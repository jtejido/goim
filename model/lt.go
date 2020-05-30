package model

import (
	"github.com/jtejido/goim/util"
	"github.com/jtejido/set"
	"github.com/lucky-se7en/grand"
	"github.com/lucky-se7en/grand/source64"
)

type LinearThreshold struct {
	base
	graph  *util.Graph
	config *util.Config
	random *grand.Rand
	src    grand.Source
	t      int
}

func NewLinearThreshold(graph *util.Graph, config *util.Config, t int) *LinearThreshold {
	ret := &LinearThreshold{graph: graph, config: config, random: grand.New(source64.NewXoShiRo256StarStar(config.Seed)), src: source64.NewMT19937(config.Seed)}
	ret.t = t
	return ret
}

func (lt *LinearThreshold) Diffuse(seeds set.Set) set.Set {
	visited := set.NewSet()
	queue := util.NewQueue()
	liveEdges := make(map[util.Node][]util.Node)
	for u := 0; u < lt.graph.Nodes().Len(); u++ {
		index := lt.graph.SampleLivingEdge(util.Node(u), lt.src)
		if index == -1 { // Unconnected node or sample with weights summing to less than 1
			continue
		}
		living_node := lt.graph.Neighbors(util.Node(u), true)[index].Target
		if liveEdges[living_node] == nil {
			liveEdges[living_node] = make([]util.Node, 0)
		}

		liveEdges[living_node] = append(liveEdges[living_node], util.Node(u))
	}

	for source := range seeds.Iter() {
		queue.Push(source.(util.Node))
		visited.Add(source.(util.Node))
	}

	for queue.Len() > 0 {
		node := queue.Peek().(util.Node)
		visited.Add(node)
		if liveEdges[node] != nil {
			for _, neighbour := range liveEdges[node] {
				if !visited.Contains(neighbour) {
					queue.Push(neighbour)
				}
			}
		}
		queue.Pop()
	}

	return visited
}
