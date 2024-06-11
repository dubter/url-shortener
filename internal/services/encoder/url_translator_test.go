package encoder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	testCases := map[int]string{
		0:  "a",
		61: "9",
		62: "ba",
		63: "bb",
	}

	encoder := New()

	for id, expectedCode := range testCases {
		assert.Equal(t, expectedCode, encoder.Encode(id))
	}
}

func TestDecode(t *testing.T) {
	testCases := map[string]int{
		"a":  0,
		"9":  61,
		"ba": 62,
		"bb": 63,
	}

	encoder := New()

	for code, expectedID := range testCases {
		assert.Equal(t, expectedID, encoder.Decode(code))
	}
}
