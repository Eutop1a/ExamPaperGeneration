package config

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
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
		// 这里需要实现登录状态检查逻辑
		// 示例：假设通过 Session 检查登录状态
		// isLoggedIn := checkLoginStatus(c)
		// if!isLoggedIn {
		//     c.JSON(http.StatusForbidden, gin.H{
		//         "status":  403,
		//         "message": "Permission denied",
		//     })
		//     c.Abort()
		//     return
		// }
		c.Next()
	}
}

// SetupWebConfig 配置 Web 服务，包括资源映射和拦截器
func SetupWebConfig(r *gin.Engine) {
	// 添加资源映射
	r.Static("/", "./static")
	r.Static("/static", "./static")

	// 添加拦截器
	r.Use(LoginHandlerInterceptor())
}
