package main

import (
	"fmt"
	"strings"
)

type Program struct {
	Variables       map[string]*Variable
	Functions       map[string]*Function
	CurrentFunction string
	Trace           bool
	StackLevel      int
}

func (program *Program) Run() {
	program.CurrentFunction = "start"
	program.Functions[program.CurrentFunction].Run(program)
}

func (program *Program) ValueOf(val interface{}) interface{} {
	switch v := val.(type) {
	case VariableReference:
		// Local variable.
		if v, ok := program.Functions[program.CurrentFunction].Variables[string(v)]; ok {
			return v.Value
		}

		// Global variable.
		if v, ok := program.Variables[string(v)]; ok {
			return v.Value
		}

		panic(fmt.Sprintf("unknown variable: %s", string(v)))

	default:
		return val
	}
}

func (program *Program) SentenceForSyntax(syntax string, args []interface{}) *Sentence {
	if _, ok := program.Functions[syntax]; ok {
		return &Sentence{
			Syntax: syntax,
			Handler: func(program *Program, args []interface{}) {
				program.CurrentFunction = syntax

				fn := program.Functions[program.CurrentFunction]
				for _, variable := range fn.Variables {
					variable.Value = args[variable.Position]
				}

				fn.Run(program)
			},
			Args: args,
		}
	}

	return nil
}

func (program *Program) PrintTrace(line string, args []interface{}) {
	if program.Trace {
		for len(args) > 0 {
			line = strings.Replace(line, "?", fmt.Sprintf("%v", args[0]), 1)
			args = args[1:]
		}

		fmt.Println("# " + strings.Repeat("  ", program.StackLevel) + line)
	}
}
