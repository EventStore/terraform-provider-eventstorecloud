package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"path"
)

type DeleteManagedClusterRequest struct {
	OrganizationID string
	ProjectID      string
	ClusterID      string
}

func (c *Client) ManagedClusterDelete(ctx context.Context, req *DeleteManagedClusterRequest) diag.Diagnostics {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("mesdb", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "clusters", req.ClusterID)

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
		return translateStatusCode(resp.StatusCode, "deleting managed cluster", resp.Body)
	}

	return nil
}
