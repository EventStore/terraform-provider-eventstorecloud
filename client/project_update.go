package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type UpdateProjectRequest struct {
	OrganizationID string
	ProjectID      string
	Name           string `json:"name"`
}

func (c *Client) ProjectUpdate(ctx context.Context, req *UpdateProjectRequest) diag.Diagnostics {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error marshalling request: %w", err))
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join(
		"resources",
		"v1",
		"organizations",
		req.OrganizationID,
		"projects",
		req.ProjectID,
	)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		requestURL.String(),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error constructing request: %w", err))
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error sending request: %w", err))
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "updating project", resp.Body)
	}

	return nil
}
