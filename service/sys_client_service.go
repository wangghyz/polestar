package service

import (
	"github.com/jinzhu/gorm"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/persistence/model"
)

type (
	// SysClientService
	SysClientService struct {
		db *gorm.DB
	}
)

var (
	_sysClientService = &SysClientService{
		db: db.NewMySQLConnectionInstance(),
	}
)

func NewSysClientServiceInstance() *SysClientService {
	return _sysClientService
}

func (s *SysClientService) CreateClient(client *model.SysClient) (*model.SysClient, error) {
	rst := s.db.Create(client)
	return client, rst.Error
}

func (s *SysClientService) DeleteClient(clientId string) error {
	return s.db.Delete(&model.SysClient{}, "client_id = ?", clientId).Error
}

func (s *SysClientService) GetClient(clientId string) (*model.SysClient, error) {
	client := &model.SysClient{}
	rst := s.db.First(client, "status = '00' and client_id = ?", clientId)
	return client, rst.Error
}

func (s *SysClientService) GetSysClients() ([]model.SysClient, error) {
	var clients []model.SysClient
	rst := s.db.Find(&clients)
	return clients, rst.Error
}
