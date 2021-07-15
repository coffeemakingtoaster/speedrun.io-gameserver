package userHelper

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

//JWT struct
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//Validate Usertoken by their token as well as their send username
//Check if identities match
func ValidateJWSToken(Usertoken string, secretKey []byte, userName string) (bool, error) {
	//Guests dont have a valid token
	if len(userName) >= 5 {
		if userName[:5] == "Guest" {
			return true, nil
		}
	}
	tokenData := &Claims{}
	_, err := jwt.ParseWithClaims(Usertoken, tokenData, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}
	if tokenData.Username != userName {
		return false, errors.New("usernames do not match")
	}
	return true, nil
}
