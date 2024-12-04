data "eventstorecloud_network" "example" {
  name       = "Example Network"
  project_id = var.project_id
}

output "network_cidr" {
  value = data.eventstorecloud_network.example.cidr_block
}
