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
Say that we randomly choose the set of candidates to be [1.2, 4.3, 9.7, 72.0].
PSO method will attempt to find the minimum (0.0 in this case) by successive
evaluation of the function, and the interaction between the set of these
candidates.
*/
package pso

import (
	"fmt"
	"math"
	"math/rand"
)

/*
Candidate represents a potential solution to a particular problem which
we want to optimize.

The candidate must provide a function that takes a []float64 representing
the position of the candidate at each dimension, and return a float64 which
is the fitness of the candidate at that position in configuration space.

The candidate must also provide a function that returns the upper and lower
boundaries for each one of the dimensions of the configuration space
of the problem. The boundaries are represented by two float64 slices,
The first of which is the upper limits of each dimension, and the
second is the lower limits of the same dimensions.
*/
type Candidate interface {
	EvalFitness(position []float64) (fitness float64)
	Bounderies() (upper []float64, lower []float64)
}

/*
Swarm is the primary data structure of the PSO package. It represents a set
of candidates (potential solutions), which it randomly scatters around the
configuration space of the problem.

The Swarm is also responsible for all the bookkeeping needed in a typical
PSO run, such as the current and best position of each candidate solution,
the current global best solution, and so on.

Finally, the Swarm also contains the various settings and configurations for
the implementation of the PSO algorithm to use, the neighborhood topology,
the social and cognitive acceleration coefficients, and so on. All of the
various settings and methods of the Swarm can be adjusted to implement a
wide range of the various PSO algorithms with minimal work.
*/
type Swarm struct {
	candids  []Candidate
	pos      [][]float64
	fit      []float64
	bPos     [][]float64
	bFit     []float64
	v        [][]float64
	gBestFit float64
	gBestPos []float64
	gBestID  int
	target   []int

	c1               float64
	c2               float64
	w                float64 // w = (0.9 - 0.4) * ((maxiter-iter)/maxiter) + 0.4
	psoType          string
	topology         string
	numIterations    int
	currentIteration int
	verbose          bool
}

/*
The DefaultSolver is a collection of sensible preset configurations for a PSO
implementation. For a large number of cases, using this solver will be
sufficient. For the cases where higher performance is needed, the user can
tinker with the various settings themselves.

The settings used in this solver are as follows:

- psoType: "Constriction"
- topology: "Global"
- social acceleration weight: 2.05
- cognitive acceleration weight: 2.05
- initial velocity: 0.0 in all dimensions.

*/
func DefaultSolver(sol Candidate, nCandids, nIters int) (float64, []float64) {
	var c []Candidate
	for i := 0; i < nCandids; i++ {
		c = append(c, sol)
	}
	s := InitSwarm(c, nIters)
	s.RunIterations()
	fmt.Println("==============================================================")
	fmt.Println("==============================================================")
	fmt.Println("==============================================================")
	fmt.Println("The minimum fitness found is", s.gBestFit)
	fmt.Println("The location of the minimum is as follows:")
	for i := range s.bPos[s.gBestID] {
		fmt.Println("In dimension\t", i, "location\t", s.bPos[s.gBestID][i])
	}
	return s.gBestFit, s.gBestPos
}

func InitSwarm(c []Candidate, numIterations int) *Swarm {
	s := new(Swarm)

	s.candids = make([]Candidate, len(c))
	s.pos = make([][]float64, len(c))
	s.bPos = make([][]float64, len(c))
	s.v = make([][]float64, len(c))
	s.fit = make([]float64, len(c))
	s.bFit = make([]float64, len(c))
	s.target = make([]int, len(c))

	for i, candidate := range c {
		s.candids[i] = candidate
		upper, lower := candidate.Bounderies()
		if len(upper) != len(lower) {
			panic("aw shucks")
		}

		pos := make([]float64, len(upper))
		s.pos[i] = make([]float64, len(upper))
		s.bPos[i] = make([]float64, len(upper))
		s.v[i] = make([]float64, len(upper))

		for i := range pos {
			if upper[i] < lower[i] {
				panic("aw shucks")
			}
			pos[i] = rand.Float64()*(lower[i]-upper[i]) + upper[i]
		}
		fitness := candidate.EvalFitness(pos)
		copy(s.pos[i], pos)
		s.fit[i] = fitness
		copy(s.bPos[i], pos)
		s.bFit[i] = fitness
	}
	s.FindGBest()
	s.gBestPos = make([]float64, len(s.bPos[s.gBestID]))
	for i := range s.bPos[s.gBestID] {
		s.gBestPos[i] = s.bPos[s.gBestID][i]
	}
	s.c1 = 2.05
	s.c2 = 2.05
	s.w = 0.9
	s.psoType = "Constriction"
	s.topology = "Global"
	s.numIterations = numIterations
	s.currentIteration = 0
	s.verbose = true
	return s
}

func (s *Swarm) FindGBest() {
	s.gBestID = 0
	s.gBestFit = s.bFit[0]
	for i := range s.candids {
		if s.bFit[i] < s.gBestFit {
			s.gBestID = i
			s.gBestFit = s.bFit[i]
		}
	}
	copy(s.gBestPos, s.bPos[s.gBestID])
}

func (s *Swarm) RunIterations() {
	for s.currentIteration < s.numIterations {
		s.Iterate()
		s.currentIteration++
	}
}

func (s *Swarm) Iterate() {
	s.UpdateTargets()
	s.UpdateVelocity()
	s.UpdatePos()
	s.CheckBoundaries()
	s.GetFitness()
	s.UpdatePersonalBests()
	s.FindGBest()
	if s.verbose {
		x1 := 0.0
		x2 := 0.0
		for i := range s.fit {
			x1 += s.fit[i]
			x2 += s.bFit[i]
		}
		x1 /= float64(len(s.fit))
		x2 /= float64(len(s.bFit))
		fmt.Println("Finished with iteration", s.currentIteration)
		fmt.Println("The global best is", s.gBestID, "with a fitness of", s.gBestFit)
		fmt.Println("The average fitness in this iteration is", x1)
		fmt.Println("The average best fitness over all iterations is", x2)
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
}

func (s *Swarm) UpdateTargets() {
	switch s.topology {
	case "Global":
		for i := range s.target {
			s.target[i] = s.gBestID
		}
	case "Ring":
		for i := 0; i < len(s.target)-1; i++ {
			s.target[i] = i + 1
		}
		s.target[len(s.target)-1] = 0
	case "Von Neuman":
		panic("Von Neumann topology not yet implemented")
	case "Random":
		for i := range s.target {
			// Find a random target, redo if the target is the candidate itself.
			redo := true
			for redo {
				target := rand.Intn(len(s.candids))
				if target == i {
					continue
				}
				// if the fitness of the target is higher, then we accelerate
				// away from it, and not toward it as the normal PSO algorithm.
				// For this reason, we will assign the target to be negative
				// so that we know to move away from it.
				if s.bFit[target] > s.bFit[i] {
					s.target[i] = -target
				} else {
					s.target[i] = target
				}
				redo = false
			}
		}
		panic("Random topology not yet implemented")
	default:
		panic("Unknown topology requested")
	}
}

func (s *Swarm) UpdateVelocity() {
	switch s.psoType {
	case "Constriction":
		phi := s.c1 + s.c2
		chi := (2.0 / math.Abs(2.0-phi-math.Sqrt((phi*phi)-(4.0*phi))))
		for i := range s.candids {
			for j := range s.v[i] {
				t := s.target[i]
				if t < 0 {
					t = -t // set it back to positive for indexing
					s.v[i][j] = chi * (s.v[i][j] +
						(rand.Float64() * s.c1 * (s.bPos[i][j] - s.pos[i][j])) +
						(rand.Float64() * s.c1 * (s.bPos[t][j] + s.pos[i][j])))
				} else {
					s.v[i][j] = chi * (s.v[i][j] +
						(rand.Float64() * s.c1 * (s.bPos[i][j] - s.pos[i][j])) +
						(rand.Float64() * s.c1 * (s.bPos[t][j] - s.pos[i][j])))
				}
			}
		}
	case "Standard":
		panic("Standard PSO algorithm not yet implemented")
	default:
		panic("Requested PSO type is not implemented")
	}
}

func (s *Swarm) UpdatePos() {
	for i := range s.candids {
		for j := range s.pos[i] {
			s.pos[i][j] += s.v[i][j]
		}
	}
}

func (s *Swarm) CheckBoundaries() {
	for i := range s.candids {
		upper, lower := s.candids[i].Bounderies()
		for j := range s.pos[i] {
			if s.pos[i][j] > upper[j] {
				s.pos[i][j] = upper[j]
				s.v[i][j] = 0.0
			}
			if s.pos[i][j] < lower[j] {
				s.pos[i][j] = lower[j]
				s.v[i][j] = 0.0
			}
		}
	}
}

func (s *Swarm) GetFitness() {
	for i := range s.candids {
		s.fit[i] = s.candids[i].EvalFitness(s.pos[i])
	}
}

func (s *Swarm) UpdatePersonalBests() {
	for i := range s.candids {
		if s.fit[i] < s.bFit[i] {
			s.bFit[i] = s.fit[i]
			copy(s.bPos[i], s.pos[i])
		}
	}
}
