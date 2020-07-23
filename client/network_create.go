package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type CreateNetworkRequest struct {
	OrganizationID   string
	ProjectID        string
	ResourceProvider string `json:"provider"`
	CidrBlock        string `json:"cidrBlock"`
	Name             string `json:"description"`
	Region           string `json:"region"`
}

type CreateNetworkResponse struct {
	NetworkID string `json:"id"`
}

func (c *Client) NetworkCreate(ctx context.Context, req *CreateNetworkRequest) (*CreateNetworkResponse, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join("infra", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "networks")

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
		return nil, translateStatusCode(resp.StatusCode, "creating network", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateNetworkResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
