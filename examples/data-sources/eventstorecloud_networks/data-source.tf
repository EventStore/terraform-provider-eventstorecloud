data "eventstorecloud_networks" "example" {
  name       = "Example Network"
  project_id = var.project_id
}

output "first_network" {
  value = data.eventstorecloud_networks.example.networks[0]
}
