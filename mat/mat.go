/*
Package mat implements a "mat" object, which behaves like a 2-dimensional array
or list in other programming languages. Under the hood, the mat object is a
flat slice, which provides for optimal performance in Go, while the methods
and constructors provide for a higher level of performance and abstraction
when compared to the "2D" slices of go (slices of slices).

All errors encountered in this package, such as attempting to access an
element out of bounds are treated as critical error, and thus, the code
immediately exits with signal 1. In such cases, the function/method in
which the error was encountered is printed to the screen, in addition
to the full stack trace, in order to help fix the issue rapidly.
*/
package mat

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
)

/*
Mat is the main struct of this library. Mat is a essentially a 1D slice
(a []float64) that contains two integers, representing rows and columns,
which allow it to behave as if it was a 2D slice. This allows for higher
performance and flexibility for the users of this library, at the expense
of some bookkeeping that is done here.

The fields of this struct are not directly accessible, and they may only
change by the use of the various methods in this library.
*/
type Mat struct {
	r, c int
	vals []float64
}

type elementFunc func(*float64)
type booleanFunc func(*float64) bool
type reducerFunc func(accumulator, next *float64)

/*
New is the primary constructor for the "Mat" object. New is a veradic function,
expecting 0 to 3 ints, with differing behavior as follows:

m := New()

m is now an empty &Mat{}, where the number of rows,
columns and the length and capacity of the underlying
slice are all zero. This is mostly for internal use.

m := New(x)

m is a x by x (square) matrix, with the underlying
slice of length x, and capacity 2x.

m := New(x, y)

m is an x by y matrix, with the underlying slice of
length rc, and capacity of 2rc. This is a good case
for when your matrix is going to expand in the
future. There is a negligible hit to performance
and a larger memory usage of your code. But in case
expanding matrices, many reallocations are avoided.

m := New(x, y, z)

m is a x by u matrix, with the underlying slice of
length rc, and capacity z. This is a good choice for
when the size of the matrix is static, or when the
application is memory constrained.

For most cases, we recommend using the New(x) or New(x, y) options, and
almost never the New() option.
*/
func New(dims ...int) *Mat {
	m := &Mat{}
	switch len(dims) {
	case 0:
	case 1:
		r := dims[0]
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
		m = &Mat{
			r,
			r,
			make([]float64, r*r, 2*r*r),
		}
	case 2:
		r := dims[0]
		c := dims[1]
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
		m = &Mat{
			r,
			c,
			make([]float64, r*c, 2*r*c),
		}
	case 3:
		r := dims[0]
		c := dims[1]
		capacity := dims[2]
		if r <= 0 {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s, the number of rows must be greater than '0',\n"
			s += "but recieved %d. "
			s = fmt.Sprintf(s, "New", r)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		if c <= 0 {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s, the number of columns must be greater than '0',\n"
			s += "but recieved %d. "
			s = fmt.Sprintf(s, "New", c)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		m = &Mat{
			r,
			c,
			make([]float64, r*c, capacity),
		}
	default:
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s expected 0 to 3 arguments, but recieved %d"
		s = fmt.Sprintf(s, "New", len(dims))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	return m
}

/*
FromSlice generated a mat object from a [][]float64 slice. The slice is
checked for being a non-jagged slice, where each row contains the same
number of elements. The creation of a mat object from jagged 2D slices
is not supported as on yet.
*/
func FromSlice(s [][]float64) *Mat {
	if s == nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, an unallocated 2D slice was recieved, where the\n"
		s += "slice is equal to nil."
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
		s += "creation of a *Mat from jagged slices is not supported.\n"
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

/*
From1DSlice creates a mat object from a slice of float64s. The created mat
object has one row, and the number of columns equal to the length of the
1D slice from which it was created.
*/
func From1DSlice(s []float64) *Mat {
	if s == nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s, an unallocated 1D slice was recieved, where the\n"
		s += "slice is equal to nil."
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

/*
FromCSV creates a mat object from a CSV (comma separated values) file. Here, we
assume that the number of rows of the resultant mat object is equal to the
number of lines, and the number of columns is equal to the number of entries
in each line. As before, we make sure that each line contains the same number
of elements.

The file to be read is assumed to be very large, and hence it is read one line
at a time. This results in some major inefficiencies, and it is recommended
that this function be used sparingly, and not as a major component of your
library/executable.

Unline other mat creation methods in this package, the mat object created here,
the capacity of the mat opject created here is the same as its length since we
assume the mat to be very large.
*/
func FromCSV(filename string) *Mat {
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
	// of columns based on the number of comma separated strings in that line.
	// Then we will read the rest of the lines one at a time, checking that the
	// number of entries in each line is the same as the first line, and
	// incrementing the number of rows each time.
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
	m := New()
	// Start with one row, and set the number of enteries per row
	m.r = 1
	m.c = len(str)
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
		line++
		if len(str) != len(row) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, line %d in %s has %d entries. The first line\n"
			s += "(line 1) has %d entries.\n"
			s += "Creation of a *Mat from jagged slices is not supported.\n"
			s = fmt.Sprintf(s, "Load", filename, err)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		m.r++
	}
	return m
}

/*
Reshape changes the row and the columns of the mat object as long as the total
number of values contained in the mat object remains constant. The order and
the values of the mat does not change with this function.
*/
func (m *Mat) Reshape(rows, cols int) *Mat {
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

/*
Dims returns the number of rows and columns of a mat object.
*/
func (m *Mat) Dims() (int, int) {
	return m.r, m.c
}

/*
Vals returns the values contained in a mat object as a 1D slice of float64s.
*/
func (m *Mat) Vals() []float64 {
	s := make([]float64, m.r*m.c)
	copy(s, m.vals)
	return s
}

/*
ToSlice returns the values of a mat object as a 2D slice of float64s.
*/
func (m *Mat) ToSlice() [][]float64 {
	s := make([][]float64, m.r)
	for i := 0; i < m.r; i++ {
		s[i] = make([]float64, m.c)
		for j := 0; j < m.c; j++ {
			s[i][j] = m.vals[i*m.c+j]
		}
	}
	return s
}

/*
ToCSV creates a file with the passed name, and writes the content of a mat
object to it, by putting each row in a single comma separated line. The
number of entries in each line is equal to the columns of the mat object.
*/
func (m *Mat) ToCSV(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "ToCSV", fileName, err)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	defer f.Close()
	str := ""
	idx := 0
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			str += strconv.FormatFloat(m.vals[idx], 'e', 14, 64)
			if j+1 != m.c {
				str += ","
			}
			idx++
		}
		if i+1 != m.r {
			str += "\n"
		}
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot write to %s due to error: %v.\n"
		s = fmt.Sprintf(s, "ToCSV", fileName, err)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
}

/*
At returns a pointer to the float64 stored in the given row and column.
*/
func (m *Mat) At(r, c int) float64 {
	if (r >= m.r) || (r < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested row %d is outside of bounds [0, %d)\n"
		s = fmt.Sprintf(s, "At", r, m.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if (c >= m.c) || (c < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [0, %d)\n"
		s = fmt.Sprintf(s, "At", r, m.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	return m.vals[r*m.c+c]
}

/*
Map applies a given function to each element of a mat object. The given
function must take a pointer to a float64, and return nothing.
*/
func (m *Mat) Map(f elementFunc) *Mat {
	for i := 0; i < m.r*m.c; i++ {
		f(&m.vals[i])
	}
	return m
}

/*
Ones sets all values of a mat to be equal to 1.0
*/
func (m *Mat) Ones() *Mat {
	f := func(i *float64) {
		*i = 1.0
		return
	}
	return m.Map(f)
}

/*
Inc takes each element of a mat object, and starting from 0.0, sets their
value to be the value of the previous entry plus 1.0. In other words, the
first few values of a mat object after this operation would be 0.0, 1.0,
2.0, ...
*/
func (m *Mat) Inc() *Mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = float64(i)
	}
	return m
}

/*
Reset sets all values of a mat object to 0.0
*/
func (m *Mat) Reset() *Mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = 0.0
	}
	return m
}

/*
SetAllTo sets all values of a mat to the passed float64 value.
*/
func (m *Mat) SetAllTo(val float64) *Mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = val
	}
	return m
}

/*
Rand sets the values of a mat to random numbers. The range from which the random
numbers are selected is determined based on the arguments passed.

- For no arguments, the range is [0, 1)
- For 1 argument, the range is [0, arg) for arg > 0, or (arg, 0] is arg < 0.
- For 2 arguments, the range is [arg1, arg2).
*/
func (m *Mat) Rand(args ...float64) *Mat {
	switch len(args) {
	case 0:
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64()
		}
	case 1:
		to := args[0]
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64() * to
		}
	case 2:
		from := args[0]
		to := args[1]
		if !(from < to) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the first argument, %f, is not less than the\n"
			s += "second argument, %f. The first argument must be strictly\n"
			s += "less than the second.\n"
			s = fmt.Sprintf(s, "Rand", from, to)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		for i := 0; i < m.r*m.c; i++ {
			m.vals[i] = rand.Float64()*(to-from) + from
		}
	default:
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s expected 0 to 2 arguments, but recieved %d."
		s = fmt.Sprintf(s, "Rand", len(args))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	return m
}

/*
Col returns a new mat object whole values are equal to a column of the original
mat object. The number of Rows of the returned mat object is equal to the
number of rows of the original mat, and the number of columns is equal to 1.
*/
func (m *Mat) Col(x int) *Mat {
	if (x >= m.c) || (x < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [0, %d)\n"
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

/*
Row returns a new mat object whose values are equal to a row of the original
mat object. The number of Rows of the returned mat object is equal to 1, and
the number of columns is equal to the number of columns of the original mat.
*/
func (m *Mat) Row(x int) *Mat {
	if (x >= m.r) || (x < 0) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested row %d is outside of the bounds [0, %d)\n"
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

/*
Equals checks to see if two mat objects are equal. That mean that the two mats
have the same number of rows, same number of columns, and have the same float64
in each entry at a given index.
*/
func (m *Mat) Equals(n *Mat) bool {
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

/*
Copy returns a duplicate of a mat object. The returned copy is "deep", meaning
that the object can be manipulated without effecting the original mat object.
*/
func (m *Mat) Copy() *Mat {
	n := New(m.r, m.c)
	copy(n.vals, m.vals)
	return n
}

/*
T returns the transpose of the original matrix. The transpose of a mat object
is defined in the usual manner, where every value at row x, and column y is
placed at row y, and column x. The number of rows and column of the transposed
mat are equal to the number of columns and rows of the original matrix,
respectively. This method creates a new mat object, and the original is
left intact.
*/
func (m *Mat) T() *Mat {
	n := New(m.c, m.r)
	idx := 0
	for i := 0; i < m.c; i++ {
		for j := 0; j < m.r; j++ {
			n.vals[idx] = m.vals[j*m.c+i]
			idx++
		}
	}
	return n
}

/*
Filter applies a function to each element of a mat object, creating a new
mat object from all elements for which the function returned true. For
example consider the following function:

f := func(i *float64) bool {
	if i > 0.0 {
		return true
	} else {
		return false
	}
}

then calling

m.Filter(f)

will create a new mat object which contains the positive elements of the
original matrix. If no elements return true for a given function, nil is
returned. It is expected that the caller of this method checks the
returned value to ensure that it is not nil.
*/
func (m *Mat) Filter(f booleanFunc) *Mat {
	var res []float64
	for i := 0; i < m.r*m.c; i++ {
		if f(&m.vals[i]) {
			res = append(res, m.vals[i])
		}
	}
	if len(res) == 0 {
		return nil
	}
	return From1DSlice(res)
}

/*
All checks if a supplied function is true for all elements of a mat object.
For instance, consider

positive := func(i *float64) bool {
	if i > 0.0 {
		return true
	} else {
		return false
	}
}

Then calling

m.All(positive)

will return true if and only if all elements in m are positive.
*/
func (m *Mat) All(f booleanFunc) bool {
	for i := 0; i < m.r*m.c; i++ {
		if !f(&m.vals[i]) {
			return false
		}
	}
	return true
}

/*
Any checks if a supplied function is true for one elements of a mat object.
For instance,


positive := func(i *float64) bool {
	if i > 0.0 {
		return true
	} else {
		return false
	}
}

Then calling

m.Any(positive)

would be true if at least one element of the mat object is positive.
*/
func (m *Mat) Any(f booleanFunc) bool {
	for i := 0; i < m.r*m.c; i++ {
		if f(&m.vals[i]) {
			return true
		}
	}
	return false
}

/*
Mul is the element-wise multiplication of a mat object by another which is
passed to this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise multiplication is stored in the original
mat on which the method was invoked.
*/
func (m *Mat) Mul(n *Mat) *Mat {
	if m.r != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first mat is %d\n"
		s += "but the number of rows of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Mul", m.r, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if m.c != n.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first mat is %d\n"
		s += "but the number of columns of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Mul", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] *= n.vals[i]
	}
	return m
}

/*
Add is the element-wise addition of a mat object with another which is passed
to this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise addition is stored in the original
mat on which the method was invoked.
*/
func (m *Mat) Add(n *Mat) *Mat {
	if m.r != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first mat is %d\n"
		s += "but the number of rows of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Add", m.r, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if m.c != n.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first mat is %d\n"
		s += "but the number of columns of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Add", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] += n.vals[i]
	}
	return m
}

/*
Sub is the element-wise subtraction of a mat object which is passed
to this method from the original mat which called the method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise subtraction is stored in the original
mat on which the method was invoked.
*/
func (m *Mat) Sub(n *Mat) *Mat {
	if m.r != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first mat is %d\n"
		s += "but the number of rows of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Sub", m.r, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if m.c != n.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first mat is %d\n"
		s += "but the number of columns of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Sub", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] -= n.vals[i]
	}
	return m
}

/*
Div is the element-wise dicition of a mat object by another which is passed
to this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise divition is stored in the original
mat on which the method was invoked. The dividing mat object (passed to this
method) must not contain any elements which are equal to 0.0.
*/
func (m *Mat) Div(n *Mat) *Mat {
	if m.r != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first mat is %d\n"
		s += "but the number of rows of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Div", m.r, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if m.c != n.c {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first mat is %d\n"
		s += "but the number of columns of the second mat is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Div", m.c, n.c)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
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

/*
Scale is the element-wise multiplication of a mat object by a scalar.

The results of the element-wise multiplication is stored in the original
mat on which the method was invoked.
*/
func (m *Mat) Scale(f float64) *Mat {
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] *= f
	}
	return m
}

/*
Sum returns the sum of the elements along a specific row or specific column.
The first argument selects the row or column (0 or 1), and the second argument
selects which row or column for which we want to calculate the sum. For
example:

m.Sum(0, 2)

will calculate the sum of the 3rd row of mat m.
*/
func (m *Mat) Sum(axis, slice int) float64 {
	if axis != 0 && axis != 1 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the first argument must be 0 or 1, however %d "
		s += "was recieved.\n"
		s = fmt.Sprintf(s, "Sum", axis)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if axis == 0 {
		if (slice >= m.r) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Sum", slice, m.r)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	} else if axis == 1 {
		if (slice >= m.c) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Sum", slice, m.c)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	}
	return m.precheckedSum(axis, slice)
}

func (m *Mat) precheckedSum(axis, slice int) float64 {
	x := 0.0
	if axis == 0 {
		for i := 0; i < m.c; i++ {
			x += m.vals[slice*m.c+i]
		}
	} else if axis == 1 {
		for i := 0; i < m.r; i++ {
			x += m.vals[i*m.c+slice]
		}
	}
	return x
}

/*
Average returns the average of the elements along a specific row or specific
column.
The first argument selects the row or column (0 or 1), and the second argument
selects which row or column for which we want to calculate the average. For
example:

m.Average(0, 2)

will calculate the average of the 3rd row of mat m.
*/
func (m *Mat) Average(axis, slice int) float64 {
	if axis != 0 && axis != 1 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the first argument must be 0 or 1, however %d "
		s += "was recieved.\n"
		s = fmt.Sprintf(s, "Average", axis)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if axis == 0 {
		if (slice >= m.r) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Average", slice, m.r)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	} else if axis == 1 {
		if (slice >= m.c) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Average", slice, m.c)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	}
	return m.precheckedAverage(axis, slice)
}

func (m *Mat) precheckedAverage(axis, slice int) float64 {
	sum := m.precheckedSum(axis, slice)
	if axis == 0 {
		return sum / float64(m.c)
	}
	return sum / float64(m.r)
}

/*
Prod returns the product of the elements along a specific row or specific
column.
The first argument selects the row or column (0 or 1), and the second argument
selects which row or column for which we want to calculate the product. For
example:

m.Prod(1, 2)

will calculate the product of the 3rd column of mat m.
*/
func (m *Mat) Prod(axis, slice int) float64 {
	if axis != 0 && axis != 1 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the first argument must be 0 or 1, however %d "
		s += "was recieved.\n"
		s = fmt.Sprintf(s, "Prod", axis)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if axis == 0 {
		if (slice >= m.r) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the requested row %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Prod", slice, m.r)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	} else if axis == 1 {
		if (slice >= m.c) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Prod", slice, m.c)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	}
	x := 1.0
	if axis == 0 {
		for i := 0; i < m.c; i++ {
			x *= m.vals[slice*m.c+i]
		}
	} else if axis == 1 {
		for i := 0; i < m.r; i++ {
			x *= m.vals[i*m.c+slice]
		}
	}
	return x
}

/*
Std returns the standard deviation of the elements along a specific row
or specific column. The standard deviation is defined as the square root of
the mean distance of each element from the mean. Look at:
http://mathworld.wolfram.com/StandardDeviation.html

For example:

m.Std(1, 0)

will calculate the standard deviation of the first column of mat m.
*/
func (m *Mat) Std(axis, slice int) float64 {
	if axis != 0 && axis != 1 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the first argument must be 0 or 1, however %d "
		s += "was recieved.\n"
		s = fmt.Sprintf(s, "Std", axis)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	if axis == 0 {
		if (slice >= m.r) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Std", slice, m.r)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	} else if axis == 1 {
		if (slice >= m.c) || (slice < 0) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
			s = fmt.Sprintf(s, "Std", slice, m.c)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
	}
	avg := m.precheckedAverage(axis, slice)
	var s []float64
	if axis == 0 {
		s = make([]float64, m.c)
		for i := 0; i < m.c; i++ {
			s[i] = avg - m.vals[slice*m.c+i]
			s[i] *= s[i]
		}
	} else {
		s = make([]float64, m.r)
		for i := 0; i < m.r; i++ {
			s[i] = avg - m.vals[i*m.c+slice]
			s[i] *= s[i]
		}
	}
	sum := 0.0
	for i := range s {
		sum += s[i]
	}
	return math.Sqrt(sum / float64(len(s)))
}

/*
Dot is the matrix multiplication of two mat objects. Consider the following two
mats:

m := New(5, 6)
n := New(6, 10)

then

o := m.Dot(n)

is a 5 by 10 mat whose element at row i and column j is given by:

Sum(m.Row(i).Mul(n.col(j))
*/
func (m *Mat) Dot(n *Mat) *Mat {
	if m.c != n.r {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the number of columns of the first mat is %d\n"
		s += "which is not equal to the number of row of the second mat,\n"
		s += "which is %d. They must be equal.\n"
		s = fmt.Sprintf(s, "Dot", m.c, n.r)
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
	o := New(m.r, n.c)
	for i := 0; i < m.r; i++ {
		for j := 0; j < n.c; j++ {
			for k := 0; k < m.c; k++ {
				o.vals[i*o.c+j] += (m.vals[i*m.c+k] * n.vals[k*n.c+j])
			}
		}
	}
	return o
}

/*
Print displays the mat to the stdout.
*/
func (m *Mat) Print() {
	var str string
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			str += strconv.FormatFloat(m.vals[i*m.c+j], 'f', 14, 64)
			str += " "
		}
		if i+1 <= m.r {
			str += "\n"
		}
	}
	fmt.Println(str)
}

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
