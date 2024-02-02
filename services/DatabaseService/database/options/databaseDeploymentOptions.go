package database

import (
	"errors"
	"services/lib/log"
)

type DatabaseDeploymentOption int

const (
	NONE DatabaseDeploymentOption = iota
	DYNAMO
	MONGO
)

func ParseDatabaseDeploymentOption(deploymentOption string) (DatabaseDeploymentOption, error) {
	log.Info("Parsing database deployment option: " + deploymentOption)
	switch deploymentOption {
	case "DYNAMO":
		return DYNAMO, nil

	case "MONGO":
		return MONGO, nil

	default:
		return NONE, errors.ErrUnsupported
	}
}
