package service

import (
	"errors"
	"fmt"
	"github.com/bytedance/gg/gslice"
	"gorm.io/gorm"
	"log"
	"strconv"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/request"
)

func UpdateScenicAreaScenic(req *request.ScenicAreaScenicPointUpdateRequest) error {
	if req.ScenicAreaCId == "" {
		return errors.New("scenic area id is required")
	}

	// 在事务中执行
	err := global.Db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询景区信息
		ssi, err := GetScenicSpotListByScenicAreaCId(req.ScenicAreaCId)
		if err != nil {
			return errors.New("scenic area found error")
		}
		originScenicSpotCIds := gslice.Map[model.ScenicSpotSummary, string](ssi.ScenicSpotList,
			func(s model.ScenicSpotSummary) string {
				return s.ScenicID
			})

		// 2. 先移除旧关联
		removeResult := tx.Model(&model.ScenicSpot{}).
			Where("c_id IN ?", originScenicSpotCIds).
			Update("belong_scenic_area_c_id", "0")

		if removeResult.Error != nil {
			return errors.New("failed to remove old scenic area associations")
		}

		// 3. 再更新新的景区关联
		updateResult := tx.Model(&model.ScenicSpot{}).
			Where("c_id IN ?", req.ScenicPointCIds).
			Update("belong_scenic_area_c_id", req.ScenicAreaCId)

		if updateResult.Error != nil {
			return errors.New("failed to update new scenic area associations")
		}

		// 检查是否真的更新了数据
		if updateResult.RowsAffected == 0 {
			return errors.New("no records updated")
		}

		return nil // 返回 nil 提交事务
	})

	return err
}

func GetScenicSpotInfoByManage(cId string) (*model.ScenicSpot, error) {
	return GetScenicSpotInfoByCId(cId), nil
}

func SaveScenicSpotByManage(ss *model.ScenicSpot) error {
	// 获取城市详情
	cityInfo, err := GetCityByCityCode(ss.QCityId)
	if err != nil {
		return fmt.Errorf("获取城市详情失败")
	}
	ss.QCityId = cityInfo.QCityId

	// 查询当前最大的 cid 然后自增 + 1
	scenicSpots := model.ScenicSpot{}
	err = global.Db.Where("is_delete = ?  order by c_id desc limit 1", 0).
		Find(&scenicSpots).Error
	if err != nil {
		log.Printf("SaveScenicSpotByManage error, %v\n", err)
		return err
	}
	// 先将字符串 ID 转换为数字
	idNum, err := strconv.ParseInt(scenicSpots.CId, 10, 64)
	if err != nil {
		// 处理转换错误，可能是 ID 不是数字
		return err
	}

	// 数字加 1 后再转回字符串
	ss.CId = strconv.FormatInt(idNum+1, 10)
	return SaveScenicSpot(ss)
}

func UpdateScenicSpotInfoByManage(ss *model.ScenicSpot) error {
	// 获取城市详情
	cityInfo, err := GetCityByCityCode(ss.QCityId)
	if err != nil {
		return fmt.Errorf("获取城市详情失败")
	}
	ss.QCityId = cityInfo.QCityId

	return global.Db.
		Where("c_id = ?", ss.ScenicId).
		Updates(ss).
		Error
}
