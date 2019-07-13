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
	"SetText": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText(""), NewText("foo"),
					},
					Instructions: []Instruction{
						{
							Call: "set ? to ?",
							Args: []int{0, 1},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("foo"), // start
		},
	},
	"SetNumber": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewNumber("0"), NewNumber("1.23"),
					},
					Instructions: []Instruction{
						{
							Call: "set ? to ?",
							Args: []int{0, 1},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewNumber("1.23"), NewNumber("1.23"), // start
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
