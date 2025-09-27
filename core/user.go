package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	AppId           = "wxa37d8e399e7affa4"
	AppSecret       = "f3f1095ea4610f81e1b23be1fd4360d0"
	Code2SessionUrl = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func Code2Session(code string) string {
	url := fmt.Sprintf(Code2SessionUrl, AppId, AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("WxInvoker.code2Session error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get successful response: %d\n", resp.StatusCode)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("WxInvoker.code2Session read body error: %v\n", err)
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("WxInvoker.code2Session json unmarshal error: %v\n", err)
		return ""
	}

	if errCode, ok := result["errcode"]; ok {
		fmt.Printf("WxInvoker.code2Session error, code: %s,errorCode : %s response: %s\n", code, errCode, string(body))
		return ""
	}

	if openid, ok := result["openid"].(string); ok {
		return openid
	}

	return ""
}
