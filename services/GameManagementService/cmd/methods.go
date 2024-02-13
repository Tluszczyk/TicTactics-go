package cmd

import (
	"fmt"
	"net/http"
	"os"

	"services/lib/log"
	messageTypes "services/lib/types"
	"services/lib/utils"

	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseOptions "services/DatabaseService/database/options"
)

func HandleRequest(request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started UserManagement HandleRequest")

	log.Info("Getting environment variables")
	// Get environment variables
	databaseDeploymentOption, err := databaseOptions.ParseDatabaseDeploymentOption(os.Getenv("DATABASE_DEPLOYMENT_OPTION"))
	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error parsing database deployment option",
		}, err
	}

	log.Info("Getting database service")
	// Get database service
	databaseService, err := database.GetDatabaseService(databaseDeploymentOption)
	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error getting database service",
		}, err
	}

	switch request.HTTPMethod {
	case "GET":
		return handleGetRequest(databaseService, request)

	case "POST":
		return handlePostRequest(databaseService, request)

	case "PUT":
		return handlePutRequest(databaseService, request)

	case "DELETE":
		return handleDeleteRequest(databaseService, request)

	default:
		return messageTypes.Response{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}
}

func handleGetRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Get")

	log.Info(fmt.Sprintf("Request path: %s", request.Path))

	switch request.Path {
	case "/game":
		return handleGetGameRequest(databaseService, request)

	case "/game/listGames":
		return handleListGamesRequest(databaseService, request)

	default:
		return messageTypes.Response{
			StatusCode: http.StatusNotFound,
			Body:       "Not found",
		}, nil
	}
}

func handleGetGameRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started GetGameRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       CreateGameResponse{Status: status},
		}, err
	}

	gamesTableName := tableNames[0]

	log.Info("Got environment variables")

	// Parse request body
	var getGameRequest GetGameRequest
	err = messageTypes.ParseRequestBody(request.Body, &getGameRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: GetGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Get game
	game, err := GetGame(databaseService, gamesTableName, getGameRequest.GID)

	if err == databaseErrors.ErrNoDocuments {
		return messageTypes.Response{
			StatusCode: http.StatusNotFound,
			Body: GetGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusNotFound,
					Message: "Game not found",
				},
			},
		}, err
	} else if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: GetGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error getting game",
				},
			},
		}, err
	}

	log.Info("Got game")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: GetGameResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Game found",
			},
			Game: game,
		},
	}, nil
}

func handleListGamesRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started ListGamesRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME", "USER_GAME_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       ListGamesResponse{Status: status},
		}, err
	}

	gamesTableName, userGameMappingTable := tableNames[0], tableNames[1]

	log.Info("Got environment variables")

	// Parse request body
	var listGamesRequest ListGamesRequest
	err = messageTypes.ParseRequestBody(request.Body, &listGamesRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: ListGamesResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// List games
	games, err := ListGames(databaseService, gamesTableName, userGameMappingTable, listGamesRequest.Filter)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: ListGamesResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error listing games",
				},
			},
		}, err
	}

	log.Info("Listed games")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: ListGamesResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Games listed",
			},
			Games: games,
		},
	}, nil
}

func handlePostRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started HandlePOSTRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME", "USER_GAME_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       CreateGameResponse{Status: status},
		}, err
	}

	gamesTableName, userGameMappingTable := tableNames[0], tableNames[1]

	log.Info("Got environment variables")

	// Parse request body
	var createGameRequest CreateGameRequest
	err = messageTypes.ParseRequestBody(request.Body, &createGameRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: CreateGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Create game
	err = CreateGame(databaseService, gamesTableName, userGameMappingTable, createGameRequest.Session.UID, createGameRequest.Settings)

	if err == databaseErrors.ErrItemAlreadyExists {
		return messageTypes.Response{
			StatusCode: http.StatusConflict,
			Body: CreateGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusConflict,
					Message: "Game already exists",
				},
			},
		}, err
	} else if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error creating game",
				},
			},
		}, err
	}

	log.Info("Created game")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: CreateGameResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Game created",
			},
		},
	}, nil
}

func handlePutRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Put")

	log.Info(fmt.Sprintf("Request path: %s", request.Path))

	switch request.Path {
	case "/game/join":
		return handleJoinRequest(databaseService, request)

	case "/game/leave":
		return handleLeaveRequest(databaseService, request)

	case "/game/leaveAll":
		return handleLeaveAllRequest(databaseService, request)

	case "/game":
		return handleUpdateRequest(databaseService, request)

	default:
		return messageTypes.Response{
			StatusCode: http.StatusNotFound,
			Body:       "Not found",
		}, nil
	}
}

func handleJoinRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Join")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME", "USER_GAME_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       CreateGameResponse{Status: status},
		}, err
	}

	gamesTableName, userGameMappingTable := tableNames[0], tableNames[1]

	log.Info("Got environment variables")

	// Parse request body
	var joinGameRequest JoinGameRequest
	err = messageTypes.ParseRequestBody(request.Body, &joinGameRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: JoinGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Join game
	err = JoinGame(databaseService, gamesTableName, userGameMappingTable, joinGameRequest.Session.UID, joinGameRequest.GID)

	if err == databaseErrors.ErrNoDocuments {
		return messageTypes.Response{
			StatusCode: http.StatusNotFound,
			Body: JoinGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusNotFound,
					Message: "Game not found",
				},
			},
		}, err
	} else if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: JoinGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error joining game",
				},
			},
		}, err
	}

	log.Info("Joined game")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: JoinGameResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Game joined",
			},
		},
	}, nil
}

func handleLeaveRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Leave")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME", "USER_GAME_MAPPING_TABLE_NAME", "USERS_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       LeaveGameResponse{Status: status},
		}, err
	}

	gamesTableName, userGameMappingTable, usersTableName := tableNames[0], tableNames[1], tableNames[2]

	log.Info("Got environment variables")

	// Parse request body
	var leaveGameRequest LeaveGameRequest
	err = messageTypes.ParseRequestBody(request.Body, &leaveGameRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: LeaveGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Leave game
	err = LeaveGame(databaseService, gamesTableName, userGameMappingTable, usersTableName, leaveGameRequest.Session.UID, leaveGameRequest.GID)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: LeaveGameResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error leaving game",
				},
			},
		}, err
	}

	log.Info("Left game")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: LeaveGameResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Game left",
			},
		},
	}, nil
}

func handleLeaveAllRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started LeaveAll")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("GAMES_TABLE_NAME", "USER_GAME_MAPPING_TABLE_NAME", "USERS_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       LeaveAllGamesResponse{Status: status},
		}, err
	}

	gamesTableName, userGameMappingTable, usersTableName := tableNames[0], tableNames[1], tableNames[2]

	log.Info("Got environment variables")

	// Parse request body
	var leaveAllGamesRequest LeaveAllGamesRequest
	err = messageTypes.ParseRequestBody(request.Body, &leaveAllGamesRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: LeaveAllGamesResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Leave game
	err = LeaveAllGames(databaseService, gamesTableName, userGameMappingTable, usersTableName, leaveAllGamesRequest.Session.UID)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: LeaveAllGamesResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error leaving game",
				},
			},
		}, err
	}

	log.Info("Left all games")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: LeaveAllGamesResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "All games left",
			},
		},
	}, nil
}

func handleUpdateRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Update")

	return messageTypes.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handleDeleteRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Delete")

	return messageTypes.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}
