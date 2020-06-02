package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"polestar/auth/server/server"
	"polestar/common/db"
	"polestar/common/util"
)

func main() {
	// 系统配置
	appConfig := util.ApplicationConfig()

	// 开启数据库
	dbConn := db.NewMySQLConnectionInstance()
	defer func() {
		if dbConn != nil {
			// 释放数据库链接
			dbConn.Close()
		}
	}()

	// web server
	g := gin.Default()

	// 初始化认证服务器
	err := server.InitGinAuthServer(g, func() (basePath string, allowedMethods []string, handlers gin.HandlersChain) {
		return "/token", []string{http.MethodPost, http.MethodGet}, nil
	})
	if err != nil {
		log.Println(err)
		return
	}

	// 启动web服务
	g.Run(appConfig.Server.Addr)
}
