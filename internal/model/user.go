package model

type UserInfo struct {
	UserName  string `gorm:"column:user_name" json:"userName"`
	NickName  string `gorm:"column:nick_name" json:"nickName"`
	AvatarUrl string `gorm:"column:avatar_url" json:"avatarUrl"`
	OpenID    string `gorm:"column:open_id" json:"openId"`
}

// TableName 指定表名
func (UserInfo) TableName() string {
	return "tm_user_info"
}
