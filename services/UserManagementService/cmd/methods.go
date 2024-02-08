package cmd

import (
	"net/http"
	"os"
	"services/lib/log"
	messageTypes "services/lib/types"
	"services/lib/utils"

	"services/DatabaseService/database"
	databaseErrors "services/DatabaseService/database/errors"
	databaseOptions "services/DatabaseService/database/options"
)

func HandleRequest(request messageTypes.Request) (response messageTypes.Response, err error) {
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

	log.Info("Proxing request")
	switch request.HTTPMethod {
	case "GET":
		return handleGetRequest(databaseService, request)

	case "POST":
		return handlePostRequest(databaseService, request)

	default:
		return messageTypes.Response{
			Body:       "Method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

func handleGetRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started HandleGETRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("USERS_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body: GetUserResponse{
				Status: status,
				User:   messageTypes.User{},
			},
		}, err
	}

	usersTableName := tableNames[0]

	log.Info("Got environment variables")

	// Parse request body
	var getUserRequest GetUserRequest
	err = messageTypes.ParseRequestBody(request.Body, &getUserRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: GetUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
				User: messageTypes.User{},
			},
		}, err
	}

	log.Info("Parsed request body")

	user, err := GetUser(databaseService, usersTableName, getUserRequest.Username)

	if err == databaseErrors.ErrNoDocuments {
		return messageTypes.Response{
			StatusCode: http.StatusNotFound,
			Body: GetUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusNotFound,
					Message: "User not found",
				},
				User: messageTypes.User{},
			},
		}, nil
	} else if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: GetUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error getting user",
				},
				User: user,
			},
		}, err
	}

	log.Info("Got user")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: GetUserResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "User found",
			},
			User: user,
		},
	}, nil
}

func handlePostRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started HandlePOSTRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("USERS_TABLE_NAME", "PASSWORD_HASHES_TABLE_NAME", "USER_PASSWORD_HASH_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       CreateUserResponse{Status: status},
		}, err
	}

	usersTableName := tableNames[0]
	passwordHashesTableName := tableNames[1]
	userPasswordHashMappingTable := tableNames[2]

	log.Info("Got environment variables")

	// Parse request body
	var createUserRequest CreateUserRequest
	err = messageTypes.ParseRequestBody(request.Body, &createUserRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: CreateUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Parsed request body")

	// Check if user already exists
	userAlreadyExists, err := DoesUserAlreadyExist(databaseService, usersTableName, createUserRequest.Credentials)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error checking if user already exists",
				},
			},
		}, err
	}

	log.Info("Checked if user already exists")

	if userAlreadyExists {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: CreateUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "User already exists",
				},
			},
		}, nil
	}

	log.Info("User does not already exist")

	// Create user
	err = CreateUser(databaseService, usersTableName, passwordHashesTableName, userPasswordHashMappingTable, createUserRequest.Credentials)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateUserResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error creating user",
				},
			},
		}, err
	}

	log.Info("Created user")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: CreateUserResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "User created",
			},
		},
	}, nil
}
