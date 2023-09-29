package main

import (
	"net/http"
	"services/lib"

	"github.com/aws/aws-lambda-go/events"
)

func HandleGet(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lib.CreateApiResponse(http.StatusNotImplemented, "Not implemented"), nil
}

func HandlePost(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lib.CreateApiResponse(http.StatusNotImplemented, "Not implemented"), nil
}

func HandlePut(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lib.CreateApiResponse(http.StatusNotImplemented, "Not implemented"), nil
}

func HandleDelete(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lib.CreateApiResponse(http.StatusNotImplemented, "Not implemented"), nil
}
