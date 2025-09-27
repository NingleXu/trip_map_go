package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
	"trip-map/internal/service"
)

func UpdateScenicAreaScenic(c *gin.Context) {
	var req request.ScenicAreaScenicPointUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	err := service.UpdateScenicAreaScenic(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func ManageGetScenicSpotListByCityCode(c *gin.Context) {
	cityCode := c.Query("cityCode")
	if cityCode == "" {
		response.FailWithMessage("参数异常", c)
		return
	}

	res, err := service.GetScenicSpotListByCityCode(cityCode, false)
	if err != nil {
		response.FailWithMessage("获取景点列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

func SaveScenicSpot(c *gin.Context) {
	var ss model.ScenicSpot
	if err := c.ShouldBindJSON(&ss); err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	err := service.SaveScenicSpotByManage(&ss)
	if err != nil {
		response.FailWithMessage("获取景点列表失败", c)
		return
	}
	response.Ok(c)

}

func GetScenicAreaListByCityCode(c *gin.Context) {
	cityCode := c.Query("cityCode")
	if cityCode == "" {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	scenicAreaList, err := service.GetScenicAreaListByCityCode(cityCode)
	if err != nil {
		response.FailWithMessage("获取景点列表失败", c)
		return
	}
	response.OkWithData(scenicAreaList, c)
}

func GetScenicSpotInfoByManage(c *gin.Context) {
	scenicId := c.Query("scenicId")
	if scenicId == "" {
		response.FailWithMessage("请求参数异常", c)
		return
	}

	scenicInfo, err := service.GetScenicSpotInfoByManage(scenicId)
	if err != nil {
		response.FailWithMessage("获取景点详情失败", c)
		return
	}
	response.OkWithData(scenicInfo, c)
}

func UpdateScenicSpotInfoByManage(c *gin.Context) {
	var ss model.ScenicSpot
	if err := c.ShouldBindJSON(&ss); err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}
	err := service.UpdateScenicSpotInfoByManage(&ss)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
