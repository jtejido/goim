package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/uhini0201/goim/evaluator"
	"github.com/uhini0201/goim/util"
)

var (
	conf       *util.Config
	cpuprofile string
	logFile    string
	confFile   string
)

func init() {
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to location")
	flag.StringVar(&logFile, "log", "", "write log to location")
	flag.StringVar(&confFile, "conf", "config.toml", "config file location")

	conf, _ = util.LoadConfig(confFile)
	flag.StringVar(&conf.OutputDir, "output", conf.OutputDir, "Path for output files.")
	flag.StringVar(&conf.GraphPath, "graph", conf.GraphPath, "Path of graph file.")
	flag.Int64Var(&conf.Seed, "seed", conf.Seed, "Seed of rng.")
	flag.IntVar(&conf.Trials, "trials", conf.Trials, "Number of trials.")
	flag.StringVar(&conf.Algorithm, "algorithm", conf.Algorithm, "Seed-selection algorithm.")
	flag.IntVar(&conf.Seeds, "seeds", conf.Seeds, "Number of seeds in each trial.")
	flag.StringVar(&conf.Model, "model", conf.Model, "Diffusion model to use.")
}

func main() {
	flag.Parse()

	if logFile != "" {
		lf, err := os.Create(logFile)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(lf)
	}

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("fatal: \n", err)
		}
	}()

	run()
}

func run() {
	graph, err := util.NewGraph(conf.GraphPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	logFileName := conf.LogFileName()
	os.Mkdir(conf.OutputDir, os.ModeDir)
	f, err := os.Create(logFileName)
	bw := bufio.NewWriter(f)
	eval := evaluator.NewEvaluator(conf, graph, bw)
	log.Println("Running Evaluator")
	if err := eval.Run(); err != nil {
		log.Fatal(err.Error())
	}
	f.Close()
}
