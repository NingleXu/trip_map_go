package response

import "trip-map/internal/model"

type CityScenicSpotRes struct {
	City           *CityLocation             `json:"city"`
	ShowDefault    bool                      `json:"showDefault"`
	ScenicSpotList []model.ScenicSpotSummary `json:"scenicSpotList"`
}

type ScenicAreaScenicSpotRes struct {
	City           *CityLocation             `json:"city"`
	ScenicArea     model.ScenicSpotSummary   `json:"scenicArea"`
	ScenicSpotList []model.ScenicSpotSummary `json:"scenicSpotList"`
}
