package userHelper

import (
	"fmt"
	"testing"
)

func TestValidateJWSToken(t *testing.T) {

	//test valid config
	ok, err := ValidateJWSToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGUifQ.3dVU1BrCRHDx6OEqc5K1KvpkuJvyf4j3u6_dOdnNFDM", []byte("your-256-bit-secret"), "example")
	if !ok || err != nil {
		fmt.Println(err)
		t.Errorf("False alarm")
	}

	//test incorrect username
	ok, err = ValidateJWSToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGUifQ.3dVU1BrCRHDx6OEqc5K1KvpkuJvyf4j3u6_dOdnNFDM", []byte("your-256-bit-secret"), "not example")
	if err == nil || ok {
		fmt.Println(err)
		t.Errorf("Wrong username was accepted")
	}

	//test incorrect token => one character changed
	ok, err = ValidateJWSToken("eyJhbGciOiJIUzI1NiIsInR5cCI6kpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGUifQ.3dVU1BrCRHDx6OEqc5K1KvpkuJvyf4j3u6_dOdnNFDM", []byte("your-256-bit-secret"), "example")
	if err == nil || ok {
		fmt.Println(err)
		t.Errorf("Invalid token should throw an error")
	}

	//Guests should always get validated
	ok, err = ValidateJWSToken("", []byte("your-256-bit-secret"), "Guest123")
	if err != nil || !ok {
		fmt.Println(err)
		t.Errorf("Invalid token should throw an error")
	}

}
