package distilgo

import (
	"bufio"
	"bytes"
	"net"
	"sort"
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

type DistilRecords []DistilRecord

func (r DistilRecords) Len() int           { return len(r) }
func (r DistilRecords) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r DistilRecords) Less(i, j int) bool { return r[i].ip32 < r[j].ip32 }

var distilDelim = []byte{'\x01'}

func LoadIPs(r *bufio.Reader) (DistilRecords, error) {
	var (
		ln   []byte
		err  error
		recs DistilRecords
	)

	for {
		ln, err = r.ReadBytes('\n')
		if err != nil {
			return nil, err
		} else if len(ln) == 0 {
			break
		}

		// 0.0.0.0\x012015-10-01 16:18:52.667\x010.3108108108108108\n
		// we have a line, and we should expect 2 \x01s according to the documentation.
		i := bytes.Index(ln, distilDelim)
		i2 := i + 1 + bytes.Index(ln[i+1:], distilDelim)
		if i == -1 || i2 == -1 || i >= 15 {
			break
		}

		ip := net.ParseIP(string(ln[0:i]))
		ip4 := ip.To4()

		date := ln[i+1 : i2]
		t, err := time.Parse("2006-01-02 15:04:05", string(date))
		if err != nil {
			return nil, err
		}

		pct, err := strconv.ParseFloat(string(ln[i2+1:len(ln)-1]), 64)
		if err != nil {
			return nil, err
		}

		recs = append(recs, DistilRecord{
			IP:                ip,
			LastDetected:      t,
			PercentBotTraffic: pct * 100,
			ip32:              (uint32(ip4[0]) * 16777216) + (uint32(ip4[1]) * 65536) + (uint32(ip4[2]) * 256) + uint32(ip4[3]),
		})
	}

	sort.Sort(recs)

	return recs, nil
}
