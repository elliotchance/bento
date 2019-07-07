package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for testName, test := range map[string]struct {
		bento    string
		expected *Program
	}{
		"Empty": {
			bento: "",
			expected: &Program{
				Variables: map[string]*Variable{},
			},
		},
		"Display": {
			bento: `Display "Hello, World!"`,
			expected: &Program{
				Variables: map[string]*Variable{},
				Sentences: []*Sentence{
					System.SentenceForSyntax("display ?", []interface{}{
						"Hello, World!",
					}),
				},
			},
		},
		"DisplayTwice": {
			bento: "Display \"hello\"\ndisplay \"twice!\"",
			expected: &Program{
				Variables: map[string]*Variable{},
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
		"Declare1": {
			bento: "declare some-variable is text",
			expected: &Program{
				Variables: map[string]*Variable{
					"some-variable": {
						Type:  "text",
						Value: "",
					},
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
				Sentences: []*Sentence{
					System.SentenceForSyntax("display ?", []interface{}{
						VariableReference("foo"),
					}),
				},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			actual, err := Parse(strings.NewReader(test.bento))
			require.NoError(t, err)

			diff := cmp.Diff(test.expected, actual,
				cmpopts.IgnoreTypes((SentenceHandler)(nil)))
			assert.Empty(t, diff)
		})
	}
}
