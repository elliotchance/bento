package main

import (
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

const (
	TokenKindEndOfFile    = "end of file"
	TokenKindEndOfLine    = "end of line"
	TokenKindWord         = "word"
	TokenKindNumber       = "number"
	TokenKindText         = "text"
	TokenKindColon        = ":"
	TokenKindOpenBracket  = "("
	TokenKindCloseBracket = ")"
	TokenKindComma        = ","
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
		case ',':
			tokens = append(tokens, Token{TokenKindComma, ""})

		case '(':
			tokens = append(tokens, Token{TokenKindOpenBracket, ""})

		case ')':
			tokens = append(tokens, Token{TokenKindCloseBracket, ""})

		case ':':
			tokens = append(tokens, Token{TokenKindColon, ""})

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
			for start := i; i < len(entire); i++ {
				if !isNumberCharacter(entire[i]) {
					tokens = append(tokens,
						Token{TokenKindNumber, entire[start:i]})
					break
				}

				if i == len(entire)-1 {
					tokens = append(tokens,
						Token{TokenKindNumber, entire[start : i+1]})
					break
				}
			}

			// TODO: Check invalid numbers like 1.2.3

		case ' ', '\t':
			// Ignore whitespace.

		default:
			// Consume a possibly hyphenated word. We have to consume all
			// characters here because the '-' is ambiguous between numbers and
			// words.
			for start := i; i < len(entire); i++ {
				if !isWordCharacter(entire[i]) {
					tokens = append(tokens,
						Token{TokenKindWord, strings.ToLower(entire[start:i])})
					i--
					break
				}

				if i == len(entire)-1 {
					tokens = append(tokens,
						Token{TokenKindWord, strings.ToLower(entire[start : i+1])})
					break
				}
			}
		}
	}

	tokens = appendEndOfLine(tokens)
	tokens = append(tokens, Token{TokenKindEndOfFile, ""})

	return
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
