package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/server"
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

	// 开启数据库
	dbConn := db.NewMySQLConnectionInstance()
	defer func() {
		if dbConn != nil {
			// 释放数据库链接
			dbConn.Close()
		}
	}()

	engine := web.NewPolestarWebEngine()

	// 初始化认证服务器
	server.InitGinAuthServer(engine, func() (basePath string, allowedMethods []string, handlers gin.HandlersChain) {
		return "/token", []string{http.MethodPost, http.MethodGet}, nil
	})

	// 启动web服务
	engine.Run(appConfig.Server.Addr)
}
