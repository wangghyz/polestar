package main

import (
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/server"
	"github.com/wangghyz/polestar/server/generator"
	"github.com/wangghyz/polestar/server/store"
	"github.com/wangghyz/polestar/web"
	"log"
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
