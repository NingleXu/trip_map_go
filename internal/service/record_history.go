package service

import (
	"encoding/json"
	"fmt"
	"time"
	"trip-map/core/wait_time"
	"trip-map/global"
	"trip-map/internal/model"
)

// BizType 枚举
type BizType int

const (
	BizTypeUnknown BizType = 0
	BizTypeScenic  BizType = 1
)

// SelectHistory 查询某天的记录，支持泛型返回
func SelectHistory[T any](date string, cId string, bizType BizType, timeLayout string) ([]T, error) {
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	end := start.Add(24 * time.Hour)

	var records []model.RecordHistory
	if err := global.Db.Where("biz_id = ? AND record_time >= ? AND record_time < ?", cId, start, end).
		Order("record_time ASC").
		Find(&records).Error; err != nil {
		return nil, err
	}

	var result []T
	for _, r := range records {
		var obj T
		switch bizType {
		case BizTypeScenic:
			if err := json.Unmarshal([]byte(r.JSONData), &obj); err != nil {
				return nil, fmt.Errorf("json unmarshal error: %w", err)
			}
			// 如果 T 是 ScenicSpotWaitTime，可以额外处理时间
			if v, ok := any(&obj).(*wait_time.ScenicSpotWaitTime); ok && timeLayout != "" {
				v.StatsTime = r.RecordTime.Format(timeLayout)
			}
		default:
			// 如果用户传的 T 不是 string，这里会 panic，所以要求 T=string 才能接收原始 json
			if v, ok := any(&obj).(*string); ok {
				*v = r.JSONData
			} else {
				return nil, fmt.Errorf("unsupported type for BizType=%d", bizType)
			}
		}
		result = append(result, obj)
	}
	return result, nil
}

// RecordHistory 记录日志
func RecordHistory(bizID string, bizType BizType, data interface{}) error {
	// 转 JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	record := model.RecordHistory{
		BizID:    bizID,
		BizType:  int8(bizType),
		JSONData: string(jsonBytes),
		IsDelete: 0,
	}

	if err := global.Db.Create(&record).Error; err != nil {
		return fmt.Errorf("insert record error: %w", err)
	}
	return nil
}

// BatchRecord 批量记录日志
// dataMap: key 是 bizID，value 是相同类型的数据切片
func BatchRecord[T any](bizType BizType, dataMap map[string]T) error {
	var records []model.RecordHistory

	for bizID, d := range dataMap {
		bs, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("json marshal error for bizID=%s: %w", bizID, err)
		}

		records = append(records, model.RecordHistory{
			BizID:    bizID,
			BizType:  int8(bizType),
			JSONData: string(bs),
			IsDelete: 0,
		})
	}

	if len(records) == 0 {
		return nil
	}

	if err := global.Db.Create(&records).Error; err != nil {
		return fmt.Errorf("batch insert error: %w", err)
	}
	return nil
}
