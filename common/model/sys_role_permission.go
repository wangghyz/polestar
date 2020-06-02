package model

import (
	base2 "polestar/common/model/base"
)

type SysRolePermission struct {
	base2.PolestarModel
	RoleId       int `gorm:"type:int(10)"`
	PermissionId int `gorm:"type:int(10)"`
}