package esc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Manages VPC (network) resources in Event Store Cloud",

		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Project ID",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"resource_provider": {
				Description:  "Cloud Provider in which to provision the network.",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validProviders, true),
				StateFunc: func(val interface{}) string {
					// Normalize to lower case
					return strings.ToLower(val.(string))
				},
			},
			"region": {
				Description: "Provider region in which to provision the network",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"cidr_block": {
				Description:  "Address space of the network in CIDR block notation",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.IsCIDRNetwork(8, 24),
			},
			"name": {
				Description: "Human-friendly name for the network",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)

	request := &client.CreateNetworkRequest{
		OrganizationID:   c.organizationId,
		ProjectID:        projectId,
		ResourceProvider: d.Get("resource_provider").(string),
		CidrBlock:        d.Get("cidr_block").(string),
		Name:             d.Get("name").(string),
		Region:           d.Get("region").(string),
	}

	resp, err := c.client.NetworkCreate(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(resp.NetworkID)

	if err := c.client.NetworkWaitForState(ctx, &client.WaitForNetworkStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      resp.NetworkID,
		State:          "available",
	}); err != nil {
		return err
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	if d.HasChange("name") {
		projectId := d.Get("project_id").(string)
		networkId := d.Id()

		request := &client.UpdateNetworkRequest{
			OrganizationID: c.organizationId,
			ProjectID:      projectId,
			NetworkID:      networkId,
			Name:           d.Get("name").(string),
		}

		err := c.client.NetworkUpdate(ctx, request)
		if err != nil {
			return err
		}
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.GetNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
	}

	resp, err := c.client.NetworkGet(ctx, request)
	if err != nil || resp.Network.Status == client.StateDeleted {
		d.SetId("")
		return diags
	}

	if err := d.Set("project_id", resp.Network.ProjectID); err != nil {
		diags = append(diags, diag.Errorf("Unable to set project_id", err)...)
	}
	if err := d.Set("resource_provider", resp.Network.Provider); err != nil {
		diags = append(diags, diag.Errorf("Unable to set resource_provider", err)...)
	}
	if err := d.Set("region", resp.Network.Region); err != nil {
		diags = append(diags, diag.Errorf("Unable to set region", err)...)
	}
	if err := d.Set("cidr_block", resp.Network.CIDRBlock); err != nil {
		diags = append(diags, diag.Errorf("Unable to set cidr_block", err)...)
	}
	if err := d.Set("name", resp.Network.Name); err != nil {
		diags = append(diags, diag.Errorf("Unable to set name", err)...)
	}

	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.DeleteNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
	}

	if err := c.client.NetworkDelete(ctx, request); err != nil {
		return err
	}

	if err := c.client.NetworkWaitForState(ctx, &client.WaitForNetworkStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
		State:          "deleted",
	}); err != nil {
		return err
	}

	return diags
}
