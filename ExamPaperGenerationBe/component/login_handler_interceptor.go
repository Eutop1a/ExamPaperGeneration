package component

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginHandlerInterceptor 登录拦截器
func LoginHandlerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 排除的路径
		excludePaths := []string{
			"/hello",
			"/getLoginStatus",
			"/login",
			"/permission_denied",
			"/registered",
			"/",
		}
		excludeExtensions := []string{
			".html",
			".css",
			".js",
			".ico",
		}
		path := c.Request.URL.Path
		// 检查是否为排除路径
		for _, excludePath := range excludePaths {
			if path == excludePath || strings.HasPrefix(path, "/static/") {
				c.Next()
				return
			}
		}
		// 检查文件扩展名
		ext := filepath.Ext(path)
		for _, excludeExt := range excludeExtensions {
			if ext == excludeExt {
				c.Next()
				return
			}
		}
		// 检查登录状态
		session := sessions.Default(c)
		username := session.Get("username")
		if username == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "Permission denied",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
