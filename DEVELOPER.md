# Developing `terraform-provider-eventstorecloud`

Starting from Terraform v0.14 there is an [ability to override provider path](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) in development purposes:

```
provider_installation {
  dev_overrides {
    "EventStore/eventstorecloud" = "path/to/go/bin"
  }
  direct {}
}
```

This configuration can be saved to `~/.terraformrc` to have a global effect or to custom path. `TF_CLI_CONFIG_FILE` environment variable should be set according to the path to file.

After that it's just enough to build provider binary to use it

```
make install
```
