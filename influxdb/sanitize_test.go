package influxdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sanitizeSData = []struct {
	Expected, Given string
}{
	{"foo", "foo"},
	{"Foo", "Foo"},
	{"Foo", "\tFoo   \n"},
	{"F__o_o", "F  o\to"},
	{"F.o.o", "F.o.o"},
}

func TestSanitizeS(t *testing.T) {
	for _, datum := range sanitizeSData {
		t.Run(fmt.Sprint(datum), func(t *testing.T) {
			assert.Equal(t, datum.Expected, SanitizeS(datum.Given))
		})
	}
}
