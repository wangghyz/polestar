package filter

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangghyz/polestar/common"
	"io/ioutil"
	"net/http"
	"strconv"
)

type (
	ExtendAuthFunc func(tokenData map[string]interface{}, ctx *gin.Context) error
)

// NewAuthFilterHandler
func NewAuthFilterHandler(extendAuthFunc ExtendAuthFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := common.ApplicationConfig()

		// 跳过非验证的uri
		skipUris := config.Auth.SkipUris
		reqMethod := ctx.Request.Method
		var skipRequest = false
		for _, su := range skipUris {
			uri := su.Uri
			suMethods := su.Methods

			// uri不匹配
			if !common.AntPathMatcher.Match(uri, ctx.FullPath()) {
				continue
			}

			for _, suMethod := range suMethods {
				if suMethod == "ALL" {
					skipRequest = true
					break
				}
				if suMethod == reqMethod {
					skipRequest = true
					break
				}
			}

			if skipRequest {
				continue
			}
		}

		// 跳过非验证的请求
		if skipRequest {
			return
		}

		// 获取token
		token := common.GetTokenFromRequest(ctx)
		tokenData := common.ParseJwtToken(token, config.Auth.Jwt.Secret)

		if common.TokenType(fmt.Sprintf("%s", tokenData[common.ClaimsType])) != common.TokenTypeAccessToken {
			common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, "无效的token类型！")
		}

		// 认证服务器验证token
		if common.ApplicationConfig().Auth.TokenCheck.CheckAtServer {
			endPoint := common.ApplicationConfig().Auth.TokenCheck.CheckEndpoint
			endPoint += token

			resp, err := http.Get(endPoint)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "服务器端token验证异常！"+err.Error())
			}
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "服务器端token验证异常！"+err.Error())
			}
			if resp.StatusCode != 200 {
				if resp.StatusCode == 900 {
					body := &common.PolestarError{}
					common.HandleErrorToPanicPolestarError(json.Unmarshal(bodyBytes, body), common.ERR_SYS_ERROR)
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, body.Error())
				} else {
					common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "服务器端token验证异常！"+err.Error())
				}
			}

			rst, err := strconv.ParseBool(string(bodyBytes))
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "服务器端token验证异常！"+err.Error())
			}
			if !rst {
				// 验证不通过
				common.PanicPolestarError(common.ERR_HTTP_AUTH_FAILED, "服务器端token验证：无效的token")
			}
		}

		// 权限验证
		// session权限信息
		sessionAuthorities := tokenData[common.TokenDataAuthorities].([]interface{})
		if sessionAuthorities == nil || len(sessionAuthorities) <= 0 {
			common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, "没有权限！")
		}

		// 进行权限验证
		// 非跳过验证的请求，进行权限验证
		var isOK = false

		// 权限配置信息
		authUris := config.Auth.AuthUris
		for _, uriInfo := range authUris {
			uri := uriInfo.Uri
			uiMethods := uriInfo.Methods
			uiAuthorities := uriInfo.Authorities

			// uri不匹配
			if !common.AntPathMatcher.Match(uri, ctx.FullPath()) {
				continue
			}

			var methodMatch = false
			for _, uiMethod := range uiMethods {
				if uiMethod == reqMethod {
					methodMatch = true
					break
				}
			}
			if !methodMatch {
				continue
			}

			for _, uiAuthority := range uiAuthorities {
				for _, sessionAuthority := range sessionAuthorities {
					if uiAuthority == sessionAuthority.(string) {
						isOK = true
						break
					}
				}
				if isOK {
					break
				}
			}
			if isOK {
				break
			}
		}

		if !isOK {
			common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, "没有权限！")
		}

		// 自定义扩展验证
		if extendAuthFunc != nil {
			err := extendAuthFunc(tokenData, ctx)
			if err != nil {
				common.PanicPolestarError(common.ERR_HTTP_REQUEST_ERROR, err.Error())
			}
		}
	}
}
