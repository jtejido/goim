package algorithm

import (
	"github.com/jtejido/goim/util"
	"github.com/jtejido/grand"
	"github.com/jtejido/grand/source64"
	"github.com/jtejido/set"
	"sort"
)

const (
	default_r = 250
)

// N. Ohsaka et al., Fast and Accurate Influence Maximization on Large Networks with Pruned Monte-Carlo Simulations, AAAI 2014.
// https://github.com/todo314/pruned-monte-carlo
type PMC struct {
	base
	graph      *util.Graph
	config     *util.Config
	es1, rs1   []util.Node
	at_e, at_r []int
}

func NewPMC(graph *util.Graph, config *util.Config, t int) *PMC {
	c := new(PMC)
	c.graph = graph
	c.config = config
	c.es1 = make([]util.Node, graph.NumEdges())
	c.rs1 = make([]util.Node, graph.NumEdges())
	c.at_e = make([]int, graph.Nodes().Len()+1)
	c.at_r = make([]int, graph.Nodes().Len()+1)
	return c
}

type pair struct {
	x, y util.Node
}

type pairs []pair

func (a pairs) Len() int           { return len(a) }
func (a pairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a pairs) Less(i, j int) bool { return a[i].x < a[j].x }

func (c *PMC) Select(activated set.Set) set.Set {
	seeds := set.NewSet()
	infs := make([]*prunedEstimator, default_r)

	for t := 0; t < default_r; t++ {
		xs := grand.New(source64.NewXoShiRo256StarStar(int64(t) + c.config.Seed))
		mp := 0               // Number of living edges
		ps := make([]pair, 0) // List of reversed living edges
		c.at_e = make([]int, c.graph.Nodes().Len()+1)
		c.at_r = make([]int, c.graph.Nodes().Len()+1)
		for i := 0; i < c.graph.Nodes().Len(); i++ {
			if c.graph.Neighbors(util.Node(i), false) == nil || len(c.graph.Neighbors(util.Node(i), false)) == 0 {
				continue
			}
			for _, edge := range c.graph.Neighbors(util.Node(i), false) {
				if xs.Float64() < edge.Dist {
					c.es1[mp] = edge.Target // Lists of activated nodes (targets)
					mp++
					c.at_e[int(edge.Src)+1]++
					ps = append(ps, pair{edge.Target, edge.Src})
				}
			}
		}
		c.at_e[0] = 0
		sort.Sort(pairs(ps))

		for i := 0; i < mp; i++ {
			c.rs1[util.Node(i)] = ps[i].y
			c.at_r[int(ps[i].x)+1]++
		}
		for i := 1; i <= c.graph.Nodes().Len(); i++ {
			c.at_e[i] += c.at_e[i-1]
			c.at_r[i] += c.at_r[i-1]
		}
		/**
		  Here at_e is the sum of the number of activations from each node.
		  Ex. [0, 2, 5, 8] if node 0 activated two edges, 1 three edges and 2 three
		  edges.
		  Similar thing for at_r;
		*/
		comp := make([]int, c.graph.Nodes().Len())

		nscc := c.scc(comp)

		es2 := set.NewSet() // List of edges in SCC graph, maintains unique pairs
		for u := 0; u < c.graph.Nodes().Len(); u++ {
			a := comp[u]
			for i := c.at_e[u]; i < c.at_e[u+1]; i++ {
				b := comp[int(c.es1[i])]
				if a != b {
					es2.Add(pair{util.Node(a), util.Node(b)})
				}
			}
		}

		es2_b := make([]pair, es2.Len())
		for i, item := range es2.ToSlice() {
			es2_b[i] = item.(pair)
		}

		sort.Sort(pairs(es2_b))

		infs[t] = newPrunedEstimator(nscc, es2_b, comp, activated)
	}

	gain := make([]int64, c.graph.Nodes().Len())
	S := make([]int, 0)

	// Selects greedily seeds
	for t := 0; t < c.config.Seeds; t++ {
		for j := 0; j < default_r; j++ {
			infs[j].update(gain)
		}
		next := 0
		for i := 0; i < c.graph.Nodes().Len(); i++ {
			if gain[i] > gain[next] {
				next = i
			}
		}

		S = append(S, next)
		for j := 0; j < default_r; j++ {
			infs[j].add(next)
		}
		seeds.Add(util.Node(next))
	}

	return seeds
}

func (c *PMC) scc(comp []int) int {
	vis := make([]bool, c.graph.Nodes().Len())
	S := util.NewStack()
	lis := make([]util.Node, 0)
	var k int
	for i := 0; i < c.graph.Nodes().Len(); i++ {
		S.Push(pair{util.Node(i), util.Node(0)})
	}

	for S.Len() > 0 {
		cp := S.Pop().(pair)
		v := cp.x
		state := cp.y
		if state == 0 {
			if vis[v] {
				continue
			}
			vis[v] = true
			S.Push(pair{v, 1})
			for i := c.at_e[v]; i < c.at_e[v+1]; i++ {
				u := c.es1[i]
				S.Push(pair{util.Node(u), util.Node(0)})
			}
		} else {
			lis = append(lis, v)
		}
	}
	for i := 0; i < c.graph.Nodes().Len(); i++ {
		S.Push(pair{lis[i], util.Node(-1)})
	}
	vis = make([]bool, c.graph.Nodes().Len())
	for S.Len() > 0 {
		cp := S.Pop().(pair)
		v := cp.x
		arg := cp.y

		if vis[v] {
			continue
		}

		vis[v] = true

		if int(arg) == -1 {
			comp[v] = k
			k++
		} else {
			comp[v] = int(arg)
		}
		for i := c.at_r[v]; i < c.at_r[v+1]; i++ {
			u := c.rs1[i]
			S.Push(pair{util.Node(u), util.Node(comp[v])})
		}
	}

	return k
}

type prunedEstimator struct {
	n, n1                int
	weight, comp, sigmas []int
	pmoc                 []int
	at_p                 []int
	up                   []int
	memo, removed        []bool
	es, rs               []int
	at_e, at_r           []int
	visited              []bool
	hub                  int
	descendant, ancestor []bool
	flag                 bool
}

func newPrunedEstimator(_n int, _es []pair, _comp []int, activated set.Set) *prunedEstimator {
	pe := new(prunedEstimator)
	pe.flag = true
	pe.n = _n
	pe.n1 = len(_comp)

	pe.visited = make([]bool, pe.n)

	m := len(_es)
	outdeg := make([]int, pe.n)
	indeg := make([]int, pe.n)
	pe.es = make([]int, m)
	pe.rs = make([]int, m)
	for i := 0; i < m; i++ {
		a := _es[i].x
		b := _es[i].y
		outdeg[a]++
		indeg[b]++
		pe.es[i] = -1
		pe.rs[i] = -1
	}

	pe.at_e = make([]int, pe.n+1)
	pe.at_r = make([]int, pe.n+1)

	pe.at_e[0] = 0
	pe.at_r[0] = 0
	for i := 1; i <= pe.n; i++ {
		pe.at_e[i] = pe.at_e[i-1] + outdeg[i-1]
		pe.at_r[i] = pe.at_r[i-1] + indeg[i-1]
	}

	for i := 0; i < m; i++ {
		a := _es[i].x
		b := _es[i].y
		pe.es[pe.at_e[a]] = int(b)
		pe.at_e[a]++
		pe.rs[pe.at_r[b]] = int(a)
		pe.at_r[b]++
	}

	pe.at_e[0] = 0
	pe.at_r[0] = 0
	for i := 1; i <= pe.n; i++ {
		pe.at_e[i] = pe.at_e[i-1] + outdeg[i-1]
		pe.at_r[i] = pe.at_r[i-1] + indeg[i-1]
	}

	pe.sigmas = make([]int, pe.n)

	pe.comp = append([]int{}, _comp...)
	ps := make([]pair, 0)
	for i := 0; i < pe.n1; i++ {
		ps = append(ps, pair{util.Node(pe.comp[i]), util.Node(i)})
	}
	sort.Sort(pairs(ps))
	pe.at_p = make([]int, pe.n+1)
	pe.pmoc = make([]int, 0)
	for i := 0; i < pe.n1; i++ {
		pe.pmoc = append(pe.pmoc, int(ps[i].y))
		pe.at_p[int(ps[i].x)+1]++
	}
	for i := 1; i <= pe.n; i++ {
		pe.at_p[i] += pe.at_p[i-1]
	}

	pe.memo = make([]bool, pe.n)
	pe.removed = make([]bool, pe.n)
	pe.weight = make([]int, pe.n1)

	for i := 0; i < pe.n1; i++ {
		if !activated.Contains(util.Node(i)) {
			pe.weight[pe.comp[i]]++
		}
	}

	pe.up = make([]int, 0)
	pe.first()
	return pe
}

func (pe *prunedEstimator) first() {
	pe.hub = 0
	for i := 0; i < pe.n; i++ {
		if (pe.at_e[i+1]-pe.at_e[i])+(pe.at_r[i+1]-pe.at_r[i]) > (pe.at_e[pe.hub+1]-pe.at_e[pe.hub])+(pe.at_r[pe.hub+1]-pe.at_r[pe.hub]) {
			pe.hub = i
		}
	}

	pe.descendant = make([]bool, pe.n)
	Q := util.NewQueue()
	Q.Push(pe.hub)
	for Q.Len() > 0 {
		// forall v, !remove[v]
		v := Q.Pop().(int)
		pe.descendant[v] = true
		for i := pe.at_e[v]; i < pe.at_e[v+1]; i++ {
			u := pe.es[i]
			if !pe.descendant[u] {
				pe.descendant[u] = true
				Q.Push(u)
			}
		}
	}

	pe.ancestor = make([]bool, pe.n)
	Q.Push(pe.hub)
	for Q.Len() > 0 {
		v := Q.Pop().(int)
		pe.ancestor[v] = true
		for i := pe.at_r[v]; i < pe.at_r[v+1]; i++ {
			u := pe.rs[i]
			if !pe.ancestor[u] {
				pe.ancestor[u] = true
				Q.Push(u)
			}
		}
	}
	pe.ancestor[pe.hub] = false

	for i := 0; i < pe.n; i++ {
		pe.sigma(i)
	}
	pe.ancestor = make([]bool, pe.n)
	pe.descendant = make([]bool, pe.n)

	for i := 0; i < pe.n1; i++ {
		pe.up = append(pe.up, i)
	}
}

func (pe *prunedEstimator) sigma1(v int) int {
	return pe.sigma(pe.comp[v])
}

func (pe *prunedEstimator) add(v0 int) {
	v0 = pe.comp[v0]
	Q := util.NewQueue()
	Q.Push(v0)
	pe.removed[v0] = true
	rm := make([]int, 0)
	for Q.Len() > 0 {
		v := Q.Pop().(int)
		rm = append(rm, v)
		for i := pe.at_e[v]; i < pe.at_e[v+1]; i++ {
			u := pe.es[i]
			if !pe.removed[u] {
				Q.Push(u)
				pe.removed[u] = true
			}
		}
	}

	pe.up = make([]int, 0)

	vec := make([]int, 0)
	for i := 0; i < len(rm); i++ {
		v := rm[i]
		pe.memo[v] = false // for update()
		for j := pe.at_p[v]; j < pe.at_p[v+1]; j++ {
			pe.up = append(pe.up, pe.pmoc[j])
		}
		for j := pe.at_r[v]; j < pe.at_r[v+1]; j++ {
			u := pe.rs[j]
			if !pe.removed[u] && !pe.visited[u] {
				pe.visited[u] = true
				vec = append(vec, u)
				Q.Push(u)
			}
		}
	}
	// reachable to removed node
	for Q.Len() > 0 {
		v := Q.Pop().(int)
		pe.memo[v] = false
		for j := pe.at_p[v]; j < pe.at_p[v+1]; j++ {
			pe.up = append(pe.up, pe.pmoc[j])
		}
		for i := pe.at_r[v]; i < pe.at_r[v+1]; i++ {
			u := pe.rs[i]
			if !pe.visited[u] {
				pe.visited[u] = true
				vec = append(vec, u)
				Q.Push(u)
			}
		}
	}
	for i := 0; i < len(vec); i++ {
		pe.visited[vec[i]] = false
	}
}

func (pe *prunedEstimator) update(sums []int64) {
	for i := 0; i < len(pe.up); i++ {
		v := pe.up[i]
		if !pe.flag {
			sums[v] -= int64(pe.sigmas[pe.comp[v]])
		}
	}
	for i := 0; i < len(pe.up); i++ {
		v := pe.up[i]
		sums[v] += int64(pe.sigma1(v))
	}

	pe.flag = false
}

func (pe *prunedEstimator) sigma(v0 int) int {
	if pe.memo[v0] {
		return pe.sigmas[v0]
	}
	pe.memo[v0] = true
	if pe.removed[v0] {
		pe.sigmas[v0] = 0
		return 0
	} else {
		child := pe.unique_child(v0)
		if child == -1 {
			pe.sigmas[v0] = pe.weight[v0]
			return pe.sigmas[v0]
		} else if child >= 0 {
			pe.sigmas[v0] = pe.sigma(child) + pe.weight[v0]
			return pe.sigmas[v0]
		} else {
			delta := 0
			vec := make([]int, 0)
			pe.visited[v0] = true
			vec = append(vec, v0)
			Q := util.NewQueue()
			Q.Push(v0)
			prune := pe.ancestor[v0]

			if prune {
				delta += pe.sigma(pe.hub)
			}

			for Q.Len() > 0 {
				v := Q.Pop().(int)

				if pe.removed[v] {
					continue
				}
				if prune && pe.descendant[v] {
					continue
				}
				delta += pe.weight[v]
				for i := pe.at_e[v]; i < pe.at_e[v+1]; i++ {
					u := pe.es[i]
					if pe.removed[u] {
						continue
					}
					if !pe.visited[u] {
						pe.visited[u] = true
						vec = append(vec, u)
						Q.Push(u)
					}
				}
			}
			for i := 0; i < len(vec); i++ {
				pe.visited[vec[i]] = false
			}
			pe.sigmas[v0] = delta
			return pe.sigmas[v0]
		}
	}
}

func (pe *prunedEstimator) unique_child(v int) int {
	outdeg := 0
	child := -1
	for i := pe.at_e[v]; i < pe.at_e[v+1]; i++ {
		u := pe.es[i]
		if !pe.removed[u] {
			outdeg++
			child = u
		}
	}
	if outdeg == 0 {
		return -1
	} else if outdeg == 1 {
		return child
	} else {
		return -2
	}
}
