package model

import (
	"github.com/jtejido/goim/util"
	"github.com/jtejido/set"
	"github.com/lucky-se7en/grand"
	"github.com/lucky-se7en/grand/source64"
)

type IndependentCascade struct {
	base
	graph   *util.Graph
	config  *util.Config
	random  *grand.Rand
	samples int
	trials  []util.TrialType
}

func NewIndependentCascade(graph *util.Graph, config *util.Config, t int) *IndependentCascade {
	ret := &IndependentCascade{graph: graph, config: config, random: grand.New(source64.NewXoShiRo256StarStar(config.Seed)), trials: make([]util.TrialType, 0)}
	ret.t = t
	return ret
}

func (ic *IndependentCascade) Sample(activated, seeds set.Set) float64 {
	return ic.sample(activated, seeds, false, false)
}

func (ic *IndependentCascade) Trial(activated, seeds set.Set, inv bool) float64 {
	return ic.sample(activated, seeds, true, inv)
}

func (ic *IndependentCascade) Trials() []util.TrialType {
	return ic.trials
}

func (ic *IndependentCascade) sample(activated, seeds set.Set, trial, inv bool) float64 {
	var outspread float64
	ic.trials = make([]util.TrialType, 0)
	samples := ic.config.Simulations
	if trial {
		samples = 1
	}

	for sample := 1; sample <= samples; sample++ {
		var reached_round float64
		queue := util.NewQueue()
		active := set.NewSet()
		for source := range seeds.Iter() {
			ss := source.(util.Node)
			queue.Push(ss)
			active.Add(ss)
		}

		for queue.Len() > 0 {
			node_id := queue.Peek().(util.Node)
			ic.sampleOutGoingEdges(node_id, queue, active, trial, inv)
			queue.Pop()
			if !activated.Contains(node_id) {
				reached_round++
			}
		}

		outspread += reached_round
	}

	return outspread / float64(samples)
}

func (ic *IndependentCascade) Diffuse(seeds set.Set) set.Set {
	active := set.NewSet()
	queue := util.NewQueue()

	for source := range seeds.Iter() {
		queue.Push(source)
		active.Add(source.(util.Node))
	}

	for queue.Len() > 0 {
		node_id := queue.Peek().(util.Node)
		ic.sampleOutGoingEdges(node_id, queue, active, false, false)
		queue.Pop()
	}

	return active
}

func (ic *IndependentCascade) sampleOutGoingEdges(node util.Node, queue *util.Queue, active set.Set, trial, inv bool) {
	neighborList := ic.graph.Neighbors(node, inv)
	if neighborList != nil {
		for _, edge := range neighborList {
			act := 0
			if !active.Contains(edge.Target) {
				if ic.random.Float64() <= edge.Dist {
					active.Add(edge.Target)
					queue.Push(edge.Target)
					act = 1
				}

				if trial {
					ic.trials = append(ic.trials, util.TrialType{node, edge.Target, act})
				}
			}
		}
	}

}
