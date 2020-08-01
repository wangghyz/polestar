package web

import (
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"github.com/wangghyz/polestar/web/recovery"
)

// NewPolestarWebEngine 新建web容器
func NewPolestarWebEngine() *gin.Engine {
	engine := gin.New()
	config := common.ApplicationConfig()

	if config.Server.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(gin.Logger(), recovery.PolestarWebRecovery())

	return engine
}
