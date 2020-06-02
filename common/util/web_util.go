package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vibrantbyte/go-antpath/antpath"
	_const "polestar/auth/const"
	"strings"
)

var (
	AntPathMatcher = antpath.New()
)

// GetTokenFromRequest 从请求中获取 access token (param、form、header)
func GetTokenFromRequest(ctx *gin.Context) (string, error) {
	token := ctx.Query(_const.AccessToken)
	if len(token) > 0 {
		return token, nil
	}

	token = ctx.PostForm(_const.AccessToken)
	if len(token) > 0 {
		return token, nil
	}

	token = ctx.GetHeader(_const.Authorization)
	if len(token) <= 0 {
		return "", errors.New("token不存在！")
	}
	if !strings.HasPrefix(strings.ToUpper(token), "BEARER") {
		return "", errors.New("不支持的token类型！")
	}

	return token[7:], nil
}
