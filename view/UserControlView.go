package view

import (
	"github.com/gin-gonic/gin"
	"webgo/log"
	"webgo/entity"
	"webgo/utils"
	"crypto/sha1"
	"strconv"
	"time"
	"net/http"
	"fmt"
)

// @Summary Register a new user
// @Description register a new user and secret password
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param   phone     formData    int     true        "register phone code"
// @Param   password     formData    string     true        "register password"
// @Success 200 {string} string	"ok"
// @Router /api/v1/register [POST]
func Register(c *gin.Context) {
	code := -1
	msg := ""
	regphone := c.DefaultPostForm("phone", "")
	var logger log.Logging
	logger.GetLogger()
	defer logger.Close()

	if regphone == "" || len(regphone) != 11 {

		msg = "register username is phone,the length has to be 11"
		logger.Trace(msg)
	} else{
		user_profile := new(entity.UserProfile)
		user_profile.Username = regphone
		phone,_ := strconv.ParseInt(regphone, 10, 64)
		orm := entity.GetDbEngine("default")
		b, e := orm.Where("phone = ?", phone).Get(user_profile)
		fmt.Println(b, e)
		if e != nil || b == true{
			msg = "username: " + regphone + " does exists"
			logger.Trace(msg)
		} else {
			password := c.DefaultPostForm("password", "")
			if len(password) < 6 {
				msg = "password is not long enough"
				logger.Trace(msg)
			} else {
				hashKey := utils.NewHashKey(sha1.New, 64, 64, 15000)
				hashedPassword := hashKey.HashPassword(password)
				cipherText := hashedPassword.CipherText
				//salt := hashed.Salt
				user_profile.Password = cipherText
				user_profile.Phone = phone
				local, _ := time.LoadLocation("Local")
				birthday := c.DefaultPostForm("birthday", time.Now().Format("2006-01-02"))
				user_profile.Birthday,_ = time.ParseInLocation("2006-01-02", birthday, local)
				user_profile.Email = c.DefaultPostForm("email", "")
				user_profile.Address = c.DefaultPostForm("address", "")
				user_profile.Nickname = c.DefaultPostForm("nickname", "")
				user_profile.CreatedAt = time.Now()
				_, err :=orm.Insert(user_profile)
				if err != nil {
					panic(err)
					msg = string(err.Error())
					logger.Trace(msg)
				}else {

					logger.Trace("register success")
				}
			}
		}
	}
	if msg == "" {
		code = 0
		msg = "SUCC"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": msg,
	})

	}

