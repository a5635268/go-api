package model

type SysUser struct {
	Uid int `json:"uid" gorm:"primary_key;AUTO_INCREMENT"`
	Nickname string `json:"nickname"`
	Openid string `json:"openid"`
	Status int `json:"status"`
}
