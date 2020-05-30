package algorithm

import (
	"github.com/jtejido/set"
)

type Algorithm interface {
	Select(activated set.Set) set.Set
}

type base struct {
	Incremental bool
}
