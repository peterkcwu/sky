package api

import "github.com/gin-gonic/gin"

func (client *ApiClient) LoadRouter(server *gin.Engine) {
	openRouter := server.Group("/api/v1")
	{
		openRouter.GET("/test", client.Hello)

	}
}