package encoder

import "strings"

const (
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Base     = len(Alphabet)
)

type Encoder struct {
	Alphabet string
	Base     int
}

func New() *Encoder {
	return &Encoder{
		Alphabet: Alphabet,
		Base:     Base,
	}
}

func (e *Encoder) Encode(id int) string {
	if id == 0 {
		return string(Alphabet[0])
	}

	buf := strings.Builder{}

	for id > 0 {
		buf.WriteByte(Alphabet[id%Base])
		id /= Base
	}

	return Reverse(buf.String())
}

func (e *Encoder) Decode(code string) int {
	id := 0

	for _, c := range code {
		id = id*Base + strings.Index(Alphabet, string(c))
	}

	return id
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
