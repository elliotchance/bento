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
			bento:    "",
			expected: &Program{},
		},
		"Display": {
			bento: `Display "Hello, World!"`,
			expected: &Program{
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
