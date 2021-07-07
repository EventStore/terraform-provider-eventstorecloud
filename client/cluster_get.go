package client

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"path"
)

type ManagedCluster struct {
	OrganizationID  string `json:"organizationId"`
	ProjectID       string `json:"projectId"`
	NetworkID       string `json:"networkId"`
	ClusterID       string `json:"id"`
	Name            string `json:"description"`
	Provider        string `json:"provider"`
	Region          string `json:"region"`
	Topology        string `json:"topology"`
	InstanceType    string `json:"instanceType"`
	DiskSizeGB      int32  `json:"diskSizeGb"`
	DiskType        string `json:"diskType"`
	ServerVersion   string `json:"serverVersion"`
	ProjectionLevel string `json:"projectionLevel"`
	Status          string `json:"status"`
	Created         string `json:"created"`
}

type GetManagedClusterRequest struct {
	OrganizationID string
	ProjectID      string
	ClusterID      string
}

type GetManagedClusterResponse struct {
	ManagedCluster ManagedCluster `json:"cluster"`
}

func (c *Client) ManagedClusterGet(ctx context.Context, req *GetManagedClusterRequest) (*GetManagedClusterResponse, diag.Diagnostics) {
	requestURL := *c.apiURL
	requestURL.Path = path.Join("mesdb", "v1", "organizations", req.OrganizationID, "projects", req.ProjectID, "clusters", req.ClusterID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, diag.Errorf("error constructing request: %w", err)
	}
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.Errorf("error sending request: %w", err)
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, translateStatusCode(resp.StatusCode, "getting managed cluster", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := GetManagedClusterResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}
