package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type CreateIntegrationData struct {
	CreateOpsGenieIntegrationData *CreateOpsGenieIntegrationData
	CreateSlackIntegrationData    *CreateSlackIntegrationData
}

type CreateIntegrationRequest struct {
	Data        map[string]interface{} `json:"data"`
	Description string                 `json:"description"`
}

type CreateIntegrationResponse struct {
	Id string `json:"id"`
}

type CreateOpsGenieIntegrationData struct {
	// API key used with the Ops Genie integration API
	ApiKey string `json:"apiKey"`
	// Required. Must be set to \"opsGenie\"
	Sink   string  `json:"sink"`
	Source *string `json:"source,omitempty"`
}

type CreateSlackIntegrationData struct {
	// Slack Channel to send messages to
	ChannelId string `json:"channelId"`
	// API token for the Slack bot
	Token string `json:"token"`
	// Required. Must be set to \"slack\"
	Sink   string  `json:"sink"`
	Source *string `json:"source,omitempty"`
}

func (c *Client) CreateIntegration(
	ctx context.Context,
	organizationId string,
	projectId string,
	createIntegrationRequest CreateIntegrationRequest,
) (*CreateIntegrationResponse, diag.Diagnostics) {
	requestBody, err := json.Marshal(createIntegrationRequest)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error marshalling request: %w", err))
	}

	url := *c.apiURL
	url.Path = "/integrate/v1/organizations/{organizationId}/projects/{projectId}/integrations"
	url.Path = strings.Replace(url.Path, "{"+"organizationId"+"}", organizationId, -1)
	url.Path = strings.Replace(url.Path, "{"+"projectId"+"}", projectId, -1)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url.String(),
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, diag.FromErr(
			fmt.Errorf("error constructing request for CreateIntegration: %w", err),
		)
	}
	request.Header.Add("Content-Type", "application/json")
	if err := c.addAuthorizationHeader(request); err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error sending request for CreateIntegration: %w", err))
	}
	defer closeIgnoreError(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, translateStatusCode(resp.StatusCode, "CreateIntegration", resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	result := CreateIntegrationResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, diag.FromErr(fmt.Errorf("error parsing response: %w", err))
	}

	return &result, nil
}
