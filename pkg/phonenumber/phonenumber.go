package phonenumber

import "regexp"

func IsValid(phoneNumber string) bool {
	var iranMobileRE = regexp.MustCompile(`^(?:\+98|0098|0)?9\d{9}$`)
	// TODO: we can use regular expression to support +98 patterns
	if iranMobileRE.MatchString(phoneNumber) {
		return true
	}

	return false
}
