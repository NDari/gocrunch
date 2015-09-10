package mat64

import (
	"log"
)

var (
	errColInx = "Mat64 Error: Column index %d is out of range"
	errRowInx = "Mat64 Error: Row index %d is out of range"
)

type mat64 struct {
	numRows uint
	numCols uint
	vals    []float64
}

type entry struct {
	value float64
	index uint
}

// Function New returns a mat64 object with the given rows and cols
func New(r, c uint) *mat64 {
	return &mat64{
		numRows: r,
		numCols: c,
		vals:    make([]float64, r*c),
	}
}

func I(r uint) *mat64 {
	identity := New(r, r)
	for i := 0; i < int(r); i++ {
		identity.vals[i*int(r)+i] = 1.0
	}
	return identity
}

// Function Col returns a slice representing a coloumn
// of a mat64 object.
func (m *mat64) Col(c uint) []float64 {
	if c >= m.numCols {
		log.Fatalf(errColInx, c)
	}
	vec := make([]float64, m.numRows)
	for i := 0; i < int(m.numRows); i++ {
		vec[i] = m.vals[i*int(m.numCols)+int(c)]
	}
	return vec
}

// Function Row returns a slice representing a row
// of a mat64 object.
func (m *mat64) Row(r uint) []float64 {
	if r >= m.numRows {
		log.Fatalf(errRowInx, r)
	}
	vec := make([]float64, m.numCols)
	for i := 0; i < int(m.numCols); i++ {
		vec[i] = m.vals[int(r)*int(m.numCols)+i]
	}
	return vec
}

// Function At returns the values of the entry in an mat64 object at
// the specified row and col. It throws errors if the indeces are out
// of range.
func (m *mat64) At(r, c uint) float64 {
	if r >= m.numRows {
		log.Fatalf(errRowInx, r)
	}
	if c > m.numCols {
		log.Fatalf(errColInx, c)
	}
	return m.vals[r*m.numCols+c]
}
