package cmd

import (
	"fmt"
	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseTypes "services/DatabaseService/database/types"
	"services/lib/log"
	"services/lib/types"
)

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
