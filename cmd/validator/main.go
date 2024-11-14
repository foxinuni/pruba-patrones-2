package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/foxinuni/patrones-prueba/internal"
)

var filepath string
var output string

func init() {
	flag.StringVar(&filepath, "file", "", "path to file")
	flag.StringVar(&output, "output", "", "path to output file")
	flag.Parse()

	if filepath == "" {
		flag.PrintDefaults()
		panic("-file most be specified")
	}
}

func main() {
	// open output file
	file, err := os.OpenFile("errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// start parser
	parser := internal.NewExcelParser(10)
	_, errors, err := parser.ParseFile(filepath)
	if err != nil {
		panic(err)
	}

	// write string to output file
	for err := range errors {
		file.WriteString(fmt.Sprintf("%s:%d: %v\n", err.Page, err.Line, err.Err))
	}
}
