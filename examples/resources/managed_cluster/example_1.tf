# Example for AWS

data "eventstorecloud_project" "example" {
  name = "Example Project"
}

resource "eventstorecloud_network" "example" {
  name = "Example Network"

  project_id = eventstorecloud_project.example.id

  resource_provider = "aws"
  region            = "us-west-2"
  cidr_block        = "172.21.0.0/16"
}

resource "eventstorecloud_managed_cluster" "example" {
  name = "Example Cluster"

  project_id = eventstorecloud_network.example.project_id
  network_id = eventstorecloud_network.example.id

  topology        = "three-node-multi-zone"
  instance_type   = "F1"
  disk_size       = 24
  disk_type       = "gp3"
  disk_iops       = 3000
  disk_throughput = 125
  server_version  = "23.10"
}
