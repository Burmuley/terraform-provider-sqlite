package main

import (
	"github.com/Burmuley/terraform-provider-sqlite/sqlite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sqlite.Provider,
	})
}
