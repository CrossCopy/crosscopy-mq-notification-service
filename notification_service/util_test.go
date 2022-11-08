package notification_service

import (
	"fmt"
	"regexp"
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

// brute force testing Digit code generator with 100 lengths, each 100 times
func TestGetRandomDigitCode(t *testing.T) {
	for l := 1; l < 100; l++ {
		for i := 1; i < 100; i++ {
			code := GetRandomDigitCode(6)
			if len(code) != 6 {
				t.Errorf("Wrong Length")
			}
			for _, char := range code {
				if match, err := regexp.MatchString("\\d", fmt.Sprintf("%c", char)); !match || err != nil {
					t.Errorf("Wrong Character")
				}
			}
		}
	}

}
