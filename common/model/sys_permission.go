package model

import (
	base2 "polestar/common/model/base"
)

type SysPermission struct {
	base2.PolestarModel
	EnName  string `gorm:"type:varchar(45);index;not null"`
	Name    string `gorm:"type:varchar(45)"`
	Comment string `gorm:"type:varchar(400)"`
	Status  string `gorm:"type:varchar(2);default:\"00\""`
}
