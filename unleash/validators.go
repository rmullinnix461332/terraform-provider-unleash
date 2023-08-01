package unleash

import (
	"fmt"
	"regexp"
)

func validateProjectId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_~.-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Only alphanumeric characters or '-' allowed in %q", k))
	}
	return
}

func validateEnvironmentType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "production" && value != "development" && value != "test" && value != "preproduction" {
		errors = append(errors, fmt.Errorf("Environment Type must be one of following: development, test, preproduction or production"))
	}
	return
}
