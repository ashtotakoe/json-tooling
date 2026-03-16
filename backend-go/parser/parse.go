package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func Parse(payload string) error {
	parsed := parseToTokens(payload)
	trimmed := trim(parsed)

	root := createKeyVal("root")

	err := parseVal(trimmed, &root)

	if err != nil {
		return err
	}

	str, _ := json.MarshalIndent(root, "", "\t")

	fmt.Println("final", string(str))
	return nil
}

func parseObj(payload []rune, parent addable) error {
	if payload[0] != '{' || payload[len(payload)-1] != '}' {
		return errors.New("obj should be wrapped with {}" + string(payload))
	}

	// stripped from outer brackets
	stripped := (string(payload[1 : len(payload)-1]))
	_ = stripped
	fmt.Println(stripped)

	// objNode := createObjVal()

	// for _, keyVal := range trimmed {

	// 	err := parseKeyVal([]rune(keyVal), &objNode)
	// 	if err != nil {
	// 		return err
	// 	}

	// }

	// parent.add(objNode)
	return nil
}

func parseArr(payload []rune, parent addable) error {
	// trim from outer brackets
	vals := (strings.Split(
		string(payload[1:len(payload)-1]),
		","))

	for _, val := range vals {
		err := parseVal(val, parent)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseKeyVal(keyVal []rune, parent addable) error {
	parsed := strings.Split(string(keyVal), ":")

	if len(parsed) < 2 {
		return errors.New("error parsing keypair " + string(keyVal))
	}

	// run from the beginning to count 2 unindented quotes
	totalQuotesAtBeginning := 0

	// final index of key part
	keyPartFinalIndex := 0

	for i, keyValChunk := range parsed {
		totalQuotesAtBeginning += searchForUnindentedQuotes([]rune(keyValChunk))

		// there should not be more than 2 quotes in one chunk
		if totalQuotesAtBeginning == 2 {
			keyPartFinalIndex = i
			break
		}

	}

	// we must stop at 2 (number of unindented quotes at the key part), no more, no less
	if totalQuotesAtBeginning != 2 {
		return errors.New("error parsing keypair " + string(keyVal))
	}

	key := strings.Join([]string(parsed[:keyPartFinalIndex+1]), ":")
	val := strings.Join([]string(parsed[keyPartFinalIndex+1:]), ":")

	keyValNode := createKeyVal(key)

	err := parseVal(val, &keyValNode)
	if err != nil {
		return err
	}

	parent.add(keyValNode)

	return nil
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

func parseVal(valTokens []Token, parent addable) error {
	val := stringFromTokens(valTokens)

	// bool
	if val == "true" || val == "false" {
		child, err := createPrimitiveVal(
			valTypeBool,
			val,
		)
		if err != nil {
			return err
		}

		parent.add(child)
		return nil
	}

	valRunes := []rune(val)

	// string
	if valRunes[0] == '"' && valRunes[len(valRunes)-1] == '"' {

		unindentedQuotesAmountInValue := searchForUnindentedQuotes([]rune(val))
		// it must equal either 0 or 2
		if !(unindentedQuotesAmountInValue == 2 || unindentedQuotesAmountInValue == 0) {
			return errors.New("error parsing val " + val)
		}

		child, err := createPrimitiveVal(
			valTypeString,
			val,
		)
		if err != nil {
			return err
		}

		parent.add(child)
		return nil
	}

	// obj
	if valRunes[0] == '{' && valRunes[len(valRunes)-1] == '}' {
		return parseObj([]rune(val), parent)
	}

	// arr
	if valRunes[0] == '[' && valRunes[len(valRunes)-1] == ']' {
		return parseArr([]rune(val), parent)
	}

	// possibly float
	if strings.Contains(val, ".") {
		parsedNum, err := strconv.ParseFloat(val, 64)
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
	parsedNum, err := strconv.ParseInt(val, 10, 64)
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

func trim(payload []Token) []Token {
	ignoredVals := []rune{' ', '\n'}

	start, end := 0, len(payload)-1

	for i, elem := range payload {

		// if token is NOT ignored val
		if !slices.ContainsFunc(ignoredVals, func(ignoredVal rune) bool {
			return elem.val == ignoredVal
		}) {
			break
		}

		start = i
	}

	for i := end; i > start; i-- {

		// if token is NOT ignored val
		if !slices.ContainsFunc(ignoredVals, func(ignoredVal rune) bool {
			return payload[i].val == ignoredVal
		}) {
			break
		}

		end = i
	}

	finalLen := end - start
	res := make([]Token, finalLen)

	copy(res, payload[start:end])

	return res
}
