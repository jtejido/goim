package util

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
	"time"
)

const (
	str_celf string = "celf"
	str_tim  string = "tim"
	str_md   string = "maxdegree"
	str_dd   string = "discountdegree"
	str_pmc  string = "pmc"
	str_ic   string = "ic"
	str_lt   string = "lt"
)

func ToAlgorithm(a string) Algorithm {
	switch strings.ToLower(a) {
	case str_celf:
		return CELF
	case str_tim:
		return TIM
	case str_md:
		return MAX_DEGREE
	case str_dd:
		return DISCOUNT_DEGREE
	case str_pmc:
		return PMC
	default:
		panic("not supported")
	}
}

func ToDiffusionModel(a string) DiffusionModel {
	switch strings.ToLower(a) {
	case str_ic:
		return IC
	case str_lt:
		return LT
	default:
		panic("not supported")
	}
}

type (
	Algorithm      int
	DiffusionModel int
)

const (
	CELF Algorithm = iota
	TIM
	MAX_DEGREE
	DISCOUNT_DEGREE
	PMC
)

const (
	IC DiffusionModel = iota
	LT
)

func (a Algorithm) String() string {
	switch a {
	case CELF:
		return strings.ToUpper(str_celf)
	case TIM:
		return strings.ToUpper(str_tim)
	case MAX_DEGREE:
		return strings.ToUpper(str_md)
	case DISCOUNT_DEGREE:
		return strings.ToUpper(str_dd)
	case PMC:
		return strings.ToUpper(str_pmc)
	default:
		panic("not supported")
	}
}

func (a DiffusionModel) String() string {
	switch a {
	case IC:
		return strings.ToUpper(str_ic)
	case LT:
		return strings.ToUpper(str_lt)
	default:
		panic("not supported")
	}
}

// This is the base Config type for the API. Extend as needed.
type Config struct {
	OutputDir   string `toml:"outputDir"`
	GraphPath   string `toml:"graphPath"`
	Trials      int    `toml:"trials"`
	Algorithm   string `toml:"algorithm"`
	Seeds       int    `toml:"seeds"`
	Model       string `toml:"model"`
	Simulations int    `toml:"simulations"`
	Seed        int64  `toml:"seed"`
}

func LoadConfig(filename string) (*Config, error) {
	var c Config

	if filename == "" {
		return nil, fmt.Errorf("cannot open an empty config filepath")
	}

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) LogFileName() (s string) {
	s += c.OutputDir + "/" // put the log file under the output path
	s += c.GraphPath[strings.LastIndexAny(c.GraphPath, "/")+1:strings.LastIndexAny(c.GraphPath, ".")] + "_"
	s += strings.ToLower(ToAlgorithm(c.Algorithm).String()) + "_"
	s += fmt.Sprintf("%d", c.Trials) + "_"
	s += fmt.Sprintf("%d", c.Seeds) + "_"
	s += fmt.Sprintf("%d", c.Seed) + "_"
	s += makeTimestampStr() + ".log"
	return
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeTimestampStr() string {
	return fmt.Sprintf("%d", makeTimestamp())
}
