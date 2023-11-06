package main

import (
	"net/http"

	"services/lib/log"
	"services/lib/types"
)

func HandleRequest(request types.Request) (types.Response, error) {
	switch request.HTTPMethod {
	case "GET":
		return handleGet(request)

	case "POST":
		return handlePost(request)

	case "PUT":
		return handlePut(request)

	case "DELETE":
		return handleDelete(request)

	default:
		return types.Response{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}
}

func handleGet(request types.Request) (types.Response, error) {
	log.Info("Started Get")

	return types.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handlePost(request types.Request) (types.Response, error) {
	log.Info("Started Post")

	return types.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handlePut(request types.Request) (types.Response, error) {
	log.Info("Started Put")

	return types.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}

func handleDelete(request types.Request) (types.Response, error) {
	log.Info("Started Delete")

	return types.Response{
		StatusCode: http.StatusNotImplemented,
		Body:       "Not implemented",
	}, nil
}
