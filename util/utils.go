package util

import (
	"bufio"
	"fmt"
	"math"

	"github.com/uhini0201/set"
)

func LogSeed(round, activated int, roundtime, timetotal float64, seeds set.Set, config *Config, bufferedWriter *bufio.Writer) {
	seedStr := SeedToLog(round, activated, roundtime, seeds)
	bufferedWriter.WriteString(seedStr + "\n")
}

func SeedToLog(round, activated int, roundtime float64, seeds set.Set) (s string) {
	s += fmt.Sprintf("%d", round) + "\t"
	s += fmt.Sprintf("%d", activated) + "\t"
	s += fmt.Sprintf("%.5f", roundtime) + "\t"
	s += "["
	i := 1
	for ss := range seeds.Iter() {
		s += fmt.Sprintf("%d", ss.(Node))
		if i != seeds.Len() {
			s += ", "
		}
		i++

	}
	s += "]"
	return
}

func equalWithinAbs(a, b, tol float64) bool {
	return a == b || math.Abs(a-b) <= tol
}

const minNormalFloat64 = 2.2250738585072014e-308

func equalWithinRel(a, b, tol float64) bool {
	if a == b {
		return true
	}
	delta := math.Abs(a - b)
	if delta <= minNormalFloat64 {
		return delta <= tol*minNormalFloat64
	}

	return delta/math.Max(math.Abs(a), math.Abs(b)) <= tol
}

func equalWithinAbsOrRel(a, b, absTol, relTol float64) bool {
	if equalWithinAbs(a, b, absTol) {
		return true
	}
	return equalWithinRel(a, b, relTol)
}
