package service

import (
	"math/rand"
	"strconv"

	"github.com/phedde/luhn-algorithm"
)

// createCardNumber формирует случайный 16-ти значный валидный номер карты.
func createCardNumber() string {
	cardNumber := rand.Int63n(int64(999999999999999-100000000000000)) + 100000000000000
	validCreditCard, _ := luhn.FullNumber(cardNumber)

	return strconv.FormatInt(validCreditCard, 10)
}

// cardNumberIsValid проверяет номер карты на валидность.
func cardNumberIsValid(number string) bool {
	cardNumber, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return false
	}

	return luhn.IsValid(cardNumber)
}
