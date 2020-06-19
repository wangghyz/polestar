package service

import (
	"github.com/jinzhu/gorm"
	"github.com/wangghyz/polestar/common/db"
	"github.com/wangghyz/polestar/common/model"
)

type (
	SysUserRoleService struct {
		db *gorm.DB
	}
)

var (
	_sysUserRoleService = &SysUserRoleService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysUserRoleService() *SysUserRoleService {
	return _sysUserRoleService
}

func (s *SysUserRoleService) CreateUserRole(userRole *model.SysUserRole) (*model.SysUserRole, error) {
	rst := s.db.Create(userRole)
	return userRole, rst.Error
}

func (s *SysUserRoleService) DeleteUserRole(userId int, roleId int) error {
	rst := s.db.Delete(&model.SysUserRole{}, "user_id = ? and role_id = ?", userId, roleId)
	return rst.Error
}
