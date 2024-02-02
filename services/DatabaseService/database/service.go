package database

import (
	"errors"
	"services/DatabaseService/database/dynamo"
	"services/DatabaseService/database/mongo"
	options "services/DatabaseService/database/options"
	"services/DatabaseService/database/types"
)

type DatabaseService interface {
	// GetItemFromDatabase returns an item from the database with the given key
	GetItemFromDatabase(*types.DatabaseGetItemInput) (*types.DatabaseGetItemOutput, error)

	// PutItemInDatabase puts an item in the database with the given key
	PutItemInDatabase(*types.DatabasePutItemInput) (*types.DatabasePutItemOutput, error)

	// DeleteItemFromDatabase deletes an item from the database with the given key
	DeleteItemFromDatabase(*types.DatabaseDeleteItemInput) (*types.DatabaseDeleteItemOutput, error)

	// UpdateItemInDatabase updates an item in the database with the given key
	UpdateItemInDatabase(*types.DatabaseUpdateItemInput) (*types.DatabaseUpdateItemOutput, error)

	// QueryDatabase queries the database with the given query
	QueryDatabase(*types.DatabaseQueryInput) (*types.DatabaseQueryOutput, error)

	// MarshallDatabaseItem marshalls a DatabaseItem into a database-specific item
	MarshallDatabaseItem(*types.DatabaseItem) (interface{}, error)

	// UnmarshallDatabaseItem unmarshalls a database-specific item into a DatabaseItem
	UnmarshallDatabaseItem(interface{}, *types.DatabaseItem) error
}

// GetDatabaseService returns a DatabaseService based on the given deployment option
func GetDatabaseService(deploymentOption options.DatabaseDeploymentOption) (DatabaseService, error) {
	switch deploymentOption {
	case options.DYNAMO:
		return dynamo.NewDynamoDatabaseService()
	case options.MONGO:
		return mongo.NewMongoDatabaseService()
	default:
		return nil, errors.ErrUnsupported
	}
}
