---
subcategory: ""
page_title: "Provision Event Store Cloud resources in AWS"
description: |-
    A sample Terraform project to provision all the Event Store Cloud resources in AWS.
---

# Event Store Cloud in AWS

The sample project creates the following resources in Event Store Cloud:
- Project
- Network
- Network peering
- Managed EventStoreDB using single F1 node with 16GB disk

From the AWS side, you still need to accept the peering request and configure the route as described in the [documentation](https://developers.eventstore.com/cloud/provision/aws/#network-peering).
This step can be also automated using the AWS Terraform provider.

```terraform
terraform {
  required_providers {
    eventstorecloud = {
      source = "EventStore/eventstorecloud"
    }
  }
}

provider "aws" {

}

provider "eventstorecloud" {
}

data "aws_caller_identity" "example" {
}

resource "aws_vpc" "example" {
  cidr_block = "172.250.0.0/24"

  tags = {
    Name = "eventstore-example"
  }
}

resource "eventstorecloud_project" "chicken_window" {
  name = "Improved Chicken Window"
}

resource "eventstorecloud_network" "chicken_window" {
  name = "Chicken Window Net"

  project_id = eventstorecloud_project.chicken_window.id

  resource_provider = "aws"
  region            = "us-west-2"
  cidr_block        = "172.21.0.0/16"
}

resource "eventstorecloud_peering" "peering" {
  name = "Example Peering"

  project_id = eventstorecloud_network.chicken_window.project_id
  network_id = eventstorecloud_network.chicken_window.id

  peer_resource_provider = eventstorecloud_network.chicken_window.resource_provider
  peer_network_region    = eventstorecloud_network.chicken_window.region

  peer_account_id = data.aws_caller_identity.example.account_id
  peer_network_id = aws_vpc.example.id
  routes          = [aws_vpc.example.cidr_block]
}

resource "aws_vpc_peering_connection_accepter" "peer" {
  vpc_peering_connection_id = eventstorecloud_peering.peering.provider_metadata.aws_peering_link_id
  auto_accept               = true

  tags = {
    Side   = "Accepter"
    Source = "Event Store"
  }
}

resource "aws_route" "peering" {
  route_table_id            = aws_vpc.example.main_route_table_id
  destination_cidr_block    = eventstorecloud_network.chicken_window.cidr_block
  vpc_peering_connection_id = aws_vpc_peering_connection_accepter.peer.id
}

resource "eventstorecloud_managed_cluster" "wings" {
  name = "Wings Cluster"

  project_id = eventstorecloud_network.chicken_window.project_id
  network_id = eventstorecloud_network.chicken_window.id

  topology        = "single-node"
  instance_type   = "F1"
  disk_size       = 16
  disk_type       = "gp3"
  disk_iops       = 3000
  disk_throughput = 125
  server_version  = "23.10"
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
