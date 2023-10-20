package main

import (
	"flag"
	"gobridge/generator"
	"gobridge/reader"
	"math/rand"
)

var (
	inputFile     = flag.String("api", "", "Target file to read")
	moduleName    = flag.String("mod", "", "Name of the go module being used for the project")
	tsOutFile     = flag.String("ts", "", "Target location to generate file to read")
	tsServiceName = flag.String("ts_service", "", "Target location to generate file to read")
	goServerFile  = flag.String("server", "", "")
)

func main() {
	flag.Parse()
	rand.Seed(123456789)

	if *inputFile == "" || *moduleName == "" {
		return
	}

	d, err := reader.ParseFile(*inputFile, *moduleName)
	if err != nil {
		panic(err)
	}

	if *tsOutFile != "" {
		err := generator.TSClient(*tsOutFile, *tsServiceName, d)
		if err != nil {
			panic(err)
		}
	}

	if *goServerFile != "" {
		err = generator.Server(*goServerFile, *moduleName, d)
		if err != nil {
			panic(err)
		}
	}
}
