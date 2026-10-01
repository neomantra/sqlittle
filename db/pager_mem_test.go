package db

import (
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

// every test database must read the same from memory as from its file
func TestOpenBytes(t *testing.T) {
	files, err := filepath.Glob("./../testdata/*.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test databases")
	}
	read := func(db *Database) (map[string][]Record, error) {
		tables, err := db.Tables()
		if err != nil {
			return nil, err
		}
		rows := map[string][]Record{}
		for _, name := range tables {
			if db.withoutRowid(name) {
				continue
			}
			table, err := db.Table(name)
			if err != nil {
				return nil, err
			}
			rows[name] = []Record{}
			if err := table.Scan(func(rowid int64, rec Record) bool {
				rows[name] = append(rows[name], append(Record{rowid}, rec...))
				return false
			}); err != nil {
				return nil, err
			}
		}
		return rows, nil
	}
	opened := 0
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		fromFile, fileErr := OpenFile(f)
		if fileErr == ErrHotJournal || len(b) == 0 {
			// a journal is a matter between files, and an empty file
			// fails in mmap
			continue
		}
		fromBytes, bytesErr := OpenBytes(b)
		if have, want := bytesErr, fileErr; have != want {
			t.Errorf("%s: have %v, want %v", f, have, want)
		}
		if fileErr != nil || bytesErr != nil {
			continue
		}
		want, wantErr := read(fromFile)
		have, haveErr := read(fromBytes)
		if !reflect.DeepEqual(haveErr, wantErr) {
			t.Errorf("%s: have %v, want %v", f, haveErr, wantErr)
		}
		if !reflect.DeepEqual(have, want) {
			t.Errorf("%s: tables read from memory differ from those read from the file", f)
		}
		fromFile.Close()
		if err := fromBytes.Close(); err != nil {
			t.Errorf("%s: %v", f, err)
		}
		opened++
	}
	if opened < 10 {
		t.Errorf("only %d databases were compared", opened)
	}
}

func TestOpenBytesShort(t *testing.T) {
	for _, b := range [][]byte{nil, {}, []byte("SQLite format 3\x00")} {
		if _, err := OpenBytes(b); err != io.EOF {
			t.Errorf("have %v, want %v", err, io.EOF)
		}
	}
}

func TestOpenBytesPastTheEnd(t *testing.T) {
	b, err := ioutil.ReadFile("./../testdata/single.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	db, err := OpenBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, id := range []int{0, -1, len(b)/db.header.PageSize + 1} {
		if _, err := db.l.page(id, db.header.PageSize); err != io.EOF {
			t.Errorf("page %d: have %v, want %v", id, err, io.EOF)
		}
	}
}
