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
			bento: "",
			expected: []Token{
				{TokenKindEndOfFile, ""},
			},
		},
		"Word": {
			bento: "hello",
			expected: []Token{
				{TokenKindWord, "hello"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"TwoWords": {
			bento: "hello world",
			expected: []Token{
				{TokenKindWord, "hello"},
				{TokenKindWord, "world"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Mix1": {
			bento: `display "hello"`,
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Mix2": {
			bento: `display "hello" ok`,
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindWord, "ok"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"AlwaysLowerCase": {
			bento: `Words in MIXED "Case"`,
			expected: []Token{
				{TokenKindWord, "words"},
				{TokenKindWord, "in"},
				{TokenKindWord, "mixed"},
				{TokenKindText, "Case"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"MultipleSpaces": {
			bento: `  foo  bar  " baz  qux"  quux   `,
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindWord, "bar"},
				{TokenKindText, " baz  qux"},
				{TokenKindWord, "quux"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Newlines": {
			bento: "foo\nbar\n\nbaz\n",
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "bar"},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "baz"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"BeginNewline": {
			bento: "\n\nfoo\nbar",
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "bar"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"DisplayTwice": {
			bento: "Display \"hello\"\ndisplay \"twice!\"",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindText, "hello"},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "display"},
				{TokenKindText, "twice!"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Comment1": {
			bento: "# comment",
			expected: []Token{
				{TokenKindEndOfFile, ""},
			},
		},
		"Comment2": {
			bento: "# comment\ndisplay",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Comment3": {
			bento: "display #comment\ndisplay",
			expected: []Token{
				{TokenKindWord, "display"},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "display"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Function1": {
			bento: "do something:\nsomething else",
			expected: []Token{
				{TokenKindWord, "do"},
				{TokenKindWord, "something"},
				{TokenKindColon, ""},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "something"},
				{TokenKindWord, "else"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Function2": {
			bento: "do something: something else",
			expected: []Token{
				{TokenKindWord, "do"},
				{TokenKindWord, "something"},
				{TokenKindColon, ""},
				{TokenKindEndOfLine, ""},
				{TokenKindWord, "something"},
				{TokenKindWord, "else"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"Tabs": {
			bento: `	foo	bar "baz	"	`,
			expected: []Token{
				{TokenKindWord, "foo"},
				{TokenKindWord, "bar"},
				{TokenKindText, "baz	"},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"FunctionWithArgument": {
			bento: `greet persons-name now (persons-name is text):`,
			expected: []Token{
				{TokenKindWord, "greet"},
				{TokenKindWord, "persons-name"},
				{TokenKindWord, "now"},
				{TokenKindOpenBracket, ""},
				{TokenKindWord, "persons-name"},
				{TokenKindWord, "is"},
				{TokenKindWord, "text"},
				{TokenKindCloseBracket, ""},
				{TokenKindColon, ""},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
			},
		},
		"FunctionWithArguments": {
			bento: `say greeting to persons-name (persons-name is text, greeting is text):`,
			expected: []Token{
				{TokenKindWord, "say"},
				{TokenKindWord, "greeting"},
				{TokenKindWord, "to"},
				{TokenKindWord, "persons-name"},
				{TokenKindOpenBracket, ""},
				{TokenKindWord, "persons-name"},
				{TokenKindWord, "is"},
				{TokenKindWord, "text"},
				{TokenKindComma, ""},
				{TokenKindWord, "greeting"},
				{TokenKindWord, "is"},
				{TokenKindWord, "text"},
				{TokenKindCloseBracket, ""},
				{TokenKindColon, ""},
				{TokenKindEndOfLine, ""},
				{TokenKindEndOfFile, ""},
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
