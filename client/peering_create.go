package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type CreatePeeringRequest struct {
	OrganizationID        string
	ProjectID             string
	NetworkId             string   `json:"networkId"`
	Name                  string   `json:"description"`
	PeerAccountIdentifier string   `json:"peerAccountId"`
	PeerNetworkIdentifier string   `json:"peerNetworkId"`
	PeerNetworkRegion     string   `json:"peerNetworkRegion"`
	Routes                []string `json:"routes"`
}

type CreatePeeringResponse struct {
	PeeringID string `json:"id"`
}

func (c *Client) PeeringCreate(ctx context.Context, req *CreatePeeringRequest) (*CreatePeeringResponse, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join("infra", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "peerings")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "creating peering", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreatePeeringResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
