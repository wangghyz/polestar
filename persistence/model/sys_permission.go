package model

import (
	"github.com/wangghyz/polestar/persistence/model/base"
)

type SysPermission struct {
	base.PolestarModel
	EnName  string `gorm:"type:varchar(45);index;not null"`
	Name    string `gorm:"type:varchar(45)"`
	Comment string `gorm:"type:varchar(400)"`
	Status  string `gorm:"type:varchar(2);default:\"00\""`
}
