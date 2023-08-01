package unleashclient

import (
	"context"
	"net/http"
)

// UnleashClient is the wrapper with authentication
type UnleashClient struct {
	Context          context.Context
	unleashEndpoint  string
	api_key          string
	ignoreCertErrors bool

	httpClient *http.Client

	// Maximum time to wait (when polling) for a task to become completed.
	retryTimeout int
}

type Project struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	FeatureLimit  int    `json:"featureLimit"`
	Mode          string `json:"mode"`
	DefStickiness string `json:"defaultStickiness"`
}

type ProjectRead struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	FeatureLimit  int      `json:"featureLimit"`
	Mode          string   `json:"mode"`
	DefStickiness string   `json:"defaultStickiness"`
	Environments  []string `json:"environments"`
}

type ProjectEnvironment struct {
	Environment           string `json:"environment"`
	ChangeRequestsEnabled bool   `json:"changeRequestsEnabled"`
}

type ProjectUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Environment struct {
	Name      string `json:"name"`
	EnvType   string `json:"type"`
	Enabled   bool   `json:"enabled"`
	SortOrder int    `json:"sortOrder"`
}

type EnvironmentUpdate struct {
	EnvType   string `json:"type"`
	SortOrder int    `json:"sortOrder"`
}

type Strategy struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  []StrategyParam `json:"parameters"`
}

type StrategyRead struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  []StrategyParam `json:"parameters"`
	Enabled     bool            `json:"deprecated"`
	Editable    bool            `json:"editable"`
}

type StrategyParam struct {
	Name        string `json:"name"`
	ParamType   string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type StrategyUpdate struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  []StrategyParam `json:"parameters"`
}
