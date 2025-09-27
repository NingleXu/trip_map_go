package model

import (
	"time"
	"trip-map/core/wait_time"
	"trip-map/internal/model/common"
)

type ScenicSpot struct {
	Id                  int
	QCityId             string
	CId                 string
	CName               string
	CIntro              string
	CAddr               string
	CImg                string
	CRatingScore        string
	CPriceNumber        string
	CTag                string
	CSightLevel         string
	TId                 string
	TSightOpenTime      string
	GoogleMapPoint      string
	BaiduMapPoint       string
	TSightLevel         string
	TSightIntroduction  string
	TAddress            string
	TTraffic            string
	TCoverImage         string
	TImages             []string      `gorm:"-"`
	Point               *common.Point `gorm:"-" json:"point"`
	ScenicId            string        `gorm:"-" json:"scenicId"`
	IsDelete            int
	CreateTime          time.Time `gorm:"autoCreateTime"`
	UpdateTime          time.Time `gorm:"autoUpdateTime"`
	PageViews           int64
	CityName            string
	IsScenicArea        int
	BelongScenicAreaCId string
	WaitTimeEnable      int
}

// TableName 指定表名
func (ScenicSpot) TableName() string {
	return "tm_scenic_spots_v2"
}

type ScenicSpotSummary struct {
	ScenicID         string                        `json:"scenicId"`
	ScenicName       string                        `json:"scenicName"`
	ScenicProfile    string                        `json:"scenicProfile"`
	ScenicFeature    string                        `json:"scenicFeature"`
	ScenicReputation string                        `json:"scenicReputation"`
	ScenicType       string                        `json:"scenicType"`
	ScenicCover      string                        `json:"scenicCover"`
	ScenicCoordinate *common.Point                 `json:"scenicCoordinate"`
	PageViews        int64                         `json:"pageViews"`
	IsScenicArea     int                           `json:"isScenicArea"`
	CityCode         string                        `json:"cityCode"`
	CityName         string                        `json:"cityName"`
	WaitTimeEnable   int                           `json:"waitTimeEnable"`
	WaitTime         *wait_time.ScenicSpotWaitTime `json:"waitTime"`
}
