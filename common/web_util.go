package common

import (
	"github.com/gin-gonic/gin"
	"github.com/vibrantbyte/go-antpath/antpath"
	"strings"
)

var (
	AntPathMatcher = antpath.New()
)

// GetTokenFromRequest 从请求中获取 access token (param、form、header)
func GetTokenFromRequest(ctx *gin.Context) string {
	token := ctx.Query(AccessToken)
	if len(token) > 0 {
		return token
	}

	token = ctx.PostForm(AccessToken)
	if len(token) > 0 {
		return token
	}

	token = ctx.GetHeader(Authorization)
	if len(token) <= 0 {
		PanicPolestarError(ERR_HTTP_AUTH_FAILED, "token不存在！")
	}
	if !strings.HasPrefix(strings.ToUpper(token), "BEARER") {
		PanicPolestarError(ERR_HTTP_AUTH_FAILED, "不支持的token类型！")
	}

	return token[7:]
}

// GetLoginUserNameFromToken 从请求的token中获取登录用户名
func GetLoginUserNameFromToken(ctx *gin.Context, key string) string {
	jwtToken := ParseJwtToken(GetTokenFromRequest(ctx), key)
	return jwtToken["user_name"].(string)
}
