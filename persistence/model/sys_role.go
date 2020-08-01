package model

import (
	"github.com/wangghyz/polestar/persistence/model/base"
)

type SysRole struct {
	base.PolestarModel
	EnName      string `gorm:"type:varchar(20);index;not null"`
	Name        string `gorm:"type:varchar(20)"`
	Comment     string `gorm:"type:varchar(400)"`
	Status      string `gorm:"type:varchar(2);default:\"00\""`
}
