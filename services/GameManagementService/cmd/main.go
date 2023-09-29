package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "GET":
		return HandleGet(request)

	case "POST":
		return HandlePost(request)

	case "PUT":
		return HandlePut(request)

	case "DELETE":
		return HandleDelete(request)

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
