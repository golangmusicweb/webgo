package view

import (
	"github.com/gin-gonic/gin"
	"webgo/entity"
	"webgo/utils"
	"crypto/sha1"
	"net/http"
	"encoding/json"
	"time"
	"webgo/setting"
	"strconv"
	"webgo/validator"
)

type UserOperation struct {
	UserName interface{}
	PassWord string
}

type ResponseMSG struct {
	status int
	code int
	message string
	data interface{}
}

func PassValidate(password string) (bool, string) {
	var ispass bool = true
	var msg string
	if len(password) < 6 {
		ispass = false
		msg = "Password length must be greater than 6"
	}
	return ispass, msg
}



func PassSecret(passwork string) (string, string) {
	pass := utils.NewHashKey(sha1.New, 64, 64, 15000)
	hashed := pass.HashPassword(passwork)
	cipherText := hashed.CipherText
	salt := hashed.Salt
	return cipherText, salt
}


// @Summary Register a new user
// @Description register a new username and secret password
// @Accept  json
// @Produce  json
// @Param	register	body	view.UserOperation	true "register a new username and secret password"
// @Success 201 {object} view.ResponseMSG
// @Success 400 {object} view.ResponseMSG
// @Success 422 {object} view.ResponseMSG
// @Router /api/v1/user/register [POST]
func RegisterView(c *gin.Context) {
	register := new(UserOperation)
	user := new(entity.UserProfile)
	response := new(ResponseMSG)

	response.code = -1
	orm := entity.GetDbEngine("default")

	body := c.Request.Body // io.Reader
	err := json.NewDecoder(body).Decode(register)

	if err == nil {
		switch t := register.UserName.(type) {
		case string:
			user.Email = t
		case float64:
			user.Phone = int64(t)
		}
		if isemail, _ := validator.EmailValidate(user.Email); isemail == true {
			if has, err := orm.Where("email=?",user.Email).Get(user); has == true && err == nil {
				response.status = http.StatusUnprocessableEntity
				response.message = "The email does already registed"
			} else {
				if ispass, msg := PassValidate(register.PassWord); ispass == false {
					response.message = msg
				} else {
					cipherText, salt := PassSecret(register.PassWord)
					user.Password = cipherText + salt
					if _, err := orm.Insert(register); err != nil {
						response.status = http.StatusUnprocessableEntity
						response.message = "Register failure"
					} else {
						response.status = http.StatusCreated
						response.code = 0
						response.message = "Register success"
					}
				}
			}
		} else if isphone, _ := validator.PhoneValidate(user.Phone); isphone == true {
			if has, err := orm.Where("phone=?", user.Phone).Get(user); has == true && err == nil {
				response.status = http.StatusUnprocessableEntity
				response.message = "The phone does already registed"
			} else {
				if ispass, msg := PassValidate(register.PassWord); ispass == false {
					response.message = msg
				} else {
					cipherText, salt := PassSecret(register.PassWord)
					user.Password = cipherText + salt
					if _, err := orm.Insert(user); err != nil {
						response.status = http.StatusUnprocessableEntity
						response.message = "Register failure"
					} else {
						response.status = http.StatusCreated
						response.code = 0
						response.message = "Register success"
					}
				}
			}
		} else {
			response.status = http.StatusBadRequest
			response.message = "The param input error"
		}
	} else {
		response.status = http.StatusBadRequest
		response.message = "The param input error"
	}
	c.JSON(response.status, gin.H{
		"status": response.status,
		"code": response.code,
		"msg": response.message,
	})
}

var REGEX_MOBILE string = `^1[358]\d{9}$|^147\d{8}$|^176\d{8}$`


// @Summary User login and get a authorization token
// @Description User login and get a authorization token
// @Accept  json
// @Produce  json
// @Param	login	body	view.UserOperation	true "User login and get a authorization token"
// @Success 200 {object} view.ResponseMSG
// @Success 400 {object} view.ResponseMSG
// @Success 422 {object} view.ResponseMSG
// @Router /api/v1/user/login [POST]
func LoginView(c *gin.Context) {
	login := new(UserOperation)
	user := new(entity.UserProfile)
	response := new(ResponseMSG)

	response.code = -1
	orm := entity.GetDbEngine("default")

	body := c.Request.Body // io.Reader
	err := json.NewDecoder(body).Decode(login)

	if err == nil {
		switch t := login.UserName.(type) {
		case string:
			user.Email = t
		case float64:
			user.Phone = int64(t)
		}


		if isemail, _ := validator.EmailValidate(user.Email); isemail == true {
			if has, _ := orm.Where("email=?", user.Email).Get(user); has == false {
				response.status = http.StatusUnprocessableEntity
				response.message = "The email does not exists"
			} else {
				response.status = http.StatusOK
			}
		} else if isphone, _ := validator.PhoneValidate(user.Phone); isphone == true {
			if has, _ := orm.Where("phone=?", user.Phone).Get(user); has == false {
				response.status = http.StatusUnprocessableEntity
				response.message = "The phone does not exists"
			} else {
				response.status = http.StatusOK
			}
		} else {
			response.status = http.StatusBadRequest
			response.message = "The param input error"
		}
	} else {
		response.status = http.StatusBadRequest
		response.message = "The param input error"
	}

	if response.status == http.StatusOK {
		// set token
		newJWT := utils.NewJwt()
		claims := new(utils.CustomClaims)
		claims.Id = user.Id
		claims.Username = login.UserName
		claims.NotBefore = int64(time.Now().Unix() - 1000)// 签名生效时间
		var config setting.Config
		config.LoadConfig()
		exp, _ := strconv.Atoi(config.Token["expired"])
		claims.ExpiresAt = time.Now().Add(time.Duration(exp) * time.Minute).Unix()// 过期时间
		claims.Issuer = "dongxiaoyi" //签名的发行者
		if token, err := newJWT.GenerateToken(claims); err == nil {
			cookie := http.Cookie{Name: "token", Value: "JWT " + token, Path: "/", HttpOnly: true}
			http.SetCookie(c.Writer, &cookie)
			response.code = 0
			response.status = http.StatusOK
			response.message = "Login success"
		} else {
			response.code = -1
			response.status = http.StatusBadRequest
			response.message = "Generate token error"
		}
	}
	c.JSON(response.status, gin.H{
		"status": response.status,
		"code": response.code,
		"msg": response.message,
	})
}

type DeleteAccount struct {
	Id int64
}

// @Summary Delete user account
// @Description Delete user account
// @Accept  json
// @Produce  json
// @Param	delete	body	view.DeleteAccount	true "Delete user account"
// @Success 204 {object} view.ResponseMSG
// @Success 400 {object} view.ResponseMSG
// @Success 422 {object} view.ResponseMSG
// @Router /api/v1/user/delete [POST]
func DeleteAccountView(c *gin.Context) {
	response := new(ResponseMSG)
	response.code = -1
	orm := entity.GetDbEngine("default")
	user := new(entity.UserProfile)
	account := new(DeleteAccount)
	body := c.Request.Body
	err := json.NewDecoder(body).Decode(account)
	if err == nil {
		user.Id = account.Id
		if _, err := orm.Id(user.Id).Delete(user); err == nil {
			response.status = http.StatusNoContent
			response.code = 0
			response.message = "Account deleted successfully"
		} else {
			response.status = http.StatusUnprocessableEntity
			response.message = "Account deletion failed"
		}
	} else {
		response.status = http.StatusBadRequest
		response.message = "The param input error"
	}
	c.JSON(response.status, gin.H{
		"status": response.status,
		"code": response.code,
		"msg": response.message,
	})
}


// @Summary Delete user account
// @Description Delete user account
// @Produce  json
// @Success 200 {object} view.ResponseMSG
// @Success 400 {object} view.ResponseMSG
// @Success 422 {object} view.ResponseMSG
// @Router /api/test/getdatabytime [GET]
func GetDataByTime(c *gin.Context) {
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	claims := c.MustGet("claims").(*utils.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
		})
	}
}