package udpwriter

import (
	"fmt"
	"io"
	"net"
	"time"
)

// asyncWriter is special component, that writes UDP messages
// line by line in async mode, using channel to sync
type asyncWriter struct {
	delivery chan []byte
	addr     *net.UDPAddr
	conn     *net.UDPConn

	delivered  uint64
	reconnects uint64
}

// NewS creates and returns new UDP connection
func NewS(addr string) (io.Writer, error) {
	netAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	return New(netAddr), nil
}

// New creates and returns new UDP connection
func New(addr *net.UDPAddr) io.Writer {
	r := &asyncWriter{
		addr:     addr,
		delivery: make(chan []byte),
	}

	// Initial connect attempt
	r.conn, _ = net.DialUDP("udp", nil, r.addr)

	// Goroutines
	go r.reconnect()
	go r.deliver()

	return r
}

// Write method is io.Writer implementation
// Queues delivery of provided bytes into goroutine
func (r *asyncWriter) Write(b []byte) (int, error) {
	l := len(b)
	if l > 0 {
		r.delivery <- b
	}

	return l, nil
}

// reconnect method will reconnect every 1 second
func (r *asyncWriter) reconnect() {
	for {
		if r.conn == nil {
			r.reconnects++
			r.conn, _ = net.DialUDP("udp", nil, r.addr)
		}

		time.Sleep(time.Second)
	}
}

// deliver method reads data from channel and delivers to connection
func (r *asyncWriter) deliver() {
	for bts := range r.delivery {
		if r.conn == nil {
			// Drop
			continue
		}

		r.delivered += uint64(len(bts))
		_, err := r.conn.Write(bts)
		if err != nil {
			r.conn = nil
		}
	}
}

// DeliveredBytes returns delivered bytes count
func (r *asyncWriter) DeliveredBytes() uint64 {
	return r.delivered
}

func (r *asyncWriter) String() string {
	return fmt.Sprintf(
		"UDP client connection to %s with %d bytes already delivered",
		r.addr.String(),
		r.delivered,
	)
}
