package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	for testName, test := range map[string]struct {
		bento    string
		expected *Program
	}{
		"Empty": {
			bento: "",
			expected: &Program{
				Functions: map[string]*Function{},
			},
		},
		"EmptyStart": {
			bento: "start:",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
					},
				},
			},
		},
		"Display": {
			bento: `start: Display "Hello, World!"`,
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("Hello, World!"),
								},
							},
						},
					},
				},
			},
		},
		"Display2": {
			bento: "start:\nDisplay \"Hello, World!\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("Hello, World!"),
								},
							},
						},
					},
				},
			},
		},
		"DisplayTwice": {
			bento: "start: Display \"hello\"\ndisplay \"twice!\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("hello"),
								},
							},
							{
								Tokens: []interface{}{
									"display", NewText("twice!"),
								},
							},
						},
					},
				},
			},
		},
		"Declare1": {
			bento: "start: declare some-variable is text",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "some-variable",
								Type:       "text",
								LocalScope: true,
							},
						},
					},
				},
			},
		},
		"Declare2": {
			bento: "start: declare foo is text\ndisplay foo",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "text",
								LocalScope: true,
							},
						},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", VariableReference("foo"),
								},
							},
						},
					},
				},
			},
		},
		"Function1": {
			bento: "start:  display \"hi\"\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("hi"),
								},
							},
						},
					},
					"do something": {
						Definition: &Sentence{Tokens: []interface{}{"do", "something"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("ok"),
								},
							},
						},
					},
				},
			},
		},
		"Function2": {
			bento: "start:do something\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"do", "something",
								},
							},
						},
					},
					"do something": {
						Definition: &Sentence{Tokens: []interface{}{"do", "something"}},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", NewText("ok"),
								},
							},
						},
					},
				},
			},
		},
		"FunctionWithArgument": {
			bento: "greet persons-name now (persons-name is text):",
			expected: &Program{
				Functions: map[string]*Function{
					"greet ? now": {
						Definition: &Sentence{Tokens: []interface{}{"greet", VariableReference("persons-name"), "now"}},
						Variables: []*VariableDefinition{
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
					},
				},
			},
		},
		"CallWithArgument": {
			bento: "greet persons-name now (persons-name is text):\ndisplay persons-name",
			expected: &Program{
				Functions: map[string]*Function{
					"greet ? now": {
						Definition: &Sentence{Tokens: []interface{}{"greet", VariableReference("persons-name"), "now"}},
						Variables: []*VariableDefinition{
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"display", VariableReference("persons-name"),
								},
							},
						},
					},
				},
			},
		},
		"FunctionWithArguments": {
			bento: "say greeting to persons-name (persons-name is text, greeting is text):",
			expected: &Program{
				Functions: map[string]*Function{
					"say ? to ?": {
						Definition: &Sentence{Tokens: []interface{}{
							"say",
							VariableReference("greeting"),
							"to",
							VariableReference("persons-name"),
						}},
						Variables: []*VariableDefinition{
							{
								Name:       "greeting",
								Type:       "text",
								LocalScope: false,
							},
							{
								Name:       "persons-name",
								Type:       "text",
								LocalScope: false,
							},
						},
					},
				},
			},
		},
		"DeclareNumber": {
			bento: "start: declare foo is number",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
							},
						},
					},
				},
			},
		},
		"SetNegativeNumber": {
			bento: "start: declare foo is number\nset foo to -1.23",
			expected: &Program{
				Functions: map[string]*Function{
					"start": {
						Definition: &Sentence{Tokens: []interface{}{"start"}},
						Variables: []*VariableDefinition{
							{
								Name:       "foo",
								Type:       "number",
								LocalScope: true,
							},
						},
						Sentences: []*Sentence{
							{
								Tokens: []interface{}{
									"set", VariableReference("foo"), "to", NewNumber("-1.23"),
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			parser := NewParser(strings.NewReader(test.bento))
			actual, err := parser.Parse()
			require.NoError(t, err)

			diff := cmp.Diff(test.expected, actual,
				cmpopts.AcyclicTransformer("NumberToString",
					func(number *big.Rat) string {
						return number.FloatString(6)
					}))

			assert.Empty(t, diff)
		})
	}
}
