package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	for testName, test := range map[string]struct {
		bento    string
		expected []Token
	}{
		"Empty": {
			bento:    "",
			expected: nil,
		},
		"Word": {
			bento: "hello",
			expected: []Token{
				{TokenKindWord, "hello"},
				{TokenKindEndline, ""},
			},
		},
		"TwoWords": {
			bento: "hello world",
			expected: []Token{
				{TokenKindWord, "hello"},
				{TokenKindWord, "world"},
				{TokenKindEndline, ""},
			},
		},
		"Mix1": {
			bento: `display "hello"`,
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindEndline, ""},
			},
		},
		"Mix2": {
			bento: `display "hello" ok`,
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindWord, "ok"},
				{TokenKindEndline, ""},
			},
		},
		"AlwaysLowerCase": {
			bento: `Words in MIXED "Case"`,
			expected: []Token{
				{TokenKindWord, "words"},
				{TokenKindWord, "in"},
				{TokenKindWord, "mixed"},
				{TokenKindText, "Case"},
				{TokenKindEndline, ""},
			},
		},
		"MultipleSpaces": {
			bento: `  foo  bar  " baz  qux"  quux   `,
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindWord, "bar"},
				{TokenKindText, " baz  qux"},
				{TokenKindWord, "quux"},
				{TokenKindEndline, ""},
			},
		},
		"Newlines": {
			bento: "foo\nbar\n\nbaz\n",
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindEndline, ""},
				{TokenKindWord, "bar"},
				{TokenKindEndline, ""},
				{TokenKindWord, "baz"},
				{TokenKindEndline, ""},
			},
		},
		"BeginNewline": {
			bento: "\n\nfoo\nbar",
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindEndline, ""},
				{TokenKindWord, "bar"},
				{TokenKindEndline, ""},
			},
		},
		"DisplayTwice": {
			bento: "Display \"hello\"\ndisplay \"twice!\"",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindEndline, ""},
				{TokenKindWord, "display"},
				{TokenKindText, "twice!"},
				{TokenKindEndline, ""},
			},
		},
		"Comment1": {
			bento:    "# comment",
			expected: nil,
		},
		"Comment2": {
			bento: "# comment\ndisplay",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindEndline, ""},
			},
		},
		"Comment3": {
			bento: "display #comment\ndisplay",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindEndline, ""},
				{TokenKindWord, "display"},
				{TokenKindEndline, ""},
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			actual, err := Tokenize(strings.NewReader(test.bento))
			require.NoError(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}
