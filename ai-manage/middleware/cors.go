package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 直接放行所有跨域请求并放行所有 OPTIONS 方法
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept,Origin,User-Agent,DNT,Cache-Control,X-Mx-ReqToken,X-Data-Type,X-Requested-With,X-Data-Type,X-Auth-Token,header,clientParams,platform,client,Content-Disposition,APIKEY")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS,HEAD,PUT,DELETE,PATCH")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type,platform,Authorization,Content-Disposition")
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}
