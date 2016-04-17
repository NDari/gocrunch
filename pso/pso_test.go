package pso

import "testing"

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
	r := []Candidate{
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
		an{},
	}
	Solve(r, 100)
}
