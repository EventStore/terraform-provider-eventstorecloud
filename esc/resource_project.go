package esc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Exists: resourceProjectExists,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Human-friendly name for the project",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	request := &client.CreateProjectRequest{
		OrganizationID: c.organizationId,
		Name:           d.Get("name").(string),
	}

	resp, err := c.client.ProjectCreate(context.Background(), request)
	if err != nil {
		return err
	}

	d.SetId(resp.ProjectID)
	return nil
}

func resourceProjectExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*providerContext)

	request := &client.GetProjectRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Id(),
	}

	_, err := c.client.ProjectGet(context.Background(), request)
	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	request := &client.GetProjectRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Id(),
	}

	resp, err := c.client.ProjectGet(context.Background(), request)
	if err != nil {
		return err
	}

	if err := d.Set("name", resp.Project.Name); err != nil {
		return err
	}

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	if d.HasChange("name") {
		request := &client.UpdateProjectRequest{
			OrganizationID: c.organizationId,
			ProjectID:      d.Id(),
			Name:           d.Get("name").(string),
		}

		if err := c.client.ProjectUpdate(context.Background(), request); err != nil {
			return err
		}
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	request := &client.DeleteProjectRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Id(),
	}

	return c.client.ProjectDelete(context.Background(), request)
}
