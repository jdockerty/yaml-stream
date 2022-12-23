# YAML Stream

A Go library to simplify the handling of YAML streams. These are multiple YAML files
which you will see delimited with `---`, for example

```yaml
items:
  - 'apples'
  - 'peaches'
language: 'go'
---
name: 'jack'
message: 'hello world'
```


### Usage

Using the example YAML file provided in [`testdata/simple_stream.yaml`](testdata/simple_stream.yaml).

Shown below are some use cases of the library, for further details you can view the tests provided within the
`stream_test.go` file.

```go
package main

import (
	"fmt"

	"github.com/jdockerty/yaml-stream"
)

func main() {
	ys := yamlstream.New()

	ys.ReadWithOpen("testdata/simple_stream.yaml")

	firstDocument := ys.Get(0)

	fmt.Println(firstDocument.String())
	// Outputs the YAML first document in the file
	// ---
	// stream_number: 1

	var myMap map[string]int
	firstDocument.Unmarshal(&myMap)

	fmt.Println(myMap["stream_number"]) // Outputs '1'

	var secondValue map[string]int
	ys.GetUnmarshal(1, &secondValue)
	fmt.Println(secondValue["stream_number"]) // Outputs '2'

}
```

When dealing with complex types, it is necessary to instead use a `map[string]interface{}` type or
unmarshal into a `struct` when the schema is known in advance.


#### CLI

Using this library, you can also utilise the CLI tool, available in the `cmd/ys` package, to print out the YAML document
which is contained within a stream.

For example

```
go install -v github.com/jdockerty/yaml-stream/cmd/ys@latest

ys -filename testdata/simple_stream.yaml
ys -filename testdata/simple_stream.yaml -index 2
```

By default, the first document is printed at index 0. The YAML stream is treated as an array of documents.
