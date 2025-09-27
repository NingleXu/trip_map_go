package request

import (
	"trip-map/internal/model/common"
)

type UserLocationScenicSpotRequest struct {
	Point   *common.Point `json:"point"`
	IsFirst *bool         `json:"isFirst"`
}

type ScenicAreaPageRequest struct {
	common.PageInfo
}

type ScenicAreaScenicPointUpdateRequest struct {
	ScenicAreaCId   string   `json:"scenicAreaCId"`
	ScenicPointCIds []string `json:"scenicPointCIds"`
}
