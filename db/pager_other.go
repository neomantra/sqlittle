//go:build js || wasip1
// +build js wasip1

// platforms without files to map and lock: only OpenBytes() works here

package db

import (
	"errors"
)

func newFilePager(file string) (pager, error) {
	return nil, errors.New("opening a file is not supported on this platform; use OpenBytes")
}
