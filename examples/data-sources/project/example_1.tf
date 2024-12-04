# This assumes a project with the name "Example Project" exists
data "eventstorecloud_project" "example" {
  name = "Example Project"
}

output "project_id" {
  value = data.eventstorecloud_project.example.id
}
