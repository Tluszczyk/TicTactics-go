package main

import (
	"io"
	"net/http"

	"services/lib/types"

	auth "services/AuthenticationService/cmd"
	gameManagement "services/GameManagementService/cmd"
	userManagement "services/UserManagementService/cmd"

	"github.com/gin-gonic/gin"
)

func handlerMonad(handler func(types.Request) (types.Response, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading request body")
			return
		}

		request := types.Request{
			HTTPMethod: c.Request.Method,
			Body:       string(bodyBytes),
		}

		response, err := handler(request)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error handling request")
			return
		}

		c.IndentedJSON(response.StatusCode, response.Body)
	}
}

func main() {
	router := gin.Default()

	// Authentication Service
	router.GET("/auth/session", handlerMonad(auth.HandleRequest))
	router.POST("/auth/session", handlerMonad(auth.HandleRequest))
	router.DELETE("/auth/session", handlerMonad(auth.HandleRequest))

	// User Management Service
	router.GET("/user", handlerMonad(userManagement.HandleRequest))
	router.POST("/user", handlerMonad(userManagement.HandleRequest))
	router.DELETE("/user", handlerMonad(userManagement.HandleRequest))

	// Game Management Service
	router.POST("/game/create", handlerMonad(gameManagement.HandleRequest))
	router.PUT("/game/join", handlerMonad(gameManagement.HandleRequest))
	router.PUT("/game/leave", handlerMonad(gameManagement.HandleRequest))
	router.GET("/game", handlerMonad(gameManagement.HandleRequest))
	router.PUT("/game", handlerMonad(gameManagement.HandleRequest))
	router.GET("/game/list-games", handlerMonad(gameManagement.HandleRequest))

	router.Run(":8080")
}
