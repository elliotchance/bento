package main

import (
	"io"
	"io/ioutil"
	"strings"
)

const (
	TokenKindWord = iota
	TokenKindText
	TokenKindEndline
)

type Token struct {
	Kind  int
	Value string
}

func Tokenize(r io.Reader) (tokens []Token, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	entire := string(data)
	word := ""

	for i := 0; i < len(entire); i++ {
		switch entire[i] {
		case '\n':
			tokens = appendWord(tokens, &word)
			tokens = appendEndline(tokens)

		case '"':
			i++
			start := i
			for ; i < len(entire); i++ {
				if entire[i] == '"' {
					tokens = append(tokens,
						Token{TokenKindText, entire[start:i]})
					break
				}
			}

		case ' ':
			tokens = appendWord(tokens, &word)

		default:
			word += string(entire[i])
		}
	}

	tokens = appendWord(tokens, &word)
	tokens = appendEndline(tokens)

	return
}

func appendWord(tokens []Token, word *string) []Token {
	if *word != "" {
		token := Token{TokenKindWord, strings.ToLower(*word)}
		*word = ""
		tokens = append(tokens, token)
	}

	return tokens
}

func appendEndline(tokens []Token) []Token {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind != TokenKindEndline {
		return append(tokens, Token{TokenKindEndline, ""})
	}

	return tokens
}
