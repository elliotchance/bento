package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"start": {},
				},
			},
		},
		"Display": {
			bento: `Display "Hello, World!"`,
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"start": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								"Hello, World!",
							}),
						},
					},
				},
			},
		},
		"DisplayTwice": {
			bento: "Display \"hello\"\ndisplay \"twice!\"",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"start": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								"hello",
							}),
							System.SentenceForSyntax("display ?", []interface{}{
								"twice!",
							}),
						},
					},
				},
			},
		},
		"Declare1": {
			bento: "declare some-variable is text",
			expected: &Program{
				Variables: map[string]*Variable{
					"some-variable": {
						Type:  "text",
						Value: "",
					},
				},
				Functions: map[string]*Function{
					"start": {},
				},
			},
		},
		"Declare2": {
			bento: "declare foo is text\ndisplay foo",
			expected: &Program{
				Variables: map[string]*Variable{
					"foo": {
						Type:  "text",
						Value: "",
					},
				},
				Functions: map[string]*Function{
					"start": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								VariableReference("foo"),
							}),
						},
					},
				},
			},
		},
		"Function1": {
			bento: "display \"hi\"\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"start": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								"hi",
							}),
						},
					},
					"do something": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								"ok",
							}),
						},
					},
				},
			},
		},
		"Function2": {
			bento: "do something\ndo something:\ndisplay \"ok\"",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"start": {
						Sentences: []*Sentence{
							{},
						},
					},
					"do something": {
						Sentences: []*Sentence{
							System.SentenceForSyntax("display ?", []interface{}{
								"ok",
							}),
						},
					},
				},
			},
		},
		"FunctionWithArgument": {
			bento: "greet persons-name now (persons-name is text):",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"greet ? now": {
						Variables: map[string]*Variable{
							"persons-name": {
								Type: "text",
							},
						},
					},
					"start": {},
				},
			},
		},
		"CallWithArgument": {
			bento: "greet persons-name now (persons-name is text):\ndisplay persons-name",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"greet ? now": {
						Variables: map[string]*Variable{
							"persons-name": {
								Type: "text",
							},
						},
						Sentences: []*Sentence{
							{
								Args: []interface{}{
									VariableReference("persons-name"),
								},
							},
						},
					},
					"start": {},
				},
			},
		},
		"FunctionWithArguments": {
			bento: "say greeting to persons-name (persons-name is text, greeting is text):",
			expected: &Program{
				Variables: map[string]*Variable{},
				Functions: map[string]*Function{
					"say ? to ?": {
						Variables: map[string]*Variable{
							"greeting": {
								Type:     "text",
								Position: 0,
							},
							"persons-name": {
								Type:     "text",
								Position: 1,
							},
						},
					},
					"start": {},
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			parser := NewParser(strings.NewReader(test.bento))
			actual, err := parser.Parse()
			require.NoError(t, err)

			diff := cmp.Diff(test.expected, actual,
				cmpopts.IgnoreTypes((SentenceHandler)(nil)))
			assert.Empty(t, diff)
		})
	}
}
