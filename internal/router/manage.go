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
}
