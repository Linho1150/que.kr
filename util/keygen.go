package util

import (
	"errors"

	"github.com/jaevor/go-nanoid"
)

func GenerateShortKey() (string, error) {
	s, error := nanoid.Standard(10)

	if error != nil {
		return "", errors.New("cannot generate candidate key")
	}

	return s(), nil
}

func GenerateSecretToken() (string, error) {
	s, error := nanoid.Standard(20)

	if error != nil {
		return "", errors.New("cannot generate candidate key")
	}

	return s(), nil
}
