// in-memory implementation of the `pager` interface, for a database that is
// already loaded. Nothing else can change it, so there is nothing to lock.

package db

import (
	"io"
)

type memPager struct {
	b []byte
}

// pages start counting at 1
func (m *memPager) page(id int, pagesize int) ([]byte, error) {
	off := int64(id-1) * int64(pagesize)
	if id < 1 || pagesize < 0 || off+int64(pagesize) > int64(len(m.b)) {
		return nil, io.EOF
	}
	buf := make([]byte, pagesize)
	copy(buf, m.b[off:])
	return buf, nil
}

func (m *memPager) RLock() error { return nil }

func (m *memPager) RUnlock() error { return nil }

func (m *memPager) CheckReservedLock() (bool, error) { return false, nil }

func (m *memPager) Close() error { return nil }
