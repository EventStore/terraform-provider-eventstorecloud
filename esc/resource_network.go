package esc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Exists: resourceNetworkExists,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

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

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
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

	resp, err := c.client.NetworkCreate(context.Background(), request)
	if err != nil {
		return err
	}

	d.SetId(resp.NetworkID)

	return c.client.NetworkWaitForState(context.Background(), &client.WaitForNetworkStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      resp.NetworkID,
		State:          "available",
	})
}

func resourceNetworkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.GetNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
	}

	network, err := c.client.NetworkGet(context.Background(), request)
	if err != nil {
		return false, nil
	}
	if network.Network.Status == client.StateDeleted {
		return false, nil
	}

	return true, nil
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	if !d.HasChange("name") {
		return nil
	}

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.UpdateNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
		Name:           d.Get("name").(string),
	}

	return c.client.NetworkUpdate(context.Background(), request)
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.GetNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
	}

	resp, err := c.client.NetworkGet(context.Background(), request)
	if err != nil {
		return err
	}

	if err := d.Set("project_id", resp.Network.ProjectID); err != nil {
		return err
	}
	if err := d.Set("resource_provider", resp.Network.Provider); err != nil {
		return err
	}
	if err := d.Set("region", resp.Network.Region); err != nil {
		return err
	}
	if err := d.Set("cidr_block", resp.Network.CIDRBlock); err != nil {
		return err
	}
	if err := d.Set("name", resp.Network.Name); err != nil {
		return err
	}

	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	networkId := d.Id()

	request := &client.DeleteNetworkRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
	}

	if err := c.client.NetworkDelete(context.Background(), request); err != nil {
		return err
	}

	return c.client.NetworkWaitForState(context.Background(), &client.WaitForNetworkStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		NetworkID:      networkId,
		State:          "deleted",
	})
}
