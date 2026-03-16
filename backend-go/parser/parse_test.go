package parser

import (
	"testing"
)

func BasicTest(t *testing.T) {
	err := Parse(`
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

	if err != nil {
		t.Error("should not have error but have " + err.Error())
	}
}
