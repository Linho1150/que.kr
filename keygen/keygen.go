package keygen

import (
	"github.com/jaevor/go-nanoid"
	"errors"
)


func Generate() (string, error) {
	s, error := nanoid.Standard(10)

	if error != nil {
		return "", errors.New("cannot generate candidate key")
	}

	return s(), nil
}

