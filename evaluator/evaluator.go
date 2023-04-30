package evaluator

import (
	"bufio"
	"log"
	"time"

	"github.com/uhini0201/goim/algorithm"
	"github.com/uhini0201/goim/model"
	"github.com/uhini0201/goim/util"
	"github.com/uhini0201/set"
)

const (
	INFLUENCE_MED int = iota
	INFLUENCE_UPPER
	INFLUENCE_ADAPTIVE
	INFLUENCE_UCB
	INFLUENCE_THOMPSON
)

type Evaluator struct {
	config    *util.Config
	graph     *util.Graph
	algorithm algorithm.Algorithm
	model     model.Model
	writer    *bufio.Writer
}

func NewEvaluator(config *util.Config, graph *util.Graph, bufferedWriter *bufio.Writer) *Evaluator {
	var algo algorithm.Algorithm
	if util.ToAlgorithm(config.Algorithm) == util.CELF {
		algo = algorithm.NewCELF(graph, config, INFLUENCE_MED)
	} else if util.ToAlgorithm(config.Algorithm) == util.TIM {
		algo = algorithm.NewTIM(graph, config, INFLUENCE_MED)
	} else if util.ToAlgorithm(config.Algorithm) == util.DISCOUNT_DEGREE {
		algo = algorithm.NewDiscountDegree(graph, config, INFLUENCE_MED)
	} else if util.ToAlgorithm(config.Algorithm) == util.MAX_DEGREE {
		algo = algorithm.NewMaxDegree(graph, config, INFLUENCE_MED)
	} else if util.ToAlgorithm(config.Algorithm) == util.PMC {
		algo = algorithm.NewPMC(graph, config, INFLUENCE_MED)
	}

	var m model.Model
	if util.ToDiffusionModel(config.Model) == util.IC {
		m = model.NewIndependentCascade(graph, config, INFLUENCE_MED)
	} else if util.ToDiffusionModel(config.Model) == util.LT {
		m = model.NewLinearThreshold(graph, config, INFLUENCE_MED)
	}

	return &Evaluator{config, graph, algo, m, bufferedWriter}
}

func (e *Evaluator) Run() error {
	activated := set.NewSet()
	var roundtime, timetotal float64
	log.Printf("Algorithm: %s \n", util.ToAlgorithm(e.config.Algorithm).String())
	log.Printf("Model: %s \n", util.ToDiffusionModel(e.config.Model).String())
	log.Printf("Output: %s", e.config.LogFileName())
	for stage := 1; stage <= e.config.Trials; stage++ {
		t0 := makeTimestamp()
		seeds := e.algorithm.Select(activated)
		diffusion := e.model.Diffuse(seeds)

		for node := range diffusion.Iter() {
			activated.Add(node.(util.Node))
		}

		t1 := makeTimestamp()

		timetotal += float64(t1-t0) / (1000.0 * 60.0)
		roundtime = float64(t1-t0) / (1000.0 * 60.0)
		util.LogSeed(stage, activated.Len(), roundtime, timetotal, seeds, e.config, e.writer)
		if err := e.writer.Flush(); err != nil {
			return err
		}
	}

	log.Printf("Time elapsed: %.5f \n", timetotal)
	return nil
}
func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
