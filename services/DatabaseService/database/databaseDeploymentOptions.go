package database

import "errors"

type DatabaseDeploymentOption int

const (
	NONE DatabaseDeploymentOption = iota
	DYNAMO
)

func ParseDatabaseDeploymentOption(deploymentOption string) (DatabaseDeploymentOption, error) {
	switch deploymentOption {
	case "DYNAMO":
		return DYNAMO, nil

	default:
		return NONE, errors.ErrUnsupported
	}
}
