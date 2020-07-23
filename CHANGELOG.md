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
