package base

import (
	"github.com/jaevor/go-nanoid"
	"github.com/samber/oops"
)

var nanoidGenerator func() string

func init() {
	var err error
	nanoidGenerator, err = nanoid.Standard(8)
	if err != nil {
		panic(oops.Wrap(err))
	}
}

func NewNanoID() string {
	return nanoidGenerator()
}
