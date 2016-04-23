package pso

import (
	"fmt"
	"testing"
)

type an struct{}

func (c an) EvalFitness(v []float64) float64 {
	val := 0.0
	for i := range v {
		val += (v[i] * v[i])
	}
	return val
}

func (c an) Bounderies() ([]float64, []float64) {
	return []float64{5.0, 5.0}, []float64{-5.0, -5.0}
}

func TestAll(t *testing.T) {
	var sol an
	fit, pos := DefaultSolver(sol, 10, 100)
	fmt.Println(fit, pos)
}
