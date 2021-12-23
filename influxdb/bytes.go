package influxdb

import (
	"bytes"
	"io"
	"strconv"
)

// NewWriter construct metric writer that writes metrics to
// arbitrary byte writer
func NewWriter(w io.Writer, flushSize int) MetricWriter {
	return influxDBUDPWriter{Writer: w, flushSize: flushSize}
}

type influxDBUDPWriter struct {
	io.Writer
	flushSize int
}

func (i influxDBUDPWriter) Write(data interface {
	ForEachMetric(func(name string, value int64, tags map[string]string))
}) error {
	if data == nil {
		return nil
	}

	var outerError error
	buf := bytes.NewBuffer(nil)
	data.ForEachMetric(func(name string, value int64, tags map[string]string) {
		if outerError != nil {
			return
		}

		// Writing into line buffer
		line := bytes.NewBuffer(nil)
		line.Write(Sanitize(name)) // Metric name
		for k, v := range tags {   // Metric tags
			line.WriteRune(',')
			line.Write(Sanitize(k))
			line.WriteRune('=')
			line.Write(Sanitize(v))
		}

		line.WriteString(" value=")
		line.WriteString(strconv.FormatInt(value, 10))
		line.WriteRune('\n')

		// Checking buffer overflow
		if line.Len()+buf.Len() > i.flushSize {
			// Buffer overflow detected
			if buf.Len() == 0 {
				_, outerError = i.Writer.Write(line.Bytes())
			} else {
				_, outerError = i.Writer.Write(buf.Bytes())
				buf.Reset()
				buf.Write(line.Bytes())
			}
		} else {
			buf.Write(line.Bytes())
		}
	})

	// Writing remainder
	if outerError == nil && buf.Len() > 0 {
		_, outerError = i.Writer.Write(buf.Bytes())
	}

	return outerError
}
