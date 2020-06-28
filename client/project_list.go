package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type ListProjectsRequest struct {
	OrganizationID string
}

type ListProjectsResponse struct {
	Projects []Project `json:"projects"`
}

func (c *Client) ProjectList(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("resources", "v1", "organizations", req.OrganizationID, "projects")

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
		return nil, translateStatusCode(resp.StatusCode, "listing projects")
	}

	decoder := json.NewDecoder(resp.Body)
	result := ListProjectsResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
