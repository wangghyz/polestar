package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"polestar/auth/client/filter"
	_const "polestar/auth/const"
	"polestar/common/util"
)

func main() {
	// 系统配置
	appConfig := util.ApplicationConfig()

	// Web服务
	g := gin.Default()

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
	appConfig := util.ApplicationConfig()

	// 从token中获取数据
	// 取得token
	token, err := util.GetTokenFromRequest(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	// 解析token中的数据
	tokenData, err := util.ParseJwtToken(token, appConfig.Auth.Jwt.Secret)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取token中的登录用户
	userName := tokenData[_const.ClaimsUserName]
	ctx.JSON(http.StatusOK, fmt.Sprintf("Hello %s! ", userName))
}

// WithoutTokenHelloHandler 没有token的请求
func WithoutTokenHelloHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello World! ")
}
