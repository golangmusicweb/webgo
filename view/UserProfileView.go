package view

import (
	"github.com/gin-gonic/gin"
	"webgo/entity"
	"webgo/utils"
	"crypto/sha1"
	"strings"
	"net/http"
	"encoding/json"
	"fmt"
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

func EmailValidate(email string) (bool, string) {
	var isemail bool = true
	var msg string
	if index := strings.IndexAny(email, "@"); index == -1 {
		isemail = false
		msg = "Please enter the correct email addres"
	}
	return isemail, msg
}

func PhoneValidate(phone int64) (bool, string) {
	var isphone bool = false
	var msg string = "Please enter the correct phone code"

		if phone > 10000000000 {
			isphone = true
			msg = ""
		}

	return isphone, msg
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
		if isemail, _ := EmailValidate(user.Email); isemail == true {
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
		} else if isphone, _ := PhoneValidate(user.Phone); isphone == true {
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
			fmt.Println("nima")
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
// @Success 201 {object} view.ResponseMSG
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


		if isemail, _ := EmailValidate(user.Email); isemail == true {
			if has, _ := orm.Where("email=?", user.Email).Get(user); has == false {
				response.status = http.StatusUnprocessableEntity
				response.message = "The email does not exists"
			} else {
				response.code = 0
				response.status = http.StatusOK
				response.message = "Login success"
			}
		} else if isphone, _ := PhoneValidate(user.Phone); isphone == true {
			if has, _ := orm.Where("phone=?", user.Phone).Get(user); has == false {
				fmt.Println(has, user.Phone, user.Password)
				response.status = http.StatusUnprocessableEntity
				response.message = "The phone does not exists"
			} else {
				fmt.Println(user)
				response.code = 0
				response.status = http.StatusOK
				response.message = "Login success"
			}
		} else {
			response.status = http.StatusBadRequest
			response.message = "The param input error"
		}
	}

	if response.status == http.StatusOK {
		// set token
		newJWT := utils.NewJwt()
		cliams := new(utils.CustomClaims)
		if token, err := newJWT.GenerateToken(cliams); err == nil {
			response.data = map[string]interface{}{"Authorization": token, "id": user.Id}
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
		"data": response.data,
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
