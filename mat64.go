package mat64

import (
	"log"
)

var (
	errColInx   = "mat64.%s Error: Column index %d is out of range"
	errRowInx   = "mat64.%s Error: Row index %d is out of range"
	errMismatch = "mat64.%s Error: Shape mismatch of the matreces"
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

// ElementalFn is a function that takes a float64 and returns a
// float64. This function can therefore be applied to each element
// of a mat64, and can be used to contruct a new trasformed mat64.
type ElementalFn func(float64) float64

// New returns a mat64 object with the given rows and cols
func New(r, c int) *mat64 {
	return &mat64{
		numRows: r,
		numCols: c,
		vals:    make([]float64, r*c),
	}
}

// Identity returns an r by r identity matrix for a given r.
func Identity(r int) *mat64 {
	identity := New(r, r)
	for i := 0; i < r; i++ {
		identity.vals[i*r+i] = 1.0
	}
	return identity
}

// Col returns a slice representing a coloumn
// of a mat64 object.
func (m *mat64) Col(c int) []float64 {
	if c >= m.numCols {
		log.Fatalf(errColInx, "Col", c)
	}
	vec := make([]float64, m.numRows)
	for i := 0; i < m.numRows; i++ {
		vec[i] = m.vals[i*m.numCols+c]
	}
	return vec
}

// Row returns a slice representing a row
// of a mat64 object.
func (m *mat64) Row(r int) []float64 {
	if r >= m.numRows {
		log.Fatalf(errRowInx, "Row", r)
	}
	vec := make([]float64, m.numCols)
	for i := 0; i < m.numCols; i++ {
		vec[i] = m.vals[r*m.numCols+i]
	}
	return vec
}

// At returns the values of the entry in an mat64 object at
// the specified row and col. It throws errors if the indeces are out
// of range.
func (m *mat64) At(r, c int) float64 {
	if r >= m.numRows {
		log.Fatalf(errRowInx, "At", r)
	}
	if c > m.numCols {
		log.Fatalf(errColInx, "At", c)
	}
	return m.vals[r*m.numCols+c]
}

// Transpose returns a copy of a given matrix with the elements
// mirriored across the diagonal. i.e. and element At(i, j) becomes the
// element At(j, i). This function leaves the original matrix intact.
func (m *mat64) Transpose() *mat64 {
	transpose := New(m.numCols, m.numRows)
	for i := 0; i < m.numRows; i++ {
		for j := 0; j < m.numCols; j++ {
			transpose.vals[j*m.numRows+i] = m.At(i, j)
		}
	}
	return transpose
}

// Equals checks if two mat objects have the same shape and the
// same entries in each row and column.
func (m *mat64) Equals(n *mat64) bool {
	if m.numRows != n.numRows || m.numCols != m.numCols {
		log.Fatalf(errMismatch, "Equals")
	}
	for i := 0; i < m.numCols*m.numRows; i++ {
		if m.vals[i] != n.vals[i] {
			return false
		}
	}
	return true
}

// Dot returns a new matrix that is the result of
// element-wise multiplication of the matrix by another, leaving
// both original matrices intact.
func (m *mat64) Dot(n *mat64) *mat64 {
	if m.numRows != n.numRows || m.numCols != m.numCols {
		log.Fatalf(errMismatch, "Dot")
	}
	o := New(m.numRows, m.numCols)
	for i := 0; i < m.numCols*m.numRows; i++ {
		o.vals[i] = m.vals[i] * n.vals[i]
	}
	return o
}

// DotInPlace multiplies a Mat64 by another in place. This means that
// the original matrix is lost.
func (m *mat64) DotInPlace(n *mat64) *mat64 {
	if m.numRows != n.numRows || m.numCols != m.numCols {
		log.Fatalf(errMismatch, "Dot")
	}
	for i := 0; i < m.numCols*m.numRows; i++ {
		m.vals[i] *= n.vals[i]
	}
	return m
}

// Apply calls a given Elemental function on each Element
// a a matrix, returning a new trasformed matrix.
func (m *mat64) Apply(f ElementalFn) *mat64 {
	n := New(m.numRows, m.numCols)
	for i := 0; i < m.numRows*m.numCols; i++ {
		n.vals[i] = f(m.vals[i])
	}
	return n
}

// ApplyInPlace calls a given Elemental function on each Element
// a a matrix, and then returns the transformed matrix.
func (m *mat64) ApplyInPlace(f ElementalFn) *mat64 {
	for i := 0; i < m.numRows*m.numCols; i++ {
		m.vals[i] = f(m.vals[i])
	}
	return m
}
