package model

import "time"

type Image struct {
	Id           int
	ScenicSpotId string
	Url          string
	IsDelete     int
	CreateTime   time.Time `gorm:"autoCreateTime"`
	UpdateTime   time.Time `gorm:"autoUpdateTime"`
}

func (Image) TableName() string {
	return "tm_scenic_spots_image"
}
