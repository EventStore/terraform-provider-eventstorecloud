# 1.5.8 (October 12, 2021)

Enhancements:

* Import functionality has been added for all resources
* Added a data source for networks
* General documentation and example improvements
* Improved error messages for peering resources in a defunct state
* Improved error messages for project and network data sources
* Cluster names can now be updated

# 1.5.7 (September 2, 2021)

Enhancements:

* `resource/eventstorecloud_peering`: `provider_metadata` had been inadvertently changed into a set of maps. It's now back to being a map like before.

# 1.5.6 (August 4, 2021)

Enhancements:

* Allow server version "21.6" for managed clusters


# 1.5.5 (Jul 20, 2021)

Enhancements:

* Provider migrated to v2 Terraform plugin SDK
* Docs are now generated and made to fit the Terraform expected structure

# 1.5.4 (Jun 23, 2021)

Enhancements:

* Add `eventstorecloud_integration` resource

# 1.5.3 (Apr 8, 2021)

Enhancements:

* Add `eventstorecloud_scheduled_backup` resource

# 1.5.2 (Mar 26, 2021)

Fixes:

* Resources will now fail when encountering defunct state, instead of waiting for timeout

# 1.5.1 (Jan 26, 2021)

Enhancements:

* Add missing server version
* Add missing instance types
* Fix `ESC_ORG_ID` environment variable key

# 1.5.0 (Dec 11, 2020)

Enhancements:

* Publish to registry

# 1.4.0 (Dec 9, 2020)

Enhancements:

* Azure support

Fixes:

* Increase the allowed CIDR range for networks and peerings to accommodate GCP and Azure

# 1.3.0 (Nov 13, 2020)

Fixes:

* Pass the refresh token properly to ESC client 

# 1.2.0 (September 30, 2020)

Enhancements:

* `resource/eventstorecloud_managed_cluster`: `disk_size` can now be updated to expand the disk size of a cluster. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_managed_cluster`: deleted managed clusters are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_network`: deleted networks are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: `provider_metadata` now contains `gcp_project_name` and `gcp_project_id`. `gcp_project_id` is suitable for passing to `google_compute_network_peering` resources. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: deleted peerings are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))

Fixes:

* ESC Client auth token directory gets created if it does not exist

## 1.1.0 (July 23, 2020)

Enhancements:

* `resource/eventstorecloud_peering`: Add computed property `provider_metadata`. ([#4](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/4))
* `resource/eventstorecloud_managed_cluster`: Correct `three-node-multi-zone` topology validation. ([#2](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/2))
* `resource/eventstorecloud_managed_cluster`: Add `projection_level` field. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))
* Meaningful error messages from the Event Store Cloud API are now presented. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))

## 1.0.0 (June 28, 2020)

Features:

* Initial release of the provider
* **New Data Source:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_network`
* **New Resource:** `eventstorecloud_peering`
* **New Resource:** `eventstorecloud_managedcluster`
