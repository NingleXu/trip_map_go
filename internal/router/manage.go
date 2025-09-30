package router

import (
	"github.com/gin-gonic/gin"
	v1 "trip-map/internal/api/v1"
)

var ManageRouterApi = new(ManageRouter)

type ManageRouter struct {
}

func (mr *ManageRouter) InitRouter(r *gin.RouterGroup) {
	r.POST("/updateScenicAreaScenic", v1.UpdateScenicAreaScenic)
	r.GET("/getScenicSpotListByCityCode", v1.ManageGetScenicSpotListByCityCode)
	r.POST("/saveScenicSpot", v1.SaveScenicSpot)
	r.GET("/getScenicAreaListByCityCode", v1.GetScenicAreaListByCityCode)
	r.GET("/getScenicSpotInfoByManage", v1.GetScenicSpotInfoByManage)
	r.POST("/updateScenicSpotInfoByManage", v1.UpdateScenicSpotInfoByManage)

	// 等待时间缓存管理接口
	r.GET("/waitTime/cacheStatus", v1.GetWaitTimeCacheStatus)
	r.POST("/waitTime/forceRefresh", v1.ForceRefreshWaitTimeCache)
	r.DELETE("/waitTime/clearCache", v1.ClearWaitTimeCache)
	r.GET("/waitTime/direct", v1.GetWaitTimeDirectly)
	r.GET("/waitTime/rawData", v1.GetScenicAreaRawData)
}
