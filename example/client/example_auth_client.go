package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/client/filter"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/web"
	"log"
	"net/http"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Fatalf("%v\n", r)
		}
	}()

	// 获取配置对象
	appConfig := common.ApplicationConfig()

	// Web服务
	g := web.NewPolestarWebEngine()

	// 开放请求
	openGroup := g.Group("/open")
	openGroup.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "This is accessed without access token!")
	})

	// 认证拦截
	authGroup := g.Group("/api", filter.NewAuthFilterHandler(func(tokenData map[string]interface{}, ctx *gin.Context) error {
		// TODO: 自定义扩展验证
		return nil
	}))

	// 业务请求
	helloGroup := authGroup.Group("/hello")
	helloGroup.GET("", WithTokenHelloHandler)
	helloGroup.POST("", WithoutTokenHelloHandler)

	// 启动服务
	g.Run(appConfig.Server.Addr)
}

// WithTokenHelloHandler 含有token的请求
func WithTokenHelloHandler(ctx *gin.Context) {
	appConfig := common.ApplicationConfig()
	// 获取token中的登录用户
	userName := common.GetLoginUserNameFromToken(ctx, appConfig.Auth.Jwt.Secret)
	ctx.JSON(http.StatusOK, fmt.Sprintf("Hello %s! ", userName))
}

// WithoutTokenHelloHandler 没有token的请求
func WithoutTokenHelloHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello World! ")
}
