package service

import (
	"github.com/jinzhu/gorm"
	"polestar/common/db"
	"polestar/common/model"
)

type (
	SysUserService struct {
		db *gorm.DB
	}
)

var (
	_sysUserService = &SysUserService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysUserService() *SysUserService {
	return _sysUserService
}

func (s *SysUserService) CreateUser(user *model.SysUser) (*model.SysUser, error) {
	rst := s.db.Create(user)
	return user, rst.Error
}

func (s *SysUserService) GetUserByUserName(userName string) (*model.SysUser, error) {
	user := &model.SysUser{}
	rst := s.db.First(user, "user_name = ? and status = '00'", userName)
	return user, rst.Error
}
