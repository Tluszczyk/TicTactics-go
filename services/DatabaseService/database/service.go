package database

import (
	"errors"
	"services/DatabaseService/database/dynamo"
	"services/DatabaseService/database/types"
)

type DatabaseService interface {
	// GetItemFromDatabase returns an item from the database with the given key
	GetItemFromDatabase(tableName string, key string) (types.DatabaseItem, error)

	// PutItemInDatabase puts an item in the database with the given key
	PutItemInDatabase(tableName string, key string, item types.DatabaseItem) error

	// DeleteItemFromDatabase deletes an item from the database with the given key
	DeleteItemFromDatabase(tableName string, key string) error

	// UpdateItemInDatabase updates an item in the database with the given key
	UpdateItemInDatabase(tableName string, key string, item types.DatabaseItem) error

	// QueryDatabase queries the database with the given query
	QueryDatabase(tableName string, query types.DatabaseQuery) (types.DatabaseQueryResult, error)
}

// GetDatabaseService returns a DatabaseService based on the given deployment option
func GetDatabaseService(deploymentOption DatabaseDeploymentOption) (DatabaseService, error) {
	switch deploymentOption {
	case DYNAMO:
		return dynamo.NewDynamoDatabaseService()
	default:
		return nil, errors.ErrUnsupported
	}
}
