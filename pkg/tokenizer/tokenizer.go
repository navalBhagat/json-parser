package tokenizer

import (
	"bufio"
	"fmt"
	"io"
)

type TokenType int

const (
	TokenLeftBrace TokenType = iota
	TokenRightBrace
	TokenColon
	TokenComma
	TokenString
	TokenKeyValue
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string // "" for no value
}

type Tokenizer struct {
	scanner *bufio.Scanner
}

func NewTokenizerFromReader(r io.Reader) *Tokenizer {
	scanner := bufio.NewScanner(r)
	scanner.Split(byteByByteSplitter)
	return &Tokenizer{
		scanner: scanner,
	}
}

func byteByByteSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) > 0 {
		return 1, data[:1], nil
	}

	return 0, nil, nil
}

func (t *Tokenizer) NextToken() (Token, error) {
	if !t.scanner.Scan() {
		if err := t.scanner.Err(); err != nil {
			return Token{}, err
		}
		return Token{Type: TokenEOF}, nil
	}

	char := t.scanner.Bytes()[0]

	switch char {
	case '{':
		return Token{Type: TokenLeftBrace, Value: "{"}, nil
	case '}':
		return Token{Type: TokenRightBrace, Value: "}"}, nil
	case ':':
		return Token{Type: TokenColon, Value: ":"}, nil
	case ',':
		return Token{Type: TokenComma, Value: ","}, nil
	case '"':
		return t.ReadString()
	}

	if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
		return t.NextToken()
	}

	return Token{}, fmt.Errorf("unexpected character: %c", char)
}

func (t *Tokenizer) ReadString() (Token, error) {
	var str string
	for t.scanner.Scan() {
		char := t.scanner.Bytes()[0]
		if char == '"' {
			return Token{Type: TokenString, Value: str}, nil
		}
		str += string(char)
	}
	return Token{}, fmt.Errorf("unterminated string")
}
