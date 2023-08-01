package unleash

import (
	"testing"
)

func TestValidateProjectId(t *testing.T) {
	validNames := []string{
		"ValidName",
		"validname",
		"valid-name",
	}
	for _, v := range validNames {
		_, errors := validateProjectId(v, "application")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Application name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid:id",
		"invalid id",
		"invalid/name",
		"",
	}
	for _, v := range invalidNames {
		_, errors := validateProjectId(v, "application")
		if len(errors) == 0 {
			t.Fatalf("%q should be a valid Application name", v)
		}
	}
}
