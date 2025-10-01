package wait_time

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ShanghaiDisney struct {
	BaseHandler
}

// ShanghaiDisneyAPIResponse 定义上海迪士尼API响应结构
type ShanghaiDisneyAPIResponse struct {
	Entries []struct {
		ID       string `json:"id"`
		WaitTime struct {
			FastPass struct {
				Available bool   `json:"available"`
				StartTime string `json:"startTime,omitempty"`
			} `json:"fastPass"`
			Status            string `json:"status"`
			SingleRider       bool   `json:"singleRider"`
			PostedWaitMinutes int    `json:"postedWaitMinutes,omitempty"`
		} `json:"waitTime"`
	} `json:"entries"`
}

func (h *ShanghaiDisney) GetWaitTime() (map[string]ScenicSpotWaitTime, error) {
	// 1. 创建 HTTP 请求
	req, err := http.NewRequest(
		"GET",
		"https://app.apigw.shanghaidisneyresort.com/explorer-service/public/wait-times/shdr;entityType=destination;destination=shdr?region=CN",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 2. 设置请求头
	req.Header.Set("Host", "app.apigw.shanghaidisneyresort.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

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
	var apiResponse ShanghaiDisneyAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 6. 转换为目标结构
	result := make(map[string]ScenicSpotWaitTime)
	for _, entry := range apiResponse.Entries {
		// 提取ID（去除分号后的部分）
		entryID := strings.Split(entry.ID, ";")[0]

		// 尝试转换ID
		scenicSpotId, err := h.ConvertID(entryID)
		if err != nil {
			// ID未映射时记录日志并跳过
			//marshal, _ := json.Marshal(entry)
			//log.Printf("上海迪士尼景点等待时间item:%s,找不到对应的景点id,%v", marshal, err)
			continue
		}

		waitTime := ScenicSpotWaitTime{
			ScenicSpotId: scenicSpotId,
		}

		// 解析等待时间和开放状态
		waitTime.WaitMinute, waitTime.OpenState = h.parseWaitTimeAndStatus(
			entry.WaitTime.PostedWaitMinutes,
			entry.WaitTime.Status,
		)

		result[scenicSpotId] = waitTime
	}

	return result, nil
}

// parseWaitTimeAndStatus 解析等待时间和开放状态
func (h *ShanghaiDisney) parseWaitTimeAndStatus(postedWaitMinutes int, status string) (int64, OpenState) {
	// 根据status判断状态
	switch strings.ToLower(status) {
	case "closed":
		return 0, OpenStateClosed
	case "renewal", "refurbishment", "maintenance":
		return 0, OpenStateMaintain
	case "operating":
		// 如果是开放状态，返回等待时间
		return int64(postedWaitMinutes), OpenStateOpen
	default:
		return 0, OpenStateUnknown
	}
}

func NewShanghaiDisney() *ShanghaiDisney {
	return &ShanghaiDisney{
		BaseHandler: BaseHandler{
			// ID映射 - 格式: "迪士尼景点ID": "系统内部景点ID"
			IDMapping: map[string]string{
				// 游乐设施 (Attractions)
				"attBuzzLightyearPlanetRescue": "100000053", // 巴斯光年星际营救
				"attDumboFlyingElephant":       "100000054", // 小飞象(奇想花园店)
				"attPeterPansFlight":           "100000055", // 小飞侠天空奇遇
				"attPiratesOfCaribbean":        "100000056", // 加勒比海盗沉落宝藏之战
				"attRexsRCRacer":               "100000057", // 抱抱龙冲天赛车
				"attRoaringRapids":             "100000058", // 雷鸣山漂流
				"attSevenDwarfsMineTrain":      "100000059", // 七个小矮人矿山车
				"attSoaringOverHorizon":        "100000060", // 翱翔飞越地平线
				"attAdventuresWinniePooh":      "100000061", // 小熊维尼历险记
				"attTronLightcyclePowerRun":    "100000062", // 创极速光轮
				"attVoyageToCrystalGrotto":     "100000063", // 晶彩奇航
				"attWoodysRoundUp":             "100000064", // 胡迪牛仔嘉年华
				"attChallengeTrails":           "100000065", // 古迹探索营
				"attZootopiaHotPursuit":        "100000066", // 热力追踪
				"attAliceWonderlandMaze":       "100000067", // 爱丽丝梦游仙境迷宫
				"attMarvelUniverse":            "100000068", // 漫威英雄总部
				"attOnceUponTimeAdventure":     "100000074", // 漫游童话时光
				"attExplorerCanoes":            "100000070", // 探险家独木舟
				"attFantasiaCarousel":          "100000071", // 幻想曲旋转木马
				"attHunnyPotSpin":              "100000072", // 旋转疯蜜罐
				"attJetPacks":                  "100000073", // 喷气背包飞行器
				"attSlinkyDogSpin":             "100000075", // 弹簧狗团团转

				// 角色见面 (Character Meet & Greet)
				"charCharactersMickeyMouseGardensImagination":                           "100000076", // 与米奇见面
				"charSelfieSpotLinaBell":                                                "100000077", // 与玲娜贝儿见面
				"charCharactersDisneyPrincessesStorybookCourt":                          "100000078", // 与迪士尼王室见面
				"charCharactersHeroicEncounterAtTheMarvelUniverse":                      "100000079", // 与复仇者联盟见面
				"charCharactersMeetChipAndDale":                                         "100000080", // 与奇奇蒂蒂见面
				"charCharactersSelfieSpotDuffy":                                         "100000081", // 与达菲和朋友们见面
				"attJackSparrowTreasureCove":                                            "100000082", // 与杰克船长见面
				"charCharactersMeetJudyNickAtZootopiaPoliceDepartmentRecruitmentCenter": "100000083", // 与朱迪和尼克见面
				"charCharactersMinnieMouseFriends":                                      "100000084", // 与米妮和朋友们见面
				"charCharactersMeetPixarCharactersAtTomorrowland":                       "100000085", // 和皮克斯朋友们见面
				"charCharactersMarvelUniverseSpiderMan":                                 "100000086", // 与蜘蛛侠见面
				"charMeetTheCharacterFromUPAtHappyCircle":                               "100000087", // 与《飞屋环游记》中的朋友见面
				"charCharactersMeetTheFriendsFromLiloStitchAtGalaxyGathering":           "100000088", // 与《星际宝贝》中的朋友见面
				"charCharactersMeetPooh":                                                "100000089", // 与小熊维尼见面
				"charCharactersMeetingPost":                                             "100000090", // 与您喜爱的玩具朋友见面
			},
		},
	}
}
