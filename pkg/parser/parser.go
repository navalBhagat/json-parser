package parser

import (
	"fmt"
	"json-parser/pkg/tokenizer"
)

func Parse(t *tokenizer.Tokenizer) bool {
	stack := []tokenizer.TokenType{}
	count := 0
	for {
		token, err := t.NextToken()
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}

		switch token.Type {
		case tokenizer.TokenLeftBrace:
			if len(stack) > 0 {
				fmt.Println("Error: Reached Left brace prematurely.")
				return false
			}
			stack = append(stack, tokenizer.TokenLeftBrace)
			count++

		case tokenizer.TokenRightBrace:
			if len(stack) == 0 || (stack[len(stack)-1] != tokenizer.TokenKeyValue && stack[len(stack)-1] != tokenizer.TokenLeftBrace) {
				fmt.Println("Error: Reached Right brace prematurely")
				return false
			}

			for len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenKeyValue {
				stack = stack[:len(stack)-1]
			}

			if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenLeftBrace {
				fmt.Println("Error: cannot find corresponding left brace")
				return false
			}

			stack = stack[:len(stack)-1]
			count++

		case tokenizer.TokenString:
			if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftBrace) {
				if stack[len(stack)-1] == tokenizer.TokenComma {
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, tokenizer.TokenString)
			}

			if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
					fmt.Println("Error: could not find corresponding key for value: ", token.Value)
					return false
				}
				stack = stack[:len(stack)-1]
				stack = append(stack, tokenizer.TokenKeyValue)
			}

		case tokenizer.TokenColon:
			if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
				fmt.Println("Error: found Colon prematurely")
				return false
			}
			stack = append(stack, tokenizer.TokenColon)

		case tokenizer.TokenComma:
			if stack[len(stack)-1] != tokenizer.TokenKeyValue {
				fmt.Println("Error: did not expect comma")
				return false
			}
			stack = append(stack, tokenizer.TokenComma)

		case tokenizer.TokenEOF:
			if count != 0 && len(stack) == 0 {
				return true
			} else {
				fmt.Println("Error: Reached EOF prematurely.")
				return false
			}

		default:
			fmt.Printf("Unexpected token type: %d\n", token.Type)
		}
	}
}
