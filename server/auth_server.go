package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/server/handler"
	"net/http"
	"strings"
)

type (
	TokenEndpointGenerator func() (basePath string, allowedMethods []string, handlers gin.HandlersChain)
)

// InitGinAuthServer 初始化Web认证服务
func InitGinAuthServer(engine *gin.Engine, tokenEndpoint TokenEndpointGenerator) {
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
				common.PanicPolestarError(common.ERR_BUSINESS_ERROR, fmt.Sprintf("不支持的授权端点请求类型！[%s]", method))
			}
		}
	} else {
		isPost = true
	}
	if tokenHandlers == nil {
		tokenHandlers = []gin.HandlerFunc{}
	}
	tokenGroup := engine.Group(tokenBasePath, tokenHandlers...)
	if isGet {
		tokenGroup.GET("", handler.NewDefaultTokenHandlerFunc())
	}
	if isPost {
		tokenGroup.POST("", handler.NewDefaultTokenHandlerFunc())
	}
}
