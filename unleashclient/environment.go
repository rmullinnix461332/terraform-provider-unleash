package unleashclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (m *UnleashClient) CreateEnvironment(environment Environment) error {
	envUrl := m.UnleashEndpoint() + "/api/admin/environments"

	buf, _ := json.Marshal(environment)

	req, err := http.NewRequest("POST", envUrl, bytes.NewBufferString(string(buf)))

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
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return errors.New("Unleash environment creation failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) GetEnvironment(environmentName string) (Environment, error) {
	var env Environment

	envUrl := m.UnleashEndpoint() + "/api/admin/environments/" + environmentName
	req, err := http.NewRequest("GET", envUrl, nil)
	if err != nil {
		return env, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return env, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return env, errors.New("Unleash environment get failed: " + string(body))
	}

	err = json.Unmarshal(body, &env)
	if err != nil {
		return env, err
	}

	return env, nil
}

func (m *UnleashClient) UpdateEnvironment(environmentName, envType string) error {
	envUrl := m.UnleashEndpoint() + "/api/admin/environments/update/" + environmentName

	envUpdate := EnvironmentUpdate{EnvType: envType, SortOrder: 9999}
	buf, _ := json.Marshal(envUpdate)

	req, err := http.NewRequest("PUT", envUrl, bytes.NewBufferString(string(buf)))

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
	if resp.StatusCode != 204 {
		return errors.New("Unleash environment update failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) EnableEnvironment(environmentName string, enabled bool) error {
	toggle := "off"
	if enabled {
		toggle = "on"
	}

	envUrl := m.UnleashEndpoint() + "/api/admin/environments/" + environmentName + "/" + toggle

	req, err := http.NewRequest("POST", envUrl, nil)

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
	if resp.StatusCode != 204 {
		return errors.New("Unleash environment update failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) DeleteEnvironment(environmentName string) error {
	envUrl := m.UnleashEndpoint() + "/api/admin/environments/" + environmentName

	req, err := http.NewRequest("DELETE", envUrl, nil)

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
	if resp.StatusCode != 204 {
		return errors.New("Unleash environment delete failed: " + string(body))
	}

	return nil
}
