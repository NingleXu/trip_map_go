package service

import (
	"log"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
)

func GetScenicAreaListByPage(req request.ScenicAreaPageRequest) (*response.PageResponseData[model.ScenicSpotSummary], error) {
	var scenicSpots []model.ScenicSpot
	var total int64

	err := global.Db.
		Where("is_delete = 0 and is_scenic_area = 1").
		Model(&model.ScenicSpot{}).
		Count(&total).
		Error
	if err != nil {
		log.Printf("查询分页查询景区总数失败%v\n", err)
		return response.PageResponseWithEmpty[model.ScenicSpotSummary](&req.PageInfo), err
	}
	err = global.Db.
		Where("is_delete = 0 and is_scenic_area = 1").
		Offset(req.Offset).
		Limit(req.PageSize).
		Find(&scenicSpots).
		Error
	if err != nil {
		log.Printf("查询分页查询景区失败%v\n", err)
		return response.PageResponseWithEmpty[model.ScenicSpotSummary](&req.PageInfo), err
	}

	return response.PageResponseWithData[model.ScenicSpotSummary](convertScenicSpots(scenicSpots), total, &req.PageInfo), err
}
