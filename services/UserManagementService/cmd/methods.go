package cmd

import (
	"fmt"
	"net/http"
	"os"
	"services/lib/log"
	messageTypes "services/lib/types"

	"services/DatabaseService/database"
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

	return messageTypes.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handlePostRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started HandlePOSTRequest")

	// Get environment variables
	usersTableName := os.Getenv("USERS_TABLE_NAME")
	passwordHashesTableName := os.Getenv("PASSWORDHASH_TABLE_NAME")
	userPasswordHashMappingTable := os.Getenv("USER_PASSWORD_HASH_MAPPING_TABLE")

	if usersTableName == "" || passwordHashesTableName == "" || userPasswordHashMappingTable == ""{
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateUserResponse{
				Status: messageTypes.Status{
					Code: http.StatusInternalServerError,
					Message: fmt.Sprintf(
						"Environment variables not set: USERS_TABLE_NAME=%s, PASSWORDHASH_TABLE_NAME=%s, USER_PASSWORD_HASH_MAPPING_TABLE=%s",
						usersTableName, passwordHashesTableName, userPasswordHashMappingTable,
					),
				},
			},
		}, nil
	}

	log.Info("Got environment variables")

	// Parse request body
	var createUserRequest CreateUserRequest
	err := messageTypes.ParseRequestBody(request.Body, &createUserRequest)

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
