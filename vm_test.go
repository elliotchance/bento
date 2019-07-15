package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var vmTests = map[string]struct {
	program        *CompiledProgram
	expectedMemory []interface{}
	expectedOutput string
}{
	"Simple": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("hello"),
					},
					Instructions: []Instruction{
						&CallInstruction{
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
		expectedOutput: "hello\n",
	},
	"Call1": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("Bob"),
					},
					Instructions: []Instruction{
						&CallInstruction{
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
						&CallInstruction{
							Call: "display ?",
							Args: []int{1},
						},
						&CallInstruction{
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
		expectedOutput: "hi\nBob\n",
	},
	"SetText": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText(""), NewText("foo"),
					},
					Instructions: []Instruction{
						&CallInstruction{
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
						NewNumber("0", 6), NewNumber("1.23", 6),
					},
					Instructions: []Instruction{
						&CallInstruction{
							Call: "set ? to ?",
							Args: []int{0, 1},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewNumber("1.23", 6), NewNumber("1.23", 6), // start
		},
	},
	"InlineIfTrue": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("foo"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     1,
							False:    2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("foo"), NewText("match!"), NewText("done"), // start
		},
		expectedOutput: "match!\ndone\n",
	},
	"InlineIfFalse": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     1,
							False:    2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"), // start
		},
		expectedOutput: "done\n",
	},
	"InlineIfElseTrue": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("foo"), NewText("match!"), NewText("no match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     1,
							False:    2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&JumpInstruction{
							Forward: 2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{4},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("foo"), NewText("match!"), NewText("no match!"), NewText("done"), // start
		},
		expectedOutput: "match!\ndone\n",
	},
	"InlineUnlessTrue": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("foo"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     2,
							False:    1,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("foo"), NewText("match!"), NewText("done"), // start
		},
		expectedOutput: "done\n",
	},
	"InlineUnlessFalse": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     2,
							False:    1,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("bar"), NewText("match!"), NewText("done"), // start
		},
		expectedOutput: "match!\ndone\n",
	},
	"InlineUnlessElseTrue": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewText("foo"), NewText("foo"), NewText("match!"), NewText("no match!"), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorEqual,
							True:     2,
							False:    1,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{2},
						},
						&JumpInstruction{
							Forward: 2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{4},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewText("foo"), NewText("foo"), NewText("match!"), NewText("no match!"), NewText("done"), // start
		},
		expectedOutput: "done\n",
	},
	"InlineWhile": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewNumber("0", 6), NewNumber("5", 6), NewNumber("1", 6), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorLessThan,
							True:     1,
							False:    3,
						},
						&CallInstruction{
							Call: "add ? and ? into ?",
							Args: []int{0, 2, 0},
						},
						&JumpInstruction{
							Forward: -2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewNumber("5", 6), NewNumber("5", 6), NewNumber("1", 6), NewText("done"), // start
		},
		expectedOutput: "done\n",
	},
	"InlineUntil": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewNumber("0", 6), NewNumber("5", 6), NewNumber("1", 6), NewText("done"),
					},
					Instructions: []Instruction{
						&ConditionJumpInstruction{
							Left:     0,
							Right:    1,
							Operator: OperatorGreaterThan,
							True:     3,
							False:    1,
						},
						&CallInstruction{
							Call: "add ? and ? into ?",
							Args: []int{0, 2, 0},
						},
						&JumpInstruction{
							Forward: -2,
						},
						&CallInstruction{
							Call: "display ?",
							Args: []int{3},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewNumber("6", 6), NewNumber("5", 6), NewNumber("1", 6), NewText("done"), // start
		},
		expectedOutput: "done\n",
	},
	"SetNumberRequiresRounding": {
		program: &CompiledProgram{
			Functions: map[string]*CompiledFunction{
				"start": {
					Variables: []interface{}{
						NewNumber("0", 1), NewNumber("1.23", 6),
					},
					Instructions: []Instruction{
						&CallInstruction{
							Call: "set ? to ?",
							Args: []int{0, 1},
						},
					},
				},
			},
		},
		expectedMemory: []interface{}{
			NewNumber("1.2", 1), NewNumber("1.23", 6), // start
		},
	},
}

var vmConditionTests = map[string]interface{}{
	`"foo" = "foo"`: true, // text
	`"foo" = "bar"`: false,
	`1.230 = 1.23`:  true, // number
	`1.23 = 2.23`:   false,
	`1.23 = "1.23"`: "cannot compare: number = text", // mixed
	`"1.23" = 1.23`: "cannot compare: text = number",

	`"foo" != "foo"`: false, // text
	`"foo" != "bar"`: true,
	`1.230 != 1.23`:  false, // number
	`1.23 != 2.23`:   true,
	`1.23 != "1.23"`: "cannot compare: number != text", // mixed
	`"1.23" != 1.23`: "cannot compare: text != number",

	`"foo" < "foo"`: false, // text
	`"foo" < "bar"`: false,
	`1.230 < 1.23`:  false, // number
	`1.23 < 2.23`:   true,
	`1.23 < "1.23"`: "cannot compare: number < text", // mixed
	`"1.23" < 1.23`: "cannot compare: text < number",

	`"foo" <= "foo"`: true, // text
	`"foo" <= "bar"`: false,
	`1.230 <= 1.23`:  true, // number
	`1.23 <= 2.23`:   true,
	`1.23 <= "1.23"`: "cannot compare: number <= text", // mixed
	`"1.23" <= 1.23`: "cannot compare: text <= number",

	`"foo" > "foo"`: false, // text
	`"foo" > "bar"`: true,
	`1.230 > 1.23`:  false, // number
	`1.23 > 2.23`:   false,
	`1.23 > "1.23"`: "cannot compare: number > text", // mixed
	`"1.23" > 1.23`: "cannot compare: text > number",

	`"foo" >= "foo"`: true, // text
	`"foo" >= "bar"`: true,
	`1.230 >= 1.23`:  true, // number
	`1.23 >= 2.23`:   false,
	`1.23 >= "1.23"`: "cannot compare: number >= text", // mixed
	`"1.23" >= 1.23`: "cannot compare: text >= number",
}

func TestVirtualMachine_Run(t *testing.T) {
	for testName, test := range vmTests {
		t.Run(testName, func(t *testing.T) {
			vm := NewVirtualMachine(test.program)
			vm.out = bytes.NewBuffer(nil)

			err := vm.Run()
			require.NoError(t, err)

			assert.Equal(t, test.expectedMemory, vm.memory)
			assert.Equal(t, test.expectedOutput, vm.out.(*bytes.Buffer).String())
		})
	}
}

func TestVirtualMachine_ConditionTests(t *testing.T) {
	for test, expected := range vmConditionTests {
		t.Run(test, func(t *testing.T) {
			parser := NewParser(strings.NewReader(
				"start: if " + test + ", display \"yes\"",
			))
			program, err := parser.Parse()
			require.NoError(t, err)

			compiler := NewCompiler(program)
			compiledProgram := compiler.Compile()

			vm := NewVirtualMachine(compiledProgram)
			vm.out = bytes.NewBuffer(nil)
			err = vm.Run()

			switch expected {
			case true:
				assert.Equal(t, "yes\n", vm.out.(*bytes.Buffer).String())
				assert.NoError(t, err)

			case false:
				assert.Equal(t, "", vm.out.(*bytes.Buffer).String())
				assert.NoError(t, err)

			default:
				assert.EqualError(t, err, expected.(string))
			}
		})
	}
}
