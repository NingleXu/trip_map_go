package service

import (
	"fmt"
	"github.com/bytedance/gg/gptr"
	"log"
	"trip-map/core/wait_time"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/response"
)

func GetWaitTimeScenicAreaList() ([]model.ScenicSpotSummary, error) {
	var scenicSpots []model.ScenicSpot

	err := global.Db.
		Where("is_delete = 0 and is_scenic_area = 1 and wait_time_enable = 1").
		Find(&scenicSpots).
		Error
	if err != nil {
		log.Printf("查询带等待时间的景区列表失败%v\n", err)
		return nil, err
	}

	return convertScenicSpots(scenicSpots), err
}

func GetScenicAreaWaitTime(scId string) (*response.ScenicAreaScenicSpotRes, error) {
	// 检查是否支持等待时间
	_, ok := wait_time.ScenicAreaWaitTimeHandlerMap[scId]
	if !ok {
		log.Printf("景区:%s不存在等待时间\n", scId)
		return nil, fmt.Errorf("景区 %s 不存在", scId)
	}

	// 获取缓存管理器
	cacheManager := wait_time.GetCacheManager()

	// 从缓存获取等待时间数据
	waitTimeData, exists := cacheManager.GetWaitTimeData(scId)
	if !exists {
		log.Printf("景区:%s的等待时间数据不在缓存中，尝试直接获取\n", scId)

		// 如果缓存中没有数据，尝试直接获取
		handler, _ := wait_time.ScenicAreaWaitTimeHandlerMap[scId]
		directData, err := handler.GetWaitTime()
		if err != nil {
			log.Printf("直接获取景区:%s等待时间数据失败: %v\n", scId, err)
			return nil, fmt.Errorf("景区 %s 的等待时间数据暂不可用，请稍后重试", scId)
		}

		// 将直接获取的数据存入缓存
		cacheManager.SetWaitTimeData(scId, directData)
		waitTimeData = directData
		log.Printf("景区:%s等待时间数据直接获取成功并已缓存\n", scId)
	}

	// 同步查数据库
	res, err := GetScenicSpotListByScenicAreaCId(scId)
	if err != nil {
		log.Printf("获取景区:%s的景点信息失败...异常%v\n", scId, err)
		return nil, err
	}

	// 合并结果
	originList := res.ScenicSpotList
	var scenicSpotList = make([]model.ScenicSpotSummary, 0)
	for i := range res.ScenicSpotList {
		scenicID := originList[i].ScenicID
		waitTime, ok := waitTimeData[scenicID]
		if !ok {
			continue
		}
		originList[i].WaitTime = gptr.Of(waitTime)
		scenicSpotList = append(scenicSpotList, originList[i])
	}

	res.ScenicSpotList = scenicSpotList

	// 记录缓存信息用于调试
	if cacheInfo, exists := cacheManager.GetCacheInfo(scId); exists {
		log.Printf("景区:%s等待时间数据获取成功，缓存更新时间:%s，景点数量:%d",
			scId, cacheInfo.UpdatedAt.Format("2006-01-02 15:04:05"), len(scenicSpotList))
	}

	return res, nil
}

func GetScenicSpotWaitTimeList(scId string, date string) ([]wait_time.ScenicSpotWaitTime, error) {
	records, err := SelectHistory[wait_time.ScenicSpotWaitTime](date, scId, BizTypeScenic, "15:04")
	if err != nil {
		log.Printf("获取景点排队时间异常！scId:%s,dateStr:%s,%v\n", scId, date, err)
		return nil, err
	}
	return records, nil
}
