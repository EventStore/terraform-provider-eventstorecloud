package esc

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourcePeering() *schema.Resource {
	return &schema.Resource{
		Create: resourcePeeringCreate,
		Exists: resourcePeeringExists,
		Read:   resourcePeeringRead,
		Update: resourcePeeringUpdate,
		Delete: resourcePeeringDelete,

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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_peering_link_id": {
							Description: "AWS Peering link ID for the peering. Empty if the peering Provider is not AWS.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"gcp_project_id": {
							Description: "GCP Project ID for the peering. Empty if the peering Provider is not GCP.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"gcp_network_name": {
							Description: "GCP Network Name for the peering. Empty if the peering Provider is not GCP.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"gcp_network_id": {
							Description: "GCP Network ID in URL format. Can be passed to google_compute_network_peering resources. Empty if the peering Provider is not GCP.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourcePeeringSetProviderMetadata(d *schema.ResourceData, provider string, metadata map[string]string) error {
	providerPeeringMetadata := map[string]string{}

	switch provider {
	case "aws":
		if val, hasVal := metadata["peeringLinkId"]; hasVal {
			providerPeeringMetadata["aws_peering_link_id"] = val
		} else {
			return errors.New("AWS peering link missing remote peering link identifier")
		}
	case "gcp":
		if val, hasVal := metadata["projectId"]; hasVal {
			providerPeeringMetadata["gcp_project_id"] = val
		} else {
			return errors.New("GCP peering link missing remote peering link project identifier")
		}
		if val, hasVal := metadata["networkId"]; hasVal {
			providerPeeringMetadata["gcp_network_name"] = val
		} else {
			return errors.New("GCP peering link missing remote peering link network identifier")
		}
		providerPeeringMetadata["gcp_network_id"] = path.Join(
			"projects",
			providerPeeringMetadata["gcp_project_id"],
			"global",
			"networks",
			providerPeeringMetadata["gcp_network_name"])
	case "azure":
		break
	default:
		return fmt.Errorf("Unknown provider %q from Event Store Cloud API", provider)
	}

	return d.Set("provider_metadata", providerPeeringMetadata)
}

func resourcePeeringCreate(d *schema.ResourceData, meta interface{}) error {
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

	resp, err := c.client.PeeringCreate(context.Background(), request)
	if err != nil {
		return err
	}

	d.SetId(resp.PeeringID)

	peering, err := c.client.PeeringWaitForState(context.Background(), &client.WaitForPeeringStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      resp.PeeringID,
		State:          "initiated",
	})
	if err != nil {
		return err
	}

	return resourcePeeringSetProviderMetadata(d, peering.Provider, peering.ProviderPeeringMetadata)
}

func resourcePeeringExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	peeringId := d.Id()

	request := &client.GetPeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
	}

	peering, err := c.client.PeeringGet(context.Background(), request)
	if err != nil {
		return false, nil
	}
	if peering.Peering.Status == client.StateDeleted {
		return false, nil
	}

	return true, nil
}

func resourcePeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	if !d.HasChange("name") {
		return nil
	}

	request := &client.UpdatePeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      d.Get("project_id").(string),
		PeeringID:      d.Id(),
		Name:           d.Get("name").(string),
	}

	if err := c.client.PeeringUpdate(context.Background(), request); err != nil {
		return err
	}

	return nil
}

func resourcePeeringRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	peeringId := d.Id()

	request := &client.GetPeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
	}

	resp, err := c.client.PeeringGet(context.Background(), request)
	if err != nil {
		return err
	}

	if err := d.Set("project_id", resp.Peering.ProjectID); err != nil {
		return err
	}
	if err := d.Set("network_id", resp.Peering.NetworkID); err != nil {
		return err
	}
	if err := d.Set("peer_resource_provider", resp.Peering.Provider); err != nil {
		return err
	}
	if err := d.Set("peer_network_region", resp.Peering.PeerNetworkRegion); err != nil {
		return err
	}
	if err := d.Set("peer_account_id", resp.Peering.PeerAccountIdentifier); err != nil {
		return err
	}
	if err := d.Set("peer_network_id", resp.Peering.PeerNetworkIdentifier); err != nil {
		return err
	}
	if err := d.Set("name", resp.Peering.Name); err != nil {
		return err
	}
	if err := d.Set("routes", resp.Peering.Routes); err != nil {
		return err
	}

	return resourcePeeringSetProviderMetadata(d, resp.Peering.Provider, resp.Peering.ProviderPeeringMetadata)
}

func resourcePeeringDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	peeringId := d.Id()

	request := &client.DeletePeeringRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
	}

	if err := c.client.PeeringDelete(context.Background(), request); err != nil {
		return err
	}

	_, err := c.client.PeeringWaitForState(context.Background(), &client.WaitForPeeringStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		PeeringID:      peeringId,
		State:          "deleted",
	})
	return err
}
