// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file
package util

import (
	"github.com/lucky-se7en/grand"
)

type Weighted struct {
	weights []float64
	// heap is a weight heap.
	heap []float64
}

func NewWeighted(w []float64) Weighted {
	s := Weighted{
		weights: make([]float64, len(w)),
		heap:    make([]float64, len(w)),
	}

	s.reweightAll(w)
	return s
}

func (s Weighted) Len() int { return len(s.weights) }

// Take returns an index from the Weighted with probability proportional
// to the weight of the item. The weight of the item is then set to zero.
// Take returns false if there are no items remaining.
func (s Weighted) Take(src grand.Source) (idx int, ok bool) {
	if src == nil {
		panic("Source cannot be nil")
	}
	const small = 1e-12
	if equalWithinAbsOrRel(s.heap[0], 0, small, small) {
		return -1, false
	}

	rnd := grand.New(src)
	r := s.heap[0] * rnd.Float64()

	i := 1
	last := -1
	left := len(s.weights)
	for {
		if r -= s.weights[i-1]; r <= 0 {
			break // Fall within item i-1.
		}
		i <<= 1 // Move to left child.
		if d := s.heap[i-1]; r > d {
			r -= d
			// If enough r to pass left child
			// move to right child state will
			// be caught at break above.
			i++
		}
		if i == last || left < 0 {
			// No progression.
			return -1, false
		}
		last = i
		left--
	}

	w, idx := s.weights[i-1], i-1

	s.weights[i-1] = 0
	for i > 0 {
		s.heap[i-1] -= w
		// The following condition is necessary to
		// handle floating point error. If we see
		// a heap value below zero, we know we need
		// to rebuild it.
		if s.heap[i-1] < 0 {
			s.reset()
			return idx, true
		}
		i >>= 1
	}

	return idx, true
}

func (s Weighted) Reweight(idx int, w float64) {
	w, s.weights[idx] = s.weights[idx]-w, w
	idx++
	for idx > 0 {
		s.heap[idx-1] -= w
		idx >>= 1
	}
}

func (s Weighted) reweightAll(w []float64) {
	if len(w) != s.Len() {
		panic("floats: length of the slices do not match")
	}
	copy(s.weights, w)
	s.reset()
}

func (s Weighted) reset() {
	copy(s.heap, s.weights)
	for i := len(s.heap) - 1; i > 0; i-- {
		s.heap[((i+1)>>1)-1] += s.heap[i]
	}
}
