package client

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

type DeleteManagedClusterRequest struct {
	OrganizationID string
	ProjectID      string
	ClusterID      string
}

func (c *Client) ManagedClusterDelete(ctx context.Context, req *DeleteManagedClusterRequest) error {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("mesdb", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "clusters", req.ClusterID)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error constructing request: %w", err)
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "deleting managed cluster", resp.Body)
	}

	return nil
}
