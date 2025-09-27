package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/internal/service"
)

func GetWaitTimeScenicAreaList(c *gin.Context) {
	res, err := service.GetWaitTimeScenicAreaList()
	if err != nil {
		response.FailWithMessage("查询带等待时间的景区列表异常", c)
		return
	}
	response.OkWithData(res, c)
}

func GetScenicAreaWaitTimeList(c *gin.Context) {
	cId := c.Query("cId")
	if cId == "" {
		response.FailWithMessage("参数异常", c)
		return
	}

	res, err := service.GetScenicAreaWaitTime(cId)
	if err != nil {
		response.FailWithMessage("查询景区景点等待时间失败", c)
		return
	}
	response.OkWithData(res, c)
}
