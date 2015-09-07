package mat64

import (
	"log"
)

var (
	errColInx = "Mat64 Error: Column index %d is out of range"
	errRowInx = "Mat64 Error: Row index %d is out of range"
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
		log.Fatalf(errColInx, c)
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
		log.Fatalf(errRowInx, r)
	}
	vec := make([]float64, m.numCols)
	for i := 0; i < m.numCols; i++ {
		vec[i] = m.vals[r*m.numCols+i]
	}
	return vec
}

// Function At returns the values of the entry in an mat64 object at
// the specified row and col. It throws errors if the indeces are out
// of range.
func (m *mat64) At(r, c int) float64 {
	if r >= m.numRows {
		log.Fatalf(errRowInx, r)
	}
	if c > m.numCols {
		log.Fatalf(errColInx, c)
	}
	return m.vals[r*m.numCols+c]
}
