package yamlstream

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Stream struct {
	Stream [][]byte
	Count  int
}

func (s *Stream) Bytes() []byte {
	b := &bytes.Buffer{}
	for _, data := range s.Stream {
		b.Write(data)
	}

	return b.Bytes()
}

func (s *Stream) String() string {
	return string(s.Bytes())
}

func (s *Stream) Read(r io.Reader) error {
	rd := bufio.NewReader(r)

	// Initialise count as 1 when passed a Reader, this covers
	// the case where there are no delimiters since this is
	// considered a "stream", i.e. a single file.
	s.Count = 1

	b := &bytes.Buffer{}
	for {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// Append the current buffer as we won't
				// hit another delimiter to add it with the EOF.
				s.Stream = append(s.Stream, b.Bytes())
				break
			}
			return fmt.Errorf("unable to read: %w", err)
		}

		if bytes.Equal(line, []byte("---\n")) {
			s.Count += 1
			s.Stream = append(s.Stream, b.Bytes())

			// Clear buffer for next delimiter
			b = &bytes.Buffer{}
		}

		_, err = b.Write(line)
		if err != nil {
			return fmt.Errorf("unable to write data: %w", err)
		}

	}

	return nil

}

func New() *Stream {
	return &Stream{
		Stream: nil,
		Count:  0,
	}
}
