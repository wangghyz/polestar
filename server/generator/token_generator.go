package generator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/server/store"
	"github.com/wangghyz/polestar/service"
	"net/http"
)

type (
	// TokenHandlerForm token请求数据体
	TokenHandlerForm struct {
		GrantType    common.GrantType `form:"grant_type" json:"grant_type"`
		ClientId     string           `form:"client_id" json:"client_id"`
		ClientSecret string           `form:"client_secret" json:"client_secret"`
		UserName     string           `form:"userName" json:"userName"`
		Password     string           `form:"password" json:"password"`
		RefreshToken string           `form:"refresh_token" json:"refresh_token"`
	}

	// Jwt中的角色信息
	JwtRolesGenerator func(clientId, userName string) ([]string, error)
	// Jwt中的权限信息
	JwtAuthoritiesGenerator func(clientId, userName string) ([]string, error)
	// JwtCustomPayloadGenerator Jwt Token Payload 自定义内容生成器
	JwtCustomPayloadGenerator func(clientId, userName string) (map[string]interface{}, error)
)

// JwtTokenGenerator 生成JWT token
func JwtTokenGenerator(
	clientStore store.ClientStore,
	tokenStore store.TokenStore,
	rolesFunc JwtRolesGenerator,
	authoritiesFunc JwtAuthoritiesGenerator,
	payloadFunc JwtCustomPayloadGenerator) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		// 取得请求参数
		form := &TokenHandlerForm{}
		err := ctx.Bind(form)
		if err != nil {
			common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
		}
		formClientId := form.ClientId
		formClientSecret := form.ClientSecret

		// 获取clientInfo
		clientInfo, err := clientStore.GetClient(formClientId)
		if err != nil {
			common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "客户端信息不存在!")
		}

		// 请求GrantType判断
		switch form.GrantType {
		case common.GrantTypePasswordCredentials:
			// 密码模式
			flg := false
			for _, v := range clientInfo.GrantType {
				if v == common.GrantTypePasswordCredentials {
					flg = true
				}
			}
			if !flg {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "客户端不支持的授权模式!")
			}

			// 验证Client Secret
			if !common.VerifySecretOrPassword(clientInfo.ClientSecret, formClientSecret) {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "错误的ClientSecret!")
			}

			// 查询用户数据
			formUserName := form.UserName
			formPassword := form.Password

			user, err := service.NewSysUserService().GetUserByUserName(formUserName)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "用户名或密码错误!")
			}
			if !common.VerifySecretOrPassword(user.Password, formPassword) {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "用户名或密码错误!")
			}

			// 做成JWT Payload 自定义内容
			var session map[string]interface{}
			if payloadFunc != nil {
				session, err = payloadFunc(formClientId, formUserName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}
			var roles, authorities []string
			if rolesFunc != nil {
				roles, err = rolesFunc(formClientId, formUserName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}
			if authoritiesFunc != nil {
				authorities, err = authoritiesFunc(formClientId, formUserName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}

			// Scope
			scope := clientInfo.Scope

			// 生成Token
			token, refreshToken, err := tokenStore.GenerateToken(clientInfo, formUserName, scope, roles, authorities, session)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
			}

			// 返回数据
			ctx.JSON(http.StatusOK, gin.H{
				"access_token":  token,
				"refresh_token": refreshToken,
			})
			break
		case common.GrantTypeRefreshToken:
			// token刷新模式
			flg := false
			for _, v := range clientInfo.GrantType {
				if v == common.GrantTypeRefreshToken {
					flg = true
				}
			}
			if !flg {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "客户端不支持的授权模式！")
			}

			// 解析refresh_token
			tokenBodyMap := common.ParseJwtToken(form.RefreshToken, common.ApplicationConfig().Auth.Jwt.Secret)

			// 判断token类型
			if common.TokenType(fmt.Sprintf("%s", tokenBodyMap[common.ClaimsType])) != common.TokenTypeRefreshToken {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "无效的刷新token！")
			}

			// 验证Client Secret
			if !common.VerifySecretOrPassword(clientInfo.ClientSecret, formClientSecret) {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "错误的Client Secret！")
			}

			userName := fmt.Sprintf("%s", tokenBodyMap[common.ClaimsUserName])
			if len(userName) <= 0 {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "refresh_token不包含用户信息！")
			}

			// 做成JWT Payload 自定义内容
			var session map[string]interface{}
			if payloadFunc != nil {
				session, err = payloadFunc(formClientId, userName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}
			var roles, authorities []string
			if rolesFunc != nil {
				roles, err = rolesFunc(formClientId, userName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}
			if authoritiesFunc != nil {
				authorities, err = authoritiesFunc(formClientId, userName)
				if err != nil {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
				}
			}

			// Scope
			scope := clientInfo.Scope

			// 生成Token
			token, refreshToken, err := tokenStore.RefreshToken(clientInfo, userName, scope, roles, authorities, session)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, err.Error())
			}

			// 做成返回消息
			ctx.JSON(http.StatusOK, gin.H{
				"access_token":  token,
				"refresh_token": refreshToken,
			})
			return
		default:
			// 非法模式
			common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "不支持的GrantType！")
		}
	}
}
