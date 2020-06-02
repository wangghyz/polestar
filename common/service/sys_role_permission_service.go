package service

import (
	"github.com/jinzhu/gorm"
	"polestar/common/db"
	"polestar/common/model"
)

type (
	SysRolePermissionService struct {
		db *gorm.DB
	}
)

var (
	_sysRolePermissionService = &SysRolePermissionService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysRolePermissionService() *SysRolePermissionService {
	return _sysRolePermissionService
}

func (s *SysRolePermissionService) CreateRolePermission(rolePermission *model.SysRolePermission) (*model.SysRolePermission, error) {
	rst := s.db.Create(rolePermission)
	return rolePermission, rst.Error
}

func (s *SysRolePermissionService) DeleteRolePermission(roleId, permissionId int) error {
	rst := s.db.Delete(&model.SysRolePermission{}, "role_id = ? and permission_id = ?", roleId, permissionId)
	return rst.Error
}
