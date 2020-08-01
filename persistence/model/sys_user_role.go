package model

import (
	"github.com/wangghyz/polestar/persistence/model/base"
)

type SysUserRole struct {
	base.PolestarModel
	UserId int `gorm:"type:int(10)"`
	RoleId int `gorm:"type:int(10)"`
}