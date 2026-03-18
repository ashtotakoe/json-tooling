package parser2

import (
	"back/utils"
	"errors"
	"log"
	"slices"
	"strconv"
	"strings"
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

func parseFlatExpression(payload []Token) (*val, error) {

	log.Println("-----")

	log.Println(
		stringFromTokens(payload),
	)

	log.Println("-----")

	val, _ := createPrimitiveVal(valTypeBool, true)
	return &val, nil
}

func parseVal(valTokens []Token, parent addable) error {
	valString := stringFromTokens(trim(valTokens))

	// bool
	if valString == "true" || valString == "false" {
		child, err := createPrimitiveVal(
			valTypeBool,
			valString,
		)
		if err != nil {
			return err
		}

		parent.add(child)
		return nil
	}

	valRunes := []rune(valString)

	// string
	if valRunes[0] == '"' && valRunes[len(valRunes)-1] == '"' {

		unindentedQuotesAmountInValue := searchForUnindentedQuotes(valTokens)
		// it must equal either 0 or 2
		if unindentedQuotesAmountInValue == 2 {
			return errors.New("error parsing val " + valString)
		}

		child, err := createPrimitiveVal(
			valTypeString,
			valString,
		)
		if err != nil {
			return err
		}

		parent.add(child)
		return nil
	}

	// obj
	if valRunes[0] == '{' && valRunes[len(valRunes)-1] == '}' {
		return parseObj(valTokens, parent)
	}

	// arr
	if valRunes[0] == '[' && valRunes[len(valRunes)-1] == ']' {
		return parseArr(valTokens, parent)
	}

	// possibly float
	if strings.Contains(valString, ".") {
		parsedNum, err := strconv.ParseFloat(valString, 64)
		if err != nil {
			return err
		}

		child, err := createPrimitiveVal(valTypeNumber, parsedNum)
		if err != nil {
			return err
		}

		parent.add(child)

		return nil
	}

	// int
	parsedNum, err := strconv.ParseInt(valString, 10, 64)
	if err != nil {
		return err
	}

	child, err := createPrimitiveVal(valTypeNumber, parsedNum)
	if err != nil {
		return err
	}

	parent.add(child)

	return nil
}

func parseObj(payload []Token, parent addable) error {

	if len(payload) < 2 {
		return errors.New("failed to parse object payload. length is less than 2")
	}

	if *payload[0].val != '{' || *payload[len(payload)-1].val != '}' {
		return errors.New("obj should be wrapped with {}" + (stringFromTokens(payload)))
	}

	// stripped from outer brackets
	stripped := payload[1 : len(payload)-1]

	keyVals := splitTokens(stripped, ',')

	// objNode := createObjVal()

	for _, keyVal := range keyVals {

		log.Println(stringFromTokens(keyVal))

		// err := parseKeyVal([]rune(keyVal), &objNode)
		// if err != nil {
		// 	return err
		// }
		//
	}

	// parent.add(objNode)
	return nil
}

func parseKeyVal(keyVal []Token, parent addable) error {
	parsed := splitTokens(keyVal, ':')

	if len(parsed) < 2 {
		return errors.New("error parsing keypair " + stringFromTokens(keyVal))
	}

	// run from the beginning to count 2 unindented quotes
	// totalQuotesAtBeginning := 0

	// final index of key part
	// keyPartFinalIndex := 0

	// for i, keyValChunk := range parsed {
	// 	totalQuotesAtBeginning += searchForUnindentedQuotes(keyValChunk)
	//
	// there should not be more than 2 quotes in one chunk
	// if totalQuotesAtBeginning == 2 {
	// 		keyPartFinalIndex = i
	// 		break
	// 	}
	//
	// }

	// we must stop at 2 (number of unindented quotes at the key part), no more, no less
	// if totalQuotesAtBeginning != 2 {
	// 	return errors.New("error parsing keypair " + stringFromTokens(keyVal))
	// }
	//
	// key := strings.Join([]string(parsed[:keyPartFinalIndex+1]), ":")
	// val := strings.Join([]string(parsed[keyPartFinalIndex+1:]), ":")
	//
	// keyValNode := createKeyVal(key)
	//
	// err := parseVal(val, &keyValNode)
	// if err != nil {
	// 	return err
	// }
	//
	// parent.add(keyValNode)
	//
	return nil
}

func parseArr(payload []Token, parent addable) error {
	// trim from outer brackets
	// vals := (strings.Split(
	// 	string(payload[1:len(payload)-1]),
	// 	","))
	//
	// for _, val := range vals {
	// 	err := parseVal(val, parent)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	return nil
}
