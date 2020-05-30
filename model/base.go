package model

import (
	"github.com/jtejido/set"
)

type Model interface {
	Diffuse(seeds set.Set) set.Set
	Type() int
}

type base struct {
	t int
}

func (b *base) Type() int {
	return b.t
}
