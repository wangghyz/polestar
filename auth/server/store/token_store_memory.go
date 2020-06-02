package store

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	_const "polestar/auth/const"
	"polestar/common/util"
	"time"
)

var (
	_memoryTokenStore *memoryTokenStore
)

type (
	memoryTokenStore struct {
	}
)

func NewMemoryTokenStoreInstance() *memoryTokenStore {
	if _memoryTokenStore == nil {
		_memoryTokenStore = &memoryTokenStore{}
	}
	return _memoryTokenStore
}

//GenerateToken 生成Token
func (s *memoryTokenStore) GenerateToken(clientInfo *ClientInfo, userName string, scope []string, roles []string, authorities []string, payloadExtendData map[string]interface{}) (accessToken, refreshToken string, err error) {
	if clientInfo == nil {
		return "", "", errors.New("客户端信息不能为空！")
	}

	isRefresh := false
	for _, item := range clientInfo.GrantType {
		if item == _const.GrantTypeRefreshToken {
			isRefresh = true
			break
		}
	}

	jti := uuid.NewV4().String()

	tokenSecret := util.ApplicationConfig().Auth.Jwt.Secret

	now := time.Now()
	accessToken, err = util.GenerateJwtToken(
		_const.TokenTypeAccessToken,
		clientInfo.ClientId,
		userName,
		tokenSecret,
		jti,
		now.Add(clientInfo.AccessTokenDuration).Unix(),
		scope,
		roles,
		authorities,
		payloadExtendData)

	if err != nil {
		return "", "", err
	}
	if isRefresh {
		refreshToken, err = util.GenerateJwtToken(
			_const.TokenTypeRefreshToken,
			clientInfo.ClientId,
			userName,
			tokenSecret,
			jti,
			now.Add(clientInfo.RefreshTokenDuration).Unix(),
			scope,
			nil,
			nil,
			nil)

		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}
