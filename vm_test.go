package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var vmTests = map[string]struct {
	program        *CompiledProgram
	expectedMemory []interface{}
}{
	"Simple": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("hello"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("hello"),
		},
	},
	"Call1": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("Bob"),
					},
					Instructions: []Instruction{
						{
							Call: "print ?",
							Args: []int{0},
						},
					},
				},
				"print ?": {
					Variables: []interface{}{
						nil, NewText("hi"),
					},
					Instructions: []Instruction{
						{
							Call: "display ?",
							Args: []int{1},
						},
						{
							Call: "display ?",
							Args: []int{0},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("Bob"),                // start
			NewText("Bob"), NewText("hi"), // print ?
		},
	},
}

func TestVirtualMachine_Run(t *testing.T) {
	for testName, test := range vmTests {
		t.Run(testName, func(t *testing.T) {
			vm := NewVirtualMachine(test.program)
			vm.Run()
			assert.Equal(t, test.expectedMemory, vm.memory)
		})
	}
}
