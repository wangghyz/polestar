package service

import (
	"github.com/jinzhu/gorm"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/persistence/model"
)

type (
	SysPermissionService struct {
		db *gorm.DB
	}
)

var (
	_sysPermissionService = &SysPermissionService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysPermissionService() *SysPermissionService {
	return _sysPermissionService
}

func (s *SysPermissionService) GetPermissionByEnName(enName string) (*model.SysPermission, error) {
	permission := &model.SysPermission{}
	rst := s.db.First(permission, "en_name = ? and status = '00'", enName)
	return permission, rst.Error
}

func (s *SysPermissionService) CreatePermission(permission *model.SysPermission) (*model.SysPermission, error) {
	rst := s.db.Create(permission)
	return permission, rst.Error
}

func (s *SysPermissionService) GetPermissionsByUserName(userName string) ([]model.SysPermission, error) {
	sql :=
		`SELECT DISTINCT
			p.*
		FROM
			sys_role r
				JOIN
			sys_user_role ur ON r.id = ur.role_id
				JOIN
			sys_user u ON ur.user_id = u.id
				JOIN
			sys_role_permission rp ON rp.role_id = r.id
				JOIN
			sys_permission p ON rp.permission_id = p.id
		WHERE
			u.user_name = ?
				AND u.status = '00'
				AND r.status = '00'
				AND u.deleted_at IS NULL
				AND r.deleted_at IS NULL
				AND ur.deleted_at IS NULL
				AND rp.deleted_at IS NULL
				AND p.deleted_at IS NULL`

	rst := s.db.Raw(sql, userName)
	if rst.Error != nil {
		return nil, rst.Error
	}

	permission := &[]model.SysPermission{}
	rst.Scan(permission)

	return *permission, nil
}
