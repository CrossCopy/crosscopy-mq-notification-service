package notification_service

import (
	"strings"
	"testing"
)

func TestConstructSignupEmailVerificationRedisKeySuccess(t *testing.T) {
	key, err := ConstructSignupEmailVerificationRedisKey("dev", "dev@crosscopy.io", "1024")
	if err != nil {
		t.Errorf("Don't expect error here")
	}
	expected := "signup:email-verification:dev:dev@crosscopy.io:1024"
	if key != expected {
		t.Errorf("Wrong key generated, expecting: %s, got: %s", expected, key)
	}
}

func TestConstructSignupEmailVerificationRedisKeyFailWithEmptyParameters(t *testing.T) {
	_, err := ConstructSignupEmailVerificationRedisKey("", "dev@crosscopy.io", "1024")
	if !(err != nil && strings.Contains(err.Error(), "empty username")) {
		t.Errorf("Expect 'empty username' Error")
	}
	_, err = ConstructSignupEmailVerificationRedisKey("dev", "", "1024")
	if !(err != nil && strings.Contains(err.Error(), "empty email")) {
		t.Errorf("Expect 'empty esername' Error")
	}
	_, err = ConstructSignupEmailVerificationRedisKey("dev", "dev@crosscopy.io", "")
	if !(err != nil && strings.Contains(err.Error(), "empty code")) {
		t.Errorf("Expect 'empty code' Error")
	}
}
