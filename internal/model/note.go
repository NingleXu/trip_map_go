package model

import "time"

type Note struct {
	Id              int
	QNoteId         string
	AuthorNickName  string
	AuthorHeadImg   string
	Title           string
	Body            string
	Images          string
	DetailUrl       string
	CPoiId          string
	CPoiName        string
	VideoUrl        string
	VideoCoverImage string
	VideoDuration   int
	UsefulCnt       int
	IsDelete        int
	CreateTime      time.Time `gorm:"autoCreateTime"`
	UpdateTime      time.Time `gorm:"autoUpdateTime"`
}

func (Note) TableName() string {
	return "tm_scenic_spots_note"
}
