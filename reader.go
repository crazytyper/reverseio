package reverseio

import (
	"io"
	"os"
)

type reader struct {
	r    io.ReadSeeker
	rpos int64
}

// Read reads the next chunk of bytes.
func (r *reader) Read(buf []byte) (l int, err error) {
	if r.rpos == 0 {
		return 0, io.EOF
	}

	l = len(buf)
	if r.rpos < 0 {
		r.rpos, err = r.r.Seek(0, os.SEEK_END)
	}

	if err == nil {
		rpos := r.rpos - int64(l)
		if rpos < 0 {
			rpos = 0
			l = int(r.rpos)
		}
		r.rpos, err = r.r.Seek(rpos, os.SEEK_SET)
	}

	if err == nil {
		_, err = io.ReadAtLeast(r.r, buf, l)
	}
	return
}

// NewReader creates a new reader that reads
// the underlying reader from end to start.
//
// Reading 2 byte chunks from "ABCD" will first give "CD" and then "AB".
func NewReader(r io.ReadSeeker) io.Reader {
	return &reader{r, -1}
}
