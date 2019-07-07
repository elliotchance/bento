package main

import "strings"

type Tokens []Token

func (tokens Tokens) Syntax() string {
	var words []string
	for _, token := range tokens {
		words = append(words, token.Value)
	}

	return strings.Join(words, " ")
}
