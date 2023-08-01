// Copyright (c) 2018, Google, Inc.
// Copyright (c) 2019, Noel Cower.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package unleashclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"

	"github.com/pkg/errors"
)

func (m *UnleashClient) UnleashEndpoint() string {
	if m.unleashEndpoint == "" {
		return "http://localhost:4242"
	}
	return m.unleashEndpoint
}

func (m *UnleashClient) RetryTimeout() int {
	return m.retryTimeout
}

// Create new unleash client with flag
func NewUnleashClient(unleashEndpoint string, api_key string, ignoreCertErrors bool) (*UnleashClient, error) {
	unleashClient := &UnleashClient{
		unleashEndpoint:  unleashEndpoint,
		ignoreCertErrors: ignoreCertErrors,
		api_key:          api_key,
		retryTimeout:     60,
		Context:          context.Background(),
	}

	// Api client initialization.
	err := unleashClient.InitializeHTTPClient()
	if err != nil {
		return nil, errors.New("Could not initialize http client, failing.")
	}

	err = unleashClient.Health()
	if err != nil {
		return nil, errors.New("Could not reach Unleash, please ensure it is running. Failing.")
	}

	return unleashClient, nil
}

// InitializeHTTPClient will return an *http.Client with TLS
func (m *UnleashClient) InitializeHTTPClient() error {
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:       cookieJar,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: m.ignoreCertErrors,
	}

	m.httpClient = &client

	return nil
}

type healthStatus struct {
	Health string `json:"health"`
}

func (m *UnleashClient) Health() error {
	healthUrl := m.UnleashEndpoint() + "/health"
	req, err := http.NewRequest("GET", healthUrl, nil)
	if err != nil {
		return err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New("Unleash health check failed: " + string(body))
	}

	var health healthStatus
	err = json.Unmarshal(body, &health)

	if health.Health != "GOOD" {
		return errors.New("Unleash health check failed: " + health.Health)
	}

	return nil
}
