package distilgo

import (
	"bufio"
	"os"
	"testing"
)

func testLoadDB() (*os.File, chan DistilRecord) {
	f, err := os.Open("testdata/db")
	if err != nil {
		return nil, nil
	}

	ch := LoadIPs(bufio.NewReader(f))
	return f, ch
}

func TestValidRecords(t *testing.T) {
	f, l := testLoadDB()

	c := 0
	for {
		r, more := <-l
		if !more {
			break
		}

		c++

		if _, err := r.IP.MarshalText(); err != nil {
			t.Error("Invalid IP address read")
		}

		if r.LastDetected.IsZero() {
			t.Error("Invalid date read")
		}
	}

	if c == 0 {
		t.Error("Failed to read any records from test data")
	}

	f.Close()
}
