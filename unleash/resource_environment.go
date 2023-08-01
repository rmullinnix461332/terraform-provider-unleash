package unleash

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-unleash/unleashclient"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEnvironmentType,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Update: resourceEnvironmentUpdate,
		Delete: resourceEnvironmentDelete,
		Exists: resourceEnvironmentExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnvironmentCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	var env unleashclient.Environment

	env.Name = data.Get("name").(string)
	env.EnvType = data.Get("type").(string)
	env.Enabled = data.Get("enabled").(bool)

	if err := client.CreateEnvironment(env); err != nil {
		fmt.Println("Error creating environment: ", err)
		return err
	}

	data.SetId(env.Name)

	return nil
}

func resourceEnvironmentRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	envName := data.Id()

	env, err := client.GetEnvironment(envName)
	if err != nil {
		fmt.Println("Error getting environment: ", err)
		return err
	}

	data.Set("name", env.Name)
	data.Set("type", env.EnvType)
	data.Set("enabled", env.Enabled)

	data.SetId(env.Name)

	return nil
}

func resourceEnvironmentUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	envName := data.Id()

	if data.HasChanges("type") {
		_, err := client.GetEnvironment(envName)
		if err != nil {
			fmt.Println("Error getting environment for update: ", err)
			return err
		}

		envType := data.Get("type").(string)

		if err := client.UpdateEnvironment(envName, envType); err != nil {
			fmt.Println("Error updating environment: ", err)
			return err
		}
	}

	if data.HasChanges("enabled") {
		enabled := data.Get("enabled").(bool)

		if err := client.EnableEnvironment(envName, enabled); err != nil {
			fmt.Println("Error enabling environment: ", err)
			return err
		}
	}

	data.SetId(envName)

	return resourceEnvironmentRead(data, meta)
}

func resourceEnvironmentDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	envName := data.Id()

	env, err := client.GetEnvironment(envName)
	if err != nil {
		fmt.Println("Error getting environment for delete: ", err)
		return err
	}
	if env.Enabled {
		client.EnableEnvironment(envName, false)
	}

	return client.DeleteEnvironment(envName)
}

func resourceEnvironmentExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	envName := data.Id()

	env, err := client.GetEnvironment(envName)
	if err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if env.Name == "" {
		return false, nil
	}

	return true, nil
}
