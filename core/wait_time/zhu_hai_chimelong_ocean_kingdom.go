package wait_time

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ZhuHaiChimelongOceanKingdom struct {
	BaseHandler
}

// APIResponse 定义API响应结构
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ImageData []struct {
			Code        string `json:"code"`
			EndTime     string `json:"endTime"`
			ImageUrl    string `json:"imageUrl"`
			Name        string `json:"name"`
			StartTime   string `json:"startTime"`
			WaitingTime string `json:"waitingTime"`
		} `json:"imageData"`
		Instruction string `json:"instruction"`
	} `json:"data"`
	SubCode int    `json:"subCode"`
	Time    string `json:"time"`
	TrackId string `json:"trackId"`
}

func (h *ZhuHaiChimelongOceanKingdom) GetWaitTime() (map[string]ScenicSpotWaitTime, error) {
	// 1. 准备请求数据
	requestBody := map[string]string{"code": "ZH56"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 2. 创建 HTTP 请求
	req, err := http.NewRequest(
		"POST",
		"https://zx-api.chimelong.com/v2/miniProgram/scenicFacilities/findWaitTimeMap",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 3. 设置请求头
	req.Header.Set("Host", "zx-api.chimelong.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "_hid=KEa_CvTFq2inDlZyR1VVAAA; _hid1=KEa_CvTFq2inDlZyR1VVAAA; acw_tc=0ae5a87217581237880215237e4b8970a3f7d93afc40a0dca5430ebc641280")
	req.Header.Set("deviceType", "APP_IOS")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("channelCode", "ONLINE")
	req.Header.Set("User-Agent", "Travel/7.9.6 (iPhone; iOS 18.4; Scale/3.00)")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	// 4. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 5. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 6. 解析响应数据
	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if apiResponse.Code != 0 {
		return nil, fmt.Errorf("API返回错误: %s", apiResponse.Message)
	}

	// 7. 转换为目标结构
	result := make(map[string]ScenicSpotWaitTime)
	for _, spotData := range apiResponse.Data.ImageData {
		scenicSpotId, err := h.ConvertID(spotData.Code)
		if err != nil {
			marshal, _ := json.Marshal(spotData)
			log.Printf("景点等待时间item:%s,找不到对应的景点id,%v", marshal, err)
			continue
		}

		waitTime := ScenicSpotWaitTime{
			ScenicSpotId: scenicSpotId,
			StartTime:    spotData.StartTime,
			EndTime:      spotData.EndTime,
		}

		// 解析等待时间
		waitTime.WaitMinute = parseWaitTime(spotData.WaitingTime)

		// 判断开放状态
		waitTime.OpenState = determineOpenState(spotData.WaitingTime, waitTime.WaitMinute)

		result[scenicSpotId] = waitTime
	}

	return result, nil
}

// parseWaitTime 解析等待时间字符串为分钟数
func parseWaitTime(waitingTime string) int64 {
	// 处理特殊状态
	if strings.Contains(waitingTime, "暂停接待") ||
		strings.Contains(waitingTime, "关闭") ||
		strings.Contains(waitingTime, "维护") {
		return 0
	}

	if minutes, err := strconv.ParseInt(waitingTime, 10, 64); err == nil {
		return minutes
	}

	return 0
}

// determineOpenState 根据等待时间字符串判断开放状态
func determineOpenState(waitingTime string, waitMinutes int64) OpenState {
	switch {
	case strings.Contains(waitingTime, "暂停接待"):
		return OpenStateClosed
	case strings.Contains(waitingTime, "维护"):
		return OpenStateMaintain
	case strings.Contains(waitingTime, "关闭"):
		return OpenStateClosed
	case waitMinutes > 0:
		return OpenStateOpen
	case waitingTime == "" || waitingTime == "0":
		return OpenStateOpen // 可能表示无需等待
	default:
		return OpenStateUnknown
	}
}

func NewZhuHaiChimelongOceanKingdom() *ZhuHaiChimelongOceanKingdom {
	return &ZhuHaiChimelongOceanKingdom{
		BaseHandler: BaseHandler{
			IDMapping: map[string]string{"31011001": "100000007",
				"31011002": "100000005",
				"31011003": "100000013",
				"31011004": "100000012",
				"31011005": "100000010",
				"31011006": "100000009",
				"31011007": "100000015",
				"31011008": "100000017",
				"31011009": "100000034",
				"31011010": "100000029",
				"31011011": "100000026",
				"31011012": "100000032",
				"31011013": "100000035",
				"31011014": "100000027",
				"31011015": "100000020",
				"31011020": "100000002",
				"31011021": "100000006",
				"31011022": "100000011",
				"31011023": "100000001",
				"31011024": "100000016",
				"31011025": "100000018",
				"31011026": "100000019",
				"31011032": "100000024",
				"31011033": "100000021",
			},
		},
	}
}
