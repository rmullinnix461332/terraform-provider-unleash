package unleashclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (m *UnleashClient) CreateStrategy(strategy Strategy) error {
	stratUrl := m.UnleashEndpoint() + "/api/admin/strategies"

	buf, _ := json.Marshal(strategy)

	req, err := http.NewRequest("POST", stratUrl, bytes.NewBufferString(string(buf)))

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
		return errors.New("Unleash strategy creation failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) GetStrategy(strategyName string) (StrategyRead, error) {
	var strat StrategyRead

	stratUrl := m.UnleashEndpoint() + "/api/admin/strategies/" + strategyName
	req, err := http.NewRequest("GET", stratUrl, nil)
	if err != nil {
		return strat, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", m.api_key)

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return strat, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return strat, errors.New("Unleash strategy get failed: " + string(body))
	}

	err = json.Unmarshal(body, &strat)
	if err != nil {
		return strat, err
	}

	strat.Enabled = !strat.Enabled

	return strat, nil
}

func (m *UnleashClient) UpdateStrategy(strategyName string, description string, parameters []StrategyParam) error {
	stratUrl := m.UnleashEndpoint() + "/api/admin/strategies/" + strategyName

	stratUpdate := StrategyUpdate{Name: strategyName, Description: description, Parameters: parameters}
	buf, _ := json.Marshal(stratUpdate)

	req, err := http.NewRequest("PUT", stratUrl, bytes.NewBufferString(string(buf)))

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
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return errors.New("Unleash strategy update failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) EnableStrategy(strategyName string, enabled bool) error {
	toggle := "deprecate"
	if enabled {
		toggle = "reactivate"
	}

	stratUrl := m.UnleashEndpoint() + "/api/admin/strategies/" + strategyName + "/" + toggle

	req, err := http.NewRequest("POST", stratUrl, nil)

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
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return errors.New("Unleash strategy update failed: " + string(body))
	}

	return nil
}

func (m *UnleashClient) DeleteStrategy(strategyName string) error {
	stratUrl := m.UnleashEndpoint() + "/api/admin/strategies/" + strategyName

	req, err := http.NewRequest("DELETE", stratUrl, nil)

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
	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return errors.New("Unleash strategy delete failed: " + string(body))
	}

	return nil
}
