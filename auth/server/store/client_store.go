package store

import (
	_const "polestar/auth/const"
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
		GrantType            []_const.GrantType
		Scope                []string
		AccessTokenDuration  time.Duration
		RefreshTokenDuration time.Duration
	}
)
