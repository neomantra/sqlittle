package db

import (
	"encoding/binary"
)

// cellPayload represents the payload part of a cell. If overflow is non-zero the
// cellPayload field will be truncated. Use addOverflow() to get a full payload.
type cellPayload struct {
	Length   int64
	Payload  []byte
	Overflow int
}

// overflow is stored on different pages. Load whatever is needed to complete
// the payload data.
// The result never shares memory with a page: pages are cached.
func addOverflow(db *Database, pl cellPayload) ([]byte, error) {
	if pl.Length < 0 {
		return nil, ErrCorrupted
	}
	to := append([]byte(nil), pl.Payload...)
	overflow := pl.Overflow
	seen := map[int]struct{}{}
	for {
		if overflow == 0 {
			if int64(len(to)) < pl.Length {
				return nil, ErrCorrupted
			}
			return to[:pl.Length], nil
		}
		if _, ok := seen[overflow]; ok || int64(len(to)) >= pl.Length {
			// a loop, or more pages than the payload needs
			return nil, ErrCorrupted
		}
		seen[overflow] = struct{}{}
		buf, err := db.page(overflow)
		if err != nil {
			return nil, err
		}
		if len(buf) < 4 {
			return nil, ErrCorrupted
		}
		next, buf := int(binary.BigEndian.Uint32(buf[:4])), buf[4:]
		if used := db.header.PageSize - db.header.Reserved - 4; len(buf) > used {
			// don't read the reserved space at the end of the page
			buf = buf[:used]
		}
		to = append(to, buf...)
		overflow = next
	}
}
