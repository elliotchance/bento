package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var sentenceTests = map[string]struct {
	sentence       *Sentence
	expectedSyntax string
}{
	"Display": {
		sentence: &Sentence{
			Tokens: []interface{}{
				"display", NewText("hello"),
			},
		},
		expectedSyntax: "display ?",
	},
}

func TestSentence_Syntax(t *testing.T) {
	for testName, test := range sentenceTests {
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, test.expectedSyntax, test.sentence.Syntax())
		})
	}
}
