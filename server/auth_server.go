package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/server/generator"
	"github.com/wangghyz/polestar/server/store"
	"net/http"
	"strings"
)

type (
	// Token 端点生成器
	TokenEndpointGenerator func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain)
	// Token Check 端点生成器
	CheckTokenEndpointGenerator func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain)
	// 端点生成器
	Oauth2EndpointGenerator func() (tokenEndpoint TokenEndpointGenerator, checkTokenEndpoint CheckTokenEndpointGenerator)
)

// InitGinAuthServer 初始化Web认证服务
func InitGinAuthServer(
	engine *gin.Engine,
	tokenEndpoint Oauth2EndpointGenerator,
	clientStore store.ClientStore,
	tokenStore store.TokenStore,
	rolesFunc generator.JwtRolesGenerator,
	authoritiesFunc generator.JwtAuthoritiesGenerator,
	customPayloadFunc generator.JwtCustomPayloadGenerator) {

	var teg TokenEndpointGenerator
	var cteg CheckTokenEndpointGenerator
	if tokenEndpoint == nil {
		teg = func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain) {
			return "/token", []string{http.MethodPost}, []gin.HandlerFunc{}
		}
		cteg = func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain) {
			return "/check", []string{http.MethodGet}, []gin.HandlerFunc{}
		}
	} else {
		teg, cteg = tokenEndpoint()
	}

	// Token 生成端点 默认（/token）
	serveTokenEndpoint(engine, teg, clientStore, tokenStore, rolesFunc, authoritiesFunc, customPayloadFunc)
	// Token 检查端点 默认（/check）
	serveCheckTokenEndpoint(engine, cteg, tokenStore)
}

// DefaultOauth2EndpointGenerator 默认端点生成器
func DefaultOauth2EndpointGenerator() Oauth2EndpointGenerator {
	return func() (tokenEndpoint TokenEndpointGenerator, checkTokenEndpoint CheckTokenEndpointGenerator) {
		// token 生成端点
		tokenEndpoint = func() (basePath string, allowedMethods []string, handlers gin.HandlersChain) {
			return "/token", []string{http.MethodPost, http.MethodGet}, nil
		}
		// token 检查端点
		checkTokenEndpoint = func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain) {
			return "/check", []string{http.MethodPost, http.MethodGet}, nil
		}
		return tokenEndpoint, checkTokenEndpoint
	}
}

// serveTokenEndpoint 处理Token Endpoint
func serveTokenEndpoint(
	engine *gin.Engine,
	endpoint TokenEndpointGenerator,
	clientStore store.ClientStore,
	tokenStore store.TokenStore,
	rolesFunc generator.JwtRolesGenerator,
	authoritiesFunc generator.JwtAuthoritiesGenerator,
	customPayloadFunc generator.JwtCustomPayloadGenerator) {

	var tokenBasePath string
	var tokenMethods []string
	var customMiddleware []gin.HandlerFunc

	// 参数获取
	if endpoint == nil {
		tokenBasePath = "/token"
		tokenMethods = []string{http.MethodPost}
		customMiddleware = []gin.HandlerFunc{}
	} else {
		tokenBasePath, tokenMethods, customMiddleware = endpoint()
	}

	// 认证端点
	// 验证tokenMethods，默认为POST
	var isGet, isPost bool
	if tokenMethods != nil && len(tokenMethods) > 0 {
		for _, method := range tokenMethods {
			switch strings.ToUpper(method) {
			case http.MethodGet:
				isGet = true
				continue
			case http.MethodPost:
				isPost = true
				continue
			default:
				common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, fmt.Sprintf("不支持的授权端点请求类型！[%s]", method))
			}
		}
	} else {
		isPost = true
	}
	if customMiddleware == nil {
		customMiddleware = []gin.HandlerFunc{}
	}
	tokenGroup := engine.Group(tokenBasePath, customMiddleware...)
	handlerFunc := generator.JwtTokenGenerator(clientStore, tokenStore, rolesFunc, authoritiesFunc, customPayloadFunc)
	if isGet {
		tokenGroup.GET("", handlerFunc)
	}
	if isPost {
		tokenGroup.POST("", handlerFunc)
	}
}

// serveCheckTokenEndpoint 处理Check Token Endpoint
func serveCheckTokenEndpoint(
	engine *gin.Engine,
	endpoint CheckTokenEndpointGenerator,
	tokenStore store.TokenStore) {

	var tokenBasePath string
	var tokenMethods []string
	var customMiddleware []gin.HandlerFunc

	if endpoint == nil {
		tokenBasePath = "/check"
		tokenMethods = []string{http.MethodGet}
		customMiddleware = []gin.HandlerFunc{}
	} else {
		tokenBasePath, tokenMethods, customMiddleware = endpoint()
	}

	var isGet, isPost bool
	if tokenMethods != nil && len(tokenMethods) > 0 {
		for _, method := range tokenMethods {
			switch strings.ToUpper(method) {
			case http.MethodGet:
				isGet = true
				continue
			case http.MethodPost:
				isPost = true
				continue
			default:
				common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, fmt.Sprintf("不支持的授权端点请求类型！[%s]", method))
			}
		}
	} else {
		isPost = true
	}
	if customMiddleware == nil {
		customMiddleware = []gin.HandlerFunc{}
	}

	tokenGroup := engine.Group(tokenBasePath, customMiddleware...)
	handlerFunc := func(c *gin.Context) {
		var token string
		switch c.Request.Method {
		case http.MethodGet:
			token = c.Query("access_token")
			break
		case http.MethodPost:
			token = c.Query("access_token")
			if token == "" {
				token = c.PostForm("access_token")
			}
			break
		default:
			common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, fmt.Sprintf("不支持的端点请求类型！[%s]", c.Request.Method))
		}

		rst := tokenStore.CheckAccessToken(token)
		c.JSON(http.StatusOK, rst)
	}
	if isGet {
		tokenGroup.GET("", handlerFunc)
	}
	if isPost {
		tokenGroup.POST("", handlerFunc)
	}
}
