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

		program, err := Parse(file)
		if err != nil {
			log.Fatalln(err)
		}

		program.Run()
	}
}
