package core

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"trip-map/internal/model"
)

type Response struct {
	Ret     bool   `json:"ret"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Ver     int    `json:"ver"`
	Data    Data   `json:"data"`
}

type Data struct {
	More  bool   `json:"more"`
	Count int    `json:"count"`
	List  []Item `json:"list"`
}

type Item struct {
	ID              int         `json:"id"`
	CreateTime      int64       `json:"createTime"`
	ReplyCount      int         `json:"replyCount"`
	SubType         int         `json:"subType"`
	BookID          int         `json:"bookId"`
	BookTitle       string      `json:"bookTitle"`
	PoiID           int         `json:"poiId"`
	PoiName         string      `json:"poiName"`
	VideoUrl        *string     `json:"videoUrl"`
	VideoCoverImage *string     `json:"videoCoverImage"`
	VideoDuration   int         `json:"videoDuration"`
	VideoWidth      int         `json:"videoWidth"`
	VideoHeight     int         `json:"videoHeight"`
	Body            string      `json:"body"`
	IPLocation      *string     `json:"ipLocation"`
	UserRating      int         `json:"userRating"`
	Images          []string    `json:"images"`
	DestImages      []DestImage `json:"destImages"`
	Author          Author      `json:"author"`
	Title           string      `json:"title"`
	Quality         int         `json:"quality"`
	UsefulCnt       int         `json:"usefulCnt"`
	BookType        int         `json:"bookType"`
	ElementPhoto    *string     `json:"elementPhoto"`
	HaveLike        bool        `json:"haveLike"`
	Replies         *string     `json:"replies"`
	DetailUrl       string      `json:"detailUrl"`
}

type DestImage struct {
	ID      int64   `json:"id"`
	Intro   *string `json:"intro"`
	URL     string  `json:"url"`
	UserID  *string `json:"userId"`
	ImgType *string `json:"imgType"`
	BookID  int64   `json:"bookId"`
	Source  string  `json:"source"`
	Quality int     `json:"quality"`
}

type Author struct {
	UserID   string `json:"userId"`
	HeadImg  string `json:"headImg"`
	NickName string `json:"nickName"`
}

func GetNoteListByPoiId(poiId string, poiName string) ([]model.Note, error) {
	client := &http.Client{}
	baseUrl := "https://hy.travel.qunar.com/api/comment/poi/comment"

	// 初始化结果列表
	var allNotes []model.Note
	page := 1
	maxPages := 5
	rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器

	for page <= maxPages {
		// 构建URL，每次获取16条数据
		url := fmt.Sprintf("%s?poiId=%s&page=%d&limit=16&ctrip=true&needLvtu=true&useEs=true&from_page=page_search_result",
			baseUrl, poiId, page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return allNotes, err // 返回已获取的部分数据和错误
		}

		resp, err := client.Do(req)
		if err != nil {
			return allNotes, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // 立即关闭body，而不是 defer
		if err != nil {
			return allNotes, err
		}

		var response Response
		if err := json.Unmarshal(body, &response); err != nil {
			return allNotes, err
		}

		if !response.Ret {
			return allNotes, fmt.Errorf("response returned Ret=false")
		}

		// 处理当前页数据
		list := response.Data.List
		for _, item := range list {
			note := model.Note{
				QNoteId:         strconv.Itoa(item.ID),
				AuthorNickName:  item.Author.NickName,
				AuthorHeadImg:   item.Author.HeadImg,
				Title:           item.Title,
				Body:            item.Body,
				Images:          strings.Join(item.Images, ","),
				DetailUrl:       item.DetailUrl,
				CPoiId:          poiId,
				CPoiName:        poiName,
				UsefulCnt:       item.UsefulCnt,
				VideoUrl:        StringOrEmpty(item.VideoUrl),
				VideoCoverImage: StringOrEmpty(item.VideoCoverImage),
				VideoDuration:   item.VideoDuration,
			}
			allNotes = append(allNotes, note)
		}

		// 检查是否还有更多数据，如果没有则退出循环
		if !response.Data.More {
			break
		}

		// 准备下一页
		page++

		// 如果不是最后一页，等待随机时间(100~200ms)
		//if page <= maxPages {
		//	delay := time.Duration(100) * time.Millisecond
		//	time.Sleep(delay)
		//}
	}

	return allNotes, nil
}

func StringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
