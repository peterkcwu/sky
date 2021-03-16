package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Response struct {
	RetCode int         `json:"ret"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"`
	Count   int32       `json:"count,omitempty"`
}

type LingdongResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func MustGetID(c *gin.Context, key string) (uint, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		RespWithError(c, 400, fmt.Sprintf("%s must be valid integer", key))
		return 0, err
	}
	return uint(id), nil
}

func RespWithErr(c *gin.Context, errCode int, err error) {
	logrus.WithField("code", errCode).WithError(err).Error("response with Error")
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	result := Response{
		RetCode: errCode,
		ErrMsg:  errMsg,
	}
	c.JSON(200, result)
	c.Abort()
}

func RespWithError(c *gin.Context, errcode int, errmsg string) {
	logrus.WithField("code", errcode).WithField("msg", errmsg).Error("response with error")
	result := Response{RetCode: errcode, ErrMsg: errmsg}
	c.JSON(200, result)
	c.Abort()
}

func RespWithMsg(c *gin.Context, msg string) {
	result := Response{RetCode: 0, ErrMsg: msg}
	c.JSON(200, result)
	c.Abort()
}

func RespWithErrorMsg(c *gin.Context, errCode int, msg string) {
	logrus.WithField("code", errCode).WithField("error", msg).Error("response with Error")
	result := Response{RetCode: errCode, ErrMsg: msg}
	c.JSON(200, result)
	c.Abort()

}

func RespWithData(c *gin.Context, data interface{}) {
	result := Response{RetCode: 0, Data: data}
	c.JSON(200, result)
}

func RespWithCount(c *gin.Context, data interface{}, count int32) {
	result := Response{RetCode: 0, Data: data, Count: count}
	c.JSON(200, result)
}

func LingdongResp(c *gin.Context, data interface{}, errCode int, msg string) {
	result := LingdongResponse{Code: errCode, Data: data, Msg: msg}
	c.JSON(200, result)
}
