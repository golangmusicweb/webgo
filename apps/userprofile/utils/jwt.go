package utils

import (
	"github.com/dgrijalva/jwt-go"
	"errors"
	"fmt"
	"time"
	"strconv"
	"webgo/setting"
)


var (
	TokenExpired error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed error = errors.New("That's not even a token")
	TokenInvalid error = errors.New("Couldn't handle this token:")
	MySign string="dongxy")

type CustomClaims struct {
	Id int64
	Username interface{}
	jwt.StandardClaims
}

type Jwt struct {
	SignKey []byte
}

func NewJwt() *Jwt {
	return &Jwt{
		SignKey: []byte(MySign),
	}
}

// Create token: building and signing
func (j *Jwt) GenerateToken(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)  //key: type of interface
	return token.SignedString(j.SignKey) // key: interface ====> string
}

// Parse token: parsing and validating
func (j *Jwt) ParseToken(tokenString string) (interface{}, error){
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	fmt.Println(token.Claims)
	fmt.Println(TokenInvalid)
	return nil, TokenInvalid
}

// Refresh Token
func (j *Jwt) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now

		// get expired from settting
		var config setting.Config
		config.LoadConfig()
		exp, _ := strconv.Atoi(config.Token["expired"])
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(exp) * time.Minute).Unix()
		return j.GenerateToken(claims)
	}
	return "", TokenInvalid
}

