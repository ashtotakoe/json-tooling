package parser

type Token struct {
	val  rune
	line int
	col  int
}

// func tokensFromString(payload string) []Token {
// 	acc := make([]Token, 0, len(payload))

// 	for _, token := range payload {
// 		acc = append(acc,
// 			Token{
// 				val: token,
// 			},
// 		)
// 	}

// 	return acc
// }

func stringFromTokens(payload []Token) string {
	acc := ""

	for _, token := range payload {
		acc += string(token.val)
	}

	return acc
}

func parseToTokens(payload string) []Token {
	line, col := 0, 0
	tokens := make([]Token, 0, len(payload))

	for _, elem := range payload {
		tokens = append(tokens, Token{
			val:  elem,
			line: line,
			col:  col,
		})

		col++

		// if new line
		if elem == '\n' {
			col = 0
			line++
		}
	}

	return tokens

}
