package store

import (
	"errors"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"github.com/wangghyz/polestar/common"
	"sync"
	"time"
)

var (
	_memoryTokenStore *memoryTokenStore
)

const (
	C_ACCESS_TOKEN  = "access_token"
	C_REFRESH_TOKEN = "refresh_token"
)

type (
	memoryTokenStore struct {
		// 缓存
		// token的存储是有多个键值对的，详细如下：
		// AccessToken：
		//		KEY：jti + "-" + access_token,			VAL: cachedTokenData{jti, userName, accessToken}
		//		KEY：userName + "-" + access_token,		VAL: cachedTokenData{jti, userName, accessToken}
		// RefreshToken:
		//  	KEY：jti + "-" + refresh_token,			VAL: cachedTokenData{jti, userName, refresh_token}
		//		KEY：userName + "-" + refresh_token,	VAL: cachedTokenData{jti, userName, refresh_token}
		// 分别存储jti和userName是方便几个存储值间可以互相找到
		Cache *cache.Cache

		// 同步锁
		Locker sync.Mutex
	}

	// token cache数据类型（其中TokenData可以是accessToken 也可以是refreshToken）
	cachedTokenData struct {
		Jti       string
		UserName  string
		TokenData string
	}
)

// NewMemoryTokenStoreInstance 实例化
func NewMemoryTokenStoreInstance() *memoryTokenStore {
	if _memoryTokenStore == nil {
		var d time.Duration
		if common.ApplicationConfig().Auth.Cache.CleanupInterval.Nanoseconds() <= 0 {
			d = 5 * time.Minute
		} else {
			d = common.ApplicationConfig().Auth.Cache.CleanupInterval * time.Minute
		}
		_memoryTokenStore = &memoryTokenStore{Cache: cache.New(cache.NoExpiration, d)}
	}
	return _memoryTokenStore
}

//GenerateToken 生成Token
func (s *memoryTokenStore) GenerateToken(
	clientInfo *ClientInfo,
	userName string,
	scope []string,
	roles []string,
	authorities []string,
	customPayload map[string]interface{}) (accessToken, refreshToken string, err error) {

	if clientInfo == nil {
		return "", "", errors.New("客户端信息不能为空！")
	}

	isRefresh := false
	for _, item := range clientInfo.GrantType {
		if item == common.GrantTypeRefreshToken {
			isRefresh = true
			break
		}
	}

	// 缓存查询token是否存在
	s.Locker.Lock()
	defer s.Locker.Unlock()
	if tokenData, exists := s.Cache.Get(userName + "-" + C_ACCESS_TOKEN); exists {
		accessToken = tokenData.(cachedTokenData).TokenData
	}
	if isRefresh {
		if tokenData, exists := s.Cache.Get(userName + "-" + C_REFRESH_TOKEN); exists {
			refreshToken = tokenData.(cachedTokenData).TokenData
		}
	}

	jti := uuid.NewV4().String()
	now := time.Now()

	if accessToken == "" {
		// 生成 Access Token
		accessToken, err = generateAccessToken(clientInfo, userName, scope, roles, authorities, customPayload, jti, now)
		if err != nil {
			return "", "", err
		}

		// 存储token
		data := cachedTokenData{
			Jti:       jti,
			UserName:  userName,
			TokenData: accessToken,
		}
		s.Cache.Set(userName+"-"+C_ACCESS_TOKEN, data, clientInfo.AccessTokenDuration)
		s.Cache.Set(jti+"-"+C_ACCESS_TOKEN, data, clientInfo.AccessTokenDuration)
	}
	if isRefresh && refreshToken == "" {
		// 生成 Refresh Token
		refreshToken, err = generateRefreshToken(clientInfo, userName, scope, jti, now)
		if err != nil {
			return "", "", err
		}

		// 存储token
		data := cachedTokenData{
			Jti:       jti,
			UserName:  userName,
			TokenData: refreshToken,
		}
		s.Cache.Set(userName+"-"+C_REFRESH_TOKEN, data, clientInfo.RefreshTokenDuration)
		s.Cache.Set(jti+"-"+C_REFRESH_TOKEN, data, clientInfo.RefreshTokenDuration)
	}

	return accessToken, refreshToken, nil
}

// RefreshToken 刷新Token（之前的AccessToken和RefreshToken作废，重新生成）
func (s *memoryTokenStore) RefreshToken(
	refreshToken string,
	clientInfo *ClientInfo,
	userName string,
	scope []string,
	roles []string,
	authorities []string,
	customPayload map[string]interface{}) (accessTokenRtn, refreshTokenRtn string, err error) {

	// 获取refresh token
	var jti string
	if tokenData, exists := s.Cache.Get(userName + "-" + C_REFRESH_TOKEN); exists {
		if refreshToken != tokenData.(cachedTokenData).TokenData {
			return "", "", errors.New("非法的RefreshToken！")
		}

		jti = tokenData.(cachedTokenData).Jti
	} else {
		return "", "", errors.New("非法的RefreshToken！")
	}

	s.Locker.Lock()
	s.Cache.Delete(userName + "-" + C_ACCESS_TOKEN)
	s.Cache.Delete(userName + "-" + C_REFRESH_TOKEN)
	s.Cache.Delete(jti + "-" + C_ACCESS_TOKEN)
	s.Cache.Delete(jti + "-" + C_REFRESH_TOKEN)
	s.Locker.Unlock()

	// 重新生成token
	return s.GenerateToken(clientInfo, userName, scope, roles, authorities, customPayload)
}

// generateAccessToken 生成access token
func generateAccessToken(clientInfo *ClientInfo,
	userName string,
	scope []string,
	roles []string,
	authorities []string,
	customPayload map[string]interface{},
	jti string,
	generateTime time.Time) (string, error) {
	tokenSecret := common.ApplicationConfig().Auth.Jwt.Secret

	accessToken, err := common.GenerateJwtToken(
		common.TokenTypeAccessToken,
		clientInfo.ClientId,
		userName,
		tokenSecret,
		jti,
		generateTime.Add(clientInfo.AccessTokenDuration).Unix(),
		scope,
		roles,
		authorities,
		customPayload)

	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// generateRefreshToken 生成refresh token
func generateRefreshToken(clientInfo *ClientInfo,
	userName string,
	scope []string,
	jti string,
	generateTime time.Time) (string, error) {

	tokenSecret := common.ApplicationConfig().Auth.Jwt.Secret

	refreshToken, err := common.GenerateJwtToken(
		common.TokenTypeRefreshToken,
		clientInfo.ClientId,
		userName,
		tokenSecret,
		jti,
		generateTime.Add(clientInfo.RefreshTokenDuration).Unix(),
		scope,
		nil,
		nil,
		nil)

	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

// Access Token 有效性检查
func (s *memoryTokenStore) CheckAccessToken(accessToken string) bool {
	claims := common.ParseJwtToken(accessToken, common.ApplicationConfig().Auth.Jwt.Secret)
	_, exists := s.Cache.Get(claims[common.ClaimsJTI].(string) + "-" + C_ACCESS_TOKEN)
	return exists
}
