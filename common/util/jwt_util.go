package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	_const "polestar/auth/const"
)

func GenerateJwtToken(tokenType _const.TokenType, clientId, userName, key string, jti string, exp int64, scope []string, roles []string, authorities []string, payloadExtendData map[string]interface{}) (tokenString string, err error) {
	claims := jwt.MapClaims{}

	claims[_const.ClaimsEXP] = exp
	claims[_const.ClaimsJTI] = jti
	claims[_const.ClaimsType] = tokenType
	claims[_const.ClaimsClientId] = clientId
	claims[_const.ClaimsUserName] = userName
	claims[_const.ClaimsScope] = scope

	if roles != nil {
		claims[_const.TokenDataRoles] = roles
	}
	if authorities != nil {
		claims[_const.TokenDataAuthorities] = authorities
	}
	if payloadExtendData != nil {
		claims[_const.TokenDataSession] = payloadExtendData
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// ParseJwtToken 解析jwt token
func ParseJwtToken(tokenString, key string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (k interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的加密模式: %v", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("无效的token！")
	}
}

// VerifySecretOrPassword 验证Secret或者密码
func VerifySecretOrPassword(cipher string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(cipher), []byte(plain))
	return err == nil
}