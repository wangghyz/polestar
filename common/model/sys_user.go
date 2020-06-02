package model

import (
	base2 "polestar/common/model/base"
)

type SysUser struct {
	base2.PolestarModel
	UserName    string `gorm:"type:varchar(100);index;not null"`
	Name        string `gorm:"type:varchar(45)"`
	Password    string `gorm:"type:varchar(200)"`
	HeaderImage string `gorm:"type:varchar(200)"`
	Comment     string `gorm:"type:varchar(400)"`
	Status      string `gorm:"type:varchar(2);default:\"00\""`
}
