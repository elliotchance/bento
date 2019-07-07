package main

import "fmt"

type Program struct {
	Variables       map[string]*Variable
	Functions       map[string]*Function
	CurrentFunction string
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
