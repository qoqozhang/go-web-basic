package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qoqozhang/go-web-basic/utils"
	"strings"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authJwt, ok := c.Get("jwt")
		if !ok {
			c.Next()
			return
		}
		auth := c.Request.Header.Get("Authorization")
		token := strings.Split(auth, " ")[1]
		jwt := authJwt.(*utils.MiddlewareJwt)
		valid, data, err := jwt.Parse(token)
		if !valid || err != nil {
			c.Abort()
			return
		}
		c.Set("username", data["username"])
		c.Next()

	}
}
