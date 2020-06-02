# polestar
Golang based OAuth2 server & client

# Copyright
Copyright 2020 wangghyz/kouki All Rights Reserved.

# Use the server
```go
func main() {
    // 系统配置
    // Application configuration
    appConfig := util.ApplicationConfig()

    // 开启数据库
    // Open the database(MySQL)
    dbConn := db.NewMySQLConnectionInstance()
    defer func() {
        if dbConn != nil {
            // 释放数据库链接
            // Release db connection
            dbConn.Close()
        }
    }()
    
    // Web服务器(gin)
    // web server(gin)
    g := gin.Default()

    // 初始化认证服务器
    // Init the auth server
    err := server.InitGinAuthServer(g, func() (basePath string, allowedMethods []string, handlers gin.HandlersChain) {
        return "/token", []string{http.MethodPost, http.MethodGet}, nil
    })
    if err != nil {
        log.Println(err)
        return
    }

    // 启动web服务
    // Run the web server
    g.Run(appConfig.Server.Addr)
}
```

# Use the client
```go
func main() {
    // 系统配置
    // Application config
    appConfig := util.ApplicationConfig()

    // Web服务(gin)
    // Web server(gin)
    g := gin.Default()

    // 开放请求
    // The open apis
    openGroup := g.Group("/open")
    openGroup.GET("/hello", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, "This is accessed without access token!")
    })

    // 认证拦截
    // The auth required apis
    authGroup := g.Group("/api", filter.NewAuthFilterHandler(func(tokenData map[string]interface{}, ctx *gin.Context) error {
        // TODO: 自定义扩展验证
        // TODO: The customized auth logic
        return nil
    }))
    // 业务请求
    // Business apis
    helloGroup := authGroup.Group("/hello")
    // Need auth requests（Configured in `application.yaml`: authUris）
    helloGroup.GET("", WithTokenHelloHandler)
    // Skip auth requests（Configured in `application.yaml`: skipUris）
    helloGroup.POST("", WithoutTokenHelloHandler)

    // 启动服务
    // Run the web server
    g.Run(appConfig.Server.Addr)
}

// WithTokenHelloHandler 含有token的请求
func WithTokenHelloHandler(ctx *gin.Context) {
     // 系统配置
     // Application config
    appConfig := util.ApplicationConfig()

    // 从token中获取数据
    // Get data from token
    // 取得token
    // Get the token from context
    token, err := util.GetTokenFromRequest(ctx)
    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
        return
    }
    // 解析token中的数据
    // Parse token
    tokenData, err := util.ParseJwtToken(token, appConfig.Auth.Jwt.Secret)
    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
        return
    }

    // 获取token中的登录用户
    // Get the login user's name from token data
    userName := tokenData[_const.ClaimsUserName]
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

## client
example/auth/client/example_auth_client.go

# For communication
E-mail：wangghyz@163.com