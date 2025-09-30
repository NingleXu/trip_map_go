package wait_time

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BeijingUniversalStudios struct {
	BaseHandler
}

// BeijingUniversalAPIResponse 定义北京环球影城API响应结构
type BeijingUniversalAPIResponse struct {
	Ret       int    `json:"ret"`
	Msg       string `json:"msg"`
	ErrorCode int    `json:"errorcode"`
	Data      struct {
		Pagination struct {
			Total       int `json:"total"`
			TotalPage   int `json:"total_page"`
			PageSize    int `json:"page_size"`
			CurrentPage int `json:"current_page"`
		} `json:"pagination"`
		List []struct {
			ID              string `json:"id"`
			MaterialID      string `json:"material_id"`
			MaterialType    string `json:"material_type"`
			GemsStatus      string `json:"gems_status"`
			LocationID      string `json:"location_id"`
			GemsCopywriting string `json:"gems_copywriting"`
			WaitingTime     int    `json:"waiting_time"`
			NextTime        string `json:"next_time"`
			ServiceTime     struct {
				Open  string `json:"open"`
				Close string `json:"close"`
			} `json:"service_time"`
			ShowTime         string `json:"show_time"`
			Title            string `json:"title"`
			Subtitle         string `json:"subtitle"`
			Label            string `json:"label"`
			MapLabel         string `json:"map_label"`
			CustomLabel      string `json:"custom_label"`
			CoverImage       string `json:"cover_image"`
			IsSupportExpress int    `json:"is_support_express"`
			ThrillingDegree  int    `json:"thrilling_degree"`
			ListSort         int    `json:"list_sort"`
			Area             string `json:"area"`
			IsClosed         int    `json:"is_closed"`
			CollectCount     int    `json:"collect_count"`
			Position         struct {
				Longitude string `json:"longitude"`
				Latitude  string `json:"latitude"`
				Address   string `json:"address"`
				Distance  string `json:"distance"`
				PoiID     string `json:"poi_id"`
			} `json:"position"`
			Indoor     bool   `json:"indoor"`
			ShowIndoor string `json:"show_indoor"`
			Favourited bool   `json:"favourited"`
		} `json:"list"`
	} `json:"data"`
}

func (h *BeijingUniversalStudios) GetWaitTime() (map[string]ScenicSpotWaitTime, error) {
	// 1. 创建 HTTP 请求
	req, err := http.NewRequest(
		"GET",
		"https://g.app.universalbeijingresort.com/map/attraction/list?mode=list",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 2. 设置请求头
	req.Header.Set("Host", "g.app.universalbeijingresort.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("MiniappToken", "")
	req.Header.Set("Language", "ch")
	req.Header.Set("appVersion", "4.8.3")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Lat", "22.569582248263888")
	req.Header.Set("Lng", "113.857060546875")
	req.Header.Set("USERAREA", "other")
	req.Header.Set("X-WECHAT-HOSTSIGN", `{"noncestr":"471f745ae7746231cd6256379e53315a","timestamp":1758604302,"signature":"5bac8a95a7fd6ce4c8bb3feadb8583bf9a938694"}`)
	req.Header.Set("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.57(0x18003932) NetType/WIFI Language/zh_CN")
	req.Header.Set("Referer", "https://servicewechat.com/wx3ba512d53df66a75/72/page-frame.html")

	// 3. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 4. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 5. 解析响应数据
	var apiResponse BeijingUniversalAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if apiResponse.Ret != 0 {
		return nil, fmt.Errorf("API返回错误: %s", apiResponse.Msg)
	}

	// 6. 转换为目标结构
	result := make(map[string]ScenicSpotWaitTime)
	for _, attraction := range apiResponse.Data.List {
		// 尝试转换ID
		scenicSpotId, err := h.ConvertID(attraction.ID)
		if err != nil {
			continue
		}

		waitTime := ScenicSpotWaitTime{
			ScenicSpotId: scenicSpotId,
			StartTime:    attraction.ServiceTime.Open,
			EndTime:      attraction.ServiceTime.Close,
		}

		// 解析等待时间和开放状态
		waitTime.WaitMinute, waitTime.OpenState = h.parseWaitTimeAndStatus(
			attraction.WaitingTime,
			attraction.IsClosed,
			attraction.GemsCopywriting,
		)

		result[scenicSpotId] = waitTime
	}

	return result, nil
}

// parseWaitTimeAndStatus 解析等待时间和开放状态
func (h *BeijingUniversalStudios) parseWaitTimeAndStatus(waitingTime int, isClosed int, gemsCopywriting string) (int64, OpenState) {
	// 判断是否关闭
	if isClosed == 1 {
		return 0, OpenStateClosed
	}

	// 根据gems_copywriting判断状态
	switch gemsCopywriting {
	case "关闭":
		return 0, OpenStateClosed
	case "维护", "暂停":
		return 0, OpenStateMaintain
	}

	// 解析等待时间
	if waitingTime < 0 {
		// -1 通常表示没有等待时间数据或关闭
		return 0, OpenStateUnknown
	}

	if waitingTime == 0 {
		// 0 表示无需等待
		return 0, OpenStateOpen
	}

	// 正数表示等待时间（分钟）
	return int64(waitingTime), OpenStateOpen
}

func NewBeijingUniversalStudios() *BeijingUniversalStudios {
	return &BeijingUniversalStudios{
		BaseHandler: BaseHandler{
			// 暂时设置为空映射，后续完善
			// 格式: "环球影城景点ID": "系统内部景点ID"
			IDMapping: map[string]string{
				// 变形金刚：火种源争夺战 - 匹配响应体id与自定义CID 100000037
				"5fa247615ba7f1491f6289a2": "100000037",
				// 哈利·波特与禁忌之旅{1} - 匹配响应体id与自定义CID 100000038（名称中{1}为原数据格式，不影响匹配）
				"5fa22624dcb5e53c6c72ab42": "100000038",
				// 神偷奶爸小黄人闹翻天 - 匹配响应体id与自定义CID 100000039
				"5f914afc26509774c642cf42": "100000039",
				// 奥利凡德{1} - 匹配响应体id与自定义CID 100000041（名称中{1}为原数据格式，不影响匹配）
				"6121b15723ccd253665b6822": "100000041",
				// 灯光，摄像，开拍！ - 匹配响应体id与自定义CID 100000042（响应体名称含逗号，自定义CID名称无逗号，为同一景点）
				"611f0913e291ec2b16214a88": "100000042",
				// 灯影传奇 - 匹配响应体id与自定义CID 100000043
				"60f97e7286449732e736681c": "100000043",
				// 炫转武侠 - 匹配响应体id与自定义CID 100000044
				"60f97dfceb33e120fe66fc3e": "100000044",
				// 大黄蜂回旋机 - 匹配响应体id与自定义CID 100000045
				"60f976da05482c2ea6737a9c": "100000045",
				// 飞越侏罗纪 - 匹配响应体id与自定义CID 100000046
				"5fa24b84b9227045895c7495": "100000046",
				// 超萌漩漩涡 - 匹配响应体id与自定义CID 100000047
				"5fa23cdcbbdbd416ae43dc05": "100000047",
				// 萌转过山车 - 匹配响应体id与自定义CID 100000048
				"5fa22dab40a74131c80b8a72": "100000048",
				// 功夫熊猫:神龙大侠之旅 - 匹配响应体id与自定义CID 100000049
				"5f9161f64b89d447253b02c5": "100000049",
				// 鹰马飞行{1} - 匹配响应体id与自定义CID 100000050（名称中{1}为原数据格式，不影响匹配）
				"5f8812ec14ae8a2d80450794": "100000050",
				// 奇遇迅猛龙 - 匹配响应体id与自定义CID 100000051
				"5f91438e16d92d317523b6a4": "100000051",
				// 阿宝功夫训练营 - 匹配响应体id与自定义CID 100000052
				"5f912fd8d72a471ef6359d02": "100000040",
			},
		},
	}
}
