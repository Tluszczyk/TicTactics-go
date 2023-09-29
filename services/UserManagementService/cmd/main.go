package main

import (
	"net/http"
	"os"

	"services/DatabaseService/database"

	"services/lib"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get environment variables
	databaseDeploymentOption, err := database.ParseDatabaseDeploymentOption(os.Getenv("DATABASE_DEPLOYMENT_OPTION"))
	if err != nil {
		return lib.CreateApiResponse(http.StatusInternalServerError, "Error parsing database deployment option"), err
	}

	// Get database service
	databaseService, err := database.GetDatabaseService(databaseDeploymentOption)
	if err != nil {
		return lib.CreateApiResponse(http.StatusInternalServerError, "Error getting database service"), err
	}

	switch request.HTTPMethod {
	case "GET":
		return HandleGetRequest(databaseService, request)

	case "POST":
		return HandlePostRequest(databaseService, request)

	default:
		return events.APIGatewayProxyResponse{
			Body:       "Method not allowed",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}
