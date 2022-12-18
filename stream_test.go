package yamlstream_test

import (
	"os"
	"testing"

	yamlstream "github.com/jdockerty/yaml-stream"
	"github.com/stretchr/testify/assert"
)

const (
	simpleYAML            = "testdata/simple.yaml"
	simpleYAMLStream      = "testdata/simple_stream.yaml"
	simpleYAMLStreamOne   = "testdata/simple_stream_one.yaml"
	simpleYAMLStreamTwo   = "testdata/simple_stream_two.yaml"
	simpleYAMLStreamThree = "testdata/simple_stream_three.yaml"
)

func TestNewWithDefaults(t *testing.T) {

	ys := yamlstream.New()

	assert.IsType(t, ys, &yamlstream.Stream{})
	assert.Equal(t, ys, &yamlstream.Stream{Count: 0, Stream: nil})
}

func TestStreamCount(t *testing.T) {

	tests := []struct {
		Name     string
		Filename string
		Expected int
	}{
		{
			Name:     "is correct for simple yaml file",
			Filename: simpleYAML,
			Expected: 1,
		},
		{
			Name:     "is correct for simple yaml stream file",
			Filename: simpleYAMLStream,
			Expected: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			ys := yamlstream.New()

			f, _ := os.Open(tc.Filename)
			defer f.Close()

			err := ys.Read(f)
			assert.Nil(t, err)

			assert.Equal(t, tc.Expected, ys.Count)
		})
	}
}

func TestBytes(t *testing.T) {

	ys := yamlstream.New()

	f, _ := os.Open(simpleYAMLStream)
	defer f.Close()

	_ = ys.Read(f)

	assert.IsType(t, ys.Bytes(), make([]byte, 0))
}

func TestReadEquality(t *testing.T) {

	tests := []struct {
		Name     string
		Filename string
	}{
		{
			Name:     "is equal for yaml stream",
			Filename: simpleYAMLStream,
		},
		{
			Name:     "is equal for regular yaml file",
			Filename: simpleYAML,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			ys := yamlstream.New()

			f, _ := os.Open(tc.Filename)
			defer f.Close()

			err := ys.Read(f)
			assert.Nil(t, err)

			fileAsBytes, _ := os.ReadFile(tc.Filename)

			assert.Equal(t, fileAsBytes, ys.Bytes())

		})
	}

}

func TestStreamGet(t *testing.T) {

	tests := []struct {
		Name               string
		YAMLStreamFilename string
		SegmentedFilename  string
		Index              int
	}{
		{
			Name:               "is correct for 0th index",
			YAMLStreamFilename: simpleYAMLStream,
			SegmentedFilename:  simpleYAMLStreamOne,
			Index:              0,
		},
		{
			Name:               "is correct for 1st index",
			YAMLStreamFilename: simpleYAMLStream,
			SegmentedFilename:  simpleYAMLStreamTwo,
			Index:              1,
		},
		{
			Name:               "is correct for 2nd index",
			YAMLStreamFilename: simpleYAMLStream,
			SegmentedFilename:  simpleYAMLStreamThree,
			Index:              2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			ys := yamlstream.New()

			f, _ := os.Open(tc.YAMLStreamFilename)
			defer f.Close()

			err := ys.Read(f)
			assert.Nil(t, err)

			expectedAsBytes, _ := os.ReadFile(tc.SegmentedFilename)
            yamlDoc := ys.Get(tc.Index)
			assert.Equal(t, string(expectedAsBytes), yamlDoc.String())
		})
	}
}
