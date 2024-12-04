package client

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type DeleteProjectRequest struct {
	OrganizationID string
	ProjectID      string
}

func (c *Client) ProjectDelete(ctx context.Context, req *DeleteProjectRequest) diag.Diagnostics {
	requestURL := *c.apiURL
	requestURL.Path = path.Join(
		"resources",
		"v1",
		"organizations",
		req.OrganizationID,
		"projects",
		req.ProjectID,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL.String(), nil)
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
		return translateStatusCode(resp.StatusCode, "deleting project", resp.Body)
	}

	return nil
}
