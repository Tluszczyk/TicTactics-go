package main

import (
	"net/http"
	// "os"
	"services/lib"
	// messageTypes "services/lib/types"

	"services/DatabaseService/database"

	"github.com/aws/aws-lambda-go/events"
)

func HandleGetRequest(databaseService database.DatabaseService, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lib.CreateApiResponse(http.StatusNotImplemented, "Not implemented"), nil
}

func HandlePostRequest(databaseService database.DatabaseService, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// // Get environment variables
	// usersTableName := os.Getenv("USERS_TABLE_NAME")
	// passwordHashesTableName := os.Getenv("PASSWORDHASH_TABLE_NAME")

	// // Parse request body
	// var createUserRequest CreateUserRequest
	// err := messageTypes.ParseRequestBody(request.Body, &createUserRequest)

	// if err != nil {
	// 	return lib.CreateApiResponse(
	// 		http.StatusBadRequest,
	// 		CreateUserResponse{
	// 			Status: messageTypes.Status{
	// 				Code:    http.StatusBadRequest,
	// 				Message: "Error parsing request body",
	// 			},
	// 		},
	// 	), err
	// }

	// // Check if user already exists
	// userAlreadyExists, err := CheckIfUserAlreadyExists(databaseService, usersTableName, createUserRequest.Credentials.Username)

	// if err != nil {
	// 	return lib.CreateApiResponse(
	// 		http.StatusInternalServerError, CreateUserResponse{
	// 			Status: messageTypes.Status{
	// 				Code:    http.StatusInternalServerError,
	// 				Message: "Error checking if user already exists",
	// 			},
	// 		},
	// 	), err
	// }

	// if userAlreadyExists {
	// 	return lib.CreateApiResponse(
	// 		http.StatusBadRequest, CreateUserResponse{
	// 			Status: messageTypes.Status{
	// 				Code:    http.StatusBadRequest,
	// 				Message: "User already exists",
	// 			},
	// 		},
	// 	), nil
	// }

	// // Create user
	// err = CreateUser(databaseService, usersTableName, passwordHashesTableName, createUserRequest.Credentials)

	// if err != nil {
	// 	return lib.CreateApiResponse(
	// 		http.StatusInternalServerError,
	// 		CreateUserResponse{
	// 			Status: messageTypes.Status{
	// 				Code:    http.StatusInternalServerError,
	// 				Message: "Error creating user",
	// 			},
	// 		},
	// 	), err
	// }

	// return lib.CreateApiResponse(
	// 	http.StatusOK,
	// 	CreateUserResponse{
	// 		Status: messageTypes.Status{
	// 			Code:    http.StatusOK,
	// 			Message: "User created",
	// 		},
	// 	},
	// ), nil

	return lib.CreateApiResponse(
		http.StatusNotImplemented,
		"Not implemented",
	), nil
}
