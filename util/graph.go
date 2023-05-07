package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/uhini0201/grand"
	"github.com/uhini0201/set"
)

const separator string = " "

type Edge struct {
	Src    Node
	Target Node
	Dist   float64
}

type Graph struct {
	nodes        set.Set
	neighbors    map[Node][]Edge
	invNeighbors map[Node][]Edge
	ltDist       map[Node]Weighted
}

func NewGraph(graphFilePath string) (g *Graph, err error) {
	f, err := os.Open(graphFilePath)
	if err != nil {
		return
	}

	defer f.Close()
	g = new(Graph)
	g.nodes = set.NewSet()
	g.neighbors = make(map[Node][]Edge)
	g.invNeighbors = make(map[Node][]Edge)
	g.ltDist = make(map[Node]Weighted)
	br := bufio.NewReader(f)
	var numEdges int
	maxNodeId := math.MinInt32

	log.Printf("Reading graph file from %s \n", graphFilePath)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		fields := strings.Fields(string(line))
		if len(fields) < 3 {
			return nil, fmt.Errorf("Invalid graph file format.")
		}
		// each line contains one directed edge: (u, v, p_uv)
		u, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		v, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		p, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		// TO-DO. internal ID that always starts with 0.
		nu := Node(u)
		nv := Node(v)
		numEdges++
		g.addEdge(nu, nv, p)

		if u > maxNodeId {
			maxNodeId = u
		}

		if v > maxNodeId {
			maxNodeId = v
		}
	}

	for u := 0; u < g.Nodes().Len(); u++ {
		if g.invNeighbors[Node(u)] != nil { // Only reversed edges are interesting
			neighbours := g.invNeighbors[Node(u)]
			w := make([]float64, len(neighbours)+1)
			var total float64
			for i := 0; i < len(neighbours); i++ {
				cur_weight := neighbours[Node(i)].Dist
				total += cur_weight
				w[i] = cur_weight
			}
			if total < 1 { // Weights do not sum to 1, we can sample no edge
				w[len(neighbours)] = 1 - total
			}

			g.ltDist[Node(u)] = NewWeighted(w[:len(w)-1])
		}
	}

	log.Printf("Number of nodes = %d \n", g.nodes.Len())
	log.Printf("Number of edges = %d \n", numEdges)
	log.Println("Finished reading graph file!")
	log.Printf("Max node id = %d \n", maxNodeId)
	return g, nil
}

func (g *Graph) addEdge(src, tgt Node, dist float64) {
	g.addNode(src)
	g.addNode(tgt)
	if g.neighbors[src] == nil {
		g.neighbors[src] = make([]Edge, 0)
	}

	if g.invNeighbors[tgt] == nil {
		g.invNeighbors[tgt] = make([]Edge, 0)
	}

	g.neighbors[src] = append(g.neighbors[src], Edge{src, tgt, dist})
	g.invNeighbors[tgt] = append(g.invNeighbors[tgt], Edge{tgt, src, dist})
}

func (g *Graph) addNode(n Node) {
	g.nodes.Add(n)
}

func (g *Graph) Nodes() set.Set {
	return g.nodes
}

func (g *Graph) NumEdges() int {
	return len(g.neighbors)
}

func (g *Graph) Neighbors(node Node, inv bool) []Edge {
	if inv {
		return g.invNeighbors[node]
	}

	return g.neighbors[node]
}

func (g *Graph) SampleLivingEdge(node Node, src grand.Source) int {
	if g.invNeighbors[node] != nil {
		index, ok := g.ltDist[node].Take(src)
		if ok {
			if index < len(g.invNeighbors[node]) {
				return index
			}
		}
	}

	return -1
}
