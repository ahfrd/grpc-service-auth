package helpers

import "regexp"

func RegexNumber(phoneNumber string) bool {
	sampleRegexp := regexp.MustCompile(`\d+`)
	match := sampleRegexp.MatchString(phoneNumber)
	return match
}
