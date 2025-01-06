package esc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceAcl() *schema.Resource {
	return &schema.Resource{
		Description: "Manages ACL resources in Event Store Cloud",

		CreateContext: resourceAclCreate,
		ReadContext:   resourceAclRead,
		UpdateContext: resourceAclUpdate,
		DeleteContext: resourceAclDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Project ID",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"cidr_blocks": {
				Description: "CIDR blocks allowed by this ACL",
				Required:    true,
				ForceNew:    false,
				Type:        schema.TypeList,
				Elem:         &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"name": {
				Description: "Human-friendly name for the Acl",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func mapToCidrBlock(data map[string]interface{}) client.AclCidrBlock {
	address, _ := data["address"].(string)
	comment, _ := data["comment"].(string)
	return client.AclCidrBlock{
		CidrBlock: address,
		Comment:   comment,
	}
}

func translateTfDataToCidrBlocks(data interface{}) []client.AclCidrBlock {
	maps := interfaceToMapList(data)
	result := []client.AclCidrBlock{}
	for _, m := range maps {
		result = append(result, mapToCidrBlock(m))
	}
	return result
}

func translateCidrBlocksToTf(cidrBlocks []client.AclCidrBlock) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, e := range cidrBlocks {
		result = append(result, map[string]interface{}{
			"address": e.CidrBlock,
			"comment": e.Comment,
		})
	}
	return result
}

// In terraform when we read back values they will always be of type []interface{}
// even if we passed a []map[string]interface{} originally. This takes []interface{} and builds
// a []string by casting each element individually.
func interfaceToCid(value interface{}) []client.AclCidrBlock {
	list := value.([]interface{})
	result := []client.AclCidrBlock{}
	for _, element := range list {
		result = append(result, element.(client.AclCidrBlock))
	}
	return result
}

func resourceAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)

	request := &client.CreateAclRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		CidrBlocks:     translateTfDataToCidrBlocks(d.Get("cidr_blocks")),
		Name:           d.Get("name").(string),
	}

	resp, err := c.client.AclCreate(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(resp.AclID)

	return resourceAclRead(ctx, d, meta)
}

func resourceAclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	if d.HasChange("name") || d.HasChange("cidr_blocks") {
		projectId := d.Get("project_id").(string)
		AclId := d.Id()

		request := &client.AclUpdateRequest{
			OrganizationID: c.organizationId,
			ProjectID:      projectId,
			AclID:          AclId,
			CidrBlocks:     translateTfDataToCidrBlocks(d.Get("cidr_blocks")),
			Description:    d.Get("name").(string),
		}

		err := c.client.AclUpdate(ctx, request)
		if err != nil {
			return err
		}
	}

	return resourceAclRead(ctx, d, meta)
}

func resourceAclRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	AclId := d.Id()

	request := &client.GetAclRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		AclID:          AclId,
	}

	resp, err := c.client.AclGet(ctx, request)
	if err != nil {
		return diag.Errorf("Internal Server Error, try again later")
	}
	if resp.Acl.Status == client.StateDeleted {
		d.SetId("")
		return nil
	}

	if err := d.Set("project_id", resp.Acl.ProjectID); err != nil {
		diags = append(diags, diag.Errorf("Unable to set project_id", err)...)
	}
	if err := d.Set("cidr_blocks", translateCidrBlocksToTf(resp.Acl.CidrBlocks)); err != nil {
		diags = append(diags, diag.Errorf("Unable to set cidr_blocks", err)...)
	}
	if err := d.Set("name", resp.Acl.Name); err != nil {
		diags = append(diags, diag.Errorf("Unable to set name", err)...)
	}

	return diags
}

func resourceAclDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	AclId := d.Id()

	request := &client.DeleteAclRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		AclID:          AclId,
	}

	if err := c.client.AclDelete(ctx, request); err != nil {
		return err
	}

	return diags
}
