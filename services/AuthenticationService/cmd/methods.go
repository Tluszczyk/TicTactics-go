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

func HandleRequest(request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started UserManagement HandleRequest")

	log.Info(fmt.Sprintf("Request: %v", request))

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
	log.Info("Started HandleGETRequest")

	// Get environment variables
	sessionsTableName := os.Getenv("SESSIONS_TABLE_NAME")
	userSessionMappingTableName := os.Getenv("USER_SESSION_MAPPING_TABLE_NAME")

	if sessionsTableName == "" {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: ValidateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error getting environment variables",
				},
			},
		}, nil
	}

	log.Info("Got environment variables")

	// Parse request body
	var validateSessionRequest ValidateSessionRequest
	err := messageTypes.ParseRequestBody(request.Body, &validateSessionRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: ValidateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info(fmt.Sprintf("%v", validateSessionRequest.Session))

	log.Info("Parsed request body")

	// Validate session

	isValid, err := ValidateSession(databaseService, sessionsTableName, userSessionMappingTableName, validateSessionRequest.Session)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: ValidateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error validating session",
				},
			},
		}, err
	}

	if !isValid {
		return messageTypes.Response{
			StatusCode: http.StatusUnauthorized,
			Body: ValidateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusUnauthorized,
					Message: "Invalid session",
				},
			},
		}, nil
	}

	log.Info("Validated session")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: ValidateSessionResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Valid session",
			},
		},
	}, nil
}

func handlePostRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Post")

	return messageTypes.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handlePutRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started Put")

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
