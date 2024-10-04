package tokenizer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type TokenType int

const (
	TokenLeftBrace TokenType = iota
	TokenRightBrace
	TokenColon
	TokenComma
	TokenString
	TokenKeyValue
	TokenNull
	TokenTrue
	TokenFalse
	TokenNumber
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string // "" for no value
}

type Tokenizer struct {
	scanner *bufio.Scanner
	buffer  *byte
}

func NewTokenizerFromReader(r io.Reader) *Tokenizer {
	scanner := bufio.NewScanner(r)
	scanner.Split(byteByByteSplitter)
	return &Tokenizer{
		scanner: scanner,
		buffer:  nil,
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
	var char byte

	if t.buffer != nil {
		char = *t.buffer
		t.buffer = nil
	} else {
		if !t.scanner.Scan() {
			if err := t.scanner.Err(); err != nil {
				return Token{}, err
			}
			return Token{Type: TokenEOF}, nil
		}
		char = t.scanner.Bytes()[0]
	}

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
	case 'n':
		return t.ReadNull("n")
	case 't':
		return t.ReadTrue("t")
	case 'f':
		return t.ReadFalse("f")
	}

	if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
		return t.NextToken()
	}

	if unicode.IsDigit(rune(char)) || char == '-' {
		return t.ReadNumber(char)
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

func (t *Tokenizer) ReadNull(str string) (Token, error) {
	expected := "null"
	for t.scanner.Scan() {
		char := t.scanner.Bytes()[0]
		str += string(char)
		if len(str) > len(expected) || char != expected[len(str)-1] {
			t.buffer = &char
			str = str[:len(str)-1]
			break
		}
		if len(str) == len(expected) {
			break
		}
	}

	if str == "null" {
		return Token{Type: TokenNull, Value: "null"}, nil
	}
	return Token{}, fmt.Errorf("incorrect spelling for 'null'")
}

func (t *Tokenizer) ReadTrue(str string) (Token, error) {
	expected := "true"
	for t.scanner.Scan() {
		char := t.scanner.Bytes()[0]
		str += string(char)
		if len(str) > len(expected) || char != expected[len(str)-1] {
			t.buffer = &char
			str = str[:len(str)-1]
			break
		}
		if len(str) == len(expected) {
			break
		}
	}

	if str == "true" {
		return Token{Type: TokenTrue, Value: "true"}, nil
	}
	return Token{}, fmt.Errorf("incorrect spelling for 'true'")
}

func (t *Tokenizer) ReadFalse(str string) (Token, error) {
	expected := "false"
	for t.scanner.Scan() {
		char := t.scanner.Bytes()[0]
		str += string(char)
		if len(str) > len(expected) || char != expected[len(str)-1] {
			t.buffer = &char
			str = str[:len(str)-1]
			break
		}
		if len(str) == len(expected) {
			break
		}
	}

	if str == "false" {
		return Token{Type: TokenFalse, Value: "false"}, nil
	}
	return Token{}, fmt.Errorf("incorrect spelling for 'false'")
}

func (t *Tokenizer) ReadNumber(start byte) (Token, error) {
	var numStr string
	numStr += string(start)

	for t.scanner.Scan() {
		char := t.scanner.Bytes()[0]
		if unicode.IsDigit(rune(char)) || char == '.' {
			numStr += string(char)
		} else {
			t.buffer = &char
			break
		}
	}

	return Token{Type: TokenNumber, Value: numStr}, nil
}
