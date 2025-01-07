package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Network struct {
	NetworkID    string `json:"id"`
	ProjectID    string `json:"projectId"`
	Provider     string `json:"provider"`
	Region       string `json:"region"`
	CIDRBlock    string `json:"cidrBlock,omitempty"`
	Name         string `json:"description"`
	Status       string `json:"status"`
	PublicAccess bool   `json:"publicAccess"`
}

type GetNetworkRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkID      string
}

type GetNetworkResponse struct {
	Network Network `json:"network"`
}

func (c *Client) NetworkGet(
	ctx context.Context,
	req *GetNetworkRequest,
) (*GetNetworkResponse, diag.Diagnostics) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join(
		"infra",
		"v1",
		"organizations",
		req.OrganizationID,
		"projects",
		req.ProjectID,
		"networks",
		req.NetworkID,
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
		return nil, translateStatusCode(resp.StatusCode, "getting network", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := GetNetworkResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error parsing response: %w", err))
	}

	return &result, nil
}
