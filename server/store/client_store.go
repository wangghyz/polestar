package store

import (
	"github.com/wangghyz/polestar/common"
	"time"
)

type (
	// ClientStore 接口
	ClientStore interface {
		// GetClient 获取客户端
		GetClient(clientId string) (*ClientInfo, error)
		// AddClient 追加客户端信息
		AddClient(clientInfo *ClientInfo) error
		// RemoveClient 删除客户端信息
		RemoveClient(clientId string) error
	}

	// ClientInfo 客户端信息
	ClientInfo struct {
		ClientId             string
		ClientSecret         string
		GrantType            []common.GrantType
		Scope                []string
		AccessTokenDuration  time.Duration
		RefreshTokenDuration time.Duration
	}
)
