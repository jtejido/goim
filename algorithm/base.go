package algorithm

import (
	"github.com/uhini0201/set"
)

type Algorithm interface {
	Select(activated set.Set) set.Set
}

type base struct {
	Incremental bool
}
