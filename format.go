package distilgo

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"time"
)

type DistilRecord struct {
	IP                net.IP
	LastDetected      time.Time
	PercentBotTraffic float64

	// internal; for sorting
	ip32 uint32
}

var distilDelim = []byte{'\x01'}

// LoadIPs returns a channel that you can receive DistilRecords on.
// It will not return any errors at this time.
func LoadIPs(r *bufio.Reader) chan DistilRecord {
	var (
		ln  []byte
		err error
		ch  = make(chan DistilRecord)
	)

	go func() {
		for {
			ln, err = r.ReadBytes('\n')
			if (err != nil && err != io.EOF) || len(ln) == 0 {
				break
			}

			// 0.0.0.0\x012015-10-01 16:18:52.667\x010.3108108108108108\n
			// we have a line, and we should expect 2 \x01s according to the documentation.
			i := bytes.Index(ln, distilDelim)
			i2 := i + 1 + bytes.Index(ln[i+1:], distilDelim)
			if i == -1 || i2 == -1 || i > 15 {
				continue
			}

			ip := net.ParseIP(string(ln[0:i]))
			ip4 := ip.To4()

			date := ln[i+1 : i2]
			t, err := time.Parse("2006-01-02 15:04:05", string(date))
			if err != nil {
				break
			}

			pct, err := strconv.ParseFloat(string(ln[i2+1:len(ln)-1]), 64)
			if err != nil {
				break
			}

			ch <- DistilRecord{
				IP:                ip,
				LastDetected:      t,
				PercentBotTraffic: pct * 100,
				ip32:              (uint32(ip4[0]) * 16777216) + (uint32(ip4[1]) * 65536) + (uint32(ip4[2]) * 256) + uint32(ip4[3]),
			}
		}

		close(ch)
	}()

	return ch
}
