package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type TokenID uint

const (
	TOKEN_PROGRAM TokenID = iota
	TOKEN_EXPRESSION
	TOKEN_VALUE
	TOKEN_OPENPAREN
	TOKEN_CLOSEPAREN
)

var ErrNoLex = errors.New("failed to lex token")

type Token struct {
	id      TokenID
	content string
}

var rules = map[TokenID]string{
	TOKEN_PROGRAM:    "Program",
	TOKEN_EXPRESSION: "Expression",
	TOKEN_VALUE:      "Value",
	TOKEN_OPENPAREN:  "Open Parentheses",
	TOKEN_CLOSEPAREN: "Close Parentheses",
}

func (t Token) String() string {

	var sanitised = t.content
	sanitised = strings.Replace(sanitised, "\n", "\\n", -1)

	return fmt.Sprintf("%s: %s", rules[t.id], sanitised)
}

func Lex(source string) ([]Token, error) {
	tokens := make([]Token, 0)
	_, err := LexProgram(source, &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func LexProgram(source string, tokens *[]Token) (int, error) {
	var program_token Token
	program_token.id = TOKEN_PROGRAM

	progToks := make([]Token, 0)

	var read int = 0
	for {
		size, err := LexExpression(source[read:], &progToks)
		if err != nil {
			return read, err
		}
		read += size

		size, err = LexWhiteSpace(source[read:])
		if err != nil {
			return read, err
		}
		read += size

		_, err = LexEOF(source[read:], &progToks)
		if err == nil {
			break
		} else if err != ErrNoLex {
			return read, err
		}
	}
	program_token.content = source[:read]
	(*tokens) = append((*tokens), program_token)

	(*tokens) = append((*tokens), progToks...)
	return read, nil
}

func LexEOF(source string, tokens *[]Token) (int, error) {
	if source == "" {
		return 0, nil
	} else {
		return 0, ErrNoLex
	}
}

func LexExpression(source string, tokens *[]Token) (int, error) {
	var expr_token Token
	expr_token.id = TOKEN_EXPRESSION
	var read int = 0

	exprToks := make([]Token, 0)

	size, err := LexValue(source[read:], tokens)
	if err == nil {
		read += size
		return read, nil
	} else if err != ErrNoLex {
		return read, err
	}

	size, err = LexOpenParen(source[read:], &exprToks)
	if err != nil {
		return read, err
	}
	read += size

	size, err = LexWhiteSpace(source[read:])
	if err != nil {
		return read, err
	}
	read += size

	size, err = LexExpression(source[read:], &exprToks)
	if err != nil {
		return read, err
	}
	read += size

	size, err = LexWhiteSpace(source[read:])
	if err != nil {
		return read, err
	}
	read += size

	for {
		size, err := LexExpression(source[read:], &exprToks)
		if err != nil {
			return read, err
		}
		read += size

		size, err = LexWhiteSpace(source[read:])
		if err != nil {
			return read, err
		}
		read += size

		size, err = LexCloseParen(source[read:], &exprToks)
		if err == nil {
			read += size
			break
		} else if err != ErrNoLex {
			return read, err
		}

		size, err = LexWhiteSpace(source[read:])
		if err != nil {
			return read, err
		}
		read += size
	}

	expr_token.content = source[:read]
	(*tokens) = append((*tokens), expr_token)

	(*tokens) = append((*tokens), exprToks...)

	return read, nil
}

func LexValue(source string, tokens *[]Token) (int, error) {
	var ident_token Token
	ident_token.id = TOKEN_VALUE
	var read int = 0

	for ; read < len(source); read++ {
		if !isIdentChar(rune(source[read])) {
			break
		}
	}

	if read == 0 {
		return 0, ErrNoLex
	}

	ident_token.content = source[:read]
	(*tokens) = append((*tokens), ident_token)
	return read, nil
}

var allowed_specials = "+-/*_"

func isIdentChar(c rune) bool {
	return (unicode.IsLetter(c) || unicode.IsDigit(c) || strings.Contains(allowed_specials, string(c)))
}

func LexOpenParen(source string, tokens *[]Token) (int, error) {
	if len(source) == 0 {
		return 0, io.EOF
	}
	if source[0] == '(' {
		return 1, nil
	} else {
		return 0, ErrNoLex
	}
}
func LexCloseParen(source string, tokens *[]Token) (int, error) {
	if len(source) == 0 {
		return 0, io.EOF
	}
	if source[0] == ')' {
		(*tokens) = append((*tokens), Token{TOKEN_CLOSEPAREN, ")"})
		return 1, nil
	} else {
		return 0, ErrNoLex
	}
}

func LexWhiteSpace(source string) (int, error) {
	var read int = 0

	for ; read < len(source); read++ {
		if !unicode.IsSpace(rune(source[read])) {
			break
		}
	}

	return read, nil
}
