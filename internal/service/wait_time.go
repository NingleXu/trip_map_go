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
	handler, ok := wait_time.ScenicAreaWaitTimeHandlerMap[scId]
	if !ok {
		log.Printf("景区:%s不存在等待时间\n", scId)
		return nil, fmt.Errorf("景区 %s 不存在", scId)
	}

	// channel 用来接收结果
	waitTimeCh := make(chan struct {
		data map[string]wait_time.ScenicSpotWaitTime
		err  error
	}, 1)

	// 异步获取等待时间
	go func() {
		m, err := handler.GetWaitTime()
		waitTimeCh <- struct {
			data map[string]wait_time.ScenicSpotWaitTime
			err  error
		}{data: m, err: err}
	}()

	// 同步查数据库
	res, err := GetScenicSpotListByScenicAreaCId(scId)
	if err != nil {
		log.Printf("获取景区:%s的景点信息失败...异常%v\n", scId, err)
		return nil, err
	}

	// 等待异步结果
	waitRes := <-waitTimeCh
	if waitRes.err != nil {
		log.Printf("获取景区:%s等待时间失败...异常%v\n", scId, waitRes.err)
		return nil, waitRes.err
	}

	// 合并结果
	originList := res.ScenicSpotList
	var scenicSpotList = make([]model.ScenicSpotSummary, 0)
	for i := range res.ScenicSpotList {
		scenicID := originList[i].ScenicID
		time, ok := waitRes.data[scenicID]
		if !ok {
			continue
		}
		originList[i].WaitTime = gptr.Of(time)
		scenicSpotList = append(scenicSpotList, originList[i])
	}

	res.ScenicSpotList = scenicSpotList
	return res, nil
}
