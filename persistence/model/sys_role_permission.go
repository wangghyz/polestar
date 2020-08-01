package model

import (
	"github.com/wangghyz/polestar/persistence/model/base"
)

type SysRolePermission struct {
	base.PolestarModel
	RoleId       int `gorm:"type:int(10)"`
	PermissionId int `gorm:"type:int(10)"`
}
