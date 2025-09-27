package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"trip-map/core"
	"trip-map/core/wait_time"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/request"
	"trip-map/internal/model/response"
	"trip-map/internal/utils"
)

// GetScenicSpotByQCityId
// @description: 通过q_city_id查询景点
func GetScenicSpotByQCityId(qCityId string, getTop30 bool) ([]model.ScenicSpot, error) {
	var scenicSpots []model.ScenicSpot
	err := global.Db.Where("q_city_id = ?", qCityId).Find(&scenicSpots).Error
	if err != nil {
		log.Printf("查询qCityId为%s的景点失败:\n", qCityId)
		return nil, err
	}
	return SortScenicSpots(scenicSpots, getTop30), nil
}

func GetScenicSpotInfoByCId(cId string) *model.ScenicSpot {
	info := model.ScenicSpot{}
	err := global.Db.Where("c_id = ?", cId).First(&info).Error
	if err != nil {
		log.Println("无法找到c_id为" + cId + "的景点记录")
		return nil
	}
	info.Point = utils.BaiduCoordinateStrConverter(info.BaiduMapPoint)

	go addScenicSpotsViewCount(cId)

	return &info
}

// CaptureCityScenicSpot
// 循环去抓取 然后保存数据库
func CaptureCityScenicSpot(req request.CaptureCityScenicSpotReq) error {

	for _, val := range req.CaptureCityScenicSpotList {
		// 查询 城市 如果存在，则补充上 qCityId
		if !strings.HasSuffix(val.CityName, "市") {
			val.CityName = val.CityName + "市"
		}
		city, err := GetCityByName(val.CityName)
		if err != nil {
			log.Printf("ERROR:获取城市出现异常: %v", err)
			continue
		}
		if nil == city {
			log.Printf("城市:%s不存在", val.CityName)
			continue
		}
		//
		//if err := UpdateQCityId(city.ID, val.DestId); err != nil {
		//	log.Printf("ERROR: updateQCityId failed: %v", err)
		//	continue
		//}
		// 执行成功

		log.Printf("==========开始处理城市:%s", val.CityName)
		scenicSpots, err := handleCapture(val.DestId)
		if err != nil {
			log.Printf("ERROR: handleCapture failed: %v", err)
			return err
		}
		log.Printf("城市:%s,抓取到景点数量: %d", city.CityName, len(scenicSpots))

		// 创建带随机种子的 rand 对象
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// 对每景点进行操作
		for _, scenicSpot := range scenicSpots {
			log.Printf("-----开始处理景点:%s", scenicSpot.Name)
			handleSingleScenicSpot(&scenicSpot, city)

			// 每个景点之间间隔 100~600随机
			delay := time.Duration(200+r.Intn(500)) * time.Millisecond
			time.Sleep(delay)
		}
	}

	return nil
}

type PoiResponse struct {
	Ret     bool   `json:"ret"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Ver     int    `json:"ver"`
	Data    struct {
		More      bool        `json:"more"`
		DistLevel int         `json:"distLevel"`
		List      []PoiDetail `json:"list"`
	} `json:"data"`
	ExtraParamMap map[string]string `json:"extraParamMap"`
}

type PoiDetail struct {
	DistId        int     `json:"distId"`
	Name          string  `json:"name"`
	Intro         string  `json:"intro"`
	Addr          string  `json:"addr"`
	Image         string  `json:"image"`
	RatingScore   float64 `json:"ratingScore"`
	PriceNumber   float64 `json:"priceNumber"`
	CommentCount  int     `json:"commentCount"`
	Tag           string  `json:"tag"`
	SightLevel    string  `json:"sightLevel"`
	DistName      string  `json:"distName"`
	Scheme        string  `json:"scheme"`
	Id            int     `json:"id"`
	SubType       int     `json:"subType"`
	SmartReferCnt int     `json:"smartReferCount"`
	Blat          float64 `json:"blat"`
	Blng          float64 `json:"blng"`
	Glat          float64 `json:"glat"`
	Glng          float64 `json:"glng"`
}

func handleCapture(destId string) ([]PoiDetail, error) {
	client := &http.Client{}
	baseUrl := "https://hy.travel.qunar.com/api/poi/search"
	page := 1
	allList := make([]PoiDetail, 0, 60)

	for page <= 2 { // 最多拉两页
		url := fmt.Sprintf("%s?RN=1&destType=1&destId=%s&type=4&useEs=true&needTopFeelingTag=true&from_page=city_new_rn&dir=desc&limit=30&page=%d",
			baseUrl, destId, page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var poiResp PoiResponse
		if err := json.Unmarshal(body, &poiResp); err != nil {
			return nil, err
		}

		allList = append(allList, poiResp.Data.List...)

		// 没有更多就退出
		if !poiResp.Data.More {
			break
		}
		time.Sleep(300 * time.Millisecond)
		page++
	}

	return allList, nil
}

func handleSingleScenicSpot(poiDetail *PoiDetail, city *model.City) {
	// 抓取 内容的帖子信息
	nodeList, err := core.GetNoteListByPoiId(strconv.Itoa(poiDetail.Id), poiDetail.Name)
	if err != nil {
		log.Printf("ERROR: 城市:%s,景点名称:%s,景点id:%d,内容接口抓帖子失败failed: %v",
			city.CityName,
			poiDetail.Name,
			poiDetail.Id,
			err)
	}
	if len(nodeList) == 0 {
		log.Printf("ERROR: 城市:%s,景点名称:%s,景点id:%d,内容接口抓帖子为空",
			city.CityName,
			poiDetail.Name,
			poiDetail.Id)
	} else {
		// 保存笔记
		if err := BatchSaveNodeList(&nodeList); err != nil {
			log.Printf("ERROR: 保存笔记列表失败:%v\n", err)
		} else {
			// 打印日志
			log.Printf("- 成功抓帖子数量为:%d \n", len(nodeList))
		}
	}

	// 抓取门票详情
	tScenicSpotInfo, err := core.GetTicketScenicSpotInfo(strconv.Itoa(poiDetail.Id))
	if err != nil {
		log.Printf("ERROR: 内容转门票信息失败:%v\n", err)
		return
	}
	if tScenicSpotInfo == nil {
		log.Printf("内容景点无法找到对应的门票信息\n")
		tScenicSpotInfo = &core.SightInfo{}
	}

	// 补充基本信息
	scenicSpot := model.ScenicSpot{
		QCityId:            city.QCityId,
		CId:                strconv.Itoa(poiDetail.Id),
		CName:              poiDetail.Name,
		CIntro:             poiDetail.Intro,
		CAddr:              poiDetail.Addr,
		CImg:               poiDetail.Image,
		CRatingScore:       strconv.FormatFloat(poiDetail.RatingScore, 'f', 2, 64),
		CPriceNumber:       strconv.FormatFloat(poiDetail.PriceNumber, 'f', 2, 64),
		CTag:               poiDetail.Tag,
		CSightLevel:        poiDetail.SightLevel,
		BaiduMapPoint:      fmt.Sprintf("%f,%f", poiDetail.Blng, poiDetail.Blat),
		GoogleMapPoint:     fmt.Sprintf("%f,%f", poiDetail.Glng, poiDetail.Glat),
		TId:                strconv.Itoa(tScenicSpotInfo.SightId),
		TSightOpenTime:     tScenicSpotInfo.SightOpenTime,
		TSightLevel:        tScenicSpotInfo.SightLevel,
		TSightIntroduction: tScenicSpotInfo.SightIntro,
		TAddress:           tScenicSpotInfo.Address,
		TTraffic:           tScenicSpotInfo.Traffic,
		TCoverImage:        tScenicSpotInfo.CoverImage,
	}
	if err := SaveScenicSpot(&scenicSpot); err != nil {
		log.Printf("ERROR: 保存景点信息失败:%v\n", err)
		return
	} else {
		log.Printf("- 成功保存景点信息\n")
	}

	// 保存图片信息
	tScenicSpotImg := make([]model.Image, len(tScenicSpotInfo.BigImages))
	for idx, imgUrl := range tScenicSpotInfo.BigImages {
		tScenicSpotImg[idx] = model.Image{
			ScenicSpotId: strconv.Itoa(poiDetail.Id),
			Url:          imgUrl,
		}
	}
	if len(tScenicSpotImg) == 0 {
		log.Printf("ERROR: 城市:%s,景点名称:%s,景点id:%d,景点图片为空",
			city.CityName,
			poiDetail.Name,
			poiDetail.Id)
	} else {
		// 保存图片信息
		if err := BatchSaveScenicSpotImages(&tScenicSpotImg); err != nil {
			log.Printf("ERROR: 保存景点图片失败:%v\n", err)
		} else {
			log.Printf("- 成功保存图片%d张\n", len(tScenicSpotImg))
		}
	}
}

func SaveScenicSpot(s *model.ScenicSpot) error {
	return global.Db.Save(s).Error
}

func GetScenicSpotListByUserLocation(req *request.UserLocationScenicSpotRequest) (*response.CityScenicSpotRes, error) {
	// 查询用户坐标对应的CityCode
	cityCode, cityCodeError := GetCityCodeByUserPosition(req.Point.Lat, req.Point.Lng)
	if cityCodeError != nil {
		// 默认改为北京
		cityCode = "156110000"
	}

	// 获取城市详情
	cityInfo, err := GetCityByCityCode(cityCode)
	if err != nil {
		return nil, fmt.Errorf("获取城市详情失败")
	}

	// 获取城市对应的景点
	scenicSpots, err := GetScenicSpotByQCityId(cityInfo.QCityId, true)
	if err != nil {
		return nil, fmt.Errorf("获取城市对应景点失败")
	}
	return &response.CityScenicSpotRes{
		City: Convert2CityLocation(cityInfo), ShowDefault: cityCodeError != nil, ScenicSpotList: convertScenicSpots(scenicSpots),
	}, nil
}

func GetScenicSpotListByCityCode(cityCode string, getTop30 bool) (*response.CityScenicSpotRes, error) {
	// 获取城市详情
	cityInfo, err := GetCityByCityCode(cityCode)
	if err != nil {
		return nil, fmt.Errorf("获取城市详情失败")
	}

	// 获取城市对应的景点
	scenicSpots, err := GetScenicSpotByQCityId(cityInfo.QCityId, getTop30)
	if err != nil {
		return nil, fmt.Errorf("获取城市对应景点失败")
	}
	return &response.CityScenicSpotRes{
		City: &response.CityLocation{
			Key:            cityInfo.Key,
			CityCode:       cityInfo.CityCode,
			CityName:       cityInfo.CityName,
			CityCover:      cityInfo.CityCover,
			CityCoordinate: utils.BaiduCoordinateStrConverter(cityInfo.CityCoordinate),
			IsHide:         cityInfo.IsHide,
		}, ShowDefault: false, ScenicSpotList: convertScenicSpots(scenicSpots),
	}, nil
}

// convertScenicSpots 将ScenicSpot切片转换为ScenicSpotSummary切片
func convertScenicSpots(spots []model.ScenicSpot) []model.ScenicSpotSummary {
	// 初始化目标切片，预分配内存以提高性能
	summaries := make([]model.ScenicSpotSummary, 0, len(spots))

	// 遍历原始切片，逐个转换并添加到目标切片
	for _, spot := range spots {
		summaries = append(summaries, convert2ScenicSpotSummary(&spot, nil))
	}

	return summaries
}
func convert2ScenicSpotSummary(spot *model.ScenicSpot, waitTime *wait_time.ScenicSpotWaitTime) model.ScenicSpotSummary {
	scenicCover := spot.TCoverImage
	if scenicCover == "" {
		scenicCover = spot.CImg
	}

	summary := model.ScenicSpotSummary{
		ScenicID:         spot.CId,
		ScenicName:       spot.CName,
		ScenicProfile:    spot.CIntro,
		ScenicFeature:    spot.CTag,
		ScenicReputation: spot.CSightLevel,
		ScenicType:       spot.CTag,
		ScenicCover:      scenicCover,
		ScenicCoordinate: utils.BaiduCoordinateStrConverter(spot.BaiduMapPoint),
		PageViews:        spot.PageViews,
		IsScenicArea:     spot.IsScenicArea,
		CityName:         spot.CityName,
		WaitTimeEnable:   spot.WaitTimeEnable,
		WaitTime:         waitTime,
	}
	return summary
}

func GetTop10ScenicSpotList() ([]model.ScenicSpotSummary, error) {
	var scenicSpots []model.ScenicSpot
	err := global.Db.Where("is_delete = ?", 0).
		Order("page_views desc").
		Limit(10).
		Find(&scenicSpots).Error
	if err != nil {
		return nil, err
	}
	return convertScenicSpots(scenicSpots), nil
}

func GetScenicSpotByKeyword(keyword string) ([]model.ScenicSpotSummary, error) {
	var scenicSpots []model.ScenicSpot
	err := global.Db.Where("is_delete = ? AND c_name LIKE ?", 0, "%"+keyword+"%").
		Find(&scenicSpots).Error
	if err != nil {
		return nil, err
	}

	scenicSpots = SortScenicSpots(scenicSpots, false)
	return convertScenicSpots(scenicSpots), nil
}

func SortScenicSpots(spots []model.ScenicSpot, getTop30 bool) []model.ScenicSpot {
	// 排序：先 TId 降序，其次 CSightLevel 降序，最后 CRatingScore 降序
	sort.Slice(spots, func(i, j int) bool {
		// 规则1: TId = 0 的排在最后
		if spots[i].TId == "0" && spots[j].TId != "0" {
			return false
		}
		if spots[i].TId != "0" && spots[j].TId == "0" {
			return true
		}

		if spots[i].CSightLevel != spots[j].CSightLevel {
			return spots[i].CSightLevel > spots[j].CSightLevel
		}

		if spots[i].PageViews != spots[j].PageViews {
			return spots[i].PageViews > spots[j].PageViews
		}

		return spots[i].CRatingScore > spots[j].CRatingScore
	})

	// 取前 30 个
	if getTop30 && len(spots) > 30 {
		return spots[:30]
	}
	return spots
}

func addScenicSpotsViewCount(cId string) {
	err := global.Db.Model(&model.ScenicSpot{}).
		Where("c_id = ?", cId).
		UpdateColumn("page_views", gorm.Expr("page_views + ?", 1)).Error

	if err != nil {
		log.Printf("addScenicSpotsViewCount error, %v\n", err)
	}
}

func GetScenicSpotListByScenicAreaCId(scId string) (*response.ScenicAreaScenicSpotRes, error) {
	// 查找景区详情
	scenicAreaInfo := GetScenicSpotInfoByCId(scId)
	if nil == scenicAreaInfo {
		return nil, errors.New("scenic area not found")
	}
	qCityId := scenicAreaInfo.QCityId
	// 查询城市详情
	cityInfo, err := GetCityByQCityId(qCityId)
	if err != nil {
		return nil, err
	}
	var scenicSpots []model.ScenicSpot
	// 查询附属的景点列表
	err = global.Db.
		Where("is_delete = ? and belong_scenic_area_c_id = ?", 0, scenicAreaInfo.CId).
		Find(&scenicSpots).Error
	if err != nil {
		return nil, err
	}

	return &response.ScenicAreaScenicSpotRes{
		City:           Convert2CityLocation(cityInfo),
		ScenicArea:     convert2ScenicSpotSummary(scenicAreaInfo, nil),
		ScenicSpotList: convertScenicSpots(scenicSpots),
	}, nil
}
func GetScenicAreaListByCityCode(cityCode string) ([]model.ScenicSpotSummary, error) {
	// 获取城市详情
	cityInfo, err := GetCityByCityCode(cityCode)
	if err != nil {
		return nil, fmt.Errorf("获取城市详情失败")
	}

	var scenicSpots []model.ScenicSpot
	err = global.Db.Where("is_scenic_area = 1 AND is_delete = ? AND q_city_id = ? ", 0, cityInfo.QCityId).
		Find(&scenicSpots).Error
	if err != nil {
		return nil, err
	}

	scenicSpots = SortScenicSpots(scenicSpots, false)
	return convertScenicSpots(scenicSpots), nil
}
