package parser2

import (
	"back/utils"
	"log"
	"slices"
)

func trim(payload []Token) []Token {
	ignoredVals := []rune{' ', '\t', '\n', '\r'}

	log.Println("trim", len(payload))
	log.Println("payload", stringFromTokens(payload))

	start, end := 0, len(payload)-1

	for i, elem := range payload {

		start = i

		if elem.tokenType == tokenBlank {
			continue
		}

		if elem.tokenType == tokenParsed {
			break
		}

		utils.Assert(elem.tokenType == tokenVal, "token must have type val after all checks")
		utils.Assert(elem.val != nil, "val propery must be non nil")

		// if token is NOT ignored val
		if !slices.ContainsFunc(ignoredVals, func(ignoredVal rune) bool {
			return *elem.val == ignoredVal
		}) {
			break
		}

	}

	for i := end; i > start; i-- {
		log.Println("inside loop", i)
		elem := payload[i]

		end = i

		if elem.tokenType == tokenBlank {
			continue
		}

		if elem.tokenType == tokenParsed {
			break
		}

		utils.Assert(elem.tokenType == tokenVal, "token must have type val after all checks")
		utils.Assert(elem.val != nil, "val propery must be non nil")

		// if token is NOT ignored val
		if !slices.ContainsFunc(ignoredVals, func(ignoredVal rune) bool {
			return *elem.val == ignoredVal
		}) {
			break
		}

	}

	finalLen := end - start + 1
	res := make([]Token, finalLen)
	copy(res, payload[start:end+1])

	return res
}

func insertValIntoExprSegment(val val, segment []Token) {
	// fill everything with blanks
	for i := range segment {
		segment[i] = createBlankToken()
	}

	segment[0] = createParsedToken(val)

}

func searchForUnindentedQuotes(tokens []Token) int {

	acc := 0
	// iterate on every char of keyValChunk
	for i, token := range tokens {

		if token.tokenType != tokenVal {
			continue
		}

		utils.Assert(token.val != nil, "must be nonnill")

		if *token.val == '"' {

			prevToken := tokens[i-1]
			isPrevTokenEqualToBackSlash := prevToken.tokenType == tokenVal && *prevToken.val == '\\'

			// if elem is first it is surely unindented
			// chech if unindented
			if i == 0 || !isPrevTokenEqualToBackSlash  {

				acc++
				continue
			}
		}
	}

	return acc
}

func splitTokens(tokens []Token, sep rune) [][]Token {

	res := make([][]Token, 0)
	acc := make([]Token, 0)

	for _, token := range tokens {
		if token.tokenType == tokenVal && *token.val == sep {

			res = append(res, trim(acc))

			acc = make([]Token, 0)
			continue
		}

		acc = append(acc, token)

	}

	return res

}
