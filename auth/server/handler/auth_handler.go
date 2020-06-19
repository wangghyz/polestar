package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	_const "github.com/wangghyz/polestar/auth/const"
	"github.com/wangghyz/polestar/auth/server/store"
	"github.com/wangghyz/polestar/common/service"
	"github.com/wangghyz/polestar/common/util"
)

type (
	// TokenHandlerForm token请求数据体
	TokenHandlerForm struct {
		GrantType    _const.GrantType `form:"grant_type"`
		ClientId     string           `form:"client_id"`
		ClientSecret string           `form:"client_secret"`
		UserName     string           `form:"userName"`
		Password     string           `form:"password"`
		RefreshToken string           `form:"refresh_token"`
	}
)

// NewDefaultTokenHandlerFunc 默认TokenHandler
func NewDefaultTokenHandlerFunc() gin.HandlerFunc {
	clientStore := store.NewMySQLClientStoreInstance()
	tokenStore := store.NewMemoryTokenStoreInstance()

	return NewTokenHandlerFunc(clientStore,
		tokenStore,
		func(clientId, userName string) ([]string, error) {
			// 角色信息
			roles, err := service.NewSysRoleService().GetRolesByUserName(userName)
			if gorm.IsRecordNotFoundError(err) {
				return nil, errors.New(fmt.Sprintf("用户[%s]不存在角色信息！\n", userName))
			} else if err != nil {
				return nil, err
			} else {
				var tmp []string
				for _, role := range roles {
					tmp = append(tmp, role.EnName)
				}
				return tmp, nil
			}
		},
		func(clientId, userName string) ([]string, error) {
			// 权限信息
			permissions, err := service.NewSysPermissionService().GetPermissionsByUserName(userName)
			if gorm.IsRecordNotFoundError(err) {
				return nil, errors.New(fmt.Sprintf("用户[%s]不存在权限信息！\n", userName))
			} else if err != nil {
				log.Println(err)
				return nil, err
			} else {
				var tmp []string
				for _, per := range permissions {
					tmp = append(tmp, per.EnName)
				}
				return tmp, nil
			}
		},
		func(clientId, userName string) (map[string]interface{}, error) {
			session := make(map[string]interface{})

			// 用户信息
			user, err := service.NewSysUserService().GetUserByUserName(userName)
			if gorm.IsRecordNotFoundError(err) {
				return nil, errors.New(fmt.Sprintf("用户[%s]不存在！\n", userName))
			} else if err != nil {
				log.Println(err)
				return nil, err
			} else {
				session = map[string]interface{}{
					"userId":    user.ID,
					"name":      user.Name,
					"headerImg": user.HeaderImage,
					"userName":  user.UserName,
				}
			}

			return session, nil
		},
	)
}

// NewTokenHandlerFunc 生成TokenHandlerFunc
func NewTokenHandlerFunc(
	clientStore store.ClientStore,
	tokenStore store.TokenStore,
	rolesFunc store.JwtAuthRolesGenerator,
	authoritiesFunc store.JwtAuthAuthoritiesGenerator,
	payloadFunc store.JwtPayloadExtendDataGenerator) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		// 取得请求参数
		form := &TokenHandlerForm{}
		err := ctx.Bind(form)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, util.ResponseMessage{
				Code:    util.ResponseStatusSystemError,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		formClientId := form.ClientId
		formClientSecret := form.ClientSecret

		// 获取clientInfo
		clientInfo, err := clientStore.GetClient(formClientId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "客户端信息不存在！",
				Data:    nil,
			})
			return
		}

		// 请求GrantType判断
		switch form.GrantType {
		case _const.GrantTypePasswordCredentials:
			// 密码模式

			flg := false
			for _, v := range clientInfo.GrantType {
				if v == _const.GrantTypePasswordCredentials {
					flg = true
				}
			}
			if !flg {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "客户端不支持的授权模式！",
					Data:    nil,
				})
				return
			}

			// 验证Client Secret
			if !util.VerifySecretOrPassword(clientInfo.ClientSecret, formClientSecret) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "错误的ClientSecret！",
					Data:    nil,
				})
				return
			}

			// 查询用户数据
			formUserName := form.UserName
			formPassword := form.Password

			user, err := service.NewSysUserService().GetUserByUserName(formUserName)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "用户名或密码错误！",
					Data:    nil,
				})
				return
			}
			if !util.VerifySecretOrPassword(user.Password, formPassword) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "用户名或密码错误！",
					Data:    nil,
				})
				return
			}

			// 做成JWT Payload 自定义内容
			var session map[string]interface{}
			if payloadFunc != nil {
				session, err = payloadFunc(formClientId, formUserName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}
			var roles, authorities []string
			if rolesFunc != nil {
				roles, err = rolesFunc(formClientId, formUserName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}
			if authoritiesFunc != nil {
				authorities, err = authoritiesFunc(formClientId, formUserName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}

			// Scope
			scope := clientInfo.Scope

			// 生成Token
			token, refreshToken, err := tokenStore.GenerateToken(clientInfo, formUserName, scope, roles, authorities, session)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: err.Error(),
					Data:    nil,
				})
				return
			}

			// 做成返回消息
			msg := util.ResponseMessage{
				Code: util.ResponseStatusSuccess,
				Data: map[string]string{
					"access_token":  token,
					"refresh_token": refreshToken,
				},
			}

			// 返回数据
			ctx.JSON(http.StatusOK, msg)
			break
		case _const.GrantTypeRefreshToken:
			// token刷新模式

			flg := false
			for _, v := range clientInfo.GrantType {
				if v == _const.GrantTypeRefreshToken {
					flg = true
				}
			}
			if !flg {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "客户端不支持的授权模式！",
					Data:    nil,
				})
				return
			}

			// 解析refresh_token
			tokenBodyMap, err := util.ParseJwtToken(form.RefreshToken, util.ApplicationConfig().Auth.Jwt.Secret)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: err.Error(),
					Data:    nil,
				})
				return
			}

			// 判断token类型
			if _const.TokenType(fmt.Sprintf("%s", tokenBodyMap[_const.ClaimsType])) != _const.TokenTypeRefreshToken {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "无效的刷新token！",
					Data:    nil,
				})
				return
			}

			// 验证Client Secret
			if !util.VerifySecretOrPassword(clientInfo.ClientSecret, formClientSecret) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "错误的Client Secret！",
					Data:    nil,
				})
				return
			}

			userName := fmt.Sprintf("%s", tokenBodyMap[_const.ClaimsUserName])
			if len(userName) <= 0 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: "refresh_token不包含用户信息！",
					Data:    nil,
				})
				return
			}

			// 做成JWT Payload 自定义内容
			var session map[string]interface{}
			if payloadFunc != nil {
				session, err = payloadFunc(formClientId, userName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}
			var roles, authorities []string
			if rolesFunc != nil {
				roles, err = rolesFunc(formClientId, userName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}
			if authoritiesFunc != nil {
				authorities, err = authoritiesFunc(formClientId, userName)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
						Code:    util.ResponseStatusLogicError,
						Message: err.Error(),
						Data:    nil,
					})
					return
				}
			}

			// Scope
			scope := clientInfo.Scope

			// 生成Token
			token, refreshToken, err := tokenStore.GenerateToken(clientInfo, userName, scope, roles, authorities, session)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
					Code:    util.ResponseStatusLogicError,
					Message: err.Error(),
					Data:    nil,
				})
				return
			}

			// 做成返回消息
			ctx.JSON(http.StatusOK, util.ResponseMessage{
				Code: util.ResponseStatusSuccess,
				Data: map[string]string{
					"access_token":  token,
					"refresh_token": refreshToken,
				},
			})
			return
		default:
			// 非法模式
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseMessage{
				Code:    util.ResponseStatusLogicError,
				Message: "不支持的GrantType！",
				Data:    nil,
			})
			return
		}
	}
}


