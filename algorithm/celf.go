package algorithm

import (
	"github.com/uhini0201/goim/model"
	"github.com/uhini0201/goim/util"
	"github.com/uhini0201/set"
)

const (
	default_cascade_samples = 100
)

type celfNode struct {
	Node
	mg float64
}

type CELF struct {
	base
	covQueue *util.PriorityQueue
	graph    *util.Graph
	config   *util.Config
	ic       *model.IndependentCascade
}

func NewCELF(graph *util.Graph, config *util.Config, t int) *CELF {
	c := new(CELF)
	c.graph = graph
	c.config = config
	c.covQueue = util.NewPriorityQueue(func(n1, n2 interface{}) bool {
		if n1.(*celfNode).mg != n2.(*celfNode).mg {
			return n2.(*celfNode).mg < n1.(*celfNode).mg // we want sorting in DESC order
		}

		return n1.(*celfNode).id < n2.(*celfNode).id
	})

	c.ic = model.NewIndependentCascade(graph, config, t)
	return c
}

func (c *CELF) Select(activated set.Set) set.Set {

	s := set.NewSet()

	for node := range c.graph.Nodes().Iter() {
		u := new(celfNode)
		u.id = node.(util.Node)
		seeds := set.NewSet()
		seeds.Add(node.(util.Node))
		u.mg = c.ic.Sample(activated, seeds)
		c.covQueue.Push(u)
	}

	s.Add(c.covQueue.Peek().(*celfNode).id)
	c.covQueue.Pop()
	for s.Len() < c.config.Seeds {
		var found bool
		for !found {
			u := c.covQueue.Peek().(*celfNode)
			c.covQueue.Pop()
			seeds := set.NewSet()
			for node := range s.Iter() {
				seeds.Add(node.(util.Node))
			}

			seeds.Add(u.id)
			prev_val := u.mg
			u.mg = c.ic.Sample(activated, seeds) - prev_val
			if u.mg >= c.covQueue.Peek().(*celfNode).mg {
				s.Add(u.id)
				found = true
			} else {
				c.covQueue.Push(u)
			}
		}
	}

	return s
}
