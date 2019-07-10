package main

import (
	"flag"
	"log"
	"os"
)

var (
	flagTrace bool
)

func main() {
	flag.BoolVar(&flagTrace, "trace", false,
		"Show all executed sentences and values.")
	flag.Parse()

	for _, arg := range flag.Args() {
		file, err := os.Open(arg)
		if err != nil {
			log.Fatalln(err)
		}

		parser := NewParser(file)
		program, err := parser.Parse()
		if err != nil {
			log.Fatalln(err)
		}

		program.Trace = flagTrace
		program.Run()
	}
}
