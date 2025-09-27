package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/internal/service"
)

func GetCityList(c *gin.Context) {
	list, err := service.GetCityList()
	if err != nil {
		response.FailWithMessage("请求异常", c)
		return
	}
	response.OkWithData(list, c)
}
