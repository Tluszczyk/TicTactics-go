package main

import (
	"services/DatabaseService/database"
	databaseTypes "services/DatabaseService/database/types"
	messageTypes "services/lib/types"

	"github.com/google/uuid"
)

// CheckIfUserAlreadyExists checks if a user already exists in the database
func CheckIfUserAlreadyExists(databaseService database.DatabaseService, usersTableName string, credentials messageTypes.Credentials) (bool, error) {
	item, err := databaseService.GetItemFromDatabase(
		&databaseTypes.DatabaseGetItemInput{
			TableName: usersTableName,
			Key: databaseTypes.DatabaseItem{
				PK: map[databaseTypes.FieldType]interface{}{
					databaseTypes.EMAIL: credentials.Email,
				},
			},
		},
	)

	if err != nil {
		return false, err
	}

	return !item.Item.IsNil(), nil
}

// CreateUser saves user's credentials to the database
func CreateUser(databaseService database.DatabaseService, usersTableName string, passwordHashesTableName string, credentials messageTypes.Credentials) error {
	// Create uid
	uid := uuid.New().String()

	// Create hash_id
	hashID := uuid.New().String()

	// Save password hash
	item := databaseTypes.DatabaseItem{
		PK:         map[databaseTypes.FieldType]interface{}{databaseTypes.EMAIL: credentials.Email},
		SK:         map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: hashID},
		Attributes: map[databaseTypes.FieldType]interface{}{databaseTypes.PASSWORD_HASH: credentials.PasswordHash},
	}

	putItemInput := databaseTypes.DatabasePutItemInput{
		TableName: passwordHashesTableName,
		Item:      item,
	}

	_, err := databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return err
	}

	// Save user
	item = databaseTypes.DatabaseItem{
		PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		SK: map[databaseTypes.FieldType]interface{}{databaseTypes.USERNAME: credentials.Username},
		Attributes: map[databaseTypes.FieldType]interface{}{
			databaseTypes.EMAIL: credentials.Email,
			databaseTypes.ELO:   "1000",
		},
	}

	putItemInput = databaseTypes.DatabasePutItemInput{
		TableName: usersTableName,
		Item:      item,
	}

	_, err = databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return err
	}

	// TODO: what happens if we save the hash but not the user?

	return nil
}
