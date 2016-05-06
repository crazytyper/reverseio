package reverseio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type readLines []string

func (r *readLines) Record(line []byte, fragment bool, err error) {
	var s string
	if line == nil {
		s = fmt.Sprintf("(null,%t,%v)", fragment, err)
	} else {
		s = fmt.Sprintf("(%q,%t,%v)", string(line), fragment, err)
	}
	*r = append(*r, s)
}

type readStrings []string

func (r *readStrings) Record(line string, err error) {
	*r = append(*r, fmt.Sprintf("(%q,%v)", string(line), err))
}

func expectedReadLines(s string, size int) (r readLines) {
	lines := strings.Split(s, "\n")
	if strings.HasSuffix(s, "\n") {
		lines = lines[:len(lines)-1]
	}
	r = readLines{}
	l := len(lines) - 1
	for i := l; i >= 0; i-- {
		line := []byte(strings.TrimSuffix(lines[i], "\r"))
		// split the line into chunks that do not exceed the buffer size
		for len(line) > size {
			e := len(line) - size
			r.Record(line[e:], true, nil)
			line = line[:e]
		}
		r.Record(line, false, nil)
	}
	r.Record(nil, false, io.EOF)
	return
}

func expectedReadStrings(s string, size int) (r readStrings) {
	lines := strings.Split(s, "\n")
	if strings.HasSuffix(s, "\n") {
		lines = lines[:len(lines)-1]
	}
	r = readStrings{}
	l := len(lines) - 1
	for i := l; i >= 0; i-- {
		line := strings.TrimSuffix(lines[i], "\r")
		r.Record(line, nil)
	}
	r.Record("", io.EOF)
	return r
}

func actualReadLines(s string, size int) readLines {
	r := NewLineReaderSize(strings.NewReader(s), size)
	actual := readLines{}
	for {
		line, fragment, err := r.ReadLine()
		actual.Record(line, fragment, err)
		if err != nil {
			break
		}
	}
	return actual
}

func actualReadStrings(s string, size int) readStrings {
	r := NewLineReaderSize(strings.NewReader(s), size)
	actual := readStrings{}
	for {
		line, err := r.ReadString()
		actual.Record(line, err)
		if err != nil {
			break
		}
	}
	return actual
}

var (
	testCases = []string{
		"",
		"\n",
		"\n\n",
		"\n\n\n",
		"\r\n",
		"\r\n\r\n",
		"\r\n\r\n\r\n",
		"Hello\nworld!",
		"Hello world!",
		"Hello\r\nworld!",
		"Hello\nworld!\n",
		"Hello\n\nworld!\n",
		"Hello\r\nworld!\r\n",
		"Hello\r\n\r\nworld!\r\n",
	}
)

func TestReadLine(t *testing.T) {
	for bufSize := 1; bufSize <= 4096; bufSize++ {
		for _, testCase := range testCases {
			assert.Equal(t,
				expectedReadLines(testCase, bufSize),
				actualReadLines(testCase, bufSize),
				"test case %q (buf size %d)", testCase, bufSize)
		}
	}
}

func TestReadString(t *testing.T) {
	for bufSize := 1; bufSize <= 4096; bufSize++ {
		for _, testCase := range testCases {
			assert.Equal(t,
				expectedReadStrings(testCase, bufSize),
				actualReadStrings(testCase, bufSize),
				"test case %q (buf size %d)", testCase, bufSize)
		}
	}
}

type mockReaderSeeker struct {
	mock.Mock
}

func (m *mockReaderSeeker) Read(data []byte) (int, error) {
	ret := m.Called(data)
	return ret.Int(0), ret.Error(1)
}

func (m *mockReaderSeeker) Seek(off int64, whence int) (int64, error) {
	ret := m.Called(off, whence)
	return ret.Get(0).(int64), ret.Error(1)
}

func TestReadLineError(t *testing.T) {
	expectedErr := errors.New("Deliberate Error")

	mockedReader := &mockReaderSeeker{}
	mockedReader.On("Seek", int64(0), os.SEEK_END).Return(int64(42), nil)
	mockedReader.On("Seek", int64(0), os.SEEK_SET).Return(int64(42), nil)
	mockedReader.On("Read", mock.Anything).Return(0, expectedErr)

	r := NewLineReader(mockedReader)
	line, fragment, err := r.ReadLine()
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, line)
	assert.False(t, fragment)
}
