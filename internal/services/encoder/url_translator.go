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

	buf := ""

	for id > 0 {
		buf = string(Alphabet[id%Base]) + buf
		id /= Base
	}

	return buf
}

func (e *Encoder) Decode(code string) int {
	id := 0

	for _, c := range code {
		id = id*Base + strings.Index(Alphabet, string(c))
	}

	return id
}
