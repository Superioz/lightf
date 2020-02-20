package token

import (
	"testing"
)

func TestGen(t *testing.T) {
	charset, length := "abcdefgABCDEFG0123456789", 18
	gen := G(charset, length, &Settings{})

	if gen.Charset != charset || gen.Length != length {
		t.Errorf("Unexpected values in created generator, got charset=%q, length=%q", charset, length)
	}
}

func TestDo(t *testing.T) {
	charset, length := "abcdefgABCDEFG0123456789", 18
	gen := G(charset, length, &Settings{})

	token, err := gen.Do()
	if err != nil {
		t.Error(err)
	}
	if token == nil {
		t.Errorf("Somehow the token couldn't be generated, nil value ...")
	}
	if l := len(token.Str); l != length {
		t.Errorf("Generated token with wrong length, expected %d, got %d", length, l)
	}
}
