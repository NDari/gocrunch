// Package mat64 supplies functions that create or act
// on 2D slices of float64s, for the Go language.
package mat64

import (
	"encoding/csv"
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

// Ones returns a new 2D slice where all the elements are equal to 1.0
func Ones(r, c int) [][]float64 {
	return Apply(func(i float64) float64 { return 1.0 }, New(r, c))
}

// Inc returns a 2D slice, where element [0][0] == 0.0, and each
// subsequent elemnt is incrmeneted by 1.0
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

// Dump prints the content of a mat64 object to a file, using comma as the
// delimeter between the elements of a row, and a new line between rows.
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

// Load generates a 2D slice of floats from a csv file.
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
// the same shape. This is a deep copy, unlike the builtin copy function
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
		log.Fatalf(errMismatch, "AppendCol")
	}
	for i := 0; i < len(v); i++ {
		m[i] = append(m[i], v[i])
	}
	return m
}

// Print prints a 2D slice of float64s to the std out.
func Print(m [][]float64) {
	w := csv.NewWriter(os.Stdout)
	w.Comma = rune(' ')
	w.WriteAll(ToString(m))
	if err := w.Error(); err != nil {
		log.Fatalf("Error in csv writer to std out:", err)
	}
}
