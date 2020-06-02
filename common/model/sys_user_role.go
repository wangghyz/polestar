package model

import (
	base2 "polestar/common/model/base"
)

type SysUserRole struct {
	base2.PolestarModel
	UserId int `gorm:"type:int(10)"`
	RoleId int `gorm:"type:int(10)"`
}