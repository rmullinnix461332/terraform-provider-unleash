package unleash

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rmullinnix461332/terraform-provider-unleash/unleashclient"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateProjectId,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"environments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressEquivalentEnvironemntList,
			},
		},
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	var proj unleashclient.Project

	proj.ID = data.Get("project_id").(string)
	proj.Name = data.Get("name").(string)
	proj.Description = data.Get("description").(string)

	if err := client.CreateProject(proj); err != nil {
		fmt.Println("Error creating project: ", err)
		return err
	}

	newProject, _ := client.GetProject(proj.ID)

	curEnvs := make(map[string]bool, 0)
	for _, env := range newProject.Environments {
		curEnvs[env] = false
	}

	newEnvs := data.Get("environments").([]interface{})
	for _, env := range newEnvs {
		if _, ok := curEnvs[env.(string)]; !ok {
			if err := client.AddEnvironmentToProject(proj.ID, env.(string)); err != nil {
				fmt.Println("Error adding environment to project: ", err)
				return err
			}
		}
		curEnvs[env.(string)] = true
	}

	for env, retained := range curEnvs {
		if !retained {
			if err := client.RemoveEnvironmentFromProject(proj.ID, env); err != nil {
				fmt.Println("Error removing environment from project: ", err)
				return err
			}
		}
	}

	data.SetId(proj.ID)

	return nil
	// return resourceProjectRead(data, meta)
}

func resourceProjectRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	projectId := data.Id()

	proj, err := client.GetProject(projectId)
	if err != nil {
		fmt.Println("Error getting project: ", err)
		return err
	}

	data.Set("project_id", proj.ID)
	data.Set("name", proj.Name)
	data.Set("description", proj.Description)

	environments := make([]string, 0)
	for _, env := range proj.Environments {
		environments = append(environments, env)
	}
	data.Set("environments", environments)

	data.SetId(proj.ID)

	return nil
}

func resourceProjectUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	projectId := data.Id()

	if data.HasChanges("name", "description") {
		_, err := client.GetProject(projectId)
		if err != nil {
			fmt.Println("Error getting project for update: ", err)
			return err
		}

		name := data.Get("name").(string)
		description := data.Get("description").(string)

		if err := client.UpdateProject(projectId, name, description); err != nil {
			fmt.Println("Error updating project: ", err)
			return err
		}
	}

	if data.HasChange("environments") {
		project, _ := client.GetProject(projectId)

		curEnvs := make(map[string]bool, 0)
		for _, env := range project.Environments {
			curEnvs[env] = false
		}

		newEnvs := data.Get("environments").([]interface{})
		for _, env := range newEnvs {
			if _, ok := curEnvs[env.(string)]; !ok {
				if err := client.AddEnvironmentToProject(projectId, env.(string)); err != nil {
					fmt.Println("Error adding environment to project: ", err)
					return err
				}
			}
			curEnvs[env.(string)] = true
		}

		for env, retained := range curEnvs {
			if !retained {
				if err := client.RemoveEnvironmentFromProject(projectId, env); err != nil {
					fmt.Println("Error removing environment from project: ", err)
					return err
				}
			}
		}
	}

	data.SetId(projectId)

	return resourceProjectRead(data, meta)
}

func resourceProjectDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	projectId := data.Id()

	return client.DeleteProject(projectId)
}

func resourceProjectExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(unleashConfig)
	client := clientConfig.client

	projectId := data.Id()

	proj, err := client.GetProject(projectId)
	if err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if proj.Name == "" {
		return false, nil
	}

	return true, nil
}

func suppressEquivalentEnvironemntList(k, old, new string, d *schema.ResourceData) bool {
	oldData, newData := d.GetChange("environments")
	if oldData == nil || newData == nil {
		return false
	}

	oldEnv := interfaceToList(oldData.([]interface{}))
	newEnv := interfaceToList(newData.([]interface{}))
	sort.Strings(oldEnv)
	sort.Strings(newEnv)

	return reflect.DeepEqual(oldEnv, newEnv)
}

func interfaceToList(intList []interface{}) []string {
	list := make([]string, 0)
	for _, item := range intList {
		list = append(list, item.(string))
	}
	return list
}
