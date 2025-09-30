package model

import "time"

type RecordHistory struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BizID      string    `gorm:"column:biz_id;type:varchar(255);not null;default:''" json:"biz_id"`
	BizType    int8      `gorm:"column:biz_type;type:tinyint;not null;default:0" json:"biz_type"`
	JSONData   string    `gorm:"column:json_data;type:text;not null" json:"json_data"`
	RecordTime time.Time `gorm:"column:record_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"record_time"`
	IsDelete   int       `gorm:"column:is_delete;type:int;not null;default:0" json:"is_delete"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 指定表名
func (RecordHistory) TableName() string {
	return "tm_record_history"
}
