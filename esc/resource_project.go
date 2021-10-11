package esc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Manages projects within an organization in Event Store Cloud",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Human-friendly name for the project",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	request := &client.CreateProjectRequest{
		OrganizationID: c.organizationId,
		Name:           d.Get("name").(string),
	}

	resp, err := c.client.ProjectCreate(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(resp.ProjectID)

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	request := &client.GetProjectRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Id(),
	}

	resp, err := c.client.ProjectGet(ctx, request)
	if err != nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("name", resp.Project.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	if d.HasChange("name") {
		request := &client.UpdateProjectRequest{
			OrganizationID: c.organizationId,
			ProjectID:      d.Id(),
			Name:           d.Get("name").(string),
		}

		if err := c.client.ProjectUpdate(ctx, request); err != nil {
			return err
		}
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	request := &client.DeleteProjectRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Id(),
	}

	return c.client.ProjectDelete(ctx, request)
}
