// Package mat64 contains a float64 Matrix object for Go.
package mat64

import (
	"log"
	"os"
	"strconv"
)

var (
	errMismatch = "mat64.%s Error: Shape mismatch of the slices"
)

// ElementalFn is a function that takes a float64 and returns a
// float64. This function can therefore be applied to each element
// of a 2D float64 slice, and can be used to construct a new one.
type ElementalFn func(float64) float64

// New returns a 2D slice of float64s with the given row and columns.
func New(r, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := 0; i < r; i++ {
		arr[i] = make([]float64, c)
	}
	return arr
}

// I returns an r by r identity matrix for a given r.
func I(r int) [][]float64 {
	identity := New(r, r)
	for i := 0; i < r; i++ {
		identity[i][i] = 1.0
	}
	return identity
}

// Col returns a column of a 2D slice of float64s.
func Col(c int, m [][]float64) []float64 {
	vec := make([]float64, len(m))
	for r := 0; r < len(m); r++ {
		vec[r] = m[r][c]
	}
	return vec
}

// Row returns a row of a 2D slice of float64s
func Row(r int, m [][]float64) []float64 {
	vec := make([]float64, len(m[r]))
	for c := 0; c < len(m[r]); c++ {
		vec[c] = m[r][c]
	}
	return vec
}

// T returns a copy of a given matrix with the elements
// mirrored across the diagonal. for example, the element At(i, j) becomes the
// element At(j, i). This function leaves the original matrix intact.
func T(m [][]float64) [][]float64 {
	transpose := New(len(m[0]), len(m))
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			transpose[j][i] = m[i][j]
		}
	}
	return transpose
}

// Equals checks if two mat objects have the same shape and the
// same entries in each row and column.
func Equal(m, n [][]float64) bool {
	if len(m) != len(n) {
		return false
	}
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n[i]) {
			return false
		}
		for j := 0; j < len(m[i]); j++ {
			if m[i][j] != n[i][j] {
				return false
			}
		}
	}
	return true
}

// Times returns a new 2D slice that is the result of
// element-wise multiplication of two 2D slices.
func Times(m, n [][]float64) [][]float64 {
	if len(m) != len(n) {
		log.Fatalf(errMismatch, "Times")
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n[i]) {
			log.Fatalf(errMismatch, "Times")
		}
		o[i] = make([]float64, len(m[i]))
		for j := 0; j < len(m[i]); j++ {
			o[i][j] = m[i][j] * n[i][j]
		}
	}
	return o
}

// Apply calls a given Elemental function on each Element
// of a 2D slice, returning it afterwards.
func Apply(f ElementalFn, m [][]float64) [][]float64 {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			m[i][j] = f(m[i][j])
		}
	}
	return m
}

// Dot is the matrix multiplication of two 2D slices of float64s

func Dot(m, n [][]float64) [][]float64 {
	lenm := len(m)
	// make sure that the length of the row of m matches the length of
	// each column in n.
	for i := 0; i < len(n); i++ {
		if lenm != len(n[i]) {
			log.Fatalf(errMismatch, "Dot")
		}
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n) {
			log.Fatalf(errMismatch, "Dot")
		}
		o[i] = make([]float64, len(n[0]))
		for j := 0; j < len(m[i]); j++ {
			for k := 0; k < len(n); k++ {
				o[i][j] += m[i][k] * n[k][j]
			}
		}
	}
	return o
}

// Reset puts all the elements of a mat64 values set to 0.0.
func Reset(m [][]float64) [][]float64 {
	return Apply(func(i float64) float64 { return 0.0 }, m)
}

// Dump prints the content of a mat64 object to a file, using the given
// delimeter between the elements of a row, and a new line between rows.
// For instance, giving the comma (",") as a delimiter will essentially
// creates a csv file from the mat64 object.
func Dump(m [][]float64, fileName, delimiter string) {
	var str string
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			str += strconv.FormatFloat(m[i][j], 'f', 14, 64)
			str += delimiter
		}
		if i+1 != len(m) {
			str += "\n"
		}
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Cannot open %v: %v", fileName, err)
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		log.Fatalf("Cannot write to %v: %v", fileName, err)
	}
}
