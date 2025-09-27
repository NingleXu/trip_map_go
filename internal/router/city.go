package router

import (
	"github.com/gin-gonic/gin"
	v1 "trip-map/internal/api/v1"
)

var CityRouterApi = new(CityRouter)

type CityRouter struct {
}

func (u *CityRouter) InitRouter(r *gin.RouterGroup) {
	r.GET("/cityList", v1.GetCityList)
}
