package filter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_const "github.com/wangghyz/polestar/auth/const"
	"github.com/wangghyz/polestar/common/util"
)

type (
	ExtendAuthFunc func(tokenData map[string]interface{}, ctx *gin.Context) error
)

// NewAuthFilterHandler
func NewAuthFilterHandler(extendAuthFunc ExtendAuthFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := util.ApplicationConfig()

		// 跳过非验证的uri
		skipUris := config.Auth.SkipUris
		reqMethod := ctx.Request.Method
		var skipRequest = false
		for _, su := range skipUris {
			uri := su.Uri
			suMethods := su.Methods

			// uri不匹配
			if !util.AntPathMatcher.Match(uri, ctx.FullPath()) {
				continue
			}

			for _, suMethod := range suMethods {
				if suMethod == "ALL" {
					skipRequest = true
					break
				}
				if suMethod == reqMethod {
					skipRequest = true
					break
				}
			}

			if skipRequest {
				continue
			}
		}

		// 跳过非验证的请求
		if skipRequest {
			return
		}

		// 获取token
		token, err := util.GetTokenFromRequest(ctx)
		if err != nil || len(token) <= 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "token无效！",
				Data:    nil,
			})
			return
		}

		tokenData, err := util.ParseJwtToken(token, config.Auth.Jwt.Secret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: fmt.Sprintf("token无效：%s", err.Error()),
				Data:    nil,
			})
			return
		}

		if _const.TokenType(fmt.Sprintf("%s", tokenData[_const.ClaimsType])) != _const.TokenTypeAccessToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "无效的token类型！",
				Data:    nil,
			})
			return
		}

		// 权限验证


		// session权限信息
		sessionAuthorities := tokenData[_const.TokenDataAuthorities].([]interface{})
		if sessionAuthorities == nil || len(sessionAuthorities) <= 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "没有权限！",
				Data:    nil,
			})
			return
		}

		// 进行权限验证
		// 非跳过验证的请求，进行权限验证
		var isOK = false

		// 权限配置信息
		authUris := config.Auth.AuthUris
		for _, uriInfo := range authUris {
			uri := uriInfo.Uri
			uiMethods := uriInfo.Methods
			uiAuthorities := uriInfo.Authorities

			// uri不匹配
			if !util.AntPathMatcher.Match(uri, ctx.FullPath()) {
				continue
			}

			var methodMatch = false
			for _, uiMethod := range uiMethods {
				if uiMethod == reqMethod {
					methodMatch = true
					break
				}
			}
			if !methodMatch {
				continue
			}

			for _, uiAuthority := range uiAuthorities {
				for _, sessionAuthority := range sessionAuthorities {
					if uiAuthority == sessionAuthority.(string) {
						isOK = true
						break
					}
				}
				if isOK {
					break
				}
			}
			if isOK {
				break
			}
		}

		if !isOK {
			ctx.AbortWithStatusJSON(http.StatusForbidden, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "没有权限！",
				Data:    nil,
			})
			return
		}

		// 自定义扩展验证
		if extendAuthFunc != nil {
			err := extendAuthFunc(tokenData, ctx)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: err.Error(),
					Data:    nil,
				})
				return
			}
		}
	}
}
