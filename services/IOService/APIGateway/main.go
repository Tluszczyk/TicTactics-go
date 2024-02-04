package main

import (
	"io"
	"encoding/json"
	"net/http"

	"services/lib/types"

	auth "services/AuthenticationService/cmd"
	gameManagement "services/GameManagementService/cmd"
	userManagement "services/UserManagementService/cmd"

	"github.com/gin-gonic/gin"
)

func authenticator(req types.Request) (types.Response, error) {
	var data map[string]string
	
	bodyBytes := []byte(req.Body)
	err := json.Unmarshal(bodyBytes, &data)

	if err != nil {
		return types.Response{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid request body",
		}, nil
	}

	// Check if the request has a Session field
	if sessionField, ok := data["Session"]; !ok {
		// Check if the session token is valid
		response, err := auth.HandleRequest(types.Request{
			HTTPMethod: http.MethodGet,
			Body:       sessionField,
		})

		if err != nil {
			return types.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       "error authenticating request: " + err.Error(),
			}, nil
		}

		if response.StatusCode != http.StatusOK {
			return response, nil
		}
	}

	return types.Response{
		StatusCode: http.StatusOK,
		Body:       "authenticated",
	}, nil
}

func handlerMonad(handler func(types.Request) (types.Response, error), requiresAuthentication bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "error reading request body: "+err.Error())
			return
		}

		request := types.Request{
			HTTPMethod: c.Request.Method,
			Body:       string(bodyBytes),
		}

		if requiresAuthentication {
			authResponse, err := authenticator(request)
			if err != nil {
				c.String(http.StatusInternalServerError, "error authenticating request: "+err.Error())
				return
			}

			if authResponse.StatusCode != http.StatusOK {
				c.String(authResponse.StatusCode, authResponse.Body.(string))
				return
			}
		}

		response, err := handler(request)
		if err != nil {
			c.String(http.StatusInternalServerError, "error handling request: "+err.Error())
			return
		}

		c.IndentedJSON(response.StatusCode, response.Body)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	// Authentication Service
	router.GET("/auth/session", handlerMonad(auth.HandleRequest, false))
	router.POST("/auth/session", handlerMonad(auth.HandleRequest, false))
	router.DELETE("/auth/session", handlerMonad(auth.HandleRequest, true))

	// User Management Service
	router.GET("/user", handlerMonad(userManagement.HandleRequest, false))
	router.POST("/user", handlerMonad(userManagement.HandleRequest, true))
	router.DELETE("/user", handlerMonad(userManagement.HandleRequest, true))

	// Game Management Service
	router.POST("/game/create", handlerMonad(gameManagement.HandleRequest, true))
	router.PUT("/game/join", handlerMonad(gameManagement.HandleRequest, true))
	router.PUT("/game/leave", handlerMonad(gameManagement.HandleRequest, true))
	router.GET("/game", handlerMonad(gameManagement.HandleRequest, true))
	router.PUT("/game", handlerMonad(gameManagement.HandleRequest, true))
	router.GET("/game/list-games", handlerMonad(gameManagement.HandleRequest, true))

	router.Run(":8080")
}
