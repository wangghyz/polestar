package store

type (
	// TokenStore TokenStore接口
	//
	// 目前实现：
	// 内存类型：store.memoryClientStore
	// MySQL类型: store.mySQLClientStore
	TokenStore interface {
		// 生成token
		GenerateToken(clientInfo *ClientInfo, userName string, scope []string, roles []string, authorities []string, payloadExtendData map[string]interface{}) (accessToken, refreshToken string, err error)
	}

	// JwtPayloadExtendDataGenerator Jwt Token Payload 自定义内容生成器
	// 返回：map[string]interface{}
	JwtPayloadExtendDataGenerator func(clientId, userName string) (map[string]interface{}, error)

	// Jwt中的角色信息
	JwtAuthRolesGenerator func(clientId, userName string) ([]string, error)
	// Jwt中的权限信息
	JwtAuthAuthoritiesGenerator func(clientId, userName string) ([]string, error)
)
