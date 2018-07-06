package validator

import "strings"

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