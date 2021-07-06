package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"strings"
)

func (c *Client) DeleteIntegration(ctx context.Context, organizationId string, projectId string, integrationId string) diag.Diagnostics {

	url := *c.apiURL
	url.Path = "/integrate/v1/organizations/{organizationId}/projects/{projectId}/integrations/{integrationId}"
	url.Path = strings.Replace(url.Path, "{"+"organizationId"+"}", organizationId, -1)
	url.Path = strings.Replace(url.Path, "{"+"projectId"+"}", projectId, -1)
	url.Path = strings.Replace(url.Path, "{"+"integrationId"+"}", integrationId, -1)

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return diag.Errorf("error constructing request for DeleteIntegration: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return diag.Errorf("error sending request for DeleteIntegration: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return translateStatusCode(resp.StatusCode, "DeleteIntegration", resp.Body)
	}

	return nil
}
