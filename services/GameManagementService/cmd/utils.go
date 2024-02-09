package cmd

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"

	"services/lib/log"
	messageTypes "services/lib/types"

	GameLogicService "services/GameLogicService/cmd"

	"services/DatabaseService/database"
	databaseTypes "services/DatabaseService/database/types"
)

func CreateGame(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, uid messageTypes.UserID, gameSettings messageTypes.GameSettings) error {
	log.Info("Started CreateGame")

	gameLogicService := GameLogicService.GameLogicService{}

	log.Info("Create board")
	// Create board
	board := gameLogicService.CreateBoard()

	game := messageTypes.Game{
		GID:            messageTypes.GameID(uuid.New().String()),
		Board:          gameLogicService.SerializeBoard(board),
		Turn:           string(messageTypes.O),
		Winner:         string(messageTypes.NONE),
		MoveHistory:    []string{},
		AvailableMoves: []string{},
		State:          messageTypes.WATING_FOR_OPPONENT,
		TileBoard:      gameLogicService.SerializeTileBoard(gameLogicService.CreateTileBoard()),
	}

	gameLogicService.UpdateAvailableMoves(&game)

	// Save game
	log.Info("Save game")
	_, err := databaseService.PutItemInDatabase(&databaseTypes.DatabasePutItemInput{
		TableName: gamesTableName,
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: game.GID},
			Attributes: map[databaseTypes.FieldType]interface{}{
				databaseTypes.BOARD:           game.Board,
				databaseTypes.TURN:            game.Turn,
				databaseTypes.WINNER:          game.Winner,
				databaseTypes.MOVE_HISTORY:    game.MoveHistory,
				databaseTypes.AVAILABLE_MOVES: game.AvailableMoves,
				databaseTypes.STATE:           game.State,
				databaseTypes.TILE_BOARD:      game.TileBoard,
			},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error saving game: %v", err))
		return err
	}

	// Save user game mapping
	log.Info("Save user game mapping")

	_, err = databaseService.PutItemInDatabase(&databaseTypes.DatabasePutItemInput{
		TableName: userGameMappingTableName,
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: game.GID},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error saving user game mapping: %v", err))
		return err
	}

	log.Info("Finished CreateGame")
	return nil
}

func JoinGame(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, uid messageTypes.UserID, gid messageTypes.GameID) error {
	log.Info("Started JoinGame")

	// Update game
	log.Info("Update game")

	_, err := databaseService.UpdateItemInDatabase(&databaseTypes.DatabaseUpdateItemInput{
		TableName: gamesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
			Attributes: map[databaseTypes.FieldType]interface{}{
				databaseTypes.STATE: messageTypes.IN_PROGRESS,
			},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error updating game: %v", err))
		return err
	}

	// Save user game mapping
	log.Info("Save user game mapping")

	_, err = databaseService.PutItemInDatabase(&databaseTypes.DatabasePutItemInput{
		TableName: userGameMappingTableName,
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error saving user game mapping: %v", err))
		return err
	}

	log.Info("Finished JoinGame")
	return nil
}

func GetGame(databaseService database.DatabaseService, gamesTableName string, gid messageTypes.GameID) (messageTypes.Game, error) {
	log.Info("Started GetGame")

	// Get game
	log.Info("Get game")

	item, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: gamesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error getting game: %v", err))
		return messageTypes.Game{}, err
	}

	if err != nil {
		log.Error(fmt.Sprintf("Error marshalling game: %v", err))
		return messageTypes.Game{}, err
	}

	// Convert move history to []string
	moveHistory := []string{}
	for _, move := range item.Item.Attributes[databaseTypes.MOVE_HISTORY].(primitive.A) {
		moveHistory = append(moveHistory, move.(string))
	}

	// Convert available moves to []string
	availableMoves := []string{}
	for _, move := range item.Item.Attributes[databaseTypes.AVAILABLE_MOVES].(primitive.A) {
		availableMoves = append(availableMoves, move.(string))
	}

	game := messageTypes.Game{
		GID:            gid,
		Board:          item.Item.Attributes[databaseTypes.BOARD].(string),
		Turn:           item.Item.Attributes[databaseTypes.TURN].(string),
		Winner:         item.Item.Attributes[databaseTypes.WINNER].(string),
		MoveHistory:    moveHistory,
		AvailableMoves: availableMoves,
		State:          messageTypes.GameState(item.Item.Attributes[databaseTypes.STATE].(string)),
		TileBoard:      item.Item.Attributes[databaseTypes.TILE_BOARD].(string),
	}

	log.Info("Finished GetGame")
	return game, nil
}
