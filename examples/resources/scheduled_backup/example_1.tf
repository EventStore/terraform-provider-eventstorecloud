resource "eventstorecloud_scheduled_backup" "daily" {
  project_id  = eventstorecloud_project.example.id
  schedule    = "0 12 * * */1"
  description = "Creates a backup once a day at 12:00"

  source_cluster_id  = eventstorecloud_managed_cluster.example.id
  backup_description = "{cluster} Daily Backup {datetime:RFC3339}"
  max_backup_count   = 3
}
