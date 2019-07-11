package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	flagAst bool
)

func main() {
	flag.BoolVar(&flagAst, "ast", false, "Print out the parsed AST and "+
		"exist. This is useful for debugging, but you should not assume that "+
		"the format returned will be consistent or if -ast will remain in "+
		"any future version.")
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

		if flagAst {
			data, err := json.MarshalIndent(program, "", "  ")
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(data))
			os.Exit(0)
		}

		compiledProgram := CompileProgram(program)

		vm := NewVirtualMachine(compiledProgram)
		vm.Run()
	}
}
