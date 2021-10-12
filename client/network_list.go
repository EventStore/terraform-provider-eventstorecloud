package client

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ListNetworksRequest struct {
	OrganizationID string
	ProjectID      string
}

type ListNetworksResponse struct {
	Networks []Network `json:"networks"`
}

func (c *Client) NetworkList(ctx context.Context, req *ListNetworksRequest) (*ListNetworksResponse, diag.Diagnostics) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("resources", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "networks")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, diag.Errorf("error constructing request: %w", err)
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "listing networks", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := ListNetworksResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
