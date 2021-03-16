package api

import (
	"github.com/gin-gonic/gin"
	"github.com/new_web/api/basic"
)

func (client *ApiClient) Hello(c *gin.Context) {
	basic.RespWithMsg(c, "hello")
}

