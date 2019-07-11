package main

import (
	"io"
	"io/ioutil"
	"strings"
)

const (
	TokenKindEndOfFile    = "end of file"
	TokenKindEndOfLine    = "end of line"
	TokenKindWord         = "word"
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
	word := ""

	for i := 0; i < len(entire); i++ {
		switch entire[i] {
		case ',':
			tokens = appendWord(tokens, &word)
			tokens = append(tokens, Token{TokenKindComma, ""})

		case '(':
			tokens = appendWord(tokens, &word)
			tokens = append(tokens, Token{TokenKindOpenBracket, ""})

		case ')':
			tokens = appendWord(tokens, &word)
			tokens = append(tokens, Token{TokenKindCloseBracket, ""})

		case ':':
			tokens = appendWord(tokens, &word)
			tokens = append(tokens, Token{TokenKindColon, ""})

		case '#':
			tokens = appendEndOfLine(tokens)
			for ; i < len(entire); i++ {
				if entire[i] == '\n' {
					break
				}
			}

		case '\n':
			tokens = appendWord(tokens, &word)
			tokens = appendEndOfLine(tokens)

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

		case ' ', '\t':
			tokens = appendWord(tokens, &word)

		default:
			word += string(entire[i])
		}
	}

	tokens = appendWord(tokens, &word)
	tokens = appendEndOfLine(tokens)
	tokens = append(tokens, Token{TokenKindEndOfFile, ""})

	return
}

func appendWord(tokens []Token, word *string) []Token {
	if *word != "" {
		token := Token{TokenKindWord, strings.ToLower(*word)}
		*word = ""
		return append(tokens, token)
	}

	return tokens
}

func appendEndOfLine(tokens []Token) []Token {
	if len(tokens) > 0 && tokens[len(tokens)-1].Kind != TokenKindEndOfLine {
		return append(tokens, Token{TokenKindEndOfLine, ""})
	}

	return tokens
}
