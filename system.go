package main

import "fmt"

var System = &Library{
	Sentences: map[string]*Sentence{
		"display ?": {
			Handler: display,
		},
		"set ? to ?": {
			Handler: setVariable,
		},
	},
}

func display(program *Program, args []interface{}) {
	fmt.Printf("%v\n", program.ValueOf(args[0]))
}

func setVariable(program *Program, args []interface{}) {
	name := string(args[0].(VariableReference))
	program.Variables[name].Value = args[1]
}
