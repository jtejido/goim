package algorithm

import (
	"github.com/uhini0201/goim/util"
	"github.com/uhini0201/set"
)

type maxDegreeNode struct {
	Node
	deg float64
}

type MaxDegree struct {
	base
	covQueue *util.PriorityQueue
	graph    *util.Graph
	config   *util.Config
}

func NewMaxDegree(graph *util.Graph, config *util.Config, t int) *MaxDegree {
	md := new(MaxDegree)
	md.graph = graph
	md.config = config
	md.covQueue = util.NewPriorityQueue(func(n1, n2 interface{}) bool {
		if n1.(*maxDegreeNode).deg != n2.(*maxDegreeNode).deg {
			return n2.(*maxDegreeNode).deg < n1.(*maxDegreeNode).deg
		}

		return n1.(*maxDegreeNode).deg < n2.(*maxDegreeNode).deg
	})

	return md
}

func (md *MaxDegree) Select(activated set.Set) set.Set {
	s := set.NewSet()
	seeds := set.NewSet()
	for node := range md.graph.Nodes().Iter() {
		nstruct := new(maxDegreeNode)
		nstruct.id = node.(util.Node)
		nstruct.deg = float64(len(md.graph.Neighbors(node.(util.Node), false)))
		md.covQueue.Push(nstruct)
	}

	for s.Len() < md.config.Seeds {
		nstruct := md.covQueue.Peek().(*maxDegreeNode)
		if !seeds.Contains(nstruct.id) {
			s.Add(nstruct.id)
			seeds.Add(nstruct.id)
		}
		md.covQueue.Pop()
	}
	return s
}
