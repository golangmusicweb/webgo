package middlewares

import (
	"webgo/log"
	"webgo/apps/userprofile/utils"
	"strings"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	var logger log.Logging
	logger.GetLogger()
	defer logger.Close()
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
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
		claims, err := j.ParseToken(strings.Split(token, " ")[1])
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
		newJWT := utils.NewJwt()
		newtoken, _ := newJWT.RefreshToken(strings.Split(token, " ")[1])
		cookie := &http.Cookie{
			Name: "Authorization",
			Value: "JWT" + newtoken,
			Path: "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		c.Set("isPass", true)
		c.Set("claims", claims)
	}
}

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func NewAuthorizer(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		a := &BasicAuthorizer{enforcer: e}
		claims, _ := c.Get("claims")

		if !a.CheckPermission(c.Request, claims) {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg": "用户角色无权限访问",
			})
			return
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *casbin.Enforcer
}


// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request, claims interface{}) bool {
	role := ""
	if claims == nil {
		return false
	} else {
		switch t := claims.(type) {
		case utils.CustomClaims:
			role = t.Role
		}
	}
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(role, path, method)
}
