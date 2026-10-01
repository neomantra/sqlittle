package db

import (
	"encoding/binary"
	"testing"
)

func TestOverflowCorrupt(t *testing.T) {
	// two pages; the second links to itself as its own overflow page
	pagesize := 512
	b := make([]byte, 2*pagesize)
	binary.BigEndian.PutUint32(b[pagesize:], 2)
	db := &Database{l: &memPager{b: b}, header: &header{PageSize: pagesize}}

	// a payload which claims more than there is
	if _, err := addOverflow(db, cellPayload{Length: 100, Payload: []byte("short")}); err != ErrCorrupted {
		t.Errorf("have %v, want %v", err, ErrCorrupted)
	}
	if _, err := addOverflow(db, cellPayload{Length: -1, Payload: []byte("short")}); err != ErrCorrupted {
		t.Errorf("have %v, want %v", err, ErrCorrupted)
	}
	// overflow pages in a loop
	if _, err := addOverflow(db, cellPayload{Length: 1 << 40, Payload: []byte("short"), Overflow: 2}); err != ErrCorrupted {
		t.Errorf("have %v, want %v", err, ErrCorrupted)
	}
	// the payload is a copy: pages are cached
	page := []byte("payload, and the rest of the page")
	full, err := addOverflow(db, cellPayload{Length: 7, Payload: page[:7]})
	if err != nil {
		t.Fatal(err)
	}
	full = append(full, " is changed"...)
	if have, want := string(page), "payload, and the rest of the page"; have != want {
		t.Errorf("have %q, want %q", have, want)
	}
}
