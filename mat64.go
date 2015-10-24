// Package mat64 supplies functions that create or act
// on 2D slices of float64s, for the Go language.
package mat64

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// ElementalFn is a function that takes a float64 and returns a
// float64. This function can therefore be applied to each element
// of a 2D float64 slice, and can be used to construct a new one.
type ElementalFn func(float64) float64

// New returns a 2D slice of float64s with the given number of row and columns.
// This function should be used as a convenience tool, and it is exactly
// equivalent to the normal method of allocating a uniform (non-jagged)
// 2D slice of float64s.
//
// If it is anticipated that the 2D slice will grow, use the "NewExpand"
// function below. For full details, read that function's documentation.
func New(r, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := 0; i < r; i++ {
		arr[i] = make([]float64, c)
	}
	return arr
}

// NewExpand returns a 2D slice of float64s, with the given number of rows
// and columns. The difference between this function and the "New" function
// above is that the inner slices are allocated with double the capacity,
// and hence can grow without the need for reallocation up to column * 2.
//
// Note that this extended capacity will waste memory, so the NewExtend
// should be used with care in situations where the performance gained by
// avoiding reallocation justifies the extra cost in memory.
func NewExpand(r, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := 0; i < r; i++ {
		arr[i] = make([]float64, c*2)
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

// Ones returns a new 2D slice where all the elements are equal to 1.0
func Ones(r, c int) [][]float64 {
	return Apply(func(i float64) float64 { return 1.0 }, New(r, c))
}

// Inc returns a 2D slice, where element [0][0] == 0.0, and each
// subsequent element is incremented by 1.0
//
// For example, m := Inc(3, 2) is
//
// [[1.0, 2.0], [3.0, 4.0], [5.0, 6.0]]
func Inc(r, c int) [][]float64 {
	m := New(r, c)
	iter := 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m[i][j] = float64(iter)
			iter++
		}
	}
	return m
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
// mirrored across the diagonal. For example, the element [i][j] becomes the
// element [j][i]. This function leaves the original matrix intact.
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
		log.Fatalf("mat64.%v Error: Row mismatch of the slices", "Times")
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n[i]) {
			log.Fatalf("mat64.%v Error: Col mismatch of the slices at col %v", "Times", i)
		}
		o[i] = make([]float64, len(m[i]))
		for j := 0; j < len(m[i]); j++ {
			o[i][j] = m[i][j] * n[i][j]
		}
	}
	return o
}

// Apply calls a given elemental function on each Element
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
			msg := "mat64.Dot Error: length of column %v on the second matrix\n"
			msg += "is %v, which does not match the length of the row of the \n"
			msg += "first matrix, which is %v"
			log.Fatalf(msg, i, len(n[i]), len(m))
		}
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n) {
			msg := "mat64.Dot Error: length of column %v on the first matrix\n"
			msg += "is %v, which does not match the length of the row of the \n"
			msg += "second matrix, which is %v"
			log.Fatalf(msg, i, len(m[i]), len(n))
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

// Reset sets the values of all entries in a 2D slice of float64s
// to 0.0
func Reset(m [][]float64) [][]float64 {
	return Apply(func(i float64) float64 { return 0.0 }, m)
}

// ToString converts a 2D slice of float64s into a 2D slice
// of strings.
func ToString(m [][]float64) [][]string {
	str := make([][]string, len(m))
	for i := 0; i < len(m); i++ {
		str[i] = make([]string, len(m[i]))
		for j := 0; j < len(m[i]); j++ {
			str[i][j] = strconv.FormatFloat(m[i][j], 'e', 14, 64)
		}
	}
	return str
}

// Dump prints the content of a [][]float64 object to a file, using comma as the
// delimiter between the elements of a row, and a new line between rows.
func Dump(m [][]float64, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Cannot open %v: %v", fileName, err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(ToString(m))
	if err = w.Error(); err != nil {
		log.Fatalf("Error in csv writer for file %v: %v", fileName, err)
	}
}

// FromString converts a 2D slice of strings into a 2D slice of float64s.
func FromString(str [][]string) [][]float64 {
	var err error
	m := make([][]float64, len(str))
	for i := 0; i < len(str); i++ {
		m[i] = make([]float64, len(str[i]))
		for j := 0; j < len(str[i]); j++ {
			m[i][j], err = strconv.ParseFloat(str[i][j], 64)
			if err != nil {
				log.Fatalf("Died on string to float conversion: %v", err)
			}
		}
	}
	return m
}

// Load generates a 2D slice of floats from a CSV file.
func Load(fileName string) [][]float64 {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open %v: %v", fileName, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	str, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error in csv reader for file %v: %v", fileName, err)
	}
	return FromString(str)
}

// Copy copies the content of a 2D slice of float64s into another with
// the same shape. This is a deep copy, unlike the built in copy function
// that is shallow for 2D slices.
func Copy(m [][]float64) [][]float64 {
	n := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		n[i] = make([]float64, len(m[i]))
		copy(n[i], m[i])
	}
	return n
}

// AppendCol appends a column to the right side of a 2D slice of float64s.
func AppendCol(m [][]float64, v []float64) [][]float64 {
	if len(m) != len(v) {
		log.Fatalf("mat64.%v Error: Row mismatch of the slices", "AppendCol")
	}
	for i := 0; i < len(v); i++ {
		m[i] = append(m[i], v[i])
	}
	return m
}

// Concat concatenates the inner slices of 2 2D slices of float64s.
//
// For example, if we have:
//
// m := [[1.0, 2.0], [3.0, 4.0]]
//
// n := [[5.0, 6.0], [7.0, 8.0]]
//
// and
//
// o := mat64.Concat(m, n)
//
// then:
//
// o is [[1.0, 2.0, 5.0, 6.0], [3.0, 4.0, 7.0, 8.0]]
func Concat(m, n [][]float64) [][]float64 {
	if len(m) != len(n) {
		log.Fatalf("mat64.%v Error: Row mismatch of the slices", "Concat")
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		o[i] = make([]float64, len(m[i])+len(n[i]))
		o[i] = append(m[i], n[i]...)
	}
	return o
}

// Print prints a 2D slice of float64s to the stdout.
func Print(m [][]float64) {
	w := csv.NewWriter(os.Stdout)
	w.Comma = rune(' ')
	w.WriteAll(ToString(m))
	if err := w.Error(); err != nil {
		log.Fatalf("Error in csv writer to std out:", err)
	}
}
