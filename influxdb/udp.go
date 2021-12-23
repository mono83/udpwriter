package influxdb

import (
	"github.com/mono83/udpwriter"
	"net"
)

// NewUDPWriterS constructs new metric writer that can write metrics
// to InfluxDB via UDP
func NewUDPWriterS(addr string, maxUDPSize int) (MetricWriter, error) {
	writer, err := udpwriter.NewS(addr)
	if err != nil {
		return nil, err
	}
	return NewWriter(writer, maxUDPSize), nil
}

// NewUDPWriter constructs new metric writer that can write metrics
// to InfluxDB via UDP
func NewUDPWriter(addr *net.UDPAddr, flushSize int) MetricWriter {
	writer := udpwriter.New(addr)
	return NewWriter(writer, flushSize)
}
