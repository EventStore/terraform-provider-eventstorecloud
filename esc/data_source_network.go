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

	return nil
}
