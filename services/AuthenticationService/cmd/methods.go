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
	tableNames, status, err := utils.GetEnvironmentVariables("SESSIONS_TABLE_NAME", "USER_SESSION_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       ValidateSessionResponse{Status: status},
		}, err
	}

	sessionsTableName, userSessionMappingTableName := tableNames[0], tableNames[1]

	log.Info("Got environment variables")

	// Parse request body
	var validateSessionRequest ValidateSessionRequest
	err = messageTypes.ParseRequestBody(request.Body, &validateSessionRequest)

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
	log.Info("Started HandleGETRequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("USERS_TABLE_NAME", "PASSWORD_HASHES_TABLE_NAME", "USER_PASSWORD_HASH_MAPPING_TABLE_NAME", "SESSIONS_TABLE_NAME", "USER_SESSION_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body: CreateSessionResponse{
				Status:  status,
				Session: messageTypes.Session{},
			},
		}, err
	}

	usersTableName, passwordHashesTableName, userPasswordHashMappingTableName, sessionsTableName, userSessionMappingTableName := tableNames[0], tableNames[1], tableNames[2], tableNames[3], tableNames[4]

	log.Info("Got environment variables")

	// Parse request body
	var createSessionRequest CreateSessionRequest
	err = messageTypes.ParseRequestBody(request.Body, &createSessionRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Validate Credentials")

	user, err := ValidateCredentials(databaseService, usersTableName, passwordHashesTableName, userPasswordHashMappingTableName, createSessionRequest.Credentials)

	if err == databaseErrors.ErrNoDocuments {
		return messageTypes.Response{
			StatusCode: http.StatusUnauthorized,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusUnauthorized,
					Message: "Invalid credentials",
				},
			},
		}, nil
	} else if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error validating credentials",
				},
			},
		}, err
	}

	log.Info("Check if user is already logged in")

	alreadySignedIn, err := DoesUserAlreadyHaveSession(databaseService, sessionsTableName, userSessionMappingTableName, user.UID)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error checking if user already has session",
				},
			},
		}, err
	}

	if alreadySignedIn {
		return messageTypes.Response{
			StatusCode: http.StatusConflict,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusConflict,
					Message: "User already has a session",
				},
			},
		}, nil
	}

	log.Info("Create session")

	session, err := CreateSession(databaseService, sessionsTableName, userSessionMappingTableName, user.UID)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: CreateSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error creating session",
				},
			},
		}, err
	}

	log.Info("Created session")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: CreateSessionResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Session created",
			},
			Session: session,
		},
	}, nil
}

func handleDeleteRequest(databaseService database.DatabaseService, request messageTypes.Request) (messageTypes.Response, error) {
	log.Info("Started HandleDELETERequest")

	// Get environment variables
	tableNames, status, err := utils.GetEnvironmentVariables("SESSIONS_TABLE_NAME", "USER_SESSION_MAPPING_TABLE_NAME")

	if err != nil {
		return messageTypes.Response{
			StatusCode: int(status.Code),
			Body:       DeleteSessionResponse{Status: status},
		}, err
	}

	sessionsTableName, userSessionMappingTableName := tableNames[0], tableNames[1]

	log.Info("Got environment variables")

	// Parse request body
	var deleteSessionRequest DeleteSessionRequest
	err = messageTypes.ParseRequestBody(request.Body, &deleteSessionRequest)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusBadRequest,
			Body: DeleteSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusBadRequest,
					Message: "Error parsing request body",
				},
			},
		}, err
	}

	log.Info("Delete session")

	err = DeleteSession(databaseService, sessionsTableName, userSessionMappingTableName, deleteSessionRequest.Session)

	if err != nil {
		return messageTypes.Response{
			StatusCode: http.StatusInternalServerError,
			Body: DeleteSessionResponse{
				Status: messageTypes.Status{
					Code:    http.StatusInternalServerError,
					Message: "Error deleting session",
				},
			},
		}, err
	}

	log.Info("Deleted session")

	return messageTypes.Response{
		StatusCode: http.StatusOK,
		Body: DeleteSessionResponse{
			Status: messageTypes.Status{
				Code:    http.StatusOK,
				Message: "Session deleted",
			},
		},
	}, nil
}
