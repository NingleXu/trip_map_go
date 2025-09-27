package utils

import (
	"math"
	"strconv"
	"strings"
	"trip-map/internal/model/common"
)

const XPi = math.Pi * 3000.0 / 180.0

func BaiduCoordinateStrConverter(pointStr string) *common.Point {
	// 去除字符串前后的空格
	pointStr = strings.TrimSpace(pointStr)

	// 分割字符串
	parts := strings.Split(pointStr, ",")
	if len(parts) != 2 {
		return nil
	}

	// 解析经度
	lngStr := strings.TrimSpace(parts[0])
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return nil
	}

	// 解析纬度
	latStr := strings.TrimSpace(parts[1])
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil
	}

	// 返回Point结构体指针
	return BD09ToGCJ02(&common.Point{
		Lat: lat,
		Lng: lng,
	})
}

// BD09ToGCJ02 百度坐标系(BD09)转腾讯地图坐标系(GCJ02)
func BD09ToGCJ02(bdPoint *common.Point) *common.Point {
	if bdPoint == nil {
		return nil
	}

	x := bdPoint.Lng - 0.0065
	y := bdPoint.Lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*XPi)
	ggLng := z * math.Cos(theta)
	ggLat := z * math.Sin(theta)

	return &common.Point{
		Lng: ggLng,
		Lat: ggLat,
	}
}
