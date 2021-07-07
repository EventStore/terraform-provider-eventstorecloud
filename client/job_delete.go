package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"strings"
)

func (c *Client) DeleteJob(ctx context.Context, organizationId string, projectId string, jobId string) diag.Diagnostics {

	url := *c.apiURL
	url.Path = "/orchestrate/v1/organizations/{organizationId}/projects/{projectId}/jobs/{jobId}"
	url.Path = strings.Replace(url.Path, "{"+"organizationId"+"}", organizationId, -1)
	url.Path = strings.Replace(url.Path, "{"+"projectId"+"}", projectId, -1)
	url.Path = strings.Replace(url.Path, "{"+"jobId"+"}", jobId, -1)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return diag.Errorf("error constructing request for DeleteJob: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return diag.Errorf("error sending request for DeleteJob: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "DeleteJob", resp.Body)
	}

	return nil
}
