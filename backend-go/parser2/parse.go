package parser2

import (
	"back/utils"
	"errors"
	"log"
	"slices"
)

var openinngBrackets, closingBrackets = []rune{'{', '['},
	[]rune{'}', ']'}

var correspondingBrackets = map[rune]rune{
	'{': '}',
	'[': ']',
}

func Parse(payload string) error {
	tokens := parseToTokens(payload)

	root := createKeyVal("root")

	parse(tokens, &root)

	// if err != nil
	// 	return err
	// }

	// str, _ := json.MarshalIndent(root, "", "\t")

	// fmt.Println("final", string(str))
	return nil
}

func parse(payload []Token, parent addable) (*val, error) {

	payloadCopy := make([]Token, len(payload))
	copy(payloadCopy, payload)

	trimmed := trim(payloadCopy)
	// to escape surrounding quotes
	escaped := trimmed[1 : len(trimmed)-1]

	var openingBracket *rune = nil
	stack := 0
	openingBracketIndex := 0

	for i, token := range escaped {

		utils.Assert(token.tokenType == tokenVal, "token on the read of the parse must be val(not parsed, not blank)")

		if slices.Contains(openinngBrackets, *token.val) {
			// if it is first occurancd - set bracket,
			// from now we care only about them
			if openingBracket == nil {

				openingBracketIndex = i
				openingBracket = token.val

				utils.Assert(stack == 0,
					"stack must equal 0, when first opening bracket is encountered")

				stack++

				continue
			}

			if openingBracket == token.val {
				stack++
			}

		}

		if slices.Contains(closingBrackets, *token.val) {
			// if it is first occurancd - set bracket,
			// from now we care only about them
			if openingBracket == nil {
				return nil, errors.New("closing bracket before openened " + string(*token.val))
			}

			// deduce if its corresponding closing bracket
			if *token.val == correspondingBrackets[*openingBracket] {
				stack--
			}
		}

		if openingBracket != nil && stack == 0 {

			segment := escaped[openingBracketIndex : i+1]

			// pass including surrounding brackets
			val, err := parse(segment, parent)
			if err != nil {
				return nil, err
			}

			// insert val into current slice, and overwrite remaining with blanks
			insertValIntoExprSegment(*val, segment)

			// reset
			openingBracket = nil
			openingBracketIndex = 0

			continue
		}

	}

	utils.Assert(stack == 0, "opening and closed quotes must be same in expression")

	// if there were no brackets
	// if !isBracketEncountered {
	// 	return parseFlatExpression(trimmed)
	// }

	val, err := parseFlatExpression(trimmed)

	if err != nil {
		return nil, err
	}

	return val, nil
}

func insertValIntoExprSegment(val val, segment []Token) {
	// fill everything with blanks
	for i := range segment {
		segment[i] = createBlankToken()
	}

	segment[0] = createParsedToken(val)

}

func parseFlatExpression(payload []Token) (*val, error) {

	log.Println("-----")

	log.Println(
		stringFromTokens(payload),
	)

	log.Println("-----")

	val, _ := createPrimitiveVal(valTypeBool, true)
	return &val, nil
}

func searchForUnindentedQuotes(keyValChunk []rune) int {

	acc := 0
	// iterate on every char of keyValChunk
	for i, c := range keyValChunk {

		if c == '"' {

			// if elem is first it is surely unindented
			// chech if unindented
			if i == 0 || keyValChunk[i-1] != '\\' {
				acc++
				continue
			}
		}
	}

	return acc
}

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
