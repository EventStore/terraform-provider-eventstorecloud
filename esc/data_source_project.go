package esc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*providerContext)

	resp, err := ctx.client.ProjectList(context.Background(), &client.ListProjectsRequest{
		OrganizationID: ctx.organizationId,
	})
	if err != nil {
		return err
	}

	if len(resp.Projects) == 0 {
		return fmt.Errorf("Your query returned no results. Please change " +
			"your search criteria and try again.")
	}

	found := []*client.Project{}
	desiredName := d.Get("name").(string)
	for _, project := range resp.Projects {
		if project.Name == desiredName {
			found = append(found, &project)
			break
		}
	}

	if len(found) == 0 {
		return fmt.Errorf("Your query returned no results. Please change " +
			"your search criteria and try again.")
	}
	if len(found) > 1 {
		return fmt.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	d.SetId(found[0].ProjectID)

	return nil
}
