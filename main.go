package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"trip-map/config"
	"trip-map/internal/bootstrap"
	"trip-map/schedule"
)

func main() {
	r := gin.Default()

	globalApiGroup := r.Group("api")
	globalApiGroup.GET("/test", func(c *gin.Context) {})

	initSystem(r)

	r.Run(":" + strconv.Itoa(config.SysConfig.Server.Port))
}

func initSystem(r *gin.Engine) {
	// 读取配置文件
	config.InitConfig()
	// 初始化route
	bootstrap.InitRouter(r)
	// 初始化db
	bootstrap.InitDB()
	// 启动等待时间数据定时抓取任务
	schedule.StartGlobalScheduler()
}
