---
subcategory: ""
page_title: "Provision Event Store Cloud resources in GCP"
description: |-
    A sample Terraform project to provision all the Event Store Cloud resources in Google Cloud.
---

# Event Store Cloud in GCP

The sample project creates the following resources in Event Store Cloud:
- Project
- Network
- Network peering
- Managed EventStoreDB using single F1 node with 16GB disk

From the GCP side, you still need to create an incoming peering from your GCP account towards the Event Store Cloud VPC as described in the [documentation](https://developers.eventstore.com/cloud/provision/gcp/#network-peering).
This step can be also automated using the GCP Terraform provider.

```terraform
terraform {
  required_providers {
    eventstorecloud = {
      source = "EventStore/eventstorecloud"
    }
  }
}

variable "peering_route" {
  type = string
}

provider "eventstorecloud" {
}

provider "google" {
}

data "google_project" "project" {
}

data "google_compute_network" "network" {
  name = "default"
}

resource "eventstorecloud_project" "chicken_window" {
  name = "Improved Chicken Window"
}

resource "eventstorecloud_network" "chicken_window" {
  name = "Chicken Window Net"

  project_id = eventstorecloud_project.chicken_window.id

  resource_provider = "gcp"
  region            = "us-central1"
  cidr_block        = "172.29.0.0/16"
}

resource "eventstorecloud_peering" "peering" {
  name = "Example Peering"

  project_id = eventstorecloud_network.chicken_window.project_id
  network_id = eventstorecloud_network.chicken_window.id

  peer_resource_provider = eventstorecloud_network.chicken_window.resource_provider
  peer_network_region    = eventstorecloud_network.chicken_window.region

  peer_account_id = data.google_project.project.project_id
  peer_network_id = data.google_compute_network.network.name
  routes          = [var.peering_route]
}

resource "eventstorecloud_managed_cluster" "wings" {
  name = "Wings Cluster"

  project_id = eventstorecloud_network.chicken_window.project_id
  network_id = eventstorecloud_network.chicken_window.id

  topology         = "three-node-multi-zone"
  instance_type    = "F1"
  disk_size        = 16
  disk_type        = "ssd"
  server_version   = "22.6"
  projection_level = "user"
}

resource "google_compute_network_peering" "example" {
  name                 = "peering"
  network              = data.google_compute_network.default.id
  peer_network         = eventstorecloud_peering.peering.provider_metadata.gcp_network_id
  export_custom_routes = true
  import_custom_routes = true
}

output "chicken_window_id" {
  value = eventstorecloud_project.chicken_window.id
}

output "chicken_window_net" {
  value = eventstorecloud_network.chicken_window
}

output "chicken_window_peering" {
  value = eventstorecloud_peering.peering
}

output "wings_cluster_dns_name" {
  value = eventstorecloud_managed_cluster.wings.dns_name
}
```
