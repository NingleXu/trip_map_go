package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
	"trip-map/internal/service"
)

func GetScenicAreaListByPage(c *gin.Context) {
	// 获取需要抓取的city列表
	var req request.ScenicAreaPageRequest

	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	res, err := service.GetScenicAreaListByPage(req)
	if err != nil {
		response.FailWithMessage("查询景区列表异常", c)
		return
	}

	response.OkWithData(res, c)
}
