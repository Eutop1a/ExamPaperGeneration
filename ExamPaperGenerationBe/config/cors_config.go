package config

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 关键修改：当允许凭证时，不能使用通配符*
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, "+
				"session, X_Requested_With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, "+
				"X-CustomHeader, Keep-Alive, User-Agent, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
			c.Header("Access-Control-Max-Age", "172800")
		}

		// 预检请求处理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
