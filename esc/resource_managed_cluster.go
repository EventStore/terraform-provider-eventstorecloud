package esc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/EventStore/terraform-provider-eventstorecloud/client"
)

func resourceManagedCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagedClusterCreate,
		Exists: resourceManagedClusterExists,
		Read:   resourceManagedClusterRead,
		Delete: resourceManagedClusterDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "ID of the project in which the managed cluster exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"network_id": {
				Description: "ID of the network in which the managed cluster exists",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"name": {
				Description: "Name of the managed cluster",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
			},
			"topology": {
				Description:  "Topology of the managed cluster",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validTopologies, true),
			},
			"instance_type": {
				Description:  "Instance Type of the managed cluster",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validInstanceTypes, true),
				StateFunc: func(val interface{}) string {
					// Normalize to lower case
					return strings.ToLower(val.(string))
				},
			},
			"disk_size": {
				Description:  "Size of the data disks, in gigabytes",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(8, 4096),
			},
			"disk_type": {
				Description:  "Storage class of the data disks",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validDiskTypes, true),
				StateFunc: func(val interface{}) string {
					// Normalize to lower case
					return strings.ToLower(val.(string))
				},
			},
			"server_version": {
				Description:  "Server version to provision",
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validServerVersions, true),
				StateFunc: func(val interface{}) string {
					// Normalize to lower case
					return strings.ToLower(val.(string))
				},
			},
			"projection_level": {
				Description:  "Determines whether to run no projections, system projections only, or system and user projections",
				Optional:     true,
				ForceNew:     true,
				Default:      "off",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(validProjectionLevels, true),
				StateFunc: func(val interface{}) string {
					// Normalize to lower case
					return strings.ToLower(val.(string))
				},
			},

			"resource_provider": {
				Description: "Provider in which the cluster was created. Determined by the provider of the Network.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"region": {
				Description: "Region in which the cluster was created. Determined by the region of the Network",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_name": {
				Description: "DNS address of the cluster",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceManagedClusterCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)

	request := &client.CreateManagedClusterRequest{
		OrganizationID:  c.organizationId,
		ProjectID:       projectId,
		NetworkId:       d.Get("network_id").(string),
		Name:            d.Get("name").(string),
		Topology:        strings.ToLower(d.Get("topology").(string)),
		InstanceType:    strings.ToLower(d.Get("instance_type").(string)),
		DiskSizeGB:      int32(d.Get("disk_size").(int)),
		DiskType:        strings.ToLower(d.Get("disk_type").(string)),
		ServerVersion:   strings.ToLower(d.Get("server_version").(string)),
		ProjectionLevel: strings.ToLower(d.Get("projection_level").(string)),
	}

	resp, err := c.client.ManagedClusterCreate(context.Background(), request)
	if err != nil {
		return err
	}

	d.SetId(resp.ClusterID)

	if err := c.client.ManagedClusterWaitForState(context.Background(), &client.WaitForManagedClusterStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		ClusterID:      resp.ClusterID,
		State:          "available",
	}); err != nil {
		return err
	}

	return resourceManagedClusterRead(d, meta)
}

func resourceManagedClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	clusterId := d.Id()

	request := &client.GetManagedClusterRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		ClusterID:      clusterId,
	}

	_, err := c.client.ManagedClusterGet(context.Background(), request)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func resourceManagedClusterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	clusterId := d.Id()

	request := &client.GetManagedClusterRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		ClusterID:      clusterId,
	}

	resp, err := c.client.ManagedClusterGet(context.Background(), request)
	if err != nil {
		return err
	}

	if err := d.Set("project_id", resp.ManagedCluster.ProjectID); err != nil {
		return err
	}
	if err := d.Set("network_id", resp.ManagedCluster.NetworkID); err != nil {
		return err
	}
	if err := d.Set("name", resp.ManagedCluster.Name); err != nil {
		return err
	}
	if err := d.Set("topology", resp.ManagedCluster.Topology); err != nil {
		return err
	}
	if err := d.Set("instance_type", resp.ManagedCluster.InstanceType); err != nil {
		return err
	}
	if err := d.Set("disk_size", int(resp.ManagedCluster.DiskSizeGB)); err != nil {
		return err
	}
	if err := d.Set("disk_type", resp.ManagedCluster.DiskType); err != nil {
		return err
	}
	if err := d.Set("server_version", resp.ManagedCluster.ServerVersion); err != nil {
		return err
	}
	if err := d.Set("projection_level", resp.ManagedCluster.ProjectionLevel); err != nil {
		return err
	}

	if err := d.Set("resource_provider", resp.ManagedCluster.Provider); err != nil {
		return err
	}
	if err := d.Set("region", resp.ManagedCluster.Region); err != nil {
		return err
	}
	if err := d.Set("dns_name", fmt.Sprintf("%s.mesdb.eventstore.cloud", resp.ManagedCluster.ClusterID)); err != nil {
		return err
	}

	return nil
}

func resourceManagedClusterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*providerContext)

	projectId := d.Get("project_id").(string)
	clusterId := d.Id()

	request := &client.DeleteManagedClusterRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		ClusterID:      clusterId,
	}

	if err := c.client.ManagedClusterDelete(context.Background(), request); err != nil {
		return err
	}

	return c.client.ManagedClusterWaitForState(context.Background(), &client.WaitForManagedClusterStateRequest{
		OrganizationID: c.organizationId,
		ProjectID:      projectId,
		ClusterID:      clusterId,
		State:          "deleted",
	})
}
