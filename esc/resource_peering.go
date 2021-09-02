package esc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourcePeering() *schema.Resource {
	return &schema.Resource{
		Description: "Manages peering connections between Event Store Cloud VPCs and customer own VPCs",

		CreateContext: resourcePeeringCreate,
		ReadContext:   resourcePeeringRead,
		UpdateContext: resourcePeeringUpdate,
		DeleteContext: resourcePeeringDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Project ID",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"network_id": {
				Description: "Network ID",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"peer_resource_provider": {
				Description:  "Cloud Provider in which the target network exists",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validProviders, true),
			},
			"peer_network_region": {
				Description: "Provider region in which to the peer network exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"peer_account_id": {
				Description: "Account identifier in which to the peer network exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"peer_network_id": {
				Description: "Network identifier of the peer network exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "Human-friendly name for the network",
				Type:        schema.TypeString,
				Required:    true,
			},
			"routes": {
				Description: "Routes to create from the Event Store network to the peer network",
				Type:        schema.TypeSet,
				ForceNew:    true,
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDRNetwork(8, 27),
				},
				Set: schema.HashString,
			},

			"provider_metadata": {
				Description: "Metadata about the remote end of the peering connection",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourcePeeringSetProviderMetadata(d *schema.ResourceData, provider string, metadata map[string]string) diag.Diagnostics {
	providerPeeringMetadata := map[string]interface{}{}

	var diags diag.Diagnostics

	switch provider {
	case "aws":
		if val, hasVal := metadata["peeringLinkId"]; hasVal {
			providerPeeringMetadata["aws_peering_link_id"] = val
		} else {
			diags = append(diags, Warnof("AWS peering link missing remote peering link identifier")...)
		}
	case "gcp":
		projectId, hasProjectId := metadata["projectId"]
		if hasProjectId {
			providerPeeringMetadata["gcp_project_id"] = projectId
		} else {
			diags = append(diags, diag.Errorf("GCP peering link missing remote peering link project identifier")...)
		}
		networkName, hasNetworkName := metadata["networkId"]
		if hasNetworkName {
			providerPeeringMetadata["gcp_network_name"] = networkName
		} else {
			diags = append(diags, diag.Errorf("GCP peering link missing remote peering link network identifier")...)
		}
		providerPeeringMetadata["gcp_network_id"] = path.Join(
			"projects",
			projectId,
			"global",
			"networks",
			networkName)
	case "azure":
		break
	default:
		diags = append(diags, diag.Errorf("Unknown provider %q from Event Store Cloud API", provider)...)
	}

	if err := d.Set("provider_metadata", providerPeeringMetadata); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourcePeeringCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)

	routesSet := d.Get("routes").(*schema.Set)
	routes := make([]string, routesSet.Len())
	for i, route := range routesSet.List() {
		routes[i] = route.(string)
	}

	request := &client.CreatePeeringRequest{
		OrganizationID:        c.organizationId,
		ProjectID:             projectId,
		NetworkId:             d.Get("network_id").(string),
		Name:                  d.Get("name").(string),
		PeerAccountIdentifier: d.Get("peer_account_id").(string),
		PeerNetworkIdentifier: d.Get("peer_network_id").(string),
		PeerNetworkRegion:     d.Get("peer_network_region").(string),
		Routes:                routes,
	}

	resp, err := c.client.PeeringCreate(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(resp.PeeringID)

	peering, err := c.client.PeeringWaitForState(ctx, &client.WaitForPeeringStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      resp.PeeringID,
		State:          "initiated",
	})
	if err != nil {
		return err
	}

	if err := resourcePeeringSetProviderMetadata(d, peering.Provider, peering.ProviderPeeringMetadata); err != nil {
		return err
	}

	return resourcePeeringRead(ctx, d, meta)
}

func resourcePeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	if d.HasChange("name") {
		request := &client.UpdatePeeringRequest{
			OrganizationID: c.organizationId,
			ProjectID:      d.Get("project_id").(string),
			PeeringID:      d.Id(),
			Name:           d.Get("name").(string),
		}

		if err := c.client.PeeringUpdate(ctx, request); err != nil {
			return err
		}
	}

	return resourcePeeringRead(ctx, d, meta)
}

func resourcePeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	peeringId := d.Id()

	request := &client.GetPeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
	}

	resp, err := c.client.PeeringGet(ctx, request)
	if err != nil || resp.Peering.Status == client.StateDeleted {
		d.SetId("")
		return diags
	}

	if err := d.Set("project_id", resp.Peering.ProjectID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("network_id", resp.Peering.NetworkID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("peer_resource_provider", resp.Peering.Provider); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("peer_network_region", resp.Peering.PeerNetworkRegion); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("peer_account_id", resp.Peering.PeerAccountIdentifier); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("peer_network_id", resp.Peering.PeerNetworkIdentifier); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Peering.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("routes", resp.Peering.Routes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := resourcePeeringSetProviderMetadata(d, resp.Peering.Provider, resp.Peering.ProviderPeeringMetadata); err != nil {
		diags = append(diags, err...)
	}

	return diags
}

func resourcePeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	peeringId := d.Id()

	request := &client.DeletePeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
	}

	if err := c.client.PeeringDelete(ctx, request); err != nil {
		return err
	}

	peering, err := c.client.PeeringWaitForState(ctx, &client.WaitForPeeringStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
		State:          "deleted",
	})
	if peering.Status != "deleted" {
		return diag.Errorf("Peering wait for status returned, but the state is still not correct")
	}
	return err
}

func Warnof(format string, a ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf(format, a...),
		},
	}
}
