package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/persistence/db"
	"github.com/wangghyz/polestar/server"
	"github.com/wangghyz/polestar/server/store"
	"github.com/wangghyz/polestar/service"
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

	// Web 引擎 (gin)
	engine := web.NewPolestarWebEngine()

	// 初始化认证服务器
	server.InitGinAuthServer(
		engine,
		// 认证端点生成器
		func() (tokenEndpoint server.TokenEndpointGenerator, checkTokenEndpoint server.CheckTokenEndpointGenerator) {
			// token 生成端点
			tokenEndpoint = func() (basePath string, allowedMethods []string, handlers gin.HandlersChain) {
				return "/token", []string{http.MethodPost, http.MethodGet}, nil
			}
			// token 检查端点
			checkTokenEndpoint = func() (path string, allowedMethods []string, customMiddleware gin.HandlersChain) {
				return "/check", []string{http.MethodPost, http.MethodGet}, nil
			}
			return tokenEndpoint, checkTokenEndpoint
		},
		// Client Store
		store.NewMySQLClientStoreInstance(),
		// Token Store
		store.NewMemoryTokenStoreInstance(),
		// 用户角色hook
		func(clientId, userName string) ([]string, error) {
			// 角色信息
			roles, err := service.NewSysRoleService().GetRolesByUserName(userName)
			if err != nil {
				return nil, err
			} else {
				var tmp []string
				for _, role := range roles {
					tmp = append(tmp, role.EnName)
				}
				return tmp, nil
			}
		},
		// 用户权限hook
		func(clientId, userName string) ([]string, error) {
			// 权限信息
			permissions, err := service.NewSysPermissionService().GetPermissionsByUserName(userName)
			if err != nil {
				return nil, err
			} else {
				var tmp []string
				for _, per := range permissions {
					tmp = append(tmp, per.EnName)
				}
				return tmp, nil
			}
		},
		// JWT Token自定义Payload内容hook
		func(clientId, userName string) (map[string]interface{}, error) {
			session := make(map[string]interface{})

			// 用户信息
			user, err := service.NewSysUserService().GetUserByUserName(userName)
			if err != nil {
				return nil, err
			} else {
				session = map[string]interface{}{
					"userId":    user.ID,
					"name":      user.Name,
					"headerImg": user.HeaderImage,
				}
			}

			return session, nil
		},
	)

	// 启动web服务
	engine.Run(appConfig.Server.Addr)
}
