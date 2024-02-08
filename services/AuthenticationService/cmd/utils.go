package cmd

import (
	"fmt"
	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseTypes "services/DatabaseService/database/types"
	"services/lib/log"
	"services/lib/types"
	"time"

	"github.com/google/uuid"
)

func DoesUserAlreadyHaveSession(databaseService database.DatabaseService, sessionsTableName string, userSessionMappingTableName string, uid types.UserID) (bool, error) {
	log.Info("Started DoesUserAlreadyHaveSession")

	_, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userSessionMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		return false, nil

	} else if err != nil {
		return false, err
	}

	return true, nil
}

func ValidateSession(databaseService database.DatabaseService, sessionsTableName string, userSessionMappingTableName string, session types.Session) (bool, error) {
	log.Info("Started ValidateSession")

	userSessionMappingGetItemOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userSessionMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: session.UID},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info("Session not found in user session mapping")
		return false, nil

	} else if err != nil {
		return false, err
	}

	log.Info(fmt.Sprintf("%v", *userSessionMappingGetItemOutput))

	sid := userSessionMappingGetItemOutput.Item.SK[databaseTypes.SID].(string)

	_, err = databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: sessionsTableName,
		Key: databaseTypes.DatabaseItem{
			PK:         map[databaseTypes.FieldType]interface{}{databaseTypes.SID: sid},
			Attributes: map[databaseTypes.FieldType]interface{}{databaseTypes.TOKEN: session.Token},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info("Session not found in sessions")
		return false, nil

	} else if err != nil {
		return false, err
	}

	return true, nil
}

func ValidateCredentials(databaseService database.DatabaseService, usersTableName string, passwordHashesTableName string, userPasswordHashMappingTableName string, credentials types.Credentials) (types.User, error) {
	log.Info("Started ValidateCredentials")

	userGetItemOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.USERNAME: credentials.Username},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info("User not found")
		return types.User{}, nil

	} else if err != nil {
		return types.User{}, err
	}

	userPasswordHashMappingGetItemOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userPasswordHashMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: userGetItemOutput.Item.PK[databaseTypes.UID]},
		},
	})

	if err != nil {
		return types.User{}, err
	}

	passwordHashGetItemOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: passwordHashesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: userPasswordHashMappingGetItemOutput.Item.SK[databaseTypes.HASH_ID]},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info("Password hash not found")
		return types.User{}, nil

	} else if err != nil {
		return types.User{}, err
	}

	log.Info("Compare password hashes")
	if credentials.PasswordHash != passwordHashGetItemOutput.Item.SK[databaseTypes.PASSWORD_HASH].(string) {
		log.Info("Password hashes do not match")
		return types.User{}, nil
	}

	log.Info(fmt.Sprintf("User: %v", *userGetItemOutput))

	log.Info("Password hashes match")
	return types.User{
		UID:      types.UserID(userGetItemOutput.Item.PK[databaseTypes.UID].(string)),
		Username: userGetItemOutput.Item.SK[databaseTypes.USERNAME].(string),
		Email:    userGetItemOutput.Item.Attributes[databaseTypes.EMAIL].(string),
		Elo:      int(userGetItemOutput.Item.Attributes[databaseTypes.ELO].(int32)),
	}, nil
}

func CreateSession(databaseService database.DatabaseService, sessionsTableName string, userSessionMappingTableName string, uid types.UserID) (types.Session, error) {
	log.Info("Started CreateSession")

	// Create session
	sid := types.UserID(uuid.New().String())
	token := types.Token(uuid.New().String())

	// Save session
	putItemInput := databaseTypes.DatabasePutItemInput{
		TableName: sessionsTableName,
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.SID: sid},
			Attributes: map[databaseTypes.FieldType]interface{}{
				databaseTypes.TOKEN:      token,
				databaseTypes.CREATED_AT: time.Now(),
			},
		},
	}

	_, err := databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return types.Session{}, err
	}

	log.Info("Session saved")

	// Save user session mapping
	putItemInput = databaseTypes.DatabasePutItemInput{
		TableName: userSessionMappingTableName,
		Item: databaseTypes.DatabaseItem{
			PK:         map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
			SK:         map[databaseTypes.FieldType]interface{}{databaseTypes.SID: sid},
			Attributes: map[databaseTypes.FieldType]interface{}{databaseTypes.CREATED_AT: time.Now()},
		},
	}

	_, err = databaseService.PutItemInDatabase(&putItemInput)

	if err != nil {
		return types.Session{}, err
	}

	log.Info("User session mapping saved")
	log.Info("Session created")

	// Session created
	return types.Session{
		UID:   uid,
		Token: token,
	}, nil
}

func DeleteSession(databaseService database.DatabaseService, sessionsTableName string, userSessionMappingTableName string, session types.Session) error {
	log.Info("Started DeleteSession")

	userSessionMappingGetItemOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userSessionMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: session.UID},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info("Session not found in user session mapping")
		return nil

	} else if err != nil {
		return err
	}

	sid := userSessionMappingGetItemOutput.Item.SK[databaseTypes.SID].(string)

	_, err = databaseService.DeleteItemFromDatabase(&databaseTypes.DatabaseDeleteItemInput{
		TableName: sessionsTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.SID: sid},
		},
	})

	if err != nil {
		return err
	}

	log.Info("Session deleted from sessions")

	_, err = databaseService.DeleteItemFromDatabase(&databaseTypes.DatabaseDeleteItemInput{
		TableName: userSessionMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: session.UID},
		},
	})

	if err != nil {
		return err
	}

	log.Info("Session deleted from user session mapping")
	log.Info("Session deleted")

	return nil
}
