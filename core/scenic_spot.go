package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetTicketScenicSpotInfo(cPoiId string) (*SightInfo, error) {
	client := &http.Client{}
	baseUrl := "https://hy.travel.qunar.com/api/poi/getPoiDetail"
	url := fmt.Sprintf("%s/?RN=1&poiId=%s",
		baseUrl, cPoiId)
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
	var response ConvertResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if !response.Ret {
		return nil, err
	}

	// 提取tId判断是否成功
	tId := response.Data.SightID
	if 0 == tId {
		log.Printf("内容景点转门票景点失败，找不到对应的门票景点...")
		return nil, err
	}

	// 拿着tId请求详情
	client1 := &http.Client{}
	baseUrl1 := "https://piao.qunar.com/ticket/pw/wisdom/%7Bpoi%7D/sight/detail.json"
	url1 := fmt.Sprintf("%s?sightId=%d",
		baseUrl1, tId)
	req1, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		return nil, err
	}
	resp1, err := client1.Do(req1)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}
	var sightResponse SightResponse
	if err := json.Unmarshal(body1, &sightResponse); err != nil {
		return nil, err
	}
	if !sightResponse.Ret {
		return nil, err
	}

	return &sightResponse.Data.SightInfo, nil
}

type ConvertResponse struct {
	Ret     bool   `json:"ret"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Ver     int    `json:"ver"`
	Data    struct {
		Name           string `json:"name"`
		NewComponentAb string `json:"newComponentAb"`
		ID             int    `json:"id"`
		SightID        int    `json:"sightId"`
	} `json:"data"`
}

type SightResponse struct {
	Ret  bool `json:"ret"`
	Data struct {
		SightInfo `json:"sightInfo"`
	} `json:"data"`
}

type SightInfo struct {
	ImageCount    int     `json:"imageCount"`
	TravelCount   int     `json:"travelCount"`
	CommentUrl    string  `json:"commentUrl"`
	CommentScore  float64 `json:"commentScore"`
	SightCategory string  `json:"sightCategory"`
	CommentDesc   string  `json:"commentDesc"`
	Foreign       bool    `json:"foreign"`
	CoverImage    string  `json:"coverImage"`
	Traffic       string  `json:"traffic"`
	BaiduMapPoint string  `json:"baiduMapPoint"`
	Address       string  `json:"address"`
	SightIntro    string  `json:"sightIntroduction"`
	SightLevel    string  `json:"sightLevel"`
	SightBusiness struct {
		Title       string   `json:"title"`
		TitleColor  string   `json:"titleColor"`
		OtherTitles []string `json:"otherTitles"`
	} `json:"sightBusinessStatusData"`
	SightId         int      `json:"sightId"`
	GoogleMapPoint  string   `json:"googleMapPoint"`
	CommentCount    int      `json:"commentCount"`
	SightOpenTime   string   `json:"sightOpenTime"`
	SightName       string   `json:"sightName"`
	BigImages       []string `json:"bigImages"`
	IntroductionUrl string   `json:"introductionUrl"`
	Region          string   `json:"region"`
}
