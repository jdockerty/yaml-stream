package yamlstream_test

import (
	"bytes"
	"os"
	"testing"

	yamlstream "github.com/jdockerty/yaml-stream"
	"github.com/stretchr/testify/assert"
)

const (
	simpleYamlStream = "testdata/simple_stream.yaml"
)

func TestNewWithDefaults(t *testing.T) {

	ys := yamlstream.New()

	assert.IsType(t, ys, &yamlstream.Stream{})
}

func TestStreamCounter(t *testing.T) {
	ys := yamlstream.New()

	f, _ := os.Open(simpleYamlStream)
	defer f.Close()

	err := ys.Read(f)
	assert.Nil(t, err)

	assert.Equal(t, 3, ys.Count)
}

func TestBytes(t *testing.T) {

	ys := yamlstream.New()

	f, _ := os.Open(simpleYamlStream)
	defer f.Close()

	_ = ys.Read(f)

	assert.IsType(t, ys.Bytes(), make([]byte, 0))
}

func TestReadBytesEquality(t *testing.T) {

	ys := yamlstream.New()

	f, _ := os.Open(simpleYamlStream)
	defer f.Close()

	err := ys.Read(f)
	assert.Nil(t, err)

	fileAsBytes, _ := os.ReadFile(simpleYamlStream)

    t.Logf("AsBytes: %s\n", ys.Bytes())
    t.Logf("\nGot:%s\nWant:%s", ys.String(), string(fileAsBytes))
	assert.True(t, bytes.Equal(ys.Bytes(), fileAsBytes))

}
