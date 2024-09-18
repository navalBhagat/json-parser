package parser

import (
	"fmt"
	"json-parser/pkg/tokenizer"
)

func Parse(t *tokenizer.Tokenizer) bool {
	stack := []tokenizer.TokenType{}

	for {
		token, err := t.NextToken()
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}
		switch token.Type {
		case tokenizer.TokenLeftBrace:
			if len(stack) > 0 {
				return false
			}
			stack = append(stack, tokenizer.TokenLeftBrace)
		case tokenizer.TokenRightBrace:
			if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenLeftBrace {
				return false
			}
			stack = stack[:len(stack)-1]
		case tokenizer.TokenEOF:
			if len(stack) == 0 {
				return true
			} else {
				fmt.Println("Error: Reached EOF")
				return false
			}
		default:
			// ignore for now
			fmt.Printf("Unexpected token type: %d\n", token.Type)
		}
	}
}
