package lib

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var TITLE_MAP map[int]string = map[int]string{
	200: "OK",
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	500: "Internal Server Error",
}

func CreateApiResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {

	// Create response body
	responseBody, err := json.Marshal(body)

	if err != nil {
		responseBody = []byte("Error marshalling response body")
		statusCode = http.StatusInternalServerError
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: statusCode,
	}
}
