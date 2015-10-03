package distilgo

import (
	"bufio"
	"os"
	"sort"
	"testing"
)

func testLoadDB() (*os.File, DistilRecords, error) {
	f, err := os.Open("testdata/db")
	if err != nil {
		return nil, nil, err
	}

	l, e := LoadIPs(bufio.NewReader(f))
	return f, l, e
}

func TestLoad(t *testing.T) {
	f, _, err := testLoadDB()
	if err != nil {
		t.Fatal(err)
	}

	f.Close()
}

func TestValidRecords(t *testing.T) {
	f, l, _ := testLoadDB()

	if !sort.IsSorted(l) {
		t.Error("IP List was not sorted")
	}

	if len(l) == 0 {
		t.Error("Read 0 records from test data.")
	}

	for i := range l {
		r := l[i]

		if _, err := r.IP.MarshalText(); err != nil {
			t.Error("Invalid IP address read")
		}

		if r.LastDetected.IsZero() {
			t.Error("Invalid date read")
		}
	}

	f.Close()
}
