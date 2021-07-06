package esc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	resp, err := c.client.ProjectList(ctx, &client.ListProjectsRequest{
		OrganizationID: c.organizationId,
	})
	if err != nil {
		return err
	}

	if len(resp.Projects) == 0 {
		return diag.Errorf("Your query returned no results. Please change " +
			"your search criteria and try again.")
	}

	var found []*client.Project
	desiredName := d.Get("name").(string)
	for _, project := range resp.Projects {
		if project.Name == desiredName {
			found = append(found, &project)
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

	d.SetId(found[0].ProjectID)

	return nil
}
