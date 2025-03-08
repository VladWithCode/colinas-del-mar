package internal

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/leekchan/accounting"
)

func FormatPhone(p string) (string, error) {
	stripCountryCodeExp := regexp.MustCompile(`^\+52`)
	replaceExp := regexp.MustCompile(`[ -]`)
	numExp := regexp.MustCompile(`[0-9]{10}`)

	phone := stripCountryCodeExp.ReplaceAll([]byte(p), []byte(""))
	phone = replaceExp.ReplaceAll([]byte(phone), []byte(""))

	if !numExp.Match(phone) {
		return "", fmt.Errorf("the string is not a valid phone number: %v", p)
	}

	return string(phone), nil
}

func FormatMoney(amount float64) string {
	bFloat := big.NewFloat(amount)
	ac := accounting.Accounting{Symbol: "$", Precision: 2}

	return ac.FormatMoneyBigFloat(bFloat)
}
