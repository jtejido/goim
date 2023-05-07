package algorithm

import (
	"math"

	"github.com/uhini0201/goim/model"
	"github.com/uhini0201/goim/util"
	"github.com/uhini0201/grand"
	"github.com/uhini0201/grand/source64"
	"github.com/uhini0201/set"
)

const (
	max_r = 10000000
)

type TIM struct {
	base
	covQueue   *util.PriorityQueue
	graph      *util.Graph
	config     *util.Config
	hyperId    int
	hyperGraph [][]util.Node
	totalR     int
	nodes      []util.Node
	rrSets     [][]util.Node
	seeds      set.Set
	src        grand.Source64
	activated  set.Set
	k          int
	n          int
	nMax       int
	m          int
	epsilon    float64
	t          int
}

func NewTIM(graph *util.Graph, config *util.Config, t int) *TIM {
	c := new(TIM)
	c.graph = graph
	c.config = config
	c.src = source64.NewXoShiRo256StarStar(config.Seed)
	c.seeds = set.NewUnsafeSet()
	c.nodes = make([]util.Node, 0)
	c.rrSets = make([][]util.Node, 0)
	c.hyperGraph = make([][]util.Node, 0)
	c.activated = set.NewUnsafeSet()
	c.k = config.Seeds
	c.t = t
	return c
}

func (c *TIM) Select(activated set.Set) set.Set {
	c.seeds.Clear()
	for node := range activated.Iter() {
		c.activated.Add(node.(util.Node))
	}

	c.n = c.graph.Nodes().Len()
	c.m = c.graph.NumEdges()
	c.hyperId = 0
	c.totalR = 0
	c.epsilon = 0.1

	c.nodes = make([]util.Node, 0)
	c.nMax = 0
	for src := range c.graph.Nodes().Iter() {
		source := src.(util.Node)
		if !(activated.Contains(source)) {
			if int(source) > c.nMax {
				c.nMax = int(source)
			}
			c.nodes = append(c.nodes, source)
		}
	}

	sampler_s := model.NewIndependentCascade(c.graph, c.config, c.t)
	dst := grand.New(c.src)
	var ep_step2, ep_step3 float64
	ep_step3 = c.epsilon
	ep_step2 = 5 * math.Pow(math.Sqrt(ep_step3)/float64(c.config.Seeds), 1./3)
	var ept float64

	ept = c.estimateEPT(sampler_s, dst)

	c.buildSeedSet()
	c.buildHyperGraph2(ep_step2, ept, sampler_s, dst)
	ept = c.influenceHyperGraph()
	ept /= 1 + ep_step2
	c.buildHyperGraph3(ep_step3, ept, sampler_s, dst)
	c.buildSeedSet()
	return c.seeds
}

func (c *TIM) estimateEPT(sampler *model.IndependentCascade, dst *grand.Rand) float64 {
	ept := c.estimateKPT(sampler, dst)
	ept /= 2
	return ept
}

func (c *TIM) estimateKPT(sampler *model.IndependentCascade, dst *grand.Rand) float64 {
	lb := 1. / 2
	var cc float64
	var lastR int
	ret := 1.
	steps := 1.
	for steps <= math.Log(float64(c.n))/math.Log(2)-1 {
		loop := int((6*math.Log(float64(c.n)) + 6*math.Log(math.Log(float64(c.n))/math.Log(2))) / lb)
		cc = 0
		lastR = loop

		for i := 0; i < loop; i++ {
			rr := make([]util.Node, 0)

			seeds := set.NewSet()
			u := c.nodes[dst.Intn(len(c.nodes)-1)]
			seeds.Add(u)
			rr = append(rr, u)
			sampler.Trial(c.activated, seeds, true)
			for _, tt := range sampler.Trials() {
				if tt.Trial == 1 {
					rr = append(rr, tt.Target)
				}
			}

			var mg_tu float64
			for _, node := range rr {
				mg_tu += float64(len(c.graph.Neighbors(node, true)))
			}

			pu := mg_tu / float64(c.m)
			cc += 1 - math.Pow(1-pu, float64(c.k))
		}
		cc /= float64(loop)
		if cc > lb {
			ret = cc * float64(c.n)
			break
		}
		lb /= 2
		steps++
	}
	c.buildSamples(lastR, sampler, dst)
	return ret
}

func (c *TIM) buildSamples(R int, sampler *model.IndependentCascade, dst *grand.Rand) {
	c.totalR += R

	if R > max_r {
		R = max_r
	}

	c.hyperId = R
	c.hyperGraph = make([][]util.Node, c.n)
	c.rrSets = make([][]util.Node, 0)

	var totInDegree float64

	for i := 0; i < R; i++ {

		rr := make([]util.Node, 0)
		seeds := set.NewSet()
		nd := c.nodes[dst.Intn(len(c.nodes)-1)]
		seeds.Add(nd)
		rr = append(rr, nd)
		sampler.Trial(c.activated, seeds, true)

		totInDegree += float64(len(sampler.Trials()))
		for _, tt := range sampler.Trials() {
			if tt.Trial == 1 {
				rr = append(rr, tt.Target)
			}
		}
		c.rrSets = append(c.rrSets, rr)

	}

	for i := 0; i < R; i++ {
		for _, t := range c.rrSets[i] {
			if c.hyperGraph[int(t)] == nil {
				c.hyperGraph[int(t)] = make([]util.Node, 0)
			}
			c.hyperGraph[int(t)] = append(c.hyperGraph[int(t)], util.Node(i))
		}
	}
}

func (c *TIM) buildSeedSet() {
	c.seeds.Clear()
	deg := make([]int, c.n)
	visit_local := make(map[util.Node]struct{}, len(c.rrSets))
	for i := 0; i < len(c.nodes); i++ {
		deg[int(c.nodes[i])] = len(c.hyperGraph[c.nodes[i]])
	}

	for i := 0; i < c.k; i++ {
		t := -1
		id := -1
		for j := 0; j < len(deg); j++ {
			if deg[j] > t {
				t = deg[j]
				id = j
			}
		}

		c.seeds.Add(util.Node(id))
		deg[id] = 0
		for _, t := range c.hyperGraph[id] {
			if _, ok := visit_local[t]; !ok {
				visit_local[t] = struct{}{}
				for _, item := range c.rrSets[t] {
					deg[int(item)]--
				}
			}
		}
	}
}

func (c *TIM) influenceHyperGraph() float64 {
	s := make(map[util.Node]struct{})
	for t := range c.seeds.Iter() {
		for _, tt := range c.hyperGraph[int(t.(util.Node))] {
			s[tt] = struct{}{}
		}
	}

	inf := float64(c.n * len(s) / c.hyperId)
	return inf
}

func (c *TIM) buildHyperGraph2(epsilon_, ept float64, sampler *model.IndependentCascade, dst *grand.Rand) {
	R := (8 + 2*epsilon_) * (float64(c.n)*math.Log(float64(c.n)) + float64(c.n)*math.Log(2)) / (epsilon_ * epsilon_ * ept) / 4
	c.buildSamples(int(R), sampler, dst)
}

func (c *TIM) buildHyperGraph3(epsilon_, opt float64, sampler *model.IndependentCascade, dst *grand.Rand) {
	logCnk := 0.0
	j := 1
	for i := c.n; j <= c.k; i-- {
		logCnk += math.Log10(float64(i)) - math.Log10(float64(j))
		j++
	}
	R := (8 + 2*epsilon_) * (float64(c.n)*math.Log(float64(c.n)) + float64(c.n)*math.Log(2) + float64(c.n)*logCnk) / (epsilon_ * epsilon_ * opt)
	c.buildSamples(int(R), sampler, dst)
}
