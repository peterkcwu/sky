package api

import "github.com/gin-gonic/gin"

func (client *ApiClient) LoadRouter(server *gin.Engine) {
	openRouter := server.Group("/api/v1")
	{
		openRouter.GET("/test", client.Hello)
		openRouter.POST("/upload", client.FileUpload)
		openRouter.POST("/upload_img", client.PhotoUpload)
		openRouter.GET("/get_img", client.GetImg)
	}
}
