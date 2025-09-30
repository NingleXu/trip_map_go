package router

import (
	"github.com/gin-gonic/gin"
	v1 "trip-map/internal/api/v1"
)

var ScenicSpotRouterApi = new(ScenicSpotRouter)

type ScenicSpotRouter struct {
}

func (ssr *ScenicSpotRouter) InitRouter(r *gin.RouterGroup) {
	r.POST("/captureCityScenicSpot", v1.CaptureCityScenicSpot)

	// map页接口
	r.POST("/getScenicSpotListByUserLocation", v1.GetScenicSpotListByUserLocation)
	r.GET("/getScenicSpotListByCityCode", v1.GetScenicSpotListByCityCode)
	r.GET("/getScenicSpotListByScenicAreaCId", v1.GetScenicSpotListByScenicAreaCId)

	// 详情页接口
	r.GET("/getScenicSpotInfoByCId", v1.GetScenicSpotInfoByCId)
	r.GET("/getNoteListByCIdAndPage", v1.GetNoteListByCIdAndPage)

	// 搜索页
	r.GET("/getHotScenicSpotList", v1.GetTop10ScenicSpotList)
	r.GET("/getScenicSpotByKeyword", v1.GetScenicSpotByKeyword)

	// 景区
	r.GET("/getScenicAreaListByPage", v1.GetScenicAreaListByPage)

	// 带等待时间的景区列表
	r.GET("/getWaitTimeScenicAreaList", v1.GetWaitTimeScenicAreaList)
	// 查询景区的景点等待时间
	r.GET("/getScenicAreaWaitTimeList", v1.GetScenicAreaWaitTimeList)
	// 查询景点的等待时间
	r.GET("/getScenicSpotWaitTimeList", v1.GetScenicSpotWaitTimeList)
}
