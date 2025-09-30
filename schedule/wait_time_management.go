package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"trip-map/core/wait_time"
)

// GetWaitTimeCacheStatus 获取等待时间缓存状态
func GetWaitTimeCacheStatus() map[string]interface{} {
	cacheManager := wait_time.GetCacheManager()
	scheduler := GetGlobalScheduler()

	result := make(map[string]interface{})

	// 调度器状态
	result["scheduler_running"] = false
	if scheduler != nil {
		result["scheduler_running"] = scheduler.GetStatus()
	}

	// 缓存状态
	cachedAreas := cacheManager.GetAllCachedScenicAreas()
	cacheDetails := make(map[string]interface{})

	for _, areaId := range cachedAreas {
		if cacheInfo, exists := cacheManager.GetCacheInfo(areaId); exists {
			cacheDetails[areaId] = map[string]interface{}{
				"updated_at":   cacheInfo.UpdatedAt.Format("2006-01-02 15:04:05"),
				"scenic_spots": len(cacheInfo.Data),
			}
		}
	}

	result["cached_areas"] = cacheDetails
	result["total_cached_areas"] = len(cachedAreas)

	return result
}

// ForceRefreshWaitTimeCache 强制刷新等待时间缓存
func ForceRefreshWaitTimeCache() error {
	scheduler := GetGlobalScheduler()
	if scheduler == nil {
		return fmt.Errorf("调度器未初始化")
	}

	if !scheduler.GetStatus() {
		return fmt.Errorf("调度器未运行")
	}

	log.Printf("手动触发等待时间缓存刷新")
	scheduler.ForceExecute()

	return nil
}

// ClearWaitTimeCache 清空等待时间缓存
func ClearWaitTimeCache(scenicAreaId string) error {
	cacheManager := wait_time.GetCacheManager()

	if scenicAreaId == "" {
		cacheManager.ClearAllCache()
		log.Printf("已清空所有等待时间缓存")
	} else {
		cacheManager.ClearCache(scenicAreaId)
		log.Printf("已清空景区 %s 的等待时间缓存", scenicAreaId)
	}

	return nil
}

// GetWaitTimeDirectly 直接获取等待时间数据 (不使用缓存，用于测试)
func GetWaitTimeDirectly(scId string) (map[string]wait_time.ScenicSpotWaitTime, error) {
	handler, ok := wait_time.ScenicAreaWaitTimeHandlerMap[scId]
	if !ok {
		return nil, fmt.Errorf("景区 %s 不支持等待时间", scId)
	}

	log.Printf("直接获取景区 %s 的等待时间数据", scId)
	return handler.GetWaitTime()
}

// GetScenicAreaRawData 获取景区原始数据 (用于查看第三方API原始响应，方便ID映射)
func GetScenicAreaRawData(scId string) (interface{}, error) {
	switch scId {
	case "BEIJING_UNIVERSAL":
		return getBijingUniversalRawData()
	case "7469417":
		return getZhuHaiChimelongRawData()
	default:
		return nil, fmt.Errorf("景区 %s 不支持获取原始数据", scId)
	}
}

// getBijingUniversalRawData 获取北京环球影城原始数据
func getBijingUniversalRawData() (interface{}, error) {
	return fetchBeijingUniversalRawData()
}

// getZhuHaiChimelongRawData 获取珠海长隆原始数据
func getZhuHaiChimelongRawData() (interface{}, error) {
	return fetchZhuHaiChimelongRawData()
}

// fetchBeijingUniversalRawData 直接获取北京环球影城原始数据
func fetchBeijingUniversalRawData() (interface{}, error) {
	// 创建HTTP请求
	req, err := http.NewRequest(
		"GET",
		"https://g.app.universalbeijingresort.com/map/attraction/list?mode=list",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
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

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析为通用interface{}以返回原始数据
	var rawData interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return rawData, nil
}

// fetchZhuHaiChimelongRawData 直接获取珠海长隆原始数据
func fetchZhuHaiChimelongRawData() (interface{}, error) {
	// 准备请求数据
	requestBody := map[string]string{"code": "ZH56"}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest(
		"POST",
		"https://zx-api.chimelong.com/v2/miniProgram/scenicFacilities/findWaitTimeMap",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
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

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析为通用interface{}以返回原始数据
	var rawData interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return rawData, nil
}
