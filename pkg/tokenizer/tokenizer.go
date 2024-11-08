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
	TokenLeftSquare
	TokenRightSquare
	TokenColon
	TokenComma
	TokenString
	TokenKeyValue
	TokenNull
	TokenTrue
	TokenFalse
	TokenNumber
	TokenObject
	TokenArray
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
	case '[':
		return Token{Type: TokenLeftSquare, Value: "["}, nil
	case ']':
		return Token{Type: TokenRightSquare, Value: "]"}, nil
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
		if char == '\\' {
			valid, sequence := t.ValidateEscapeString()
			if !valid {
				return Token{}, fmt.Errorf("invalid escape string")
			}
			str += string(char) + sequence
		} else {
			if char < 0x20 {
				return Token{}, fmt.Errorf("invalid character in string: control character 0x%02X", char)
			}
			str += string(char)
		}
	}
	return Token{}, fmt.Errorf("unterminated string")
}

func (t *Tokenizer) ValidateEscapeString() (bool, string) {
	if !t.scanner.Scan() {
		return false, ""
	}

	char := t.scanner.Bytes()[0]

	switch char {
	case 'b', 'f', 'n', 'r', 't', '"', '\\', '/':
		return true, string(char)
	case 'u':
		unicode := string(char)
		for i := 0; i < 4; i++ {
			if !t.scanner.Scan() {
				return false, ""
			}
			digit := t.scanner.Bytes()[0]
			if !isHexDigit(digit) {
				return false, ""
			}
			unicode += string(digit)
		}
		return true, unicode
	default:
		return false, ""
	}
}

func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
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

// func (t *Tokenizer) ReadNumber(start byte) (Token, error) {
// 	var numStr string
// 	numStr += string(start)

// 	for t.scanner.Scan() {
// 		char := t.scanner.Bytes()[0]
// 		if unicode.IsDigit(rune(char)) || char == '.' {
// 			if len(numStr) > 2 && numStr[0] == '0' && numStr[1] != '.' {
// 				return Token{Type: TokenNumber, Value: numStr}, fmt.Errorf("invalid leading 0 in number")
// 			}
// 			numStr += string(char)
// 		} else {
// 			t.buffer = &char
// 			break
// 		}
// 	}

// 	return Token{Type: TokenNumber, Value: numStr}, nil
// }

func (t *Tokenizer) ReadNumber(start byte) (Token, error) {
	var numStr string
	numStr += string(start)
	var currentChar byte
	furtherProcess := true

	for t.scanner.Scan() {
		currentChar = t.scanner.Bytes()[0]

		// Integer parsing
		if currentChar != '.' && currentChar != 'e' && currentChar != 'E' {
			if unicode.IsDigit(rune(currentChar)) {
				numStr += string(currentChar)
			} else {
				t.buffer = &currentChar
				furtherProcess = false
				break
			}
		} else {
			break
		}
	}

	// Test for leading 0
	if len(numStr) > 1 && numStr[0] == '0' {
		return Token{Type: TokenNumber, Value: numStr}, fmt.Errorf("invalid leading 0 found")
	}

	// Fraction parsing
	if furtherProcess && currentChar == '.' {
		numStr += string(currentChar)
		for t.scanner.Scan() {
			currentChar = t.scanner.Bytes()[0]

			if currentChar != 'e' && currentChar != 'E' {
				if unicode.IsDigit(rune(currentChar)) {
					numStr += string(currentChar)
				} else {
					t.buffer = &currentChar
					furtherProcess = false
					break
				}
			} else {
				break
			}
		}
	}

	// Exponent parsing
	if furtherProcess && (currentChar == 'e' || currentChar == 'E') {
		numStr += string(currentChar)

		if !t.scanner.Scan() {
			return Token{Type: TokenNumber, Value: numStr}, fmt.Errorf("couldn't parse number")
		}

		currentChar = t.scanner.Bytes()[0]

		if currentChar == '+' || currentChar == '-' || unicode.IsDigit(rune(currentChar)) {
			numStr += string(currentChar)
			for t.scanner.Scan() {
				currentChar = t.scanner.Bytes()[0]

				if unicode.IsDigit(rune(currentChar)) {
					numStr += string(currentChar)
				} else {
					t.buffer = &currentChar
					furtherProcess = false
					break
				}
			}
		} else {
			return Token{Type: TokenNumber, Value: numStr}, fmt.Errorf("invalid exponential number")
		}

		// check for invalid sign
		if numStr[len(numStr)-1] == '+' || numStr[len(numStr)-1] == '-' {
			return Token{Type: TokenNumber, Value: numStr}, fmt.Errorf("Invalid sign with no numbers")
		}
	}

	return Token{Type: TokenNumber, Value: numStr}, nil
}
