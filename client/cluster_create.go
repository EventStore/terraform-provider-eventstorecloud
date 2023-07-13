package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type CreateManagedClusterRequest struct {
	OrganizationID  string
	ProjectID       string
	NetworkId       string `json:"networkId"`
	Name            string `json:"description"`
	Topology        string `json:"topology"`
	InstanceType    string `json:"instanceType"`
	DiskSizeGB      int32  `json:"diskSizeGb"`
	DiskType        string `json:"diskType"`
	DiskIops        int32  `json:"diskIops"`
	DiskThroughput  int32  `json:"diskThroughput"`
	ServerVersion   string `json:"serverVersion"`
	ProjectionLevel string `json:"projectionLevel"`
	CloudAuth       bool   `json:"cloudIntegratedAuthentication"`
	Protected       bool   `json:"protected"`
}

type CreateManagedClusterResponse struct {
	ClusterID string `json:"id"`
}

func (c *Client) ManagedClusterCreate(ctx context.Context, req *CreateManagedClusterRequest) (*CreateManagedClusterResponse, diag.Diagnostics) {
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, diag.Errorf("error marshalling request: %w", err)
	}

	requestURL := *c.apiURL
	requestURL.Path = path.Join("mesdb", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "clusters")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return nil, diag.Errorf("error constructing request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "creating managed cluster", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateManagedClusterResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
