package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Job struct {
	Data           map[string]interface{} `json:"data"`
	Description    string                 `json:"description"`
	Id             string                 `json:"id"`
	OrganizationId string                 `json:"organizationId"`
	ProjectId      string                 `json:"projectId"`
	Schedule       string                 `json:"schedule"`
	Status         string                 `json:"status"`
	Type           string                 `json:"type"`
}

type GetJobResponse struct {
	Job Job `json:"job"`
}

func (c *Client) GetJob(
	ctx context.Context,
	organizationId string,
	projectId string,
	jobId string,
) (*GetJobResponse, diag.Diagnostics) {
	url := *c.apiURL
	url.Path = "/orchestrate/v1/organizations/{organizationId}/projects/{projectId}/jobs/{jobId}"
	url.Path = strings.Replace(url.Path, "{"+"organizationId"+"}", organizationId, -1)
	url.Path = strings.Replace(url.Path, "{"+"projectId"+"}", projectId, -1)
	url.Path = strings.Replace(url.Path, "{"+"jobId"+"}", jobId, -1)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error constructing request for GetJob: %w", err))
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error sending request for GetJob: %w", err))
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, translateStatusCode(resp.StatusCode, "GetJob", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := GetJobResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error parsing response: %w", err))
	}

	return &result, nil
}
