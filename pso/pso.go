/*
pso is an implementation of the Particle Swarm Optimization (PSO) method.

The goal of the pso method is to take a set of such candidates, and through its
algorithm, search a certain space for a better solution. Consider the function
x^2 (x squared). The minimum of this one dimensional function is at 0.0, as no
other value can be given to the function to produce a lower result. The PSO
method attempts to systematically and without any addition input, to find
the solution.

In order to do this, we start by defining the function which we wish to
optimize, as well as define the boundaries of the search space. For instance,
for our problem, we will constrain the search are on the space of real numbers
from -100.0 to 100.0.

Now, we start with a given set of (typically) random guesses of candidates for
where we think the minimum would be (pretending that we do not already know).
Say that we randomly choose the set of candiates to be [1.2, 4.3, 9.7, 72.0].
PSO method will attempt to find the minimum (0.0 in this case) by successive
evaluation of the function, and the interaction between the set of these
candidates.
*/
package pso

import (
	"fmt"
	"math/rand"
)

/*
Candidate is the primary type in this library. In an abstract sense, a candidate
represents a potential solution to a particular problem which we want to optimize.

The candidate must provide a function that takes a []float64 representing
the position of the candidate at each dimension, and return a float64 which
is the fitness of the candidate at that position in configuration space.

The candidate must provide a function that returns the upper and lower
boundaries for each one of the dimensions of the configuration space
of the problem. The boundaries are represented by two float64 slices,
The first of which is the upper limits of each dimension, and the
second is the lower limits of the same dimensions.
*/
type Candidate interface {
	EvalFitness([]float64) float64
	Bounderies() (upper []float64, lower []float64)
}

type Swarm struct {
	Candids  []Candidate
	Pos      [][]float64
	Fit      []float64
	BPos     [][]float64
	BFit     []float64
	V        [][]float64
	GBestFit float64
	GBestID  int
}

func InitSwarm(c []Candidate) *Swarm {
	s := new(Swarm)
	for i := range c {
		s.Candids = append(s.Candids, c[i])
		upper, lower := c[i].Bounderies()
		if len(upper) > len(lower) {
			panic("aw shucks")
		}
		pos := make([]float64, len(upper))
		for i := range pos {
			if upper[i] > lower[i] {
				panic("aw shucks")
			}
			pos[i] = rand.Float64()*(lower[i]-upper[i]) + lower[i]
		}
		fitness := c[i].EvalFitness(pos)
		s.Pos = append(s.Pos, pos)
		s.Fit = append(s.Fit, fitness)
		s.BPos = append(s.BPos, pos)
		s.BFit = append(s.BFit, fitness)
		s.V = append(s.V, make([]float64, len(pos)))
	}
	s.FindGBest()
	return s
}

func (s *Swarm) FindGBest() {
	s.GBestID = 0
	s.GBestFit = s.BFit[0]
	for i := range s.Candids {
		if s.BFit[i] < s.GBestFit {
			s.GBestID = i
			s.GBestFit = s.BFit[i]
		}
	}
}

func (s *Swarm) UpdatePersonalBests() {
	for i := range s.Candids {
		if s.Fit[i] < s.BFit[i] {
			s.BFit[i] = s.Fit[i]
			copy(s.BPos[i], s.Pos[i])
		}
	}
}

func (s *Swarm) UpdateVelocity() {
	for i := range s.Candids {
		for j := range s.V[i] {
			s.V[i][j] = 0.729 * (s.V[i][j] +
				(rand.Float64() * 2.05 * (s.BPos[i][j] - s.Pos[i][j])) +
				(rand.Float64() * 2.05 * (s.BPos[s.GBestID][j] - s.Pos[i][j])))
		}
	}
}

func (s *Swarm) UpdatePos() {
	for i := range s.Candids {
		for j := range s.Pos[i] {
			s.Pos[i][j] += s.V[i][j]
		}
	}
}

func (s *Swarm) CheckBoundaries() {
	for i := range s.Candids {
		upper, lower := s.Candids[i].Bounderies()
		for j := range s.Pos[i] {
			if s.Pos[i][j] > upper[j] {
				s.Pos[i][j] = upper[j]
				s.V[i][j] = 0.0
				continue // skip the checking of the lower boundaries.
			}
			if s.Pos[i][j] < lower[j] {
				s.Pos[i][j] = lower[j]
				s.V[i][j] = 0.0
			}
		}
	}
}

func (s *Swarm) GetFitness() {
	for i := range s.Candids {
		s.Fit[i] = s.Candids[i].EvalFitness(s.Pos[i])
	}
}

func (s *Swarm) Iterate(iteration int) {
	s.UpdateVelocity()
	s.UpdatePos()
	s.CheckBoundaries()
	s.GetFitness()
	s.FindGBest()
	s.UpdatePersonalBests()
	x1 := 0.0
	x2 := 0.0
	for i := range s.Fit {
		x1 += s.Fit[i]
		x2 += s.BFit[i]
	}
	x1 /= float64(len(s.Fit))
	x2 /= float64(len(s.BFit))
	fmt.Println("Finished with iteration", iteration)
	fmt.Sprintf("The global best is candidate %d, with a fitness of %f\n", s.GBestID, s.GBestFit)
	fmt.Println("The average fitness in this iteration is", x1)
	fmt.Println("The average best fitness over all iterations is", x2)
}

func Solve(c []Candidate, numIterations int) {
	fmt.Sprintf("Solving with %d potential solutions over %d iterations\n", len(c), numIterations)
	s := InitSwarm(c)
	for i := 0; i < numIterations; i++ {
		s.Iterate(i)
	}
	fmt.Println("\n=========================================================================\n")
	fmt.Println("\n=========================================================================\n")
	fmt.Println("\n=========================================================================\n")
	fmt.Println("The minimum ftness found is", s.GBestFit)
	fmt.Println("The location of the minimum is as follows:")
	for i := range s.BPos[s.GBestID] {
		fmt.Println("In dimension\t", i, "location\t", s.BPos[s.GBestID][i])
	}
}
