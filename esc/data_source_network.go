package esc

import (
	"context"

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
	c := meta.(*providerContext)

	resp, err := c.client.NetworkList(ctx, &client.ListNetworksRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Get("project_id").(string),
	})
	if err != nil {
		return err
	}

	if len(resp.Networks) == 0 {
		return diag.Errorf("Your query returned no results. Please change " +
			"your search criteria and try again.")
	}

	var found []*client.Network
	desiredName := d.Get("name").(string)
	for _, network := range resp.Networks {
		if network.Name == desiredName {
			found = append(found, &network)
			break
		}
	}

	if len(found) == 0 {
		return diag.Errorf("Your query returned no results. Please change " +
			"your search criteria and try again.")
	}
	if len(found) > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	d.SetId(found[0].NetworkID)
	d.Set("cidr_block", found[0].CIDRBlock)
	d.Set("region", found[0].Region)
	d.Set("resource_provider", found[0].Provider)

	return nil
}
