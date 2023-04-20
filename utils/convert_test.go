package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvert_ToDecimals(t *testing.T) {
	amountStr := "10000000000000000000000000"
	kai := ToDecimal(amountStr, 18).String()
	assert.Equal(t, "10000000", kai)
}
