package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/schedule"
)

// GetWaitTimeCacheStatus 获取等待时间缓存状态
func GetWaitTimeCacheStatus(c *gin.Context) {
	status := schedule.GetWaitTimeCacheStatus()
	response.OkWithData(status, c)
}

// ForceRefreshWaitTimeCache 强制刷新等待时间缓存
func ForceRefreshWaitTimeCache(c *gin.Context) {
	err := schedule.ForceRefreshWaitTimeCache()
	if err != nil {
		response.FailWithMessage("强制刷新缓存失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("强制刷新缓存任务已启动", c)
}

// ClearWaitTimeCache 清空等待时间缓存
func ClearWaitTimeCache(c *gin.Context) {
	scenicAreaId := c.Query("scenicAreaId") // 可选参数，为空则清空所有缓存

	err := schedule.ClearWaitTimeCache(scenicAreaId)
	if err != nil {
		response.FailWithMessage("清空缓存失败: "+err.Error(), c)
		return
	}

	if scenicAreaId == "" {
		response.OkWithMessage("已清空所有等待时间缓存", c)
	} else {
		response.OkWithMessage("已清空景区 "+scenicAreaId+" 的等待时间缓存", c)
	}
}

// GetWaitTimeDirectly 直接获取等待时间数据 (不使用缓存，用于测试)
func GetWaitTimeDirectly(c *gin.Context) {
	cId := c.Query("cId")
	if cId == "" {
		response.FailWithMessage("参数 cId 不能为空", c)
		return
	}

	data, err := schedule.GetWaitTimeDirectly(cId)
	if err != nil {
		response.FailWithMessage("直接获取等待时间失败: "+err.Error(), c)
		return
	}

	response.OkWithData(data, c)
}

// GetScenicAreaRawData 获取景区原始数据 (用于查看第三方API原始响应，方便ID映射)
func GetScenicAreaRawData(c *gin.Context) {
	cId := c.Query("cId")
	if cId == "" {
		response.FailWithMessage("参数 cId 不能为空", c)
		return
	}

	data, err := schedule.GetScenicAreaRawData(cId)
	if err != nil {
		response.FailWithMessage("获取原始数据失败: "+err.Error(), c)
		return
	}

	response.OkWithData(data, c)
}
