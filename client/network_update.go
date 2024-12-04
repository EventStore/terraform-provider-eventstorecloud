package client

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type UpdateNetworkRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkID      string
	Name           string `json:"description"`
}

func (c *Client) NetworkUpdate(ctx context.Context, req *UpdateNetworkRequest) diag.Diagnostics {
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

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL.String(), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error constructing request: %w", err))
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error sending request: %w", err))
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "updating network", resp.Body)
	}

	return nil
}
