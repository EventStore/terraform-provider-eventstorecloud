package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type CreateProjectRequest struct {
	OrganizationID string
	Name           string `json:"name"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"id"`
}

func (c *Client) ProjectCreate(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join("resources", "v1", "organizations", req.OrganizationID, "projects")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error constructing request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "creating project")
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateProjectResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
