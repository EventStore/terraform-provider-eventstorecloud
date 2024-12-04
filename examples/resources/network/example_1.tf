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
