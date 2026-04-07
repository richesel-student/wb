package shortener

import "testing"

func TestGenerate(t *testing.T) {
	s := Generate(6)

	if len(s) != 6 {
		t.Fatal("wrong length")
	}
}
