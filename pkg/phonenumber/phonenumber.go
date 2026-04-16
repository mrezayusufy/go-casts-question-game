package phonenumber

import (
	"regexp"
)

var iranianMobileRegex = regexp.MustCompile(`^09[0-9]{9}$`)

func IsValid(phone string) bool {
	return iranianMobileRegex.MatchString(phone)
}
