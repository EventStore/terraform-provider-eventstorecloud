package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type Project struct {
	ProjectID      string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
	Created        string `json:"created"`
}

type GetProjectRequest struct {
	OrganizationID string
	ProjectID      string
}

type GetProjectResponse struct {
	Project Project `json:"project"`
}

func (c *Client) ProjectGet(ctx context.Context, req *GetProjectRequest) (*GetProjectResponse, error) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("resources", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error constructing request: %w", err)
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "getting project", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := GetProjectResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
