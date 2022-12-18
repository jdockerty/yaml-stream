# YAML Stream

A Go library to simplify the handling of YAML streams. These are multiple YAML files
which you will see delimited with `---`, for example

```yaml
---
items:
  - 'apples'
  - 'peaches'
language: 'go'
---
name: 'jack'
message: 'hello world'
```

## Usage

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/jdockerty/yaml-stream"
)

func main() {
	someYaml := []byte(`
---
hello: world
---
some_object:
    item_one: one
    item_two: two
    item_three:
        - 1
        - 2
`)

	b := &bytes.Buffer{}

	_, err := b.Write(someYaml)
	if err != nil {
		// handle err
	}

	ys := yamlstream.New()

	err = ys.Read(b)
	if err != nil {
		// handle err
	}

	firstDoc := ys.Get(0)
	secondDoc := ys.Get(1)
	fmt.Println(firstDoc.Bytes()) // Prints []byte of first YAML document
	fmt.Println(secondDoc.String()) // Prints string representation of the 2nd
}
```
