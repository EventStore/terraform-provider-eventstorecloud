package main

import (
	"github.com/EventStore/terraform-provider-eventstorecloud/esc"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: esc.Provider,
	})
}
