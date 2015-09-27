package Mat64

import (
	"log"
)

var (
	errColInx   = "Mat64.%s Error: Column index %d is out of range"
	errRowInx   = "Mat64.%s Error: Row index %d is out of range"
	errMismatch = "Mat64.%s Error: Shape mismatch of the matreces"
)

type Mat64 struct {
	NumRows int
	NumCols int
	Vals    []float64
}

// ElementalFn is a function that takes a float64 and returns a
// float64. This function can therefore be applied to each element
// of a Mat64, and can be used to construct a new transformed Mat64.
type ElementalFn func(float64) float64

// New returns a Mat64 object with the given rows and cols
func New(r, c int) *Mat64 {
	return &Mat64{
		NumRows: r,
		NumCols: c,
		Vals:    make([]float64, r*c),
	}
}

// Identity returns an r by r identity matrix for a given r.
func Identity(r int) *Mat64 {
	identity := New(r, r)
	for i := 0; i < r; i++ {
		identity.Vals[i*r+i] = 1.0
	}
	return identity
}

// Col returns a slice representing a column.
// of a Mat64 object.
func (m *Mat64) Col(c int) []float64 {
	if c >= m.NumCols {
		log.Fatalf(errColInx, "Col", c)
	}
	vec := make([]float64, m.NumRows)
	for i := 0; i < m.NumRows; i++ {
		vec[i] = m.Vals[i*m.NumCols+c]
	}
	return vec
}

// Row returns a slice representing a row
// of a Mat64 object.
func (m *Mat64) Row(r int) []float64 {
	if r >= m.NumRows {
		log.Fatalf(errRowInx, "Row", r)
	}
	vec := make([]float64, m.NumCols)
	for i := 0; i < m.NumCols; i++ {
		vec[i] = m.Vals[r*m.NumCols+i]
	}
	return vec
}

// At returns the values of the entry in an Mat64 object at
// the specified row and col. It throws errors if an index is out
// of range.
func (m *Mat64) At(r, c int) float64 {
	if r >= m.NumRows {
		log.Fatalf(errRowInx, "At", r)
	}
	if c > m.NumCols {
		log.Fatalf(errColInx, "At", c)
	}
	return m.Vals[r*m.NumCols+c]
}

// Transpose returns a copy of a given matrix with the elements
// mirrored across the diagonal. for example, the element At(i, j) becomes the
// element At(j, i). This function leaves the original matrix intact.
func (m *Mat64) Transpose() *Mat64 {
	transpose := New(m.NumCols, m.NumRows)
	for i := 0; i < m.NumRows; i++ {
		for j := 0; j < m.NumCols; j++ {
			transpose.Vals[j*m.NumRows+i] = m.At(i, j)
		}
	}
	return transpose
}

// Equals checks if two mat objects have the same shape and the
// same entries in each row and column.
func (m *Mat64) Equals(n *Mat64) bool {
	if m.NumRows != n.NumRows || m.NumCols != m.NumCols {
		log.Fatalf(errMismatch, "Equals")
	}
	for i := 0; i < m.NumCols*m.NumRows; i++ {
		if m.Vals[i] != n.Vals[i] {
			return false
		}
	}
	return true
}

// Dot returns a new matrix that is the result of
// element-wise multiplication of the matrix by another, leaving
// both original matrices intact.
func (m *Mat64) Dot(n *Mat64) *Mat64 {
	if m.NumRows != n.NumRows || m.NumCols != m.NumCols {
		log.Fatalf(errMismatch, "Dot")
	}
	o := New(m.NumRows, m.NumCols)
	for i := 0; i < m.NumCols*m.NumRows; i++ {
		o.Vals[i] = m.Vals[i] * n.Vals[i]
	}
	return o
}

// DotInPlace multiplies a Mat64 by another in place. This means that
// the original matrix is lost.
func (m *Mat64) DotInPlace(n *Mat64) *Mat64 {
	if m.NumRows != n.NumRows || m.NumCols != m.NumCols {
		log.Fatalf(errMismatch, "Dot")
	}
	for i := 0; i < m.NumCols*m.NumRows; i++ {
		m.Vals[i] *= n.Vals[i]
	}
	return m
}

// Apply calls a given Elemental function on each Element
// of a matrix, returning a new transformed matrix.
func (m *Mat64) Apply(f ElementalFn) *Mat64 {
	n := New(m.NumRows, m.NumCols)
	for i := 0; i < m.NumRows*m.NumCols; i++ {
		n.Vals[i] = f(m.Vals[i])
	}
	return n
}

// ApplyInPlace calls a given Elemental function on each Element
// of a matrix, and then returns the transformed matrix.
func (m *Mat64) ApplyInPlace(f ElementalFn) *Mat64 {
	for i := 0; i < m.NumRows*m.NumCols; i++ {
		m.Vals[i] = f(m.Vals[i])
	}
	return m
}
