package model

import "time"

type City struct {
	ID             int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CityCode       string    `gorm:"column:city_code;type:varchar(64);not null;default:'';comment:城市code" json:"city_code"`
	CityName       string    `gorm:"column:city_name;type:varchar(32);not null;default:'';comment:城市名称" json:"city_name"`
	CityCoordinate string    `gorm:"column:city_coordinate;type:varchar(64);not null;comment:城市坐标" json:"city_coordinate"`
	CityCover      string    `gorm:"column:city_cover;type:varchar(4096);not null;comment:城市封面" json:"city_cover"`
	IsHide         int       `gorm:"column:is_hide;type:int;not null;default:0;comment:是否隐藏 1是0否" json:"is_hide"`
	IsDelete       int       `gorm:"column:is_delete;type:int;not null;default:0;comment:是否删除 1是0否" json:"is_delete"`
	CreateTime     time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"`
	Key            string    `gorm:"column:key;type:char(1);not null;default:'A';comment:拼音首字母" json:"key"`
	QCityId        string    `gorm:"column:q_city_id;type:varchar(10);not null;default:'';comment:qunar城市id" json:"q_city_id"`
}

// TableName 指定表名
func (City) TableName() string {
	return "tm_city"
}
