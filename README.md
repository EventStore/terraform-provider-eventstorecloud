# Terraform Provider for Event Store Cloud

This repository contains a [Terraform][terraform] provider for provisioning resources in [Event Store Cloud][esc].

# Installation

## Terraform 0.12

We provide binary releases for macOS, Windows and Linux via GitHub releases. In order for Terraform to find the plugin, the appropriate binary must be placed into the Terraform third-party plugin directory, the location of which varies by operating system:

- `%APPDATA%\terraform.d\plugins` on Windows
- `~/.terraform.d/plugins` on macOS or Linux

Alternatively, the binary can be placed alongside the main `terraform` binary.

On macOS and Linux, you can download the provider using the following commands:

- macOS: `curl -o ./terraform-provider-eventstorecloud.zip -L https://github.com/EventStore/terraform-provider-eventstorecloud/releases/download/v1.5.3/terraform-provider-eventstorecloud_1.5.3_darwin_amd64.zip && unzip ./terraform-provider-eventstorecloud.zip && mv ./terraform-provider-eventstorecloud_v1.5.3 ~/.terraform.d/plugins/terraform-provider-eventstorecloud_v1.5.3`
- Linux: `curl -o ./terraform-provider-eventstorecloud.zip -L https://github.com/EventStore/terraform-provider-eventstorecloud/releases/download/v1.5.3/terraform-provider-eventstorecloud_1.5.3_linux_amd64.zip && unzip ./terraform-provider-eventstorecloud.zip && mv ./terraform-provider-eventstorecloud_v1.5.3 ~/.terraform.d/plugins/terraform-provider-eventstorecloud_v1.5.3`

If you prefer to install from source, use the `make install` target in this repository. You'll need a Go 1.13+
development environment.

## Terraform 0.13+

Terraform now supports third party modules installed via the plugin registry. Add the following to your terraform module
configuration.

```
terraform {
  required_providers {
    eventstorecloud = {
      source = "EventStore/eventstorecloud"
      version = "1.5.3"
    }
  }
}
```

## Documentation

You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/EventStore/eventstorecloud/latest/docs).

# Provider Configuration

The Event Store Cloud provider must be configured with an access token, however there are several additional
options which may be useful.

Provider configuration options are:

- `token` - (`ESC_TOKEN` via the environment) - *Required* - a refresh token for Event Store Cloud. This token can be created and displayed with the esc cli tool [esc cli](https://github.com/EventStore/esc), or via the "request refresh token" button on the [Authentification Tokens page](https://console.eventstore.cloud/authentication-tokens)  in the console. The token id displayed in the cloud console is not a valid token.
- `organization_id` - (`ESC_ORG_ID` via the environment) - *Required* - the identifier of the Event Store Cloud
  organization into which to provision resources.

- `url` - (`ESC_URL` via the environment) - *Optional* - the URL of the Event Store Cloud API. This defaults
  to the public cloud instance of Event Store Cloud, but may be overridden to provision resources in another
  instance.
- `token_store` - (`ESC_TOKEN_STORE` via the environment) - *Optional* - the location on the local filesystem
  of the token cache. This is shared with the Event Store Cloud CLI.

## Contributing

The Event Store Cloud Terraform provider is released under the Mozilla Public License version 2, like most Terraform
providers. We welcome pull requests and issues! We adhere to the [Contributor Covenant][codeofconduct] Code of Conduct,
and ask that contributors familiarize themselves with it. We also have a set of [Contributing Guidelines][contributing].

[terraform]: (https://terraform.io)
[esc]: https://eventstore.com/event-store-cloud/
[codeofconduct]: https://github.com/EventStore/terraform-provider-eventstorecloud/tree/trunk/CODE-OF-CONDUCT.md
[contributing]: https://github.com/EventStore/terraform-provider-eventstorecloud/tree/trunk/CONTRIBUTING.md
