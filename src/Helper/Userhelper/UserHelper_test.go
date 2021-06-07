package userHelper

import (
	"testing"
)

func TestValidateJWSToken(t *testing.T) {

	//test valid config
	ok, err := ValidateJWSToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGUiLCJqdGkiOiIzMGMzN2U1NC00NWJlLTQzNDgtYjQ0OC03ZjA2Nzc1MmJiOWUiLCJpYXQiOjE2MjMwODcxMDAsImV4cCI6MTYyMzA5MDcwMH0._az5Q6pwkQgAMjD1f3aBQ5IpEpai71CF6V_3557pUGE", []byte("secret"), "example")
	if !ok || err != nil {
		t.Errorf("False alarm")
	}

	//test incorrect username
	ok, err = ValidateJWSToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGUiLCJqdGkiOiIzMGMzN2U1NC00NWJlLTQzNDgtYjQ0OC03ZjA2Nzc1MmJiOWUiLCJpYXQiOjE2MjMwODcxMDAsImV4cCI6MTYyMzA5MDcwMH0._az5Q6pwkQgAMjD1f3aBQ5IpEpai71CF6V_3557pUGE", []byte("secret"), "not example")
	if err == nil || ok {
		t.Errorf("Wrong username was accepted")
	}

	//test incorrect token => one character changed
	ok, err = ValidateJWSToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ8.eyJ1c2VybmFtZSI6ImV4YW1wbGUiLCJqdGkiOiIzMGMzN2U1NC00NWJlLTQzNDgtYjQ0OC03ZjA2Nzc1MmJiOWUiLCJpYXQiOjE2MjMwODcxMDAsImV4cCI6MTYyMzA5MDcwMH0._az5Q6pwkQgAMjD1f3aBQ5IpEpai71CF6V_3557pUGE", []byte("secret"), "example")
	if err == nil || ok {
		t.Errorf("Invalid token should throw an error")
	}
}
