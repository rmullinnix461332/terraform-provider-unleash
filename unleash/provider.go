package unleash

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-unleash/unleashclient"
)

func New() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL for Unleash",
				DefaultFunc: schema.EnvDefaultFunc("UNLEASH_URL", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Api Key for Unleash connectivity",
				DefaultFunc: schema.EnvDefaultFunc("UNLEASH_API_KEY", nil),
			},
			"ignore_cert_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore certificate errors from unleash server",
				Default:     true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"unleash_project":     resourceProject(),
			"unleash_environment": resourceEnvironment(),
			//"unleash_project_environment": resourceProjectEnvironment(),
			"unleash_strategy": resourceStrategy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"unleash_environment": datasourceEnvironment(),
		},
		ConfigureContextFunc: providerConfigure,
	}

	return p
}

type unleashConfig struct {
	server string
	client *unleashclient.UnleashClient
}

func providerConfigure(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server := data.Get("server").(string)
	api_key := data.Get("api_key").(string)
	cert_err := data.Get("ignore_cert_errors").(bool)

	client, err := unleashclient.NewUnleashClient(server, api_key, cert_err)

	if err != nil {
		fmt.Println("config error", err)
	}
	return unleashConfig{
		server: server,
		client: client,
	}, diag.Diagnostics{}
}
