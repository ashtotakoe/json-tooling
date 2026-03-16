package parser2

import (
	"log"
	"testing"
)

func trimTest(t *testing.T) {
	tokens := parseToTokens(`
		{
			"aaao": true,
			"bbb": 5,
			"test": [
				{
					"u": 123
				}
			]
		}
		`)

	trimed := stringFromTokens(
		trim(tokens),
	)

	log.Println(trimed[0], trimed[len(trimed) - 1])
	
}
