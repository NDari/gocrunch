// Package Mat64 contains a float64 Matrix object for Go.
package Mat64

import (
	"log"
	"os"
	"strconv"
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

// Col returns a Mat64, representing a single column of the
// original mat64 object at a given location.
func (m *Mat64) Col(c int) *Mat64 {
	if c >= m.NumCols {
		log.Fatalf(errColInx, "Col", c)
	}
	vec := New(m.NumRows, 1)
	for i := 0; i < m.NumRows; i++ {
		vec.Vals[i] = m.Vals[i*m.NumCols+c]
	}
	return vec
}

// Row returns a Mat64, representing a single row of the
// original Mat64 object at the give location.
func (m *Mat64) Row(r int) *Mat64 {
	if r >= m.NumRows {
		log.Fatalf(errRowInx, "Row", r)
	}
	vec := New(1, m.NumCols)
	for i := 0; i < m.NumCols; i++ {
		vec.Vals[i] = m.Vals[r*m.NumCols+i]
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

// Times returns a new matrix that is the result of
// element-wise multiplication of the matrix by another, leaving
// both original matrices intact.
func (m *Mat64) Times(n *Mat64) *Mat64 {
	if m.NumRows != n.NumRows || m.NumCols != m.NumCols {
		log.Fatalf(errMismatch, "Times")
	}
	o := New(m.NumRows, m.NumCols)
	for i := 0; i < m.NumCols*m.NumRows; i++ {
		o.Vals[i] = m.Vals[i] * n.Vals[i]
	}
	return o
}

// TimesInPlace multiplies a Mat64 by another in place. This means that
// the original matrix is lost.
func (m *Mat64) TimesInPlace(n *Mat64) *Mat64 {
	if m.NumRows != n.NumRows || m.NumCols != m.NumCols {
		log.Fatalf(errMismatch, "Times")
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

// Dot is the matrix multiplication of two Mat64 objects. Consider
// The following two mat64 objects, pretty printed for illusration:
//
// A = [[1, 0],
//      [0, 1]]
//
// and
//
// B = [[4, 1],
//      [2, 2]]
//
// A.Dot(B) = [[4, 1],
//             [2, 2]]
//
// The number of elements in the first matrix row, must equal the number
// elements in the second matrix column.
func (m *Mat64) Dot(n *Mat64) *Mat64 {
	if m.NumCols != n.NumRows {
		log.Fatalf(errMismatch, "Dot")
	}
	o := New(m.NumRows, n.NumCols)
	items := m.NumCols
	for i := 0; i < m.NumRows; i++ {
		for j := 0; j < n.NumCols; j++ {
			for k := 0; k < items; k++ {
				o.Vals[i*o.NumRows+j] += (m.At(i, k) * n.At(k, j))
			}
		}
	}
	return o
}

// Reset puts all the elements of a Mat64 values set to 0.0.
func (m *Mat64) Reset() *Mat64 {
	return m.ApplyInPlace(func(i float64) float64 { return 0.0 })
}

// Dump prints the content of a Mat64 object to a file, using the given
// delimeter between the elements of a row, and a new line between rows.
// For instance, giving the comma (",") as a delimiter will essentially
// creates a csv file from the Mat64 object.
func (m *Mat64) Dump(fileName, delemiter string) {
	var str string
	for i := 0; i < m.NumRows; i++ {
		for j := 0; j < m.NumCols; j++ {
			str += strconv.FormatFloat(m.Vals[i*m.NumRows+j], 'f', 14, 64)
			str += delimeter
		}
		if i+1 != m.NumRows {
			str += "\n"
		}
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf(err)
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		log.Fatalf(err)
	}
}
