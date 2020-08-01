package common

type (
	// TokenType
	TokenType string

	// GrantType
	GrantType string
)

const (
	// token类型
	// 访问token
	TokenTypeAccessToken TokenType = "access_token"
	// 刷新token
	TokenTypeRefreshToken TokenType = "refresh_token"

	// GrantTypePasswordCredentials 密码模式
	GrantTypePasswordCredentials GrantType = "password"
	// GrantTypeRefreshToken token刷新模式
	GrantTypeRefreshToken GrantType = "refresh_token"

	// 请求参数token key（param、form、header）
	AccessToken   = "access_token"
	Authorization = "Authorization"

	// token 内容
	// 角色信息
	TokenDataRoles = "roles"
	// 授权信息
	TokenDataAuthorities = "authorities"
	// session数据
	TokenDataSession = "session"

	// token claims
	// token有效期
	ClaimsEXP = "exp"
	// token唯一标示
	ClaimsJTI = "jti"
	// token类型：access_token，refresh_token
	ClaimsType = "type"
	// 客户端ID
	ClaimsClientId = "client_id"
	// 用户名
	ClaimsUserName = "user_name"
	// 客户端scope
	ClaimsScope = "scope"
)
