package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
	"trip-map/internal/service"
)

func UserLogin(c *gin.Context) {
	var req request.UserLoginRequest

	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	if req.Code == "" {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	token, err := service.DoLogin(req.Code)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(token, c)
}
