# polestar
Golang based OAuth2 server & client

# Copyright
Copyright 2020 wangghyz/kouki All Rights Reserved.

# Call Stack 
```
          web.NewPolestarWebEngine()                                            // gin engin
            |   TokenEndpointGenerator                                          // token endpoint for web requests
            |     |  store.ClientStore                                          // Client Store (Default: store.mySQLClientStore)
            |     |    |  store.TokenStore                                      // Token Store (Default: store.memoryTokenStore)
            |     |    |    |  generator.JwtRolesGenerator                      // User roles
            |     |    |    |    |  generator.JwtAuthoritiesGenerator           // User authorities
            |     |    |    |    |    |  generator.JwtCustomPayloadGenerator    // Custom payload
            |     |    |    |    |    |    |
            |-------------------------------
            ↓
  server.InitGinAuthServer()                                                    // Init auth server by gin
      |
      |----> generator.JwtTokenGenerator()                                      // Generate jwt token
               |
               |----> tokenStore.GenerateToken()                                // Call token store to generate token
                        |
                        |----> common.GenerateJwtToken()                        // Call jwt token generator uitl
```

# Use the server
```go
func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Fatalf("%v\n", r)
		}
	}()

	// 获取配置对象
	appConfig := common.ApplicationConfig()

	// 开启数据库
	dbConn := db.NewMySQLConnectionInstance()
	defer func() {
		if dbConn != nil {
			// 释放数据库链接
			dbConn.Close()
		}
	}()

	// Web 引擎 (gin)
	engine := web.NewPolestarWebEngine()

	// 初始化认证服务器
	server.InitGinAuthServer(
        engine,
        // 认证端点生成器
        server.DefaultOauth2EndpointGenerator(),
        // Client Store
        store.NewMySQLClientStoreInstance(),
        // Token Store
        store.NewMemoryTokenStoreInstance(),
        // 用户角色hook
        generator.DefaultJwtRolesGenerator(),
        // 用户权限hook
        generator.DefaultJwtAuthoritiesGenerator(),
        // JWT Token自定义Payload内容hook
        generator.DefaultJwtCustomPayloadGenerator(),
	)

	// 启动web服务
	engine.Run(appConfig.Server.Addr)
}
```

# Use the client
```go
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
```

# Example
## Server
example/auth/server/example_auth_server_mysql.go

## Client
example/auth/client/example_auth_client.go

# For communication
E-mail：wangghyz@163.com