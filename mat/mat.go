/*
Package mat implements function that create or act upon 2D slices of
`float64`. This is in essence the same concept of a matrix in other
languages.

The 2D slices acted on or created by the functions below are assumed to
be non-jagged. This means that for a given [][]float64, the inner slices
are assumed to be of the same length.
*/
package mat

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
)

type mat struct {
	r, c int
	vals []float64
	work []float64
}

type elementFunc func(float64) float64
type booleanFunc func(float64) bool
type reducerFunc func(float64, float64) float64

func New(r, c int) *mat {
	if r <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of rows must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "New", r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	if c <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of columns must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "New", c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	return &mat{
		r,
		c,
		make([]float64, r*c),
		make([]float64, r*c),
	}
}

func Bare(r, c int) *mat {
	if r <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of rows must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "Bare", r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	if c <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of columns must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "Bare", c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	return &mat{
		r,
		c,
		make([]float64, r*c),
		nil,
	}
}

func FromSlice(s [][]float64, bare bool) *mat {
	if isJagged(s) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, a 'jagged' 2D slice was recieved, where the rows\n"
		s += "(the inner slice of the 2D slice) have different lengths. The"
		s += "creation of a *mat from jagged slices is not supported.\n"
		s = fmt.Sprintf(s, "FromSlice")
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	// We start with a Bare mat. if the "bare" parameter is "false", then we
	// will manually allocate the mat.work slice.
	m := Bare(len(s), len(s[0]))
	if !bare {
		m.work = make([]float64, len(s)*len(s[0]))
	}
	for i := range s {
		for j := range s[i] {
			m.vals[i*len(s[0])+j] = s[i][j]
		}
	}
	return m
}

func isJagged(s [][]float64) bool {
	c := len(s[0])
	for i := range s {
		if len(s[i]) != c {
			return true
		}
	}
	return false
}

func FromCSV(filename string) *mat {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV", filename, err)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// I am going with the assumption that a mat loaded from a CSV is going to
	// be large. So, we are going to read one line, and determine the number
	// of coloumns based on the number of comma separated strings in that line.
	// Then we will read the rest of the lines one at a time, checking that the
	// number of entries in each line is the same as the first line, and
	// incrementing the number of rows each time.
	//
	// Another thing to note is that since we are assuming that the file is
	// large, I am going to create a Bare mat.
	str, err := r.Read()
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot read from %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV", filename, err)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	line := 1
	// In order to use "append", we must allocate an empty slice ( [] ).
	// However this is not allowed in Bare and New by design, as you
	// cannot pass zeros for mat.r nor mat.c. So here we allocate and
	// then set it equal to nil for the desired effect.
	m := Bare(1, len(str))
	m.vals = nil
	row := make([]float64, len(str))

	for {
		for i := range str {
			row[i], err = strconv.ParseFloat(str[i], 64)
			if err != nil {
				fmt.Println("\nNumgo/mat error.")
				s := "In mat.%v, item %d in line %d is %s, which cannot\n"
				s += "be converted to a float64 due to: %v"
				s = fmt.Sprintf(s, "FromCSV", i, line, str[i], err)
				fmt.Println(s)
				fmt.Println("Stack trace for this error:\n")
				debug.PrintStack()
				os.Exit(1)
			}
		}
		m.vals = append(m.vals, row...)
		// Read the next line. If there is one, increment the number of rows
		str, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, cannot read from %s due to error: %v.\n"
			s = fmt.Sprintf(s, "FromCSV", filename, err)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:\n")
			debug.PrintStack()
			os.Exit(1)
		}
		line += 1
		if len(str) != len(row) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, line %d in %s has %d entries. The first line\n"
			s += "(line 1) has %d entries.\n"
			s += "Creation of a *mat from jagged slices is not supported.\n"
			s = fmt.Sprintf(s, "Load", filename, err)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:\n")
			debug.PrintStack()
			os.Exit(1)
		}
		m.r += 1
	}
	return m
}

func (m *mat) ToSlice() [][]float64 {
	s := make([][]float64, m.r)
	for i := 0; i < m.r; i++ {
		s[i] = make([]float64, m.c)
		for j := 0; j < m.c; j++ {
			s[i][j] = m.vals[i*m.c+j]
		}
	}
	return s
}

func (m *mat) Map(f elementFunc) {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = f(m.vals[i])
	}
	return
}

// /*
// NewExpand returns a 2D slice of `float64`, with the given number of rows
// and columns. The difference between this function and the "New" function
// above is that the inner slices are allocated with double the capacity,
// and hence can grow without the need for reallocation up to column * 2.

// Note that this extended capacity will waste memory, so the NewExtend
// should be used with care in situations where the performance gained by
// avoiding reallocation justifies the extra cost in memory.
// */
// func NewExpand(r, c int) [][]float64 {
// 	arr := make([][]float64, r)
// 	for i := 0; i < r; i++ {
// 		arr[i] = make([]float64, c*2)
// 	}
// 	return arr
// }

// /*
// I returns an r by r 2D slice for a given r, where the elements along
// the diagonal (where the first and the second index are equal) is set
// to `1.0`, and all other elements are set to `0.0`.
// */
// func I(r int) [][]float64 {
// 	identity := New(r, r)
// 	for i := 0; i < r; i++ {
// 		identity[i][i] = 1.0
// 	}
// 	return identity
// }

// /*
// Ones returns a new 2D slice where all the elements are equal to `1.0`.
// */
// func Ones(r, c int) [][]float64 {
// 	f := func(i float64) float64 {
// 		return 1.0
// 	}
// 	return Map(f, New(r, c))
// }

// /*
// Inc returns a 2D slice, where element `[0][0] == 0.0`, and each
// subsequent element is incremented by `1.0`.

// For example:

// m := Inc(3, 2)

// mat.Print(m) // 1.0, 2.0
//              // 3.0, 4.0
// 			 // 5.0, 6.0
// */
// func Inc(r, c int) [][]float64 {
// 	m := New(r, c)
// 	iter := 0
// 	for i := 0; i < r; i++ {
// 		for j := 0; j < c; j++ {
// 			m[i][j] = float64(iter)
// 			iter++
// 		}
// 	}
// 	return m
// }

// /*
// Col returns a column of a 2D slice of `float64`. Col uses a 0-based index,
// hence the first column of a 2D slice, m,  is `Col(0, m)`.

// This function also allows for negative indexing. For example, `Col(-1, m)`
// is the last column of the 2D slice m, and `Col(-2, m)` is the second to
// last column of m, and so on.
// */
// func Col(c int, m [][]float64) []float64 {
// 	if (c >= len(m[0])) || (-c > len(m[0])) {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s the requested column %d is outside of the bounds [-%d, %d]\n"
// 		s = fmt.Sprintf(s, "Col", c, len(m[0]), len(m[0])-1)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	vec := make([]float64, len(m))
// 	if c >= 0 {
// 		for r := 0; r < len(m); r++ {
// 			vec[r] = m[r][c]
// 		}
// 	} else if c < 0 {
// 		lenColM := len(m[0])
// 		for r := range m {
// 			vec[r] = m[r][lenColM+c]
// 		}
// 	}
// 	return vec
// }

// /*
// Row returns a row of a 2D slice of `float64`. Row uses a 0-based index, hence
// the first row of a 2D slice, m, is Row(0, m).

// This function also allows for negative indexing. For example, Row(-1, m) is
// the last row of m.
// */
// func Row(r int, m [][]float64) []float64 {
// 	if (r >= len(m)) || (-r > len(m)) {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s the requested column %d is outside of the bounds [-%d, %d]\n"
// 		s = fmt.Sprintf(s, "Row", r, len(m[0]), len(m[0])-1)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	if r >= 0 {
// 		return m[r]
// 	} else {
// 		return m[len(m)+r]
// 	}
// }

// /*
// T returns a copy of a given 2D slice with the elements of the 2D slice
// mirrored across the diagonal. For example, the element `[i][j]` becomes the
// element `[j][i]` of the returned 2D slice. This function leaves the
// original matrix intact.
// */
// func T(m [][]float64) [][]float64 {
// 	transpose := New(len(m[0]), len(m))
// 	for i := 0; i < len(m); i++ {
// 		for j := range m[i] {
// 			transpose[j][i] = m[i][j]
// 		}
// 	}
// 	return transpose
// }

// /*
// Equal checks if two 2D slices have the same shape and the same entries in
// each row and column. If either the shape or the entries of the arguments
// are different, `false` is returned. Otherwise, the return value is `true`.
// */
// func Equal(m, n [][]float64) bool {
// 	if len(m) != len(n) {
// 		return false
// 	}
// 	for i := range m {
// 		if len(m[i]) != len(n[i]) {
// 			return false
// 		}
// 		for j := range m[i] {
// 			if m[i][j] != n[i][j] {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// /*
// Mul returns a new 2D slice that is the result of element-wise multiplication
// of two 2D slices.
// */
// func Mul(m, n [][]float64) [][]float64 {
// 	if len(m) != len(n) {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s the number of rows of the first 2D slice is %d, while\n"
// 		s += "the number of rows of the second 2D slice is %d. They must be equal\n"
// 		s = fmt.Sprintf(s, "Mul", len(m), len(n))
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	o := make([][]float64, len(m))
// 	for i := range m {
// 		o[i] = vec.Mul(m[i], n[i])
// 	}
// 	return o
// }

// /*
// MapInPlace calls a given elemental function on each Element of a 2D slice, returning
// it afterwards. This function modifies the original 2D slice.
// */
// func MapInPlace(f func(float64) float64, m [][]float64) {
// 	for i := range m {
// 		vec.MapInPlace(f, m[i])
// 	}
// 	return
// }

// /*
// Map calls a given elemental function on each Element of a 2D slice, returning
// a new 2D slice. This function does not modify the original 2D slice.
// */
// func Map(f func(float64) float64, m [][]float64) [][]float64 {
// 	n := make([][]float64, len(m))
// 	for i := range m {
// 		n[i] = vec.Map(f, m[i])
// 	}
// 	return n
// }

// /*
// Dot is the matrix multiplication of two 2D slices of `float64`.
// */
// func Dot(m, n [][]float64) [][]float64 {
// 	lenm := len(m)
// 	// make sure that the length of the row of m matches the length of
// 	// each column in n.
// 	for i := range n {
// 		if lenm != len(n[i]) {
// 			fmt.Println("\nNumgo/mat error.")
// 			s := "In mat.%s the length of column %d os the second 2D slice is %d\n"
// 			s += "which is not equal to the number of rows of the first 2D slice,\n"
// 			s += "which is %d. They must be equal.\n"
// 			s = fmt.Sprintf(s, "Mul", i, len(n[i]), lenm)
// 			fmt.Println(s)
// 			fmt.Println("Stack trace for this error:\n")
// 			debug.PrintStack()
// 			os.Exit(1)
// 		}
// 	}
// 	o := make([][]float64, len(m))
// 	for i := range m {
// 		if len(m[i]) != len(n) {
// 			fmt.Println("\nNumgo/mat error.")
// 			s := "In mat.%s the length of column %d os the first 2D slice is %d\n"
// 			s += "which is not equal to the number of rows of the 2nd 2D slice,\n"
// 			s += "which is %d. They must be equal.\n"
// 			s = fmt.Sprintf(s, "Mul", i, len(m[i]), len(n))
// 			fmt.Println(s)
// 			fmt.Println("Stack trace for this error:\n")
// 			debug.PrintStack()
// 			os.Exit(1)
// 		}
// 		o[i] = make([]float64, len(n[0]))
// 		for j := range m[i] {
// 			for k := range n {
// 				o[i][j] += m[i][k] * n[k][j]
// 			}
// 		}
// 	}
// 	return o
// }

// /*
// Reset sets the values of all entries in a 2D slice of `float64` to `0.0`.
// */
// func Reset(m [][]float64) {
// 	f := func(i float64) float64 {
// 		return 0.0
// 	}
// 	MapInPlace(f, m)
// 	return
// }

// /*
// ToString converts a `[][]float64` to `[][]string`.
// */
// func ToString(m [][]float64) [][]string {
// 	str := make([][]string, len(m))
// 	for i := range m {
// 		str[i] = make([]string, len(m[i]))
// 		for j := range m[i] {
// 			str[i][j] = strconv.FormatFloat(m[i][j], 'e', 14, 64)
// 		}
// 	}
// 	return str
// }

// /*
// Dump prints the content of a `[][]float64` slice to a file, using comma as the
// delimiter between the elements of a row, and a new line between rows.
// */
// func Dump(m [][]float64, fileName string) {
// 	f, err := os.Create(fileName)
// 	if err != nil {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%v, cannot open %s due to error: %v.\n"
// 		s = fmt.Sprintf(s, "Dump", fileName, err)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	defer f.Close()
// 	w := csv.NewWriter(f)
// 	w.WriteAll(ToString(m))
// 	if err = w.Error(); err != nil {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%v, cannot write to %s due to error: %v.\n"
// 		s = fmt.Sprintf(s, "Dump", fileName, err)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// }

// /*
// FromString converts a `[][]string` to `[][]float64`.
// */
// func FromString(str [][]string) [][]float64 {
// 	var err error
// 	m := make([][]float64, len(str))
// 	for i := range str {
// 		m[i] = make([]float64, len(str[i]))
// 		for j := 0; j < len(str[i]); j++ {
// 			m[i][j], err = strconv.ParseFloat(str[i][j], 64)
// 			if err != nil {
// 				fmt.Println("\nNumgo/mat error.")
// 				s := "In mat.%v, item %d in line %d (0-indices) due to: %v.\n"
// 				s = fmt.Sprintf(s, "FromString", i, j, err)
// 				fmt.Println(s)
// 				fmt.Println("Stack trace for this error:\n")
// 				debug.PrintStack()
// 				os.Exit(1)
// 			}
// 		}
// 	}
// 	return m
// }

// /*
// Load generates a 2D slice of floats from a CSV file.
// */
// func Load(fileName string) [][]float64 {
// 	f, err := os.Open(fileName)
// 	if err != nil {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%v, cannot open %s due to error: %v.\n"
// 		s = fmt.Sprintf(s, "Load", fileName, err)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	defer f.Close()
// 	r := csv.NewReader(f)
// 	str, err := r.ReadAll()
// 	if err != nil {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%v, cannot read from %s due to error: %v.\n"
// 		s = fmt.Sprintf(s, "Load", fileName, err)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	return FromString(str)
// }

// /*
// Copy copies the content of a 2D slice of float64 into another with
// the same shape. This is a deep copy, unlike the built in copy function
// that is shallow for 2D slices.
// */
// func Copy(m [][]float64) [][]float64 {
// 	n := make([][]float64, len(m))
// 	for i := range m {
// 		n[i] = make([]float64, len(m[i]))
// 		copy(n[i], m[i])
// 	}
// 	return n
// }

// /*
// AppendCol appends a column to the right side of a 2D slice of float64s.
// */
// func AppendCol(m [][]float64, v []float64) [][]float64 {
// 	if len(m) != len(v) {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s the number of rows of the 2D slice is %d, while\n"
// 		s += "the number of rows of the vector is %d. They must be equal.\n"
// 		s = fmt.Sprintf(s, "AppendCol", len(m), len(v))
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	for i := range v {
// 		m[i] = append(m[i], v[i])
// 	}
// 	return m
// }

// /*
// Concat concatenates the inner slices of two `[][]float64` arguments..

// For example, if we have:

// m := [[1.0, 2.0], [3.0, 4.0]]
// n := [[5.0, 6.0], [7.0, 8.0]]
// o := mat.Concat(m, n)

// mat.Print(o) // 1.0, 2.0, 5.0, 6.0
//              // 3.0, 4.0, 7.0, 8.0

// */
// func Concat(m, n [][]float64) [][]float64 {
// 	if len(m) != len(n) {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s the number of rows of the first 2D slice is %d, while\n"
// 		s += "the number of rows of the second 2D slice is %d. They must be equal.\n"
// 		s = fmt.Sprintf(s, "Concat", len(m), len(n))
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// 	o := make([][]float64, len(m))
// 	for i := range m {
// 		o[i] = make([]float64, len(m[i])+len(n[i]))
// 		o[i] = append(m[i], n[i]...)
// 	}
// 	return o
// }

// /*
// Print prints a `[][]float64` to the stdout.
// */
// func Print(m [][]float64) {
// 	w := csv.NewWriter(os.Stdout)
// 	w.Comma = rune(' ')
// 	w.WriteAll(ToString(m))
// 	if err := w.Error(); err != nil {
// 		fmt.Println("\nNumgo/mat error.")
// 		s := "In mat.%s cannot write to stdout due to: %v\n"
// 		s = fmt.Sprintf(s, "Print", err)
// 		fmt.Println(s)
// 		fmt.Println("Stack trace for this error:\n")
// 		debug.PrintStack()
// 		os.Exit(1)
// 	}
// }
