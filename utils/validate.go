package utils

import "regexp"

func ValidateEmail(email string) bool {
	regex := regexp.MustCompile(`[\w\.]+@\w+\.[a-z]{2,3}(\.[a-z]{2,3})?`)
	return regex.MatchString(email)
}
