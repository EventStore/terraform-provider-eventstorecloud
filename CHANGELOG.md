# Unreleased

ENHANCEMENTS:

* `resource/eventstorecloud_managed_cluster`: `disk_size` can now be updated to expand the disk size of a cluster. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_managed_cluster`: deleted managed clusters are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_network`: deleted networks are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: `provider_metadata` now contains `gcp_project_name` and `gcp_project_id`. `gcp_project_id` is suitable for passing to `google_compute_network_peering` resources. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: deleted peerings are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))

## 1.1.0 (July 23, 2020)

ENHANCEMENTS:

* `resource/eventstorecloud_peering`: Add computed property `provider_metadata`. ([#4](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/4))
* `resource/eventstorecloud_managed_cluster`: Correct `three-node-multi-zone` topology validation. ([#2](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/2))
* `resource/eventstorecloud_managed_cluster`: Add `projection_level` field. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))
* Meaningful error messages from the Event Store Cloud API are now presented. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))

## 1.0.0 (June 28, 2020)

FEATURES:

* Initial release of provider
* **New Data Source:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_network`
* **New Resource:** `eventstorecloud_peering`
* **New Resource:** `eventstorecloud_managedcluster`
