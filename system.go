package main

import "fmt"

var System = &Library{
	Sentences: []*Sentence{
		{
			Syntax:  "display ?",
			Handler: display,
		},
	},
}

func display(args []interface{}) {
	fmt.Printf("%v\n", args[0])
}
