package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type MiddlewareJwt struct {
	//  jwt的signingKey
	Key []byte
	// jwt算法method , possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	SigningAlgorithm string
	// token有效期
	Expires time.Duration
	// 结构内部调用的token
	token *jwt.Token
}

var defaultKey = []byte("secret")

// Init 配置部分jwt的默认参数
func (m *MiddlewareJwt) Init() *MiddlewareJwt {
	if m.Key == nil {
		m.Key = defaultKey
	}
	if m.SigningAlgorithm == "" {
		m.SigningAlgorithm = "HS256"
	}
	token := jwt.New(jwt.GetSigningMethod(m.SigningAlgorithm))
	claim := token.Claims.(jwt.MapClaims)
	if m.Expires == 0 {
		m.Expires = time.Hour * 2
	}
	claim["exp"] = time.Now().Add(m.Expires).Unix()
	m.token = token
	return m
}

// SignedString  生成jwt token
// data map[string]interface{}  传入claim 的自定义fields
func (m *MiddlewareJwt) SignedString(data map[string]interface{}) (tokenString string, err error) {
	claim := m.token.Claims.(jwt.MapClaims)
	for k, v := range data {
		claim[k] = v
	}
	claim["iat"] = time.Now().Unix()
	claim["nbf"] = time.Now().Unix()
	tokenString, err = m.token.SignedString(m.Key)
	return
}

// Parse 解析token，验证是否有效，并以 map[string]interface{}的方式返回jwt的claim信息
func (m *MiddlewareJwt) Parse(tokenString string) (valid bool, data map[string]interface{}, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.Key, nil
	})
	switch {
	case token.Valid:
		data = token.Claims.(jwt.MapClaims)
		return true, data, nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return false, nil, err
	case err != nil:
		return false, nil, err
	}
	return
}
