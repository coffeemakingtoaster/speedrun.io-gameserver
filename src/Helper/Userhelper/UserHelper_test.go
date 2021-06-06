package userHelper

import "testing"

func TestValidateUser(t *testing.T) {
	if !ValidateUser("mock") {
		t.Errorf("this can never possibly happen at the moment as there is not API calls...")
	}
}
