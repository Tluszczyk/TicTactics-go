package cmd

import (
	"services/lib/log"

	GameManagementService "services/GameManagementService/cmd"

	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseTypes "services/DatabaseService/database/types"
	messageTypes "services/lib/types"

	"github.com/google/uuid"
)

// DoesUserAlreadyExist checks if a user already exists in the database
func DoesUserAlreadyExist(databaseService database.DatabaseService, usersTableName string, credentials messageTypes.Credentials) (bool, error) {
	log.Info("Started DoesUserAlreadyExist")

	_, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.USERNAME: credentials.Username},
			// Attributes: map[databaseTypes.FieldType]interface{}{databaseTypes.EMAIL: credentials.Email}, username must bu unique
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		return false, nil

	} else if err != nil {
		return false, err
	}

	return true, nil
}

// CreateUser saves user's credentials to the database
func CreateUser(databaseService database.DatabaseService, usersTableName string, passwordHashesTableName string, userPasswordHashMappingTable string, credentials messageTypes.Credentials) error {
	log.Info("Started CreateUser")

	log.Info("Create uid")
	// Create uid
	uid := uuid.New().String()

	log.Info("Create hash_id")
	// Create hash_id
	hashID := uuid.New().String()

	log.Info("Save password hash")
	// Save password hash
	item := databaseTypes.DatabaseItem{
		PK: map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: hashID},
		SK: map[databaseTypes.FieldType]interface{}{databaseTypes.PASSWORD_HASH: credentials.PasswordHash},
	}

	log.Info("Put password hash item in the database")
	putItemInput := databaseTypes.DatabasePutItemInput{
		TableName: passwordHashesTableName,
		Item:      item,
	}

	_, err := databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return err
	}

	log.Info("Save user password hash mapping item in the database")
	// Save user password hash mapping
	item = databaseTypes.DatabaseItem{
		PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		SK: map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: hashID},
	}

	log.Info("Put user password hash mapping item in the database")
	putItemInput = databaseTypes.DatabasePutItemInput{
		TableName: userPasswordHashMappingTable,
		Item:      item,
	}

	_, err = databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return err
	}

	log.Info("Save user")
	// Save user
	item = databaseTypes.DatabaseItem{
		PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		SK: map[databaseTypes.FieldType]interface{}{databaseTypes.USERNAME: credentials.Username},
		Attributes: map[databaseTypes.FieldType]interface{}{
			databaseTypes.EMAIL: credentials.Email,
			databaseTypes.ELO:   1000,
		},
	}

	log.Info("Put user item in the database")
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

// GetUser retrieves a user from the database
func GetUser(databaseService database.DatabaseService, usersTableName string, username string) (messageTypes.User, error) {
	log.Info("Started GetUser")

	// Get user
	getItemInput := databaseTypes.DatabaseGetItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.USERNAME: username},
		},
	}

	item, err := databaseService.GetItemFromDatabase(&getItemInput)

	if err != nil {
		return messageTypes.User{}, err
	}

	user := messageTypes.User{
		UID:      messageTypes.UserID(item.Item.PK[databaseTypes.UID].(string)),
		Username: item.Item.SK[databaseTypes.USERNAME].(string),
		Email:    item.Item.Attributes[databaseTypes.EMAIL].(string),
		Elo:      int(item.Item.Attributes[databaseTypes.ELO].(int32)),
	}

	return user, nil
}

// DeleteUser deletes a user from the database
func DeleteUser(databaseService database.DatabaseService, usersTableName string, uid messageTypes.UserID) error {
	log.Info("Started DeleteUser")

	// Delete user
	deleteItemInput := databaseTypes.DatabaseDeleteItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	}

	_, err := databaseService.DeleteItemFromDatabase(&deleteItemInput)

	if err != nil {
		return err
	}

	return nil
}

// DeletePassword user password
func DeletePassword(databaseService database.DatabaseService, userPasswordHashMappingTable string, passwordHashesTableName string, uid messageTypes.UserID) error {
	log.Info("Started DeletePassword")

	// Get hash_id
	getItemInput := databaseTypes.DatabaseGetItemInput{
		TableName: userPasswordHashMappingTable,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	}

	item, err := databaseService.GetItemFromDatabase(&getItemInput)

	if err != nil {
		return err
	}

	hashID := item.Item.SK[databaseTypes.HASH_ID].(string)

	// Delete user password hash mapping
	deleteItemInput := databaseTypes.DatabaseDeleteItemInput{
		TableName: userPasswordHashMappingTable,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	}

	_, err = databaseService.DeleteItemFromDatabase(&deleteItemInput)

	if err != nil {
		return err
	}

	// Delete password hash
	deleteItemInput = databaseTypes.DatabaseDeleteItemInput{
		TableName: passwordHashesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: hashID},
		},
	}

	_, err = databaseService.DeleteItemFromDatabase(&deleteItemInput)

	if err != nil {
		return err
	}

	return nil
}

// Leave all games
func LeaveAllGames(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, usersTableName string, uid messageTypes.UserID) error {
	log.Info("Started LeaveAllGames")
	
	err := GameManagementService.LeaveAllGames(databaseService, gamesTableName, userGameMappingTableName, usersTableName, uid)

	if err != nil {
		return err
	}

	return nil
}