package utils

import (
	"fmt"
	"net/http"
	"os"
	"services/lib/types"
)

func GetEnvironmentVariables(variables ...string) ([]string, types.Status, error) {
	environmentVariables := []string{}

	for _, variable := range variables {
		environmentVariable := os.Getenv(variable)

		if environmentVariable == "" {
			return []string{}, types.Status{
				Code: http.StatusInternalServerError,
				Message:    fmt.Sprintf("Environment variable not found: %s", variable),
			}, fmt.Errorf("environment variable %s not found", variable)
		}

		environmentVariables = append(environmentVariables, environmentVariable)
	}

	return environmentVariables, types.Status{
		Code: http.StatusOK,
		Message:    "Environment variables found",
	}, nil
}