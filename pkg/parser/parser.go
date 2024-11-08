package parser

import (
	"fmt"
	"json-parser/pkg/tokenizer"
)

func Parse(t *tokenizer.Tokenizer) bool {
	stack := []tokenizer.TokenType{}
	processedCount := 0
	nestingCount := 0
	const maxNestingDepth = 19
	for {
		token, err := t.NextToken()
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}

		switch token.Type {
		case tokenizer.TokenLeftBrace:
			stack, processedCount = handleLeftBrace(stack, processedCount)
			nestingCount++
			if nestingCount > maxNestingDepth {
				fmt.Println("Reached maximum nesting limit")
				return false
			}

		case tokenizer.TokenRightBrace:
			var isSuccess bool
			stack, processedCount, isSuccess = handleRightBrace(stack, processedCount)
			if !isSuccess {
				return false
			}
			nestingCount--

		case tokenizer.TokenLeftSquare:
			stack, processedCount = handleLeftSquare(stack, processedCount)
			nestingCount++
			if nestingCount > maxNestingDepth {
				fmt.Println("Reached maximum nesting limit")
				return false
			}

		case tokenizer.TokenRightSquare:
			var isSuccess bool
			stack, processedCount, isSuccess = handleRightSquare(stack, processedCount)
			if !isSuccess {
				return false
			}
			nestingCount--

		case tokenizer.TokenString:
			var isSuccess bool
			stack, isSuccess = handleString(stack, token)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenColon:
			var isSuccess bool
			stack, isSuccess = handleColon(stack)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenComma:
			var isSuccess bool
			stack, isSuccess = handleComma(stack)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenNull:
			var isSuccess bool
			stack, isSuccess = handleNull(stack, token)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenTrue:
			var isSuccess bool
			stack, isSuccess = handleTrue(stack, token)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenFalse:
			var isSuccess bool
			stack, isSuccess = handleFalse(stack, token)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenNumber:
			var isSuccess bool
			stack, isSuccess = handleNumber(stack, token)
			if !isSuccess {
				return false
			}

		case tokenizer.TokenEOF:
			return handleEOF(stack, processedCount)

		default:
			fmt.Printf("Unexpected token type: %d\n", token.Type)
		}
	}
}

func handleLeftBrace(stack []tokenizer.TokenType, processedCount int) ([]tokenizer.TokenType, int) {
	// Opening braces are always appended directly
	stack = append(stack, tokenizer.TokenLeftBrace)
	processedCount++
	return stack, processedCount
}

func handleRightBrace(stack []tokenizer.TokenType, processedCount int) ([]tokenizer.TokenType, int, bool) {
	// Right brace can only follow a KeyValue or a LeftBrace
	if len(stack) == 0 || (stack[len(stack)-1] != tokenizer.TokenKeyValue && stack[len(stack)-1] != tokenizer.TokenLeftBrace) {
		fmt.Println("Error: Reached Right brace prematurely")
		return stack, processedCount, false
	}

	// Remove any KeyValue(s) to reach the LeftBrace
	for len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenKeyValue {
		stack = stack[:len(stack)-1]
	}

	// If no LeftBrace is found throw error
	if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenLeftBrace {
		fmt.Println("Error: cannot find corresponding left brace")
		return stack, processedCount, false
	}

	// Remove LeftBrace
	stack = stack[:len(stack)-1]

	// If a Colon is found, append a KeyValue
	if len(stack) != 0 && stack[len(stack)-1] == tokenizer.TokenColon && stack[len(stack)-2] == tokenizer.TokenString {
		stack = stack[:len(stack)-2]
		stack = append(stack, tokenizer.TokenKeyValue)
	}

	// If it's an array element, append an Object
	if len(stack) != 0 && stack[len(stack)-1] == tokenizer.TokenComma {
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenObject)
	}

	processedCount++
	return stack, processedCount, true
}

func handleLeftSquare(stack []tokenizer.TokenType, processedCount int) ([]tokenizer.TokenType, int) {
	// Opening braces are always appended directly
	stack = append(stack, tokenizer.TokenLeftSquare)
	processedCount++
	return stack, processedCount
}

func handleRightSquare(stack []tokenizer.TokenType, processedCount int) ([]tokenizer.TokenType, int, bool) {
	// Right square cannot follow a comma, left brace, right brace, colon, or key value
	if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftBrace || stack[len(stack)-1] == tokenizer.TokenRightBrace || stack[len(stack)-1] == tokenizer.TokenColon || stack[len(stack)-1] == tokenizer.TokenKeyValue) {
		fmt.Println("Error: Reached right square brace prematurely")
		return stack, processedCount, false
	}

	// Remove all objects to reach a left square
	for len(stack) > 0 && stack[len(stack)-1] != tokenizer.TokenLeftSquare {
		stack = stack[:len(stack)-1]
	}

	// If no left square is found throw error
	if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenLeftSquare {
		fmt.Println("Error: cannot find corresponding left square brace")
		return stack, processedCount, false
	}

	// Remove left square
	stack = stack[:len(stack)-1]

	// If a Colon is found, append a KeyValue
	if len(stack) != 0 && stack[len(stack)-1] == tokenizer.TokenColon && stack[len(stack)-2] == tokenizer.TokenString {
		stack = stack[:len(stack)-2]
		stack = append(stack, tokenizer.TokenKeyValue)
	}

	// If it's an array element, append an Array
	if len(stack) != 0 && stack[len(stack)-1] == tokenizer.TokenComma {
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenArray)
	}

	processedCount++
	return stack, processedCount, true
}

func handleString(stack []tokenizer.TokenType, token tokenizer.Token) ([]tokenizer.TokenType, bool) {
	// KEY String (or Array string) can follow either a comma, left brace, or left square
	if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftBrace || stack[len(stack)-1] == tokenizer.TokenLeftSquare) {
		if stack[len(stack)-1] == tokenizer.TokenComma {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, tokenizer.TokenString)
		return stack, true
		// Key-Value String can only follow a colon
	} else if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
		stack = stack[:len(stack)-1]
		if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
			fmt.Println("Error: could not find corresponding key for value: ", token.Value)
			return stack, false
		}
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenKeyValue)
		return stack, true
	} else {
		fmt.Println("Error: found invalid string:", token.Value)
		return stack, false
	}
}

func handleColon(stack []tokenizer.TokenType) ([]tokenizer.TokenType, bool) {
	// Colon can only follow a String
	if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
		fmt.Println("Error: found Colon prematurely")
		fmt.Println(stack)
		return stack, false
	}
	stack = append(stack, tokenizer.TokenColon)
	return stack, true
}

func handleComma(stack []tokenizer.TokenType) ([]tokenizer.TokenType, bool) {
	// Comma can follow anything expect right/left (square) braces
	if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenLeftBrace || stack[len(stack)-1] == tokenizer.TokenRightBrace || stack[len(stack)-1] == tokenizer.TokenLeftSquare || stack[len(stack)-1] == tokenizer.TokenRightSquare) {
		fmt.Println("Error: Invalid Comma")
		return stack, false
	}
	stack = append(stack, tokenizer.TokenComma)
	return stack, true
}

func handleNull(stack []tokenizer.TokenType, token tokenizer.Token) ([]tokenizer.TokenType, bool) {
	// Null as a Value
	if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
		stack = stack[:len(stack)-1]
		if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
			fmt.Println("Error: could not find corresponding key for value: ", token.Value)
			return stack, false
		}
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenKeyValue)
		return stack, true
		// Null as an array element
	} else if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftSquare) {
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenNull)
		return stack, true
	} else {
		fmt.Println("Error: found invalid null")
		return stack, false
	}
}

func handleTrue(stack []tokenizer.TokenType, token tokenizer.Token) ([]tokenizer.TokenType, bool) {
	// True as a Value
	if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
		stack = stack[:len(stack)-1]
		if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
			fmt.Println("Error: could not find corresponding key for value: ", token.Value)
			return stack, false
		}
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenKeyValue)
		return stack, true
		// True as an array element
	} else if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftSquare) {
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenTrue)
		return stack, true
	} else {
		fmt.Println("Error: found invalid true")
		return stack, false
	}
}

func handleFalse(stack []tokenizer.TokenType, token tokenizer.Token) ([]tokenizer.TokenType, bool) {
	// False as a Value
	if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
		stack = stack[:len(stack)-1]
		if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
			fmt.Println("Error: could not find corresponding key for value: ", token.Value)
			return stack, false
		}
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenKeyValue)
		return stack, true
		// False as an array element
	} else if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftSquare) {
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenFalse)
		return stack, true
	} else {
		fmt.Println("Error: found invalid false")
		return stack, false
	}
}

func handleNumber(stack []tokenizer.TokenType, token tokenizer.Token) ([]tokenizer.TokenType, bool) {
	// Number as a Value
	if len(stack) > 0 && stack[len(stack)-1] == tokenizer.TokenColon {
		stack = stack[:len(stack)-1]
		if len(stack) == 0 || stack[len(stack)-1] != tokenizer.TokenString {
			fmt.Println("Error: could not find corresponding key for value: ", token.Value)
			return stack, false
		}
		stack = stack[:len(stack)-1]
		stack = append(stack, tokenizer.TokenKeyValue)
		return stack, true
		// Number as an array element
	} else if len(stack) > 0 && (stack[len(stack)-1] == tokenizer.TokenComma || stack[len(stack)-1] == tokenizer.TokenLeftSquare) {
		if stack[len(stack)-1] == tokenizer.TokenComma {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, tokenizer.TokenNumber)
		return stack, true
	} else {
		fmt.Println(token.Value)
		fmt.Println("Error: found invalid number")
		return stack, false
	}
}

func handleEOF(stack []tokenizer.TokenType, processedCount int) bool {
	if processedCount != 0 && len(stack) == 0 {
		return true
	} else {
		fmt.Println("Error: Reached EOF prematurely.")
		return false
	}
}
