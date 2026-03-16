package parser2

import "fmt"

type tokenType int

const (
	tokenVal    tokenType = iota
	tokenParsed           // when child expression (the sequence of tokens) is parsed, it turns into parsed, with val becoming non nil
	tokenBlank            // all remaining tokens in parsed expression become blank and wait to be cleaned
)

type Token struct {
	line int
	col  int

	tokenType tokenType
	val       *rune // simple char representing val. available in tokenVal
	parsedVal *val  // represents already parsed value of corresponding type. available in type tokenParsed
}

func stringFromTokens(payload []Token) string {
	acc := ""

	for _, token := range payload {
		if token.tokenType == tokenBlank {
			acc += " blankToken "
		}

		if token.tokenType == tokenParsed {
			acc += " parsedToken "
		}

		if token.tokenType == tokenVal {
			acc += string(*token.val)
		}
	}

	return acc
}

func parseToTokens(payload string) []Token {
	line, col := 0, 0
	tokens := make([]Token, 0, len(payload))

	for _, elem := range payload {
		tokens = append(tokens, Token{
			tokenType: tokenVal,

			val:  &elem,
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

func createBlankToken() Token {
	return Token{
		tokenType: tokenBlank,
	}
}

func createParsedToken(val val) Token {
	return Token{
		tokenType: tokenParsed,
		parsedVal: &val,
	}
}

func logTokens(tokens []Token, comment string) {
	fmt.Println("----", comment)

	fmt.Println(stringFromTokens(tokens))

	fmt.Println("----")
}
