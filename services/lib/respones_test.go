package lib

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestCreateApiResponse(t *testing.T) {
	type args struct {
		statusCode int
		body       interface{}
	}
	tests := []struct {
		name string
		args args
		want events.APIGatewayProxyResponse
	}{
		{
			name: "Test with valid body",
			args: args{
				statusCode: http.StatusOK,
				body:       map[string]string{"message": "Hello, World!"},
			},
			want: events.APIGatewayProxyResponse{
				Body:       `{"message":"Hello, World!"}`,
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Test with invalid body",
			args: args{
				statusCode: http.StatusOK,
				body:       make(chan int), // This will cause an error when marshalling to JSON
			},
			want: events.APIGatewayProxyResponse{
				Body:       "Error marshalling response body",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateApiResponse(tt.args.statusCode, tt.args.body)
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("CreateApiResponse() StatusCode = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
			if got.Body != tt.want.Body {
				t.Errorf("CreateApiResponse() Body = %v, want %v", got.Body, tt.want.Body)
			}
		})
	}
}
