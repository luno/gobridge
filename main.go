package main

import (
	"flag"
	"gobridge/generator"
	"gobridge/reader"
)

var inputFile = flag.String("go_target", "", "Target file to read")
var moduleName = flag.String("go_mod", "", "Name of the go module being used for the project")
var tsOutFile = flag.String("ts_output", "", "Target location to generate file to read")

func main() {
	flag.Parse()

	if *inputFile == "" || *moduleName == "" {
		return
	}

	data, fs, err := reader.ParseFile(*inputFile, *moduleName)
	if err != nil {
		panic(err)
	}

	if *tsOutFile != "" {
		err := generator.TSClient(*tsOutFile, data, fs)
		if err != nil {
			panic(err)
		}
	}
}

