// Package utils
package utils

import (
	"math/big"

	"github.com/dustin/go-humanize"
	"github.com/shopspring/decimal"
)

func Uint8ToString(input [32]uint8) string {
	var name []byte
	for _, b := range input {
		if b != 0 {
			name = append(name, b)
		}
	}

	return string(name)
}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func HumanizeCurrency(value string) string {
	valueF64 := ToDecimal(value, 18).InexactFloat64()
	humanizeValue := humanize.FormatFloat("#,###.##", valueF64)
	return humanizeValue
}
