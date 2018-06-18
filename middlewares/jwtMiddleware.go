package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webgo/log"
	"webgo/utils"
)

func JWTAuth() gin.HandlerFunc {
	var logger log.Logging
	logger.GetLogger()
	defer logger.Close()
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg": "请求未携带token，无权限访问",
			})
			c.Set("isPass", false)
			return
		}
		logger.Trace("get token: ", token)

		j := utils.NewJwt()

		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg": "授权已过期",
				})
				c.Set("isPass", false)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg": err.Error(),
			})
			c.Set("isPass", false)
			return
		}
		c.Set("isPass", true)
		c.Set("claims", claims)
	}
}