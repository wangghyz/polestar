package model

import (
	"github.com/wangghyz/polestar/persistence/model/base"
)

type SysUser struct {
	base.PolestarModel
	UserName    string `gorm:"type:varchar(100);unique_index;not null"`
	Name        string `gorm:"type:varchar(45)"`
	Password    string `gorm:"type:varchar(200)"`
	HeaderImage string `gorm:"type:varchar(200)"`
	Comment     string `gorm:"type:varchar(400)"`
	Status      string `gorm:"type:varchar(2);default:\"00\""`
}
