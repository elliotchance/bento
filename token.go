package main

import (
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

const (
	TokenKindEndOfFile    = "end of file"
	TokenKindEndOfLine    = "new line"
	TokenKindWord         = "word"
	TokenKindNumber       = "number"
	TokenKindText         = "text"
	TokenKindColon        = ":"
	TokenKindOpenBracket  = "("
	TokenKindCloseBracket = ")"
	TokenKindComma        = ","
	TokenKindOperator     = "operator"
	TokenKindEllipsis     = "..."
)

type Token struct {
	Kind  string
	Value string
}

func Tokenize(r io.Reader) (tokens []Token, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	entire := string(data)

	for i := 0; i < len(entire); i++ {
		switch entire[i] {
		case '.':
			// TODO: Check len() allows this.
			if entire[i+1] == '.' && entire[i+2] == '.' {
				tokens = append(tokens, Token{TokenKindEllipsis, ""})
				i += 2
			}

		case ',':
			tokens = append(tokens, Token{TokenKindComma, ""})

		case '(':
			tokens = append(tokens, Token{TokenKindOpenBracket, ""})

		case ')':
			tokens = append(tokens, Token{TokenKindCloseBracket, ""})

		case ':':
			tokens = append(tokens, Token{TokenKindColon, ""})

		case '=', '!', '>', '<':
			var operator string
			operator, i = consumeCharacters(isOperatorCharacter, entire, i)
			tokens = append(tokens, Token{TokenKindOperator, operator})

		case '#':
			tokens = appendEndOfLine(tokens)
			for ; i < len(entire); i++ {
				if entire[i] == '\n' {
					break
				}
			}

		case '\n':
			tokens = appendEndOfLine(tokens)

		case '"':
			i++
			for start := i; i < len(entire); i++ {
				if entire[i] == '"' || i == len(entire)-1 {
					tokens = append(tokens,
						Token{TokenKindText, entire[start:i]})
					break
				}
			}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			var number string
			number, i = consumeCharacters(isNumberCharacter, entire, i)
			tokens = append(tokens, Token{TokenKindNumber, number})

			// TODO: Check invalid numbers like 1.2.3

		case ' ', '\t':
			// Ignore whitespace.

		default:
			// TODO: If nothing is consumed below it will be an infinite loop.

			var word string
			word, i = consumeCharacters(isWordCharacter, entire, i)
			tokens = append(tokens, Token{TokenKindWord, strings.ToLower(word)})
		}
	}

	tokens = appendEndOfLine(tokens)
	tokens = append(tokens, Token{TokenKindEndOfFile, ""})

	return
}

func consumeCharacters(t func(byte) bool, entire string, i int) (string, int) {
	start := i

	for ; i < len(entire); i++ {
		if !t(entire[i]) {
			break
		}
	}

	return entire[start:i], i - 1
}

func isOperatorCharacter(c byte) bool {
	return c == '=' || c == '!' || c == '<' || c == '>'
}

func isNumberCharacter(c byte) bool {
	return (c >= '0' && c <= '9') || c == '.' || c == '-'
}

func isWordCharacter(c byte) bool {
	// TODO: This will not work with unicode characters.
	return unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '-'
}

func appendEndOfLine(tokens []Token) []Token {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind != TokenKindEndOfLine {
		return append(tokens, Token{TokenKindEndOfLine, ""})
	}

	return tokens
}
