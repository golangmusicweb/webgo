package validator

import (
	"strings"
	"webgo/apps/userprofile/utils"
	"crypto/sha1"
)

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

func PasswordValidate(password, cipherText, salt string)  bool {
	pass := utils.NewHashKey(sha1.New, 64, 64, 15000)
	isPass := pass.VerifyPassword(password, cipherText, salt)
	return isPass
}