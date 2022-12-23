package main

import (
	"flag"
	"fmt"
	"os"

	yamlstream "github.com/jdockerty/yaml-stream"
)

var filename = flag.String("filename", "", "input file containing one or more YAML documents.")
var index = flag.Int("index", 0, "index of the document to access from the YAML stream. Default is 0.")

func main() {

	flag.Parse()

    if *filename == "" {
        fmt.Println("filename is a required flag.")
        return
    }


    ys := yamlstream.New()

    if err := ys.ReadWithOpen(*filename); err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        return
    }
    
    if *index > ys.Count - 1 {
        fmt.Fprintf(os.Stderr, "%d is not a valid index for the YAML stream. Max index is %d\n", *index, ys.Count - 1)
        return
    }
    doc := ys.Get(*index)

    fmt.Println(doc.String())
}
