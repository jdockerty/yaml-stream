package yamlstream

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Stream represents a stream of YAML, delimited by `---`, although there
// is no requirement that this be the case, as a single Stream, i.e. a single
// file, is still valid.
type Stream struct {

	// A YAML stream is an array of bytes arrays, each element is separated
	// by the `---` delimiter.
	Stream [][]byte

	// Count is the number of streams that are present in the provided stream.
	// Note that this is not the number of delimiters present, but how many
	// YAML files there would be if the file were to be split separately.
	//
	// For example:
	//
	// hello: 'world'
	// ---
	// library: 'yaml-stream'
	//
	// This is a Count of 2.
	Count int
}

// Bytes returns the given stream as a single byte array, this is effectively
// similar to a call to `io.ReadFull` or `os.ReadFile`.
func (s *Stream) Bytes() []byte {
	b := &bytes.Buffer{}
	for _, data := range s.Stream {
		b.Write(data)
	}

	return b.Bytes()
}

// String returns the given stream in its string representated format.
// It is expected to be the file as it looks to the user.
func (s *Stream) String() string {
	return string(s.Bytes())
}

// Read takes the provided stream and segments it into multiple
// byte arrays for direct access, wrapped as a `Stream` type.
// This makes the assumption that the delimiter to the stream is
// `---\n`, i.e. three dashes and a newline character. It is the
// initialisation point for the Stream.
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

// New will return a default Stream with default values.
// There is no requirement to use this, as `&Stream{}` is
// effectively the same. It serves as the standard entrypoint
// to the library and creating the `Stream` type.
func New() *Stream {
	return &Stream{
		Stream: nil,
		Count:  0,
	}
}
