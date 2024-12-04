package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Peering struct {
	ProjectID               string            `json:"projectId"`
	PeeringID               string            `json:"id"`
	NetworkID               string            `json:"networkId"`
	Provider                string            `json:"provider"`
	Name                    string            `json:"description"`
	PeerAccountIdentifier   string            `json:"peerAccountId"`
	PeerNetworkIdentifier   string            `json:"peerNetworkId"`
	PeerNetworkRegion       string            `json:"peerNetworkRegion"`
	ProviderPeeringMetadata map[string]string `json:"providerPeeringMetadata"`
	Routes                  []string          `json:"routes"`
	Status                  string            `json:"status"`
	Created                 string            `json:"created,omitempty"`
}

type GetPeeringRequest struct {
	OrganizationID string
	ProjectID      string
	PeeringID      string
}

type GetPeeringResponse struct {
	Peering Peering `json:"peering"`
}

func (c *Client) PeeringGet(
	ctx context.Context,
	req *GetPeeringRequest,
) (*GetPeeringResponse, diag.Diagnostics) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join(
		"infra",
		"v1",
		"organizations",
		req.OrganizationID,
		"projects",
		req.ProjectID,
		"peerings",
		req.PeeringID,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error constructing request: %w", err))
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error sending request: %w", err))
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "getting peering", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := GetPeeringResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error parsing response: %w", err))
	}

	return &result, nil
}
