# Example for AWS

resource "eventstorecloud_project" "example" {
  name = "Example Project"
}

resource "eventstorecloud_network" "example" {
  name = "Example Network"

  project_id = eventstorecloud_project.example.id

  resource_provider = "aws"
  region            = "us-west-2"
  cidr_block        = "172.21.0.0/16"
}

resource "eventstorecloud_peering" "example" {
  name = "Peering from AWS into Example Network"

  project_id = eventstorecloud_network.example.project_id
  network_id = eventstorecloud_network.example.id

  peer_resource_provider = eventstorecloud_network.example.resource_provider
  peer_network_region    = eventstorecloud_network.example.region

  peer_account_id = "<Customer AWS Account ID>"
  peer_network_id = "<Customer VPC ID>"
  routes          = ["<Address space of the customer VPC>"]
}
