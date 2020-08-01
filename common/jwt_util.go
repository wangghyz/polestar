package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJwtToken(tokenType TokenType, clientId, userName, key string, jti string, exp int64, scope []string, roles []string, authorities []string, payloadExtendData map[string]interface{}) (tokenString string, err error) {
	claims := jwt.MapClaims{}

	claims[ClaimsEXP] = exp
	claims[ClaimsJTI] = jti
	claims[ClaimsType] = tokenType
	claims[ClaimsClientId] = clientId
	claims[ClaimsUserName] = userName
	claims[ClaimsScope] = scope

	if roles != nil {
		claims[TokenDataRoles] = roles
	}
	if authorities != nil {
		claims[TokenDataAuthorities] = authorities
	}
	if payloadExtendData != nil {
		claims[TokenDataSession] = payloadExtendData
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// ParseJwtToken 解析jwt token
func ParseJwtToken(tokenString, key string) map[string]interface{} {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (k interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的加密模式: %v", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		PanicPolestarError(ERR_SYS_ERROR, "Token解析错误！"+err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		PanicPolestarError(ERR_SYS_ERROR, "无效的token！")
	}
	return claims
}

// VerifySecretOrPassword 验证Secret或者密码
func VerifySecretOrPassword(cipher string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plain))
	return err == nil
}
