package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"polestar/auth/server/handler"
	"strings"
)

type (
	TokenEndpointGenerator func() (basePath string, allowedMethods []string, handlers gin.HandlersChain)
)

// InitGinAuthServer 初始化Web认证服务
func InitGinAuthServer(g *gin.Engine, tokenEndpoint TokenEndpointGenerator) error {
	var tokenBasePath string
	var tokenMethods []string
	var tokenHandlers []gin.HandlerFunc

	// 参数获取
	if tokenEndpoint == nil {
		tokenBasePath = "/token"
		tokenMethods = []string{http.MethodPost}
		tokenHandlers = []gin.HandlerFunc{}
  	} else {
		tokenBasePath, tokenMethods, tokenHandlers = tokenEndpoint()
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
				return errors.New(fmt.Sprintf("不支持的授权端点请求类型！[%s]", method))
			}
		}
	} else {
		isPost = true
	}
	if tokenHandlers == nil {
		tokenHandlers = []gin.HandlerFunc{}
	}
	tokenGroup := g.Group(tokenBasePath, tokenHandlers...)
	if isGet {
		tokenGroup.GET("", handler.NewDefaultTokenHandlerFunc())
	}
	if isPost {
		tokenGroup.POST("", handler.NewDefaultTokenHandlerFunc())
	}

	return nil
}
