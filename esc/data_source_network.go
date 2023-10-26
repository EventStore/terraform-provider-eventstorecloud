package esc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data for an existing `Network` resource",
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_provider": {
				Description: "Cloud Provider in which to provision the network.",
				Required:    false,
				ForceNew:    false,
				Computed:    true,
				Type:        schema.TypeString,
			},
			"region": {
				Description: "Provider region in which to provision the network",
				Required:    false,
				ForceNew:    false,
				Computed:    true,
				Type:        schema.TypeString,
			},
			"cidr_block": {
				Description: "Address space of the network in CIDR block notation",
				Required:    false,
				ForceNew:    false,
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("[WARN] Your warning message here")

	c := meta.(*providerContext)

	projectID := d.Get("project_id").(string)

	resp, err := c.client.NetworkList(ctx, &client.ListNetworksRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectID,
	})
	if err != nil {
		return err
	}

	if len(resp.Networks) == 0 {
		return diag.Errorf("There are no networks in project %s", projectID)
	}

	var found *client.Network
	multipleNetworksFound := false
	count := 0

	desiredName := d.Get("name").(string)
	for _, network := range resp.Networks {
		if network.Name == desiredName && network.Status == "available" {
			count++
			if count > 1 {
				multipleNetworksFound = true
				break
			}
			found = &network
		}
	}

	if multipleNetworksFound {
		return diag.Errorf("Error: Multiple networks with the same name '%s' were found. Please specify a more unique name or check your existing resources.", desiredName)
	}

	if found == nil {
		return diag.Errorf("Network %s was not found in project %s", desiredName, projectID)
	}

	d.SetId(found.NetworkID)
	d.Set("cidr_block", found.CIDRBlock)
	d.Set("region", found.Region)
	d.Set("resource_provider", found.Provider)

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "This resource is deprecated. Please use the eventstorecloud_networks data resource instead.",
		},
	}
}

func dataSourceNetworkList() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves data for an existing `Network` resource",
		ReadContext: dataSourceNetworkListRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Elem:     dataSourceNetwork(),
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	projectID := d.Get("project_id").(string)
	desiredName := d.Get("name").(string)

	resp, err := c.client.NetworkList(ctx, &client.ListNetworksRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectID,
	})
	if err != nil {
		return err
	}

	networkResources := make([]map[string]interface{}, 0)

	for _, network := range resp.Networks {
		if network.Name == desiredName && network.Status == "available" {
			networkResource := map[string]interface{}{
				"cidr_block":        network.CIDRBlock,
				"name":              network.Name,
				"project_id":        network.ProjectID,
				"region":            network.Region,
				"resource_provider": network.Provider,
			}
			networkResources = append(networkResources, networkResource)
		}
	}

	d.SetId(projectID)
	if err := d.Set("networks", networkResources); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
