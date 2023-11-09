package mongo

import (
	"services/DatabaseService/database/types"
	"testing"
)

func TestGetItemFromDatabase(t *testing.T) {

	// Create a new mongoDBService
	mongoDBService, err := NewMongoDatabaseService()
	if err != nil {
		t.Error(err)
	}

	// Create a new item input
	getItemInput := &types.DatabaseGetItemInput{
		TableName: "users",
		Key: types.DatabaseItem{
			PK:         map[types.FieldType]interface{}{types.USERNAME: "test"},
			SK:         map[types.FieldType]interface{}{types.EMAIL: "test@test"},
			Attributes: map[types.FieldType]interface{}{types.ELO: 1000},
		},
	}

	// Try to get a record
	getItemOutput, err := mongoDBService.GetItemFromDatabase(getItemInput)
	if err != nil {
		t.Error(err)
	}

	// Print the result
	t.Log(getItemOutput)

	// Disconnect from the database
	err = mongoDBService.Disconnect()
	if err != nil {
		t.Error(err)
	}
}

func TestPutItemToDatabase(t *testing.T) {
	
	// Create a new mongoDBService
	mongoDBService, err := NewMongoDatabaseService()
	if err != nil {
		t.Error(err)
	}

	// Create a new item input
	putItemInput := &types.DatabasePutItemInput{
		TableName: "users",
		Item: types.DatabaseItem{
			PK:         map[types.FieldType]interface{}{types.USERNAME: "testput"},
			SK:         map[types.FieldType]interface{}{types.EMAIL: "test@test"},
			Attributes: map[types.FieldType]interface{}{types.ELO: 1000},
		},
	}

	// Try to put a record
	_, err = mongoDBService.PutItemInDatabase(putItemInput)
	if err != nil {
		t.Error(err)
	}

	// Verify that the record was put
	getItemInput := &types.DatabaseGetItemInput{
		TableName: "users",
		Key: types.DatabaseItem{
			PK:         map[types.FieldType]interface{}{types.USERNAME: "testput"},
			SK:         map[types.FieldType]interface{}{types.EMAIL: "test@test"},
			Attributes: map[types.FieldType]interface{}{types.ELO: 1000},
		},
	}

	// Try to get a record
	getItemOutput, err := mongoDBService.GetItemFromDatabase(getItemInput)
	if err != nil {
		t.Error(err)
	}

	// Print the result
	t.Log(getItemOutput)

	// Disconnect from the database
	err = mongoDBService.Disconnect()
	if err != nil {
		t.Error(err)
	}
}