package whale

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"testing"
)

func TestName(t *testing.T) {
	data := humanize.FormatFloat("#,###.##", 1000000.5222)
	fmt.Println(data)
}
