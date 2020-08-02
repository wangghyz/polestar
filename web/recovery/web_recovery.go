package recovery

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
)

// PolestarWebRecovery Gin Recovery 用于处理Panic
// 为避免占用http status code，追加了 900 错误码，用于表示系统业务错误
func PolestarWebRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			if err, ok := common.IsPolestarError(r); ok {
				// 自定义HTTP状态码 900：业务系统错误
				c.Error(err)
				c.AbortWithStatusJSON(900, err)
			} else {
				// 其他异常
				err := errors.New(fmt.Sprintf("%v", r))
				c.Error(err)
				c.AbortWithStatusJSON(900, err.Error())
			}
		}()
		c.Next()
	}
}
