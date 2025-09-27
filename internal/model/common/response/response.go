package response

import (
	"net/http"
	"trip-map/internal/model/common"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"message"`
	Success bool        `json:"success"`
}

const (
	ERROR   = -1
	SUCCESS = 200
)

func Result(code int, data interface{}, msg string, suc bool, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
		suc,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", true, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, true, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "成功", true, c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, true, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", false, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, false, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, false, c)
}

type PageResponseData[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Offset   int   `json:"offset"`
	PageSize int   `json:"pageSize"` // 注意这里改为大写，因为Go中可导出的字段需要大写
}

func PageResponseWithEmpty[T any](pageInfo *common.PageInfo) *PageResponseData[T] {
	return &PageResponseData[T]{
		Data:     []T{},
		Total:    0,
		Offset:   pageInfo.Offset,
		PageSize: pageInfo.PageSize,
	}
}

func PageResponseWithData[T any](data []T, total int64, pageInfo *common.PageInfo) *PageResponseData[T] {
	return &PageResponseData[T]{
		Data:     data,
		Total:    total,
		Offset:   pageInfo.Offset,
		PageSize: pageInfo.PageSize,
	}
}
