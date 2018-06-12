package view

import (
	"github.com/gin-gonic/gin"
	"webgo/restful"
	"time"
	"webgo/entity"
	"fmt"
)

// @Summary userprofile
// @version v1
// @Description get user info
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Success 200 {string} string	"ok"
// @Failure 404 {string} string	"404"
// @Router /api/v1/user/{id} [get]
// @BasePath /
func User(c *gin.Context) {
	id := c.Param("id")
	dbEngine := restful.GetDbEngine("default")
	dbEngine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	user := new(entity.UserProfile)
	dbEngine.ID(id).Get(user)
	fmt.Println(user)

	c.JSON(200, gin.H{
		"message": "",
		"status":  "posted",
		"data":    user,
	})
}