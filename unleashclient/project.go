package unleashclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
)

func (m *UnleashClient) CreateProject(project Project) error {
	projUrl := m.UnleashEndpoint() + "/api/admin/projects"

	buf, _ := json.Marshal(project)

	req, err := http.NewRequest("POST", projUrl, bytes.NewBufferString(string(buf)))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Unleash project creation failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) AddEnvironmentToProject(projectId string, env string) error {
	projUrl := m.UnleashEndpoint() + "/api/admin/projects/" + projectId + "/environments"

	var envItem ProjectEnvironment

	envItem.Environment = env
	envItem.ChangeRequestsEnabled = true

	buf, _ := json.Marshal(envItem)

	req, err := http.NewRequest("POST", projUrl, bytes.NewBufferString(string(buf)))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Unleash add environment to project failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) RemoveEnvironmentFromProject(projectId string, env string) error {
	projUrl := m.UnleashEndpoint() + "/api/admin/projects/" + projectId + "/environments/" + env

	req, err := http.NewRequest("DELETE", projUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return errors.New("Unleash remove environment from project failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) GetProject(projectId string) (ProjectRead, error) {
	var proj ProjectRead

	projUrl := m.UnleashEndpoint() + "/api/admin/projects/" + projectId
	req, err := http.NewRequest("GET", projUrl, nil)
	if err != nil {
		return proj, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return proj, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return proj, errors.New("Unleash project get failed: " + string(body))
	}

	err = json.Unmarshal(body, &proj)
	if err != nil {
		return proj, err
	}

	sort.Strings(proj.Environments)

	proj.ID = projectId

	return proj, nil
}

func (m *UnleashClient) UpdateProject(projectId string, projectName string, description string) error {
	projUrl := m.UnleashEndpoint() + "/api/admin/projects/" + projectId

	projUpdate := ProjectUpdate{Name: projectName, Description: description}
	buf, _ := json.Marshal(projUpdate)

	req, err := http.NewRequest("PUT", projUrl, bytes.NewBufferString(string(buf)))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Unleash project update failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) DeleteProject(projectId string) error {
	projUrl := m.UnleashEndpoint() + "/api/admin/projects/" + projectId

	req, err := http.NewRequest("DELETE", projUrl, nil)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Unleash project delete failed: " + string(body))
	}

	return nil
}
