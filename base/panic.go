package base

import (
	"github.com/samber/oops"
)

func Recover(err any) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(error); ok {
		return oops.Wrap(e)
	}

	return oops.Errorf("%v", err)
}

func Try(f func()) *Block {
	return NewBlock().Try(f)
}

type Block struct {
	try     func()
	catch   func(error)
	finally func()
}

func NewBlock() *Block {
	return &Block{}
}

func (b *Block) Try(f func()) *Block {
	b.try = f
	return b
}

func (b *Block) Catch(f func(error)) *Block {
	b.catch = f
	return b
}

func (b *Block) Finally(f func()) *Block {
	b.finally = f
	return b
}

func (b *Block) Do() {
	defer func() {
		if err := Recover(recover()); err != nil && b.catch != nil {
			b.catch(err)
		}

		if b.finally != nil {
			b.finally()
		}
	}()

	if b.try != nil {
		b.try()
	}
}
