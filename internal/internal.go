package internal

import (
	"errors"
	"fmt"
	"regexp"
)

func FormatPhone(p string) (string, error) {
	stripCountryCodeExp := regexp.MustCompile(`^\+52`)
	replaceExp := regexp.MustCompile(`[ -]`)
	numExp := regexp.MustCompile(`[0-9]{10}`)

	phone := stripCountryCodeExp.ReplaceAll([]byte(p), []byte(""))
	phone = replaceExp.ReplaceAll([]byte(phone), []byte(""))

	if !numExp.Match(phone) {
		return "", errors.New(fmt.Sprintf("The string is not a valid phone number: %v", p))
	}

	return string(phone), nil
}
