package main

import (
	"services/lib"
	"services/lib/log"
	"services/lib/types"

	"services/UserManagementService/cmd"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Info("Started UserManagementService")

	response, err := cmd.HandleRequest(types.Request{
		Body:       request.Body,
		HTTPMethod: request.HTTPMethod,
	})

	return lib.CreateApiResponse(
		response.StatusCode,
		response.Body,
	), err
}

func main() {
	lambda.Start(handler)
}
