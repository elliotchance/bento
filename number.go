package main

import (
	"math/big"
	"strings"
)

const (
	DefaultNumericPrecision = 6

	// TODO: This is a placeholder for when it can be handled more gracefully.
	//  We need this for constants so that we don't lose precision.
	UnlimitedPrecision = 1000
)

type Number struct {
	Rat       *big.Rat
	Precision int
}

func NewNumber(s string, precision int) *Number {
	// TODO: The number of decimal places cannot be negative.

	rat, _ := big.NewRat(0, 1).SetString(s)

	return &Number{
		Rat:       rat,
		Precision: precision,
	}
}

func (number *Number) String() string {
	s := number.Rat.FloatString(number.Precision)

	// Remove any trailing zeros after the decimal point.
	if number.Precision > 0 {
		s = strings.TrimRight(s, "0")
	}

	// For integers we also want to remove the ".".
	return strings.TrimRight(s, ".")
}

func (number *Number) Cmp(number2 *Number) int {
	return number.Rat.Cmp(number2.Rat)
}

func (number *Number) Add(a, b *Number) {
	number.Rat = big.NewRat(0, 1).Add(a.Rat, b.Rat)
}

func (number *Number) Sub(a, b *Number) {
	number.Rat = big.NewRat(0, 1).Sub(a.Rat, b.Rat)
}

func (number *Number) Mul(a, b *Number) {
	result := big.NewRat(0, 1).Mul(a.Rat, b.Rat)

	// Multiplying two rational numbers may give a decimal that is higher
	// precision than we allow, so it has to be rounded before it's stored.
	number.Rat.SetString(result.FloatString(number.Precision))
}

func (number *Number) Quo(a, b *Number) {
	result := big.NewRat(0, 1).Quo(a.Rat, b.Rat)

	// TODO: What about divide-by-zero?

	// Dividing two rational numbers may give a decimal that is higher precision
	// than we allow, so it has to be rounded before it's stored.
	number.Rat.SetString(result.FloatString(number.Precision))
}

func (number *Number) Set(x *Number) {
	number.Rat.SetString(x.Rat.FloatString(number.Precision))
}

func (number *Number) Bool() bool {
	return number.Rat.Sign() != 0
}
