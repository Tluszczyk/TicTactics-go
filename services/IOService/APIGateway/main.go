package main

import (
	"io"
	"net/http"

	"services/lib/log"
	"services/lib/types"

	auth "services/AuthenticationService/cmd"
	gameManagement "services/GameManagementService/cmd"
	userManagement "services/UserManagementService/cmd"

	"github.com/gin-gonic/gin"
)

func authenticator(req types.Request) (types.Response, error) {
	log.Info("authenticating request")

	// Check if the session token is valid
	response, err := auth.HandleRequest(types.Request{
		HTTPMethod: http.MethodGet,
		Body:       req.Body,
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

	log.Info("request authenticated")
	return types.Response{
		StatusCode: http.StatusOK,
		Body:       "authenticated",
	}, nil
}

func handlerMonad(handler func(types.Request) (types.Response, error), requiresAuthentication bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, types.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       "error reading request body: " + err.Error(),
			})
			return
		}

		request := types.Request{
			HTTPMethod: c.Request.Method,
			Body:       string(bodyBytes),
		}

		if requiresAuthentication {
			authResponse, err := authenticator(request)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, types.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       "error authenticating request: " + err.Error(),
				})
				return
			}

			if authResponse.StatusCode != http.StatusOK {
				c.JSON(authResponse.StatusCode, authResponse)
				return
			}
		}

		response, err := handler(request)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, types.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       "error handling request: " + err.Error(),
			})
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
	router.GET("/user", handlerMonad(userManagement.HandleRequest, true))
	router.POST("/user", handlerMonad(userManagement.HandleRequest, false))
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
