package main

import (
	"flag"
	"log"
	"os"
)

func main() {
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

		program.Run()
	}
}
