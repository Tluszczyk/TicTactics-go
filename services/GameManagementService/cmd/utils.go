package cmd

import (
	"fmt"

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
