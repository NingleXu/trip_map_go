package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/response"
	"trip-map/internal/utils"
)

const GeoParseUrl string = "https://apis.map.qq.com/ws/geocoder/v1/?key=W3EBZ-7YS6H-HIMDU-WPEK2-IARMJ-U7FU4&location="

// GetCityList 查询所有城市列表
func GetCityList() ([]response.CityKeyList, error) {
	var cities []model.City
	if err := global.Db.Where(" q_city_id != ''").Find(&cities).Error; err != nil {
		return nil, err
	}
	// 创建一个map来按key分组
	keyMap := make(map[string][]response.CityItem)
	// 遍历城市列表，按key分组
	for _, city := range cities {
		item := response.CityItem{
			Id:   city.CityCode,
			Name: city.CityName,
		}
		keyMap[city.Key] = append(keyMap[city.Key], item)
	}
	// 将map转换为切片
	var result []response.CityKeyList
	for key, list := range keyMap {
		result = append(result, response.CityKeyList{
			Key:  key,
			List: list,
		})
	}
	// 按key字母顺序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})

	return result, nil
}

// GetCityByName 根据城市名查询单个城市
func GetCityByName(name string) (*model.City, error) {
	var city model.City
	if err := global.Db.Where("city_name = ?", name).First(&city).Error; err != nil {
		return nil, err
	}
	return &city, nil
}

func UpdateQCityId(id int, qCityId string) error {
	return global.Db.Model(&model.City{}).
		Where("id = ?", id).
		Update("q_city_id", qCityId).Error
}

func GetCityByCityCode(cityCode string) (*model.City, error) {
	// 假设model.City结构体中CityCoordinate为string类型
	var city model.City

	err := global.Db.Table("tm_city").
		Select("*, ST_AsText(city_coordinate) as city_coordinate").
		Where("city_code = ? and q_city_id != ''", cityCode).
		First(&city).Error

	if err != nil {
		return nil, err
	}

	coordStr := convertPointStr(city.CityCoordinate)

	city.CityCoordinate = coordStr // 现在是"116.024067,40.362639"格式
	return &city, nil
}

func convertPointStr(coordStr string) string {
	re := regexp.MustCompile(`POINT\((\d+\.\d+) (\d+\.\d+)\)`)
	matches := re.FindStringSubmatch(coordStr)
	if len(matches) == 3 {
		lon := matches[1] // 经度
		lat := matches[2] // 纬度
		coordStr = lon + "," + lat
	}
	return coordStr
}

func GetCityCodeByUserPosition(lat, lng float64) (string, error) {
	url := fmt.Sprintf("%s%f,%f", GeoParseUrl, lat, lng)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get successful response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 使用结构体来解析JSON
	var response struct {
		Result struct {
			AdInfo struct {
				CityCode string `json:"city_code"`
			} `json:"ad_info"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if response.Result.AdInfo.CityCode == "" {
		return "", fmt.Errorf("city_code is empty")
	}

	return response.Result.AdInfo.CityCode, nil
}

func GetCityByQCityId(qCityId string) (*model.City, error) {
	// 假设model.City结构体中CityCoordinate为string类型
	var city model.City

	err := global.Db.Table("tm_city").
		Select("*, ST_AsText(city_coordinate) as city_coordinate").
		Where("q_city_id = ?", qCityId).
		First(&city).Error

	if err != nil {
		return nil, err
	}

	coordStr := convertPointStr(city.CityCoordinate)

	city.CityCoordinate = coordStr
	return &city, nil
}

func Convert2CityLocation(cityInfo *model.City) *response.CityLocation {
	return &response.CityLocation{
		Key:            cityInfo.Key,
		CityCode:       cityInfo.CityCode,
		CityName:       cityInfo.CityName,
		CityCover:      cityInfo.CityCover,
		CityCoordinate: utils.BaiduCoordinateStrConverter(cityInfo.CityCoordinate),
		IsHide:         cityInfo.IsHide,
	}
}
