package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"path"
)

type DeleteNetworkRequest struct {
	OrganizationID string
	ProjectID      string
	NetworkID      string
}

func (c *Client) NetworkDelete(ctx context.Context, req *DeleteNetworkRequest) diag.Diagnostics {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("infra", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "networks", req.NetworkID)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL.String(), nil)
	if err != nil {
		return diag.Errorf("error constructing request: %w", err)
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return diag.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "deleting network", resp.Body)
	}

	return nil
}
