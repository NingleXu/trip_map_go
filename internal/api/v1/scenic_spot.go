package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
	"trip-map/internal/service"
)

func CaptureCityScenicSpot(c *gin.Context) {
	// 获取需要抓取的city列表
	var req request.CaptureCityScenicSpotReq

	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage("请求参数异常", c)
		return
	}
	if len(req.CaptureCityScenicSpotList) == 0 {
		response.FailWithMessage("城市列表不能为空", c)
		return
	}

	if err := service.CaptureCityScenicSpot(req); err != nil {
		response.FailWithMessage("执行异常", c)
		return
	}

	response.Ok(c)
}

func GetScenicSpotInfoByCId(c *gin.Context) {
	cId := c.Query("cId")
	if cId == "" {
		response.FailWithMessage("参数异常", c)
		return
	}

	scenicSpotInfo := service.GetScenicSpotInfoByCId(cId)
	if scenicSpotInfo == nil {
		response.FailWithMessage("查询景点详情失败", c)
		return
	}

	// 查询对应图片列表
	images, err := service.GetScenicSpotImagesByCId(cId)
	if err != nil {
		log.Printf("[Error] 查询景点图片失败%v\n", err)
	}
	imageUrls := make([]string, 0, len(images))
	for _, image := range images {
		imageUrls = append(imageUrls, image.Url)
	}

	scenicSpotInfo.TImages = imageUrls
	response.OkWithData(&scenicSpotInfo, c)
}

func GetScenicSpotListByUserLocation(c *gin.Context) {
	var req request.UserLocationScenicSpotRequest
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage("参数异常", c)
	}
	res, err := service.GetScenicSpotListByUserLocation(&req)
	if err != nil {
		response.FailWithMessage("获取景点列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

func GetScenicSpotListByCityCode(c *gin.Context) {
	cityCode := c.Query("cityCode")
	if cityCode == "" {
		response.FailWithMessage("参数异常", c)
		return
	}

	res, err := service.GetScenicSpotListByCityCode(cityCode, true)
	if err != nil {
		response.FailWithMessage("获取景点列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

func GetTop10ScenicSpotList(c *gin.Context) {
	res, err := service.GetTop10ScenicSpotList()
	if err != nil {
		response.FailWithMessage("获取热门景点列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

func GetScenicSpotByKeyword(c *gin.Context) {
	keyword := c.Query("keyword")

	if keyword == "" {
		var data []interface{}
		response.OkWithData(data, c)
		return
	}

	res, err := service.GetScenicSpotByKeyword(keyword)
	if err != nil {
		response.FailWithMessage("查询景点列表失败", c)
		return
	}
	response.OkWithData(res, c)
}

func GetScenicSpotListByScenicAreaCId(c *gin.Context) {
	scenicAreaCId := c.Query("scenicAreaCId")

	if scenicAreaCId == "" {
		response.FailWithMessage("参数异常", c)
		return
	}

	res, err := service.GetScenicSpotListByScenicAreaCId(scenicAreaCId)
	if err != nil {
		response.FailWithMessage("查询景区景点列表失败", c)
		return
	}

	response.OkWithData(res, c)
}
