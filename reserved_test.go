package sqlittle

import (
	"fmt"
	"strings"
	"testing"
)

// testdata/reserved.sqlite has reserved bytes at the end of every page, and
// rows (and index entries) which need overflow pages.
func TestReservedSpace(t *testing.T) {
	db, err := Open("./testdata/reserved.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	n := 0
	err = db.Select("big", func(r Row) {
		var (
			id int
			s  string
		)
		if err := r.Scan(&id, &s); err != nil {
			t.Fatal(err)
		}
		n++
		if id != n {
			t.Fatalf("id: have %d, want %d", id, n)
		}
		want := fmt.Sprintf("%05d", id) + strings.Repeat("ab", 6000+id*200)
		if s != want {
			t.Fatalf("row %d: payload mismatch (len %d, want %d)", id, len(s), len(want))
		}
	}, "id", "s")
	if err != nil {
		t.Fatal(err)
	}
	if n != 12 {
		t.Fatalf("have %d rows, want 12", n)
	}

	n = 0
	err = db.IndexedSelect("big", "big_s", func(r Row) {
		n++
		var s string
		if err := r.Scan(&s); err != nil {
			t.Fatal(err)
		}
		if want := fmt.Sprintf("%05d", n); !strings.HasPrefix(s, want) {
			t.Fatalf("index row %d: have prefix %q", n, s[:5])
		}
	}, "s")
	if err != nil {
		t.Fatal(err)
	}
	if n != 12 {
		t.Fatalf("have %d index rows, want 12", n)
	}
}
