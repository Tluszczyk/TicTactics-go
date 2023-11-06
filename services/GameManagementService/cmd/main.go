package main

import (
	"services/lib"
	"services/lib/log"
	"services/lib/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info("Started GameManagementService")

	response, err := HandleRequest(types.Request{
		HTTPMethod: request.HTTPMethod,
		Body:       request.Body,
	})

	return lib.CreateApiResponse(
		response.StatusCode,
		response.Body,
	), err
}

func main() {
	lambda.Start(handler)
}
