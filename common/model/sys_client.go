package model

import (
	base2 "polestar/common/model/base"
)

type SysClient struct {
	base2.PolestarModel
	ClientId             string `gorm:"type:varchar(20);index;not null"`
	ClientSecret         string `gorm:"type:varchar(100)"`
	GrantTypes           string `gorm:"type:varchar(200)"`
	Scope                string `gorm:"type:varchar(200)"`
	AccessTokenDuration  int    `gorm:"type:int;default:20"`
	RefreshTokenDuration int    `gorm:"type:int;default:10080"`
	Comment              string `gorm:"type:varchar(200)"`
	Status               string `gorm:"type:varchar(2);default:\"00\""`
}
