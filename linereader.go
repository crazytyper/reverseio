package reverseio

import "io"

const (
	defaultBufSize = 4096
)

type lineReader struct {
	r    io.Reader // a reverse reader
	rbuf []byte    // the read buffer
	rpos int       // the position within "rbuf"
	lbuf []byte    // the line buffer
	llen int       // the current length of the line
	lf   int       // the number of trailing EOL bytes at the end of the line (0,1 or 2)
	err  error     // the error
}

// ReadString reads the next line as a string.
//
// This method will consume as many bytes as the longest line in the source.
// If you need more control over memory consuption use ReadLine and deal
// with fragments yourself.
func (r *lineReader) ReadString() (line string, err error) {
	data, fragment, readErr := r.ReadLine()
	if readErr != nil {
		return "", readErr
	}
	if !fragment {
		return string(data), nil
	}
	// concatenate fragments
	buf := make([]byte, 0, len(data)*2)
	for err == nil {
		buf = prepend(buf, data)
		if !fragment {
			break
		}
		data, fragment, err = r.ReadLine()
	}
	return string(buf), err
}

// ReadLine reads the next line.
// If the line exceeds the buffer size fragment will be true
func (r *lineReader) ReadLine() (line []byte, fragment bool, err error) {
	if r.err != nil {
		return r.emitError()
	}
	for {
		c := r.next()
		if r.err != nil {
			if r.err == io.EOF {
				return r.emitLine()
			}
			return r.emitError()
		}
		if c == '\n' && r.lineLen() > 0 {
			r.backup()
			return r.emitLine()
		}
		if !r.appendLine(c) {
			r.backup()
			return r.emitLineFragement()
		}
	}
}

func (r *lineReader) next() (b byte) {
	if r.rpos == 0 {
		l, err := r.r.Read(r.rbuf)
		if err != nil {
			r.err = err
			return
		}
		r.rpos = l
		if r.rpos == 0 {
			r.err = io.EOF
			return
		}
	}
	r.rpos--
	b = r.rbuf[r.rpos]
	return
}

func (r *lineReader) backup() {
	r.rpos++
}

func (r *lineReader) appendLine(b byte) bool {
	if r.llen >= len(r.lbuf) {
		return false // line too long
	}
	if r.lf == 0 && b == '\n' {
		r.lf++ // do not append trailing LF
	} else if r.lf == 1 && b == '\r' {
		r.lf++ // do not append trailing CRLF
	} else {
		r.llen++
		r.lbuf[len(r.lbuf)-r.llen] = b
	}
	return true
}

func (r *lineReader) emitLine() (line []byte, fragment bool, err error) {
	line = r.line()
	r.resetLine()
	return
}

func (r *lineReader) emitLineFragement() (line []byte, fragment bool, err error) {
	line = r.line()
	fragment = true
	r.resetLine()
	return
}

func (r *lineReader) emitError() ([]byte, bool, error) {
	return nil, false, r.err
}

func (r *lineReader) line() []byte {
	return r.lbuf[len(r.lbuf)-r.llen:]
}

func (r *lineReader) resetLine() {
	r.llen = 0
	r.lf = 0
}

func (r *lineReader) lineLen() int {
	return r.llen + r.lf
}

// LineReader reads lines.
type LineReader interface {
	// ReadString reads the next line as a string.
	ReadString() (line string, err error)

	// ReadLine reads the next line.
	// If the line exceeds the buffer size fragment will be true
	ReadLine() (line []byte, fragment bool, err error)
}

// NewLineReader creates a new line reader that reads
// lines beginning at the end of the reader.
//
// Lines exceeding 4096 bytes will be returned as multiple fragments.
func NewLineReader(r io.ReadSeeker) LineReader {
	return NewLineReaderSize(r, defaultBufSize)
}

// NewLineReaderSize creates a new line reader that reads
// lines beginning at the end of the reader.
//
// Lines exceeding size bytes will be returned as multiple fragments.
func NewLineReaderSize(r io.ReadSeeker, size int) LineReader {
	return &lineReader{
		r:    NewReader(r),
		rbuf: make([]byte, size),
		lbuf: make([]byte, size),
	}
}
