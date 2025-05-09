package service

import (
	"math/rand"
	"strconv"
)

// createAccountNumber формирует случайный 20-ти значный номер счёта.
func createAccountNumber() string {
	return strconv.FormatUint(uint64(rand.Intn(9999999999-1000000000)+1000000000)<<32+uint64(rand.Intn(9999999999-1000000000)+1000000000), 10)
}
