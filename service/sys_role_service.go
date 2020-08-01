package service

import (
	"github.com/jinzhu/gorm"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/persistence/model"
)

type (
	SysRoleService struct {
		db *gorm.DB
	}
)

var (
	_sysRoleService = &SysRoleService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysRoleService() *SysRoleService {
	return _sysRoleService
}

func (s *SysRoleService) CreateRole(role *model.SysRole) (*model.SysRole, error) {
	rst := s.db.Create(role)
	return role, rst.Error
}

func (s *SysRoleService) GetRoleByEnName(roleEnName string) (*model.SysRole, error) {
	role := &model.SysRole{}
	rst := s.db.First(role, "en_name = ? and status = '00'", roleEnName)
	return role, rst.Error
}

func (s *SysRoleService) GetRolesByUserName(userName string) ([]model.SysRole, error) {

	sql :=
		`SELECT
			r.*
		FROM
			sys_role r
			JOIN
			sys_user_role ur ON r.id = ur.role_id
			JOIN
			sys_user u ON ur.user_id = u.id
		WHERE
			u.user_name = ?
			AND u.status = '00'
			AND r.status = '00'
			AND u.deleted_at IS NULL
			AND r.deleted_at IS NULL
			AND ur.deleted_at IS NULL`

	rst := s.db.Raw(sql, userName)
	if rst.Error != nil {
		return nil, rst.Error
	}

	roles := &[]model.SysRole{}
	rst.Scan(roles)

	return *roles, nil
}
