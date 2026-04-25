package money

import (
	"github.com/shopspring/decimal"
)

// Default precision for Banking
const (
	DefaultScale = 2
)

// ToDecimal converts string to decimal. Returns 0 if invalid.
func ToDecimal(amount string) decimal.Decimal {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Zero
	}
	return d.Round(DefaultScale)
}

// ToString converts decimal to string with scale 2.
func ToString(amount decimal.Decimal) string {
	return amount.StringFixed(DefaultScale)
}

// Add adds two decimals and rounds the result.
func Add(a, b decimal.Decimal) decimal.Decimal {
	return a.Add(b).Round(DefaultScale)
}

// Subtract subtracts b from a and rounds the result.
func Subtract(a, b decimal.Decimal) decimal.Decimal {
	return a.Sub(b).Round(DefaultScale)
}
