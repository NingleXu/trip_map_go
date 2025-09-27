package response

import "trip-map/internal/model/common"

type CityKeyList struct {
	Key  string     `json:"key"`
	List []CityItem `json:"list"`
}
type CityItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CityLocation struct {
	Key            string        `json:"key"`
	CityCode       string        `json:"cityCode"`
	CityName       string        `json:"cityName"`
	CityCover      string        `json:"cityCover"`
	CityCoordinate *common.Point `json:"cityCoordinate"`
	IsHide         int           `json:"isHide"`
}
