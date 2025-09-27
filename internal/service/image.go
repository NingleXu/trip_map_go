package service

import (
	"trip-map/global"
	"trip-map/internal/model"
)

func GetScenicSpotImagesByCId(cId string) ([]model.Image, error) {
	images := make([]model.Image, 0)
	err := global.Db.Where("scenic_spot_id = ?", cId).Find(&images).Error
	if err != nil {
		return []model.Image{}, err
	}
	return images, nil
}

func BatchSaveScenicSpotImages(images *[]model.Image) error {
	return global.Db.Save(images).Error
}
