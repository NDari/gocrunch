package mat64

import (
	"log"
)

var (
	errColInx = "Mat64 Error: Column index out of range"
	errRowInx = "Mat64 Error: Column index out of range"
)

type mat64 struct {
	numRows int
	numCols int
	vals    []float64
}

type entry struct {
	value float64
	index int
}

func NewMat(r, c int) *mat64 {
	return &mat64{
		numRows: r,
		numCols: c,
		vals:    make([]float64, r*c),
	}

}

// Function Col returns a slice representing a coloumn
// of a mat64 object.
func (m *mat64) Col(c int) []float64 {
	if c >= m.numCols {
		log.Fatal(errColInx)
	}
	vec := make([]float64, m.numRows)
	for i := 0; i < m.numRows; i++ {
		vec[i] = m.vals[i*m.numCols+c]
	}
	return vec
}

// Function Row returns a slice representing a row
// of a mat64 object.
func (m *mat64) Row(r int) []float64 {
	if r >= m.numRows {
		log.Fatal(errRowInx)
	}
	vec := make([]float64, m.numCols)
	for i := 0; i < m.numCols; i++ {
		vec[i] = m.vals[r*m.numCols+i]
	}
	return vec

}
