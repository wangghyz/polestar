package store

type (
	// TokenStore TokenStore接口
	//
	// 目前实现：
	// 内存类型：store.memoryClientStore
	// MySQL类型: store.mySQLClientStore
	TokenStore interface {
		// 生成token
		GenerateToken(
			clientInfo *ClientInfo,
			userName string,
			scope []string,
			roles []string,
			authorities []string,
			customPayload map[string]interface{}) (accessToken, refreshToken string, err error)

		// 刷新token
		RefreshToken(
			clientInfo *ClientInfo,
			userName string,
			scope []string,
			roles []string,
			authorities []string,
			customPayload map[string]interface{}) (accessToken, refreshToken string, err error)

		// 检查AccessToken
		CheckAccessToken(accessToken string) bool
	}
)
