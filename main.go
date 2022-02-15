package main

import (
    "context"

    "github.com/Burmuley/terraform-provider-sqlite/sqlite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
    provider_name := "burmuley.com/edu/sqlite"
    opts := &plugin.ServeOpts{ProviderFunc: sqlite.Provider}
    plugin.Debug(context.Background(), provider_name, opts)

    //plugin.Serve(&plugin.ServeOpts{
	//	ProviderFunc: sqlite.Provider,
	//})
}
