package random

import "github.com/sethvargo/go-password/password"

type Random struct{}

func New() *Random {
	return &Random{}
}

func (r *Random) GenerateRandomPassword(length, lenDigits, lenSymbol int) (string, error) {
	return password.Generate(length, lenDigits, lenSymbol, false, false)
}
