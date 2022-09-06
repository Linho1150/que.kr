package keygen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	s, e := Generate()

	if e != nil {
		t.Error("fail")
	}

	t.Logf("generated key: %v", s)
}
