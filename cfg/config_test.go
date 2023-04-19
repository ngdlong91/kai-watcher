package cfg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestParseInt64(t *testing.T) {
	data, err := strconv.ParseInt("-873461799", 10, 64)
	assert.Nil(t, err)
	fmt.Println(data)
}
