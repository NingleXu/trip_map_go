package router

import (
	"github.com/gin-gonic/gin"
	v1 "trip-map/internal/api/v1"
)

var UserRouterApi = new(UserRouter)

type UserRouter struct {
}

func (u *UserRouter) InitRouter(r *gin.RouterGroup) {
	r.POST("/login", v1.UserLogin)
}
