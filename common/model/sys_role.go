package model

import (
	base2 "github.com/wangghyz/polestar/common/model/base"
)

type SysRole struct {
	base2.PolestarModel
	EnName      string `gorm:"type:varchar(20);index;not null"`
	Name        string `gorm:"type:varchar(20)"`
	Comment     string `gorm:"type:varchar(400)"`
	Status      string `gorm:"type:varchar(2);default:\"00\""`
}
