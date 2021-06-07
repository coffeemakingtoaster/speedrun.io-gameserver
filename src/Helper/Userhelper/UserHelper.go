package userHelper

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func ValidateUser(uID string) bool {
	return true
}

func ValidateJWSToken(Usertoken string, secretKey []byte, userName string) (bool, error) {

	tokenData := &Claims{}
	_, err := jwt.ParseWithClaims(Usertoken, tokenData, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}
	if tokenData.Username != userName {
		return false, errors.New("Usernames do not match")
	}
	return true, nil
}
