package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Response struct {
    Code int         `json:"code"`
    Data interface{} `json:"data"`
    Msg  string      `json:"msg"`
}

const (
    ERROR   = 999
    SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
    c.JSON(http.StatusOK, Response{
        Code: code,
        Data: data,
        Msg:  msg,
    })
}

func Ok(c *gin.Context) {
    Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
    Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
    Result(SUCCESS, data, "操作成功", c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
    Result(SUCCESS, data, message, c)
}

func FailWithMessage(message string, c *gin.Context) {
    Result(ERROR, map[string]interface{}{}, message, c)
}
