package bootstrap

import (
	"github.com/gin-gonic/gin"
	v1 "trip-map/internal/router"
)

func InitRouter(r *gin.Engine) {
	apiGroup := r.Group("api/tripMap")

	// 挂载景点路由
	v1.ScenicSpotRouterApi.InitRouter(apiGroup)

	v1.UserRouterApi.InitRouter(apiGroup)

	v1.CityRouterApi.InitRouter(apiGroup)

	manageApiGroup := apiGroup.Group("/manage")
	v1.ManageRouterApi.InitRouter(manageApiGroup)
}
