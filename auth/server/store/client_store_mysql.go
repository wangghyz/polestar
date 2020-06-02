package store

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_const "polestar/auth/const"
	"polestar/common/model"
	"polestar/common/service"
	"strings"
	"time"
)

var (
	// _mysqlClientStore
	_mysqlClientStore = &mySQLClientStore{
		service: service.NewSysClientServiceInstance(),
	}
)

// mySQLClientStore
type mySQLClientStore struct {
	service *service.SysClientService
}

// GetClient 获得客户端信息
// clientId 客户端ID
func (s *mySQLClientStore) GetClient(clientId string) (*ClientInfo, error) {
	if len(clientId) <= 0 {
		return nil, errors.New("客户端ID不能为空！")
	}

	client, err := s.service.GetClient(clientId)
	if err != nil || client == nil {
		return nil, err
	}

	info := &ClientInfo{}
	info.ClientId = client.ClientId
	info.ClientSecret = client.ClientSecret
	info.Scope = strings.Split(client.Scope, ",")
	info.GrantType = make([]_const.GrantType, 0)
	gts := strings.Split(client.GrantTypes, ",")
	for _, gt := range gts {
		info.GrantType = append(info.GrantType, _const.GrantType(gt))
	}
	info.AccessTokenDuration = time.Minute * time.Duration(client.AccessTokenDuration)
	info.RefreshTokenDuration = time.Minute * time.Duration(client.RefreshTokenDuration)

	return info, nil
}

// AddClient 追加客户端信息
func (s *mySQLClientStore) AddClient(clientInfo *ClientInfo) error {
	if len(clientInfo.ClientId) > 0 {
		_, err := s.service.GetClient(clientInfo.ClientId)
		if !gorm.IsRecordNotFoundError(err) {
			if err == nil {
				return errors.New(fmt.Sprintf("客户端[%s]已存在！", clientInfo.ClientId))
			} else {
				return err
			}
		}

		client := &model.SysClient{}
		client.ClientId = clientInfo.ClientId
		client.ClientSecret = clientInfo.ClientSecret

		var gts []string
		for _, gt := range clientInfo.GrantType {
			gts = append(gts, string(gt))
		}

		client.GrantTypes = strings.Join(gts, ",")
		client.AccessTokenDuration = int(clientInfo.AccessTokenDuration.Minutes())
		client.RefreshTokenDuration = int(clientInfo.RefreshTokenDuration.Minutes())
		client.Scope = strings.Join(clientInfo.Scope, ",")

		_, err = s.service.CreateClient(client)
		return err
	} else {
		return errors.New("ClientId不能为空！")
	}
}

// 移除客户端信息
func (s *mySQLClientStore) RemoveClient(clientId string) error {
	return s.service.DeleteClient(clientId)
}

// 获取 Client Store
func NewMySQLClientStoreInstance() ClientStore {
	return _mysqlClientStore
}
