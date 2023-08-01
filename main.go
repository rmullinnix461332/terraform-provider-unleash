package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rmullinnix461332/terraform-provider-unleash/unleash"
)

func main() {

	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: unleash.New,
			ProviderAddr: "app.terraform.io/SLUS-DCP/unleash",
		},
	)
}
