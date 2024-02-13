package cmd

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"

	"services/lib/log"
	messageTypes "services/lib/types"

	GameLogicService "services/GameLogicService/cmd"

	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseTypes "services/DatabaseService/database/types"
)

// Top level functions

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
		MoveHistory:    []messageTypes.CellPosition{},
		AvailableMoves: []messageTypes.CellPosition{},
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
	moveHistory := []messageTypes.CellPosition{}
	for _, move := range item.Item.Attributes[databaseTypes.MOVE_HISTORY].(primitive.A) {
		moveHistory = append(moveHistory, messageTypes.CellPosition(move.(string)))
	}

	// Convert available moves to []string
	availableMoves := []messageTypes.CellPosition{}
	for _, move := range item.Item.Attributes[databaseTypes.AVAILABLE_MOVES].(primitive.A) {
		availableMoves = append(availableMoves, messageTypes.CellPosition(move.(string)))
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

func ListGames(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, filter messageTypes.GameFilter) ([]messageTypes.Game, error) {
	log.Info("Started ListGames")

	// Get games
	log.Info("Get games")

	getGamesFilter := databaseTypes.DatabaseGetItemsInput{
		TableName: userGameMappingTableName,
		Key:       databaseTypes.DatabaseItem{},
	}

	if filter.UID != "" {
		getGamesFilter.Key.PK = map[databaseTypes.FieldType]interface{}{databaseTypes.UID: filter.UID}
	}

	gamesOutput, err := databaseService.GetItemsFromDatabase(&getGamesFilter)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting games: %v", err))
		return nil, err
	}

	// Get games IDs
	gamesIDs := map[messageTypes.GameID]bool{}

	for _, item := range gamesOutput.Items {
		gid := messageTypes.GameID(item.SK[databaseTypes.GID].(string))
		gamesIDs[gid] = true
	}

	// Filter games
	games := []messageTypes.Game{}

	for gid := range gamesIDs {
		game, err := GetGame(databaseService, gamesTableName, gid)

		if err != nil {
			log.Error(fmt.Sprintf("Error getting game: %v", err))
			return nil, err
		}

		if filter.Check(game) {
			games = append(games, game)
		}
	}

	log.Info("Finished ListGames")
	return games, nil
}

func LeaveGame(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, usersTableName string, uid messageTypes.UserID, gid messageTypes.GameID) error {
	log.Info("Started LeaveGame")

	// Get players' symbols
	log.Info("Get users from game")

	symbolPlayerMapping, err := getSymbolPlayerMapping(databaseService, userGameMappingTableName, gid)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting symbol player mapping: %v", err))
		return err
	}

	// Get winner
	var endState messageTypes.GameWinner
	if symbolPlayerMapping["O"] == uid {
		endState = messageTypes.X
	} else {
		endState = messageTypes.O
	}

	// Update game
	log.Info("Update game")

	_, err = databaseService.UpdateItemInDatabase(&databaseTypes.DatabaseUpdateItemInput{
		TableName: gamesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
			Attributes: map[databaseTypes.FieldType]interface{}{
				databaseTypes.WINNER: endState,
				databaseTypes.STATE:  messageTypes.FINISHED,
			},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error updating game: %v", err))
		return err
	}

	// Get players' UIDs
	player1UID := symbolPlayerMapping["O"]
	player2UID := symbolPlayerMapping["X"]

	// Get result
	var result float64
	switch endState {
	case messageTypes.X:
		result = 1
	case messageTypes.O:
		result = 0
	case messageTypes.TIE:
		result = 0.5
	}

	// Update ELO ratings
	log.Info("Update ELO ratings")

	newPlayer1ELO, newPlayer2ELO, err := getUpdatedELORatings(databaseService, usersTableName, gid, player1UID, player2UID, result)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating ELO ratings: %v", err))
		return err
	}

	// Update player 1 ELO
	log.Info("Update player 1 ELO")

	err = updatePLayerELO(databaseService, usersTableName, player1UID, newPlayer1ELO)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating player 1 ELO: %v", err))
		return err
	}

	// Update player 2 ELO
	log.Info("Update player 2 ELO")

	err = updatePLayerELO(databaseService, usersTableName, player2UID, newPlayer2ELO)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating player 2 ELO: %v", err))
		return err
	}

	log.Info("Finished ConcludeGame")
	return nil
}

func LeaveAllGames(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, usersTableName string, uid messageTypes.UserID) error {
	log.Info("Started LeaveAllGames")

	// Get games
	log.Info("Get games")

	gamesOutput, err := databaseService.GetItemsFromDatabase(&databaseTypes.DatabaseGetItemsInput{
		TableName: userGameMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error getting games: %v", err))
		return err
	}

	// Leave games
	log.Info("Leave games")

	for _, item := range gamesOutput.Items {
		gid := messageTypes.GameID(item.SK[databaseTypes.GID].(string))

		err = LeaveGame(databaseService, gamesTableName, userGameMappingTableName, usersTableName, uid, gid)

		if err != nil {
			log.Error(fmt.Sprintf("Error leaving game: %v", err))
			return err
		}
	}

	log.Info("Finished LeaveAllGames")
	return nil
}

func ConcludeGame(databaseService database.DatabaseService, gamesTableName string, userGameMappingTableName string, usersTableName string, gid messageTypes.GameID, endState messageTypes.GameWinner) error {
	log.Info("Started ConcludeGame")

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
				databaseTypes.WINNER: endState,
				databaseTypes.STATE:  messageTypes.FINISHED,
			},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error updating game: %v", err))
		return err
	}

	// Get players' symbols
	log.Info("Get users from game")

	symbolPlayerMapping, err := getSymbolPlayerMapping(databaseService, userGameMappingTableName, gid)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting symbol player mapping: %v", err))
		return err
	}

	// Get players' UIDs
	player1UID := symbolPlayerMapping["O"]
	player2UID := symbolPlayerMapping["X"]

	// Get result
	var result float64
	switch endState {
	case messageTypes.X:
		result = 1
	case messageTypes.O:
		result = 0
	case messageTypes.TIE:
		result = 0.5
	}

	// Update ELO ratings
	log.Info("Update ELO ratings")

	newPlayer1ELO, newPlayer2ELO, err := getUpdatedELORatings(databaseService, usersTableName, gid, player1UID, player2UID, result)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating ELO ratings: %v", err))
		return err
	}

	// Update player 1 ELO
	log.Info("Update player 1 ELO")

	err = updatePLayerELO(databaseService, usersTableName, player1UID, newPlayer1ELO)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating player 1 ELO: %v", err))
		return err
	}

	// Update player 2 ELO
	log.Info("Update player 2 ELO")

	err = updatePLayerELO(databaseService, usersTableName, player2UID, newPlayer2ELO)

	if err != nil {
		log.Error(fmt.Sprintf("Error updating player 2 ELO: %v", err))
		return err
	}

	log.Info("Finished ConcludeGame")
	return nil
}

func PutMove(databaseService database.DatabaseService, gamesTableName string, gid messageTypes.GameID, move messageTypes.Move) error {
	log.Info("Started PutMove")

	// Get game
	log.Info("Get game")

	game, err := GetGame(databaseService, gamesTableName, gid)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting game: %v", err))
		return err
	}

	// Put move
	log.Info("Put move")

	gameLogicService := GameLogicService.GameLogicService{}

	_, err = gameLogicService.CanPutMove(&game, move)

	if err != nil {
		log.Error(fmt.Sprintf("Error can put move: %v", err))
		return err
	}

	gameLogicService.PutMove(&game, move)

	// Save game
	log.Info("Save game")

	_, err = databaseService.UpdateItemInDatabase(&databaseTypes.DatabaseUpdateItemInput{
		TableName: gamesTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
		Item: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
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

	log.Info("Finished PutMove")
	return nil
}

func DoesThisUserPlayInThisGame(databaseService database.DatabaseService, userGameMappingTableName string, uid messageTypes.UserID, gid messageTypes.GameID) (bool, error) {
	log.Info("Started DoesThisUserPlayInThisGame")

	// Get user game mapping
	log.Info("Get user game mapping")

	_, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userGameMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
	})

	if err == databaseErrors.ErrNoDocuments {
		log.Info(fmt.Sprintf("User %v does not play in game %v", uid, gid))
		return false, nil
	}

	if err != nil {
		log.Error(fmt.Sprintf("Error getting user game mapping: %v", err))
		return false, err
	}

	log.Info(fmt.Sprintf("User %v plays in game %v", uid, gid))
	return true, nil
}

func CreateMove(databaseService database.DatabaseService, userGameMappingTableName string, uid messageTypes.UserID, gid messageTypes.GameID, cellPosition messageTypes.CellPosition) (messageTypes.Move, error) {
	log.Info("Started CreateMove")

	// Get symbol
	log.Info("Get symbol")

	userGameMappingOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: userGameMappingTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error getting symbol: %v", err))
		return messageTypes.Move{}, err
	}

	symbol := userGameMappingOutput.Item.Attributes[databaseTypes.SYMBOL].(string)

	move := messageTypes.Move{
		Cell:   cellPosition,
		Symbol: symbol,
	}

	log.Info("Finished CreateMove")
	return move, nil
}

// Helper functions

func getSymbolPlayerMapping(databaseService database.DatabaseService, userGameMappingTableName string, gid messageTypes.GameID) (map[string]messageTypes.UserID, error) {
	log.Info("Started getSymbolPlayerMapping")

	// Get users from game
	log.Info("Get users from game")

	usersOutput, err := databaseService.GetItemsFromDatabase(&databaseTypes.DatabaseGetItemsInput{
		TableName: userGameMappingTableName,
		Key: databaseTypes.DatabaseItem{
			SK: map[databaseTypes.FieldType]interface{}{databaseTypes.GID: gid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error getting users from game: %v", err))
		return nil, err
	}

	if len(usersOutput.Items) != 2 {
		log.Error(fmt.Sprintf("Error getting users from game: %v, not 2 usergamemappings", err))
		return nil, err
	}

	// Get user UIDs
	player1UID := messageTypes.UserID(usersOutput.Items[0].PK[databaseTypes.UID].(string))
	player2UID := messageTypes.UserID(usersOutput.Items[1].PK[databaseTypes.UID].(string))

	// Get players' symbols
	player1Symbol := usersOutput.Items[0].Attributes[databaseTypes.SYMBOL].(string)
	player2Symbol := usersOutput.Items[1].Attributes[databaseTypes.SYMBOL].(string)

	// Create symbol player mapping
	symbolPlayerMapping := map[string]messageTypes.UserID{
		player1Symbol: player1UID,
		player2Symbol: player2UID,
	}

	return symbolPlayerMapping, nil
}

func getUserELO(databaseService database.DatabaseService, usersTableName string, uid messageTypes.UserID) (float64, error) {
	log.Info("Started getUserELO")

	// Get user ELO rating
	log.Info("Get user ELO rating")

	userELORatingOutput, err := databaseService.GetItemFromDatabase(&databaseTypes.DatabaseGetItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error getting user ELO rating: %v", err))
		return 0, err
	}

	userELORating := float64(userELORatingOutput.Item.Attributes[databaseTypes.ELO].(int32))

	return userELORating, nil
}

func getUpdatedELORatings(databaseService database.DatabaseService, usersTableName string, gid messageTypes.GameID, player1UID messageTypes.UserID, player2UID messageTypes.UserID, result float64) (float64, float64, error) {
	log.Info("Started updateELORatings")

	// Get winner ELO rating
	log.Info("Get winner ELO rating")

	player1ELO, err := getUserELO(databaseService, usersTableName, player1UID)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting winner ELO rating: %v", err))
		return 0, 0, err
	}

	// Get loser ELO rating
	log.Info("Get loser ELO rating")

	player2ELO, err := getUserELO(databaseService, usersTableName, player2UID)

	if err != nil {
		log.Error(fmt.Sprintf("Error getting loser ELO rating: %v", err))
		return 0, 0, err
	}

	// Calculate new ELO ratings
	log.Info("Calculate new ELO ratings")

	newPlayer1ELO, newPlayer2ELO := GameLogicService.CalculateEloRating(player1ELO, player2ELO, result)

	return newPlayer1ELO, newPlayer2ELO, nil
}

func updatePLayerELO(databaseService database.DatabaseService, usersTableName string, uid messageTypes.UserID, newELO float64) error {
	log.Info("Started updatePLayerELO")

	// Update user ELO rating
	log.Info("Update user ELO rating")

	_, err := databaseService.UpdateItemInDatabase(&databaseTypes.DatabaseUpdateItemInput{
		TableName: usersTableName,
		Key: databaseTypes.DatabaseItem{
			PK: map[databaseTypes.FieldType]interface{}{databaseTypes.UID: uid},
		},
		Item: databaseTypes.DatabaseItem{
			Attributes: map[databaseTypes.FieldType]interface{}{
				databaseTypes.ELO: newELO,
			},
		},
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error updating user ELO rating: %v", err))
		return err
	}

	log.Info("Finished updatePLayerELO")
	return nil
}
