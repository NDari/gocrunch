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
}

type elementFunc func(*float64)
type booleanFunc func(*float64) bool
type reducerFunc func(accumulator *float64, next float64)

func New(r, c int) *mat {
	if r <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of rows must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "New", r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if c <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of columns must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "New", c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	return &mat{
		r,
		c,
		make([]float64, r*c),
	}
}

func FromSlice(s [][]float64) *mat {
	if s == nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, an unallocated 2D slice was recieved, where the slice\n"
		s += "is equal to nil."
		s = fmt.Sprintf(s, "FromSlice")
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if isJagged(s) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, a 'jagged' 2D slice was recieved, where the rows\n"
		s += "(the inner slice of the 2D slice) have different lengths. The\n"
		s += "creation of a *mat from jagged slices is not supported.\n"
		s = fmt.Sprintf(s, "FromSlice")
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	m := New(len(s), len(s[0]))
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

func From1DSlice(s []float64) *mat {
	if s == nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, an unallocated 1D slice was recieved, where the slice\n"
		s += "is equal to nil."
		s = fmt.Sprintf(s, "From1DSlice")
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	m := New(1, len(s))
	copy(m.vals, s)
	return m
}

func FromCSV(filename string) *mat {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV", filename, err)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
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
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	line := 1
	// In order to use "append", we must allocate an empty slice ( [] ).
	// However this is not allowed in mat.New by design, as you
	// cannot pass zeros for mat.r nor mat.c. So here we allocate and
	// then set it equal to nil for the desired effect.
	m := New(1, len(str))
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
				fmt.Println("Stack trace for this error:")
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
			fmt.Println("Stack trace for this error:")
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
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		m.r += 1
	}
	return m
}

func (m *mat) Reshape(rows, cols int) *mat {
	if rows <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of rows must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "Reshape", rows)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if cols <= 0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, the number of columns must be greater than '0', but\n"
		s += "recieved %d. "
		s = fmt.Sprintf(s, "Reshape", cols)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if rows*cols != m.r*m.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, The total number of entries of the old and new shape\n"
		s += "must match.\n"
		s = fmt.Sprintf(s, "Reshape")
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	} else {
		m.r = rows
		m.c = cols
	}
	return m
}

func (m *mat) Dims() (int, int) {
	return m.r, m.c
}

func (m *mat) Vals() []float64 {
	s := make([]float64, m.r*m.c)
	copy(s, m.vals)
	return s
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

func (m *mat) Map(f elementFunc) *mat {
	for i := 0; i < m.r*m.c; i++ {
		f(&m.vals[i])
	}
	return m
}

func (m *mat) Ones() *mat {
	f := func(i *float64) {
		*i = 1.0
		return
	}
	return m.Map(f)
}

func (m *mat) Inc() *mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = float64(i)
	}
	return m
}

func (m *mat) Reset() *mat {
	f := func(i *float64) {
		*i = 0.0
		return
	}
	return m.Map(f)
}

func (m *mat) Col(x int) *mat {
	if (x >= m.c) || (x < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of the bounds [0, %d]\n"
		s = fmt.Sprintf(s, "Col", x, m.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	v := New(m.r, 1)
	for r := 0; r < m.r; r++ {
		v.vals[r] = m.vals[r*m.c+x]
	}
	return v
}

func (m *mat) Row(x int) *mat {
	if (x >= m.r) || (x < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested row %d is outside of the bounds [0, %d]\n"
		s = fmt.Sprintf(s, "Row", x, m.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	v := New(1, m.c)
	for r := 0; r < m.c; r++ {
		v.vals[r] = m.vals[x*m.c+r]
	}
	return v
}

func (m *mat) Equals(n *mat) bool {
	if m.r != n.r {
		return false
	}
	if m.c != n.c {
		return false
	}
	for i := 0; i < m.r*m.c; i++ {
		if m.vals[i] != n.vals[i] {
			return false
		}
	}
	return true
}

func (m *mat) Copy() *mat {
	n := New(m.r, m.c)
	copy(n.vals, m.vals)
	return n
}

func (m *mat) T() *mat {
	n := New(m.c, m.r)
	idx := 0
	for i := 0; i < m.c; i++ {
		for j := 0; j < m.r; j++ {
			n.vals[idx] = m.vals[j*m.c+i]
			idx += 1
		}
	}
	return n
}

func (m *mat) Filter(f booleanFunc) *mat {
	var res []float64
	for i := 0; i < m.r*m.c; i++ {
		if f(&m.vals[i]) {
			res = append(res, m.vals[i])
		}
	}
	if len(res) == 0 {
		return nil
	} else {
		n := From1DSlice(res)
		return n
	}
}

func (m *mat) All(f booleanFunc) bool {
	for i := 0; i < m.r*m.c; i++ {
		if !f(&m.vals[i]) {
			return false
		}
	}
	return true
}

func (m *mat) Any(f booleanFunc) bool {
	for i := 0; i < m.r*m.c; i++ {
		if f(&m.vals[i]) {
			return true
		}
	}
	return false
}

func (m *mat) CombineWith(n *mat, how reducerFunc) *mat {
	if m.r != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first mat is %d\n"
		s += "but the number of rows of the second mat is %d. They must\n"
		s += "match for combination.\n"
		s = fmt.Sprintf(s, "Combine", m.r, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if m.c != n.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first mat is %d\n"
		s += "but the number of columns of the second mat is %d. They must\n"
		s += "match for combination.\n"
		s = fmt.Sprintf(s, "Combine", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	for i := 0; i < m.r*m.c; i++ {
		how(&m.vals[i], n.vals[i])
	}
	return m
}

func (m *mat) Mul(n *mat) *mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] *= n.vals[i]
	}
	return m
}

func (m *mat) Add(n *mat) *mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] += n.vals[i]
	}
	return m
}

func (m *mat) Sub(n *mat) *mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] -= n.vals[i]
	}
	return m
}

func (m *mat) Div(n *mat) *mat {
	zero := func(i *float64) bool {
		if *i == 0.0 {
			return true
		}
		return false
	}
	if n.Any(zero) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, one or more elements of the second matrix are 0.0\n"
		s += "Division by zero is not allowed.\n"
		s = fmt.Sprintf(s, "Div", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] /= n.vals[i]
	}
	return m
}

func (m *mat) Scale(f float64) *mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] *= f
	}
	return m
}

func (m *mat) Sum() {
	return
}

func (m *mat) Average() {
	return
}

func (m *mat) Prod() {
	return
}

func (m *mat) Std() {
	return
}

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
