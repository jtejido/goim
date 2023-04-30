package algorithm

import (
	"github.com/uhini0201/goim/util"
	"github.com/uhini0201/set"
)

type discountDegreeNode struct {
	Node
	deg float64
}

type DiscountDegree struct {
	base
	covQueue *util.PriorityQueue
	graph    *util.Graph
	config   *util.Config
}

func NewDiscountDegree(graph *util.Graph, config *util.Config, t int) *DiscountDegree {
	dd := new(DiscountDegree)
	dd.graph = graph
	dd.config = config
	dd.covQueue = util.NewPriorityQueue(func(n1, n2 interface{}) bool {
		if n1.(*discountDegreeNode).deg != n2.(*discountDegreeNode).deg {
			return n2.(*discountDegreeNode).deg < n1.(*discountDegreeNode).deg
		}

		return n1.(*discountDegreeNode).deg < n2.(*discountDegreeNode).deg
	})

	return dd
}

func (dd *DiscountDegree) Select(activated set.Set) set.Set {
	s := set.NewSet()
	queue_nodes := make(map[util.Node]*util.Item)
	for node := range dd.graph.Nodes().Iter() {
		nstruct := new(discountDegreeNode)
		nstruct.id = node.(util.Node)
		if !activated.Contains(nstruct.id) {
			nstruct.deg = 1.
		}
		if dd.graph.Neighbors(node.(util.Node), false) != nil {
			for _, edge := range dd.graph.Neighbors(node.(util.Node), false) {
				if !activated.Contains(edge.Target) {
					nstruct.deg += edge.Dist
				}
			}
		}

		queue_nodes[nstruct.id] = dd.covQueue.Push(nstruct)
	}

	for s.Len() < dd.config.Seeds {
		nstruct := dd.covQueue.Peek().(*discountDegreeNode)
		s.Add(nstruct.id)
		if dd.graph.Neighbors(nstruct.id, false) != nil {
			for _, edge := range dd.graph.Neighbors(nstruct.id, false) {
				if !activated.Contains(edge.Target) && !s.Contains(edge.Target) {
					newN := new(discountDegreeNode)
					newN.id = edge.Target
					newN.deg = queue_nodes[edge.Target].Value().(*discountDegreeNode).deg * (1.0 - edge.Dist)
					dd.covQueue.Update(queue_nodes[edge.Target], newN)
				}
			}
		}
		dd.covQueue.Pop()
	}

	return s
}
