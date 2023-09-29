package main

import (
	"services/DatabaseService/database"
	databaseTypes "services/DatabaseService/database/types"
	messageTypes "services/lib/types"

	"github.com/google/uuid"
)

// CheckIfUserAlreadyExists checks if a user already exists in the database
func CheckIfUserAlreadyExists(databaseService database.DatabaseService, usersTableName string, username string) (bool, error) {
	queryResult, err := databaseService.QueryDatabase(
		usersTableName,
		databaseTypes.DatabaseQuery{
			"username": username,
		},
	)

	if err != nil {
		return false, err
	}

	return len(queryResult) > 0, nil
}

// CreateUser saves user's credentials to the database
func CreateUser(databaseService database.DatabaseService, usersTableName string, passwordHashesTableName string, credentials messageTypes.Credentials) error {
	// Create uid
	uid := uuid.New().String()

	// Create hash_id
	hashID := uuid.New().String()

	// Save password hash
	putItemRequest := databaseTypes.DatabaseItem{
		"hash_id": hashID,
		"hash":    credentials.PasswordHash,
	}
	err := databaseService.PutItemInDatabase(passwordHashesTableName, hashID, putItemRequest)

	if err != nil {
		return err
	}

	// Save user
	putItemRequest = databaseTypes.DatabaseItem{
		"uid":      uid,
		"username": credentials.Username,
		"hash_id":  hashID,
	}
	err = databaseService.PutItemInDatabase(usersTableName, uid, putItemRequest)

	if err != nil {
		return err
	}

	// TODO: what happens if we save the hash but not the user?

	return nil
}
