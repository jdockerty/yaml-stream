package yamlstream

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"gopkg.in/yaml.v3"
)

var (
	yamlDelimiter = []byte("---\n")
)

// Document is a single YAML document.
type Document map[string]interface{}

// Bytes will return the byte array representation of the YAML document.
// Note that this will include the YAML delimiter of '---' even if it was not
// originally included in order to allow the separation of multiple documents.
func (d *Document) Bytes() []byte {

	b := &bytes.Buffer{}

	data, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}
	// The YAML library does not include a delimiter to separate a document
	// so it is added in by the library to ensure separation between other
	// documents that may be present.
	b.Write(yamlDelimiter)
	b.Write(data)

	return b.Bytes()
}

// String returns the string representation of the YAML document.
// This is usually expected to be how the file looks to the human eye, unless
// there are no delimiters in place, in which case they will be added by the
// library.
func (d *Document) String() string {
	return string(d.Bytes())
}

// Unmarshal will unmarshal the document into the type provided in the out value.
func (d *Document) Unmarshal(out interface{}) error {
	err := yaml.Unmarshal(d.Bytes(), out)
	if err != nil {
		return err
	}
	return nil
}

// Stream represents a stream of YAML, delimited by `---`, although there
// is no requirement that this be the case, as a single Stream, i.e. a single
// file, is still valid.
type Stream struct {

	// Stream represents a stream of YAML. It is composed of zero or more YAML
	// documents. Conceptually, this is similar to an array of YAML documents,
	// able to be accessed at a particular index.
	Stream []Document

	// Count is the number of YAML documents that are present in the provided
	// stream.
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

// Get will retrieve a YAML document at a provided index.
func (s *Stream) Get(index int) Document {
	return s.Stream[index]
}

// GetUnmarshal will unmarshal the Document at the provided index into the type
// that is provided by the 'out' value.
func (s *Stream) GetUnmarshal(index int, out interface{}) error {

	if reflect.ValueOf(out).Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer (&) for 'to' value that was passed")
	}

	doc := s.Stream[index]
	err := doc.Unmarshal(out)
	if err != nil {
		return err
	}

	return nil
}

// Bytes returns the given stream as a single byte array, this is effectively
// similar to a call to `io.ReadFull` or `os.ReadFile`.
// Note that because a Document is separated by a delimiter, if they are not
// included in the file provided to `Read`, they will be included when being
// called into a byte array.
func (s *Stream) Bytes() []byte {
	b := &bytes.Buffer{}
	for _, data := range s.Stream {
		b.Write(data.Bytes())
	}

	return b.Bytes()
}

// String returns the given stream in its string representated format.
// It is expected to be the file as it looks to the user.
func (s *Stream) String() string {
	return string(s.Bytes())
}

// Read takes a reader, such as a YAML file, that contains zero or more documents
// and collects it into a singlular Stream.
func (s *Stream) Read(r io.Reader) error {

	stream := make([]Document, 0)

	d := yaml.NewDecoder(r)
	for {

		var document Document

		err := d.Decode(&document)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		stream = append(stream, document)
	}

	s.Stream = stream
	s.Count = len(stream)
	return nil
}

// New will return a default Stream with default values.
// There is no requirement to use this, as `&Stream{}` is
// effectively the same. It serves as the standard entrypoint
// to the library and creating the `Stream` type.
func New() *Stream {
	return &Stream{}
}
