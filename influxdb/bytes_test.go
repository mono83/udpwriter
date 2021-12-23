package influxdb

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWriter(t *testing.T) {
	b := bytes.NewBuffer(nil)

	w := NewWriter(b, 1024)
	if assert.NoError(t, w.Write(stubData(true))) {
		assert.Equal(
			t,
			"foo.bar value=100500\nbar_kek,some_spaces=even_here value=300\n",
			b.String(),
		)
	}
}

type stubData bool

func (stubData) ForEachMetric(f func(name string, value int64, tags map[string]string)) {
	f("foo.bar", 100500, nil)
	f("bar kek", 300, map[string]string{"some spaces": " even\there\n"})
}
