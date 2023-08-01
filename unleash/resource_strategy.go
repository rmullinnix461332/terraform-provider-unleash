package unleash

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-unleash/unleashclient"
)

func resourceStrategy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
		Create: resourceStrategyCreate,
		Read:   resourceStrategyRead,
		Update: resourceStrategyUpdate,
		Delete: resourceStrategyDelete,
		Exists: resourceStrategyExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStrategyCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	var strat unleashclient.Strategy

	strat.Name = data.Get("name").(string)
	strat.Description = data.Get("description").(string)
	enabled := data.Get("enabled").(bool)

	paramInt, ok := data.GetOk("parameters")
	strat.Parameters = make([]unleashclient.StrategyParam, 0)
	if ok {
		for _, item := range paramInt.([]interface{}) {
			var param unleashclient.StrategyParam

			param.Name = item.(map[string]interface{})["name"].(string)
			param.ParamType = item.(map[string]interface{})["type"].(string)
			param.Description = item.(map[string]interface{})["description"].(string)
			param.Required = item.(map[string]interface{})["required"].(bool)

			strat.Parameters = append(strat.Parameters, param)
		}
	}

	if err := client.CreateStrategy(strat); err != nil {
		fmt.Println("Error creating strategy: ", err)
		return err
	}

	data.SetId(strat.Name)

	fmt.Println("Disabd: enabled: ", enabled)
	if err := client.EnableStrategy(strat.Name, enabled); err != nil {
		fmt.Println("Error enabling strategy: ", err)
		return err
	}

	return nil
}

func resourceStrategyRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	stratName := data.Id()

	strat, err := client.GetStrategy(stratName)
	if err != nil {
		fmt.Println("Error getting strategy: ", err)
		return err
	}

	fmt.Println("Strategy: ", strat)

	data.Set("name", strat.Name)
	data.Set("description", strat.Description)
	data.Set("enabled", strat.Enabled)

	paramSet := make([]interface{}, 0)
	for _, item := range strat.Parameters {
		elements := make(map[string]interface{})
		elements["name"] = item.Name
		elements["type"] = item.ParamType
		elements["description"] = item.Description
		elements["required"] = item.Required

		paramSet = append(paramSet, elements)
	}
	data.Set("parameters", paramSet)

	data.SetId(strat.Name)

	return nil
}

func resourceStrategyUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	stratName := data.Id()

	if data.HasChanges("description", "parameters") {
		_, err := client.GetStrategy(stratName)
		if err != nil {
			fmt.Println("Error getting strategy for update: ", err)
			return err
		}

		description := data.Get("description").(string)
		paramInt, ok := data.GetOk("parameters")
		parameters := make([]unleashclient.StrategyParam, 0)
		if ok {
			for _, item := range paramInt.([]interface{}) {
				var param unleashclient.StrategyParam

				param.Name = item.(map[string]interface{})["name"].(string)
				param.ParamType = item.(map[string]interface{})["type"].(string)
				param.Description = item.(map[string]interface{})["description"].(string)
				param.Required = item.(map[string]interface{})["required"].(bool)

				parameters = append(parameters, param)
			}
		}

		if err := client.UpdateStrategy(stratName, description, parameters); err != nil {
			fmt.Println("Error updating strategy: ", err)
			return err
		}
	}

	if data.HasChanges("enabled") {
		enabled := data.Get("enabled").(bool)
		fmt.Println("HasChange: enabled: ", enabled)
		if err := client.EnableStrategy(stratName, enabled); err != nil {
			fmt.Println("Error enabling strategy: ", err)
			return err
		}
	}

	data.SetId(stratName)

	return resourceStrategyRead(data, meta)
}

func resourceStrategyDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	stratName := data.Id()

	return client.DeleteStrategy(stratName)
}

func resourceStrategyExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	stratName := data.Id()

	strat, err := client.GetStrategy(stratName)
	if err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if strat.Name == "" {
		return false, nil
	}

	return true, nil
}
