package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"trip-map/core"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/utils"
)

const (
	DefaultAvatar   string = "https://edu-1010xzh.oss-cn-shenzhen.aliyuncs.com/2024/02/28/4b77f690cc1749f5ac79c465709bcff4baku.jpg"
	DefaultNickname string = "微信用户"
	RsaPublishStr   string = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA29jZfWjA3O6U+gK+KY7KUWBhixXkufbG/66y51iN4q75arKMemEnb9/QhoXLdGFNpklkp+r8KaGacIf3AFlLDKB2wOu/oRae9YUurLpN5nSsexbtOOxx5aNTs3cJFGeP7oh6ABQ+cZp4UdVIWoUbFXkdwY2N6StT4kZvedyAw+RpSSa5tWvOtonAJraqUzJIYyqlEUDr1sGnuRU6aVfe6yfA9H6VZk7kr+HVSMZo1MFvteiTbfv7jkZzTqanN7+E8gpQqisFW7SHYi3Bwp/XzyqdPnVrNmngxQZxOx45dZ4iTwcTpIOfpcymd2ly2VT90zdg6QxMK+InW5X9f1tQkwIDAQAB"
	RsaPrivateStr   string = "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDb2Nl9aMDc7pT6Ar4pjspRYGGLFeS59sb/rrLnWI3irvlqsox6YSdv39CGhct0YU2mSWSn6vwpoZpwh/cAWUsMoHbA67+hFp71hS6suk3mdKx7Fu047HHlo1OzdwkUZ4/uiHoAFD5xmnhR1UhahRsVeR3BjY3pK1PiRm953IDD5GlJJrm1a862icAmtqpTMkhjKqURQOvWwae5FTppV97rJ8D0fpVmTuSv4dVIxmjUwW+16JNt+/uORnNOpqc3v4TyClCqKwVbtIdiLcHCn9fPKp0+dWs2aeDFBnE7Hjl1niJPBxOkg5+lzKZ3aXLZVP3TN2DpDEwr4idblf1/W1CTAgMBAAECggEAZCPWjYVVtE0IlwkAzbU4+vBH/i6uzPZXlsdgvnhbyNGi0rMZwfTXHeJ4/Y2cKxrXX9M2gjZLPjtaOb/1Brels86zyRSZaSsApR1RMWR7b2nd1wOOcstg5hULX0ftXtn9ec24pKiT+PM/sybPmkvfFlzg7PUpmvgdcYhb5spF7PQZRRGnhIij+AN01NfDAZ1oqtFQ7xkQUiSQFSNPVUfFC9nCSoAebWTqkFsAWzvHd2NZgFiYYrjiNmLa40OGxKjVRsrfqN1AFB6kq722B/FGgA+yk9zig9Mot5j7agGvzOSNeb11WY/f7Xd5B2muQidGHfwHDfJUquSoedKFv7KwwQKBgQD+vf8FIZtORgZTbjo9c9+rcw1tqE7zo4v8yBruzMuJMYmCeodQPGc2tkr4i0mPvH0q7HDUcbtx34oahT+Lwww/qvsEWbwycA9wMSlPTZEW0Jbxquo4VjzmPakmZjFQa82U5yMO8eS6yEOxNpPUiIH08/A9U2+RvFw+Olz5o+We+wKBgQDc7r6hsTXciUjltAWIrg+ILtYkx7S+lBtWrf1yeKQe8ne8+2cxA958oxCYj9RvAQ3SqjwE3ZMhX+8AMwiwbibjUAu1qNtq8AvsQi66GeMrLwcKmSPBkXilOp1yUpm+0vY7BpVe0adsLW5IlW+FCUGOIHf6GH+4XfLlMrD5jCYBSQKBgEeoJEtKN8id0/u1/vX4WUt+EqHs/UB1mdQiackQnJRb9eVZGCUOyK3QO2iMrcWb7M2dMuPfli2jBtMM9mIXHKPwMan4oALEGOOjQI6JMC3twPf77uSoBXtyjtk5V9faazrehbMXghK0cK4xvwXC3GOOFt75UGH7TStH+Y1TeCzvAoGAUlbkU0zBXy0HLxzVxyfgAAg8pT6MvU5jlf2IbOZLfIEvYQ5tWhYwEFGRuNo5+RjyduYdMk8GK7UeVPuwLFkRQzys8Io7JHLMbsQHuDI3uPtw62FBsz2tMh9TWK0yQa1MOZltiAYpGKch6AlRo8pcVUUCkgIZb7QL96HZ1VeHPokCgYAGhbeGBB0P6THSeD8vymO48uixxCHMwNcLyUpEd+WpbASmUf3uko2C3S77x9ovCYY6jZqG5AJVZTghmWrjeA8o5fFMgMv8o++YLUdra8rPoBU33RUoeWTP0Wuc6TSlFf2my/rN/FfGqBUrvPr/fHFB1kmRrXMLr23mZJ4k00QVBA=="
)

func DoLogin(code string) (string, error) {
	openID := core.Code2Session(code)
	if "" == openID {
		return "", fmt.Errorf("获取openId异常")
	}
	// 通过 openID 查找用户信息
	userInfo, err := getUserInfoByOpenID(openID)
	if nil != err {
		// 查找失败
		log.Printf("getUserInfoByOpenID error: %v", err)
		return "", err
	}

	// 执行注册
	if userInfo == nil {
		registeredUserInfo, err := createUserByOpenID(openID)
		if nil != err {
			log.Printf("createUserByOpenID error: %v", err)
			return "", err
		}
		log.Printf("用户注册成功！%v", registeredUserInfo)
		userInfo = registeredUserInfo
	}

	token, err := utils.CreateToken(userInfo.UserName, RsaPrivateStr)
	if err != nil {
		log.Printf("生成token失败 %v", err)
		return "", err
	}

	return token, nil
}

func getUserInfoByOpenID(openID string) (*model.UserInfo, error) {
	userInfo := model.UserInfo{}
	err := global.Db.Where("open_id = ?", openID).First(&userInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userInfo, nil
}

func createUserByOpenID(openID string) (*model.UserInfo, error) {
	userInfo := &model.UserInfo{
		UserName:  uuid.New().String(),
		NickName:  DefaultNickname,
		AvatarUrl: DefaultAvatar,
		OpenID:    openID,
	}
	err := global.Db.Create(userInfo).Error
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
