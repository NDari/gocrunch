# pso
--
    import "github.com/NDari/gocrunch/pso"

## Usage

A simple example for the most basic usage is below:

```go
package main

import (
  "fmt"
  "gocrunch/pso"
)

// Define an empty struct, so that we can define the two methods we
// need on it.
type an struct{}

// Define the EvalFitness method on the above empty struct. This method
// will evaluate the fit of a certain location in configuration space.
// For this example, EvalFitness will evaluate the sum of the square of
// the position in each dimension.
func (c an) EvalFitness(v []float64) float64 {
	val := 0.0
	for i := range v {
		val += (v[i] * v[i])
	}
	return val
}

// Define the Boundaries function. This function return the min and max
// range for each dimension. For this example, we will go from -5 to 5
// in both dimensions.
func (c an) Bounderies() ([]float64, []float64) {
	return []float64{5.0, 5.0}, []float64{-5.0, -5.0}
}

func main() {
	var sol an
	numCandidates := 10 // How many possible solution we search in each iteration
	numIterations := 100 // The number of iterations to carry out the optimization

	// optimize this function with the default solver, which is fast, and
	// has some sensible default choices. We will print the best fitness and
	// its location after the optimization is concluded.
	fit, pos := pso.DefaultSolver(sol, numCandidates, numIterations)
	fmt.Println(fit, pos)
}
```

## Documentation

Full documentation is under badges, below.

## Badges

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
[![GoDoc](https://godoc.org/github.com/NDari/gocrunch/pso?status.svg)](https://godoc.org/github.com/NDari/gocrunch/pso)
