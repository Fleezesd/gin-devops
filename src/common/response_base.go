package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS = 0
	ERROR   = 7
)

type BaseResp struct {
	Code    int         `json:"code"` // 这里的code不是http的状态码 ，而是前后端交互定义的
	Data    interface{} `json:"result"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
}

// Result 访问成功
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

// Result400 参数错误
func Result400(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// Result401 未授权
func Result401(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// Result403 没权限
func Result403(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusForbidden, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

// Result500 服务器错误
func Result500(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, BaseResp{
		Code:    code,
		Data:    data,
		Message: msg,
		Type:    "",
	})
}

func FailWithMessage(message string, c *gin.Context) {
	Result400(ERROR, map[string]interface{}{}, message, c)
}

func FailWithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result400(ERROR, data, message, c)
}

func Req401WithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(ERROR, data, message, c)
}

func Req403WithWithMessage(message string, c *gin.Context) {
	Result403(ERROR, map[string]interface{}{}, message, c)
}

func Req500WithWithDetailed(data interface{}, message string, c *gin.Context) {
	Result500(ERROR, data, message, c)
}
