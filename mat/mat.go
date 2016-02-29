/*
Package mat implements a large set of functions which act on 2-dimensional slices
of float64.

All 2D slices created or passed to the functions in this library are assumed to
be non-jagged, meaning that all rows have the same number of elements. It is the
responsibility of the called to ensure this. An easy way to do this is to use
the various functions in this library to create 2D slices or to generate them
from files.

All errors encountered in this package, such as attempting to access an
element out of bounds are treated as critical error, and thus, the code
immediately panics. In such cases, the function/method in
which the error was encountered is printed to the screen, in addition
to the full stack trace, in order to help fix the issue rapidly.
*/
package mat

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"sync"
)

/*
ElementalFunc defines the signature of a function that takes a float64, and
returns a float64
*/
type ElementFunc func(float64) float64

/*
BooleanFunc defines the signature of a function that takes a float64, and
return a bool.
*/
type BooleanFunc func(float64) bool
type reducerFunc func(accumulator, next *float64)

/*
New is a utility function to create 2D slices. New is a variadic function,
expecting 1 or 2 ints, with differing behavior as follows:

	m := mat.New(x)

m is a x by x (square) 2D slice. Alternatively

	m := mat.New(x, y)

is a 2D slice with x rows and y columns.

*/
func New(dims ...int) [][]float64 {
	var m [][]float64
	switch len(dims) {
	case 1:
		r := dims[0]
		if r <= 0 {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s, the number of rows must be greater than '0', but\n"
			s += "recieved %d. "
			s = fmt.Sprintf(s, "New", r)
			panic(s)
		}
		m = make([][]float64, r)
		for i := range m {
			m[i] = make([]float64, r)
		}
	case 2:
		r := dims[0]
		c := dims[1]
		if r <= 0 {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s, the number of rows must be greater than '0', but\n"
			s += "recieved %d. "
			s = fmt.Sprintf(s, "New", r)
			panic(s)
		}
		if c <= 0 {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%s, the number of columns must be greater than '0', but\n"
			s += "recieved %d. "
			s = fmt.Sprintf(s, "New", c)
			panic(s)
		}
		m = make([][]float64, r)
		for i := range m {
			m[i] = make([]float64, c)
		}
	default:
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s expected 1 or 2 arguments, but recieved %d"
		s = fmt.Sprintf(s, "New", len(dims))
		panic(s)
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
FromCSV creates a mat object from a CSV (comma separated values) file. Here, we
assume that the number of rows of the resultant 2D slice is equal to the
number of lines, and the number of columns is equal to the number of entries
in each line making sure that each line contains the same number
of elements.

The file to be read is assumed to be very large, and hence it is read one line
at a time.
*/
func FromCSV(filename string) [][]float64 {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot open %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV", filename, err)
		panic(s)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// I am going with the assumption that a 2D slice loaded from a CSV is going to
	// be large. So, we are going to read one line, and determine the number
	// of columns based on the number of comma separated strings in that line.
	// Then we will read the rest of the lines one at a time, checking that the
	// number of entries in each line is the same as the first line.
	str, err := r.Read()
	if err != nil {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, cannot read from %s due to error: %v.\n"
		s = fmt.Sprintf(s, "FromCSV", filename, err)
		panic(s)
	}
	line := 1
	m := [][]float64{}
	for {
		row := make([]float64, len(str))
		for i := range str {
			row[i], err = strconv.ParseFloat(str[i], 64)
			if err != nil {
				fmt.Println("\nNumgo/mat error.")
				s := "In mat.%v, item %d in line %d is %s, which cannot\n"
				s += "be converted to a float64 due to: %v"
				s = fmt.Sprintf(s, "FromCSV", i, line, str[i], err)
				panic(s)
			}
		}
		m = append(m, row)
		// Read the next line. If there is one.
		str, err = r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, cannot read from %s due to error: %v.\n"
			s = fmt.Sprintf(s, "FromCSV", filename, err)
			panic(s)
		}
		line++
		if len(str) != len(row) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, line %d in %s has %d entries. The first line\n"
			s += "(line 1) has %d entries.\n"
			s += "Creation of a *Mat from jagged slices is not supported.\n"
			s = fmt.Sprintf(s, "FromCSV", filename, err)
			panic(s)
		}
	}
	return m
}

/*
Flatten turns a 2D slice of float64 into a 1D slice of float64. This is done
by appending all rows tip to tail.
*/
func Flatten(m [][]float64) []float64 {
	var n []float64
	for i := range m {
		n = append(n, m[i]...)
	}
	return n
}

/*
ToCSV writes the content of a passed 2D slice into a CSV file with the passed
name, by putting each row in a single comma separated line. The number of
entries in each line is equal to the length of the second dimension of the
2D slice. This function return an error, which contains any errors encounterd
or nil if it succeeded.
*/
func ToCSV(m [][]float64, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	str := ""
	r, c := len(m), len(m[0])
	for i := range m {
		for j := range m[i] {
			str += strconv.FormatFloat(m[i][j], 'e', 14, 64)
			if j+1 != c {
				str += ","
			}
		}
		if i+1 != r {
			str += "\n"
		}
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

/*
Map applies a given function to each element of a 2D slice of float64. The
passed function must satisfy the signature of an ElementalFunc.
*/
func Map(f ElementFunc, m [][]float64) {
	for i := range m {
		for j := range m[i] {
			m[i][j] = f(m[i][j])
		}
	}
}

/*
SetAll sets all elements of a 2D slice to the passed value.
*/
func SetAll(m [][]float64, val float64) {
	for i := range m {
		for j := range m[i] {
			m[i][j] = val
		}
	}
}

/*
Mul multiples all elements of a 2D slice by the passed value. The passed value can be
a float64, a 1D slice of float64, or a 2D slice of float64.

When the passed value is a float64, then each element of the 2D slice are multiplied
by the passed value.

If the passed value is a []float64, then each row of the 2D slice is elementally
multipled by the corresponding entry in the passed 1D slice.


Finally, if the passed value is a [][]float64, then Mul takes each element of the
first 2D slice passed to it, and multiples that element by the corresponding element
in the second 2D slice passed to this function.
The shape of the 2D slices must be the same (same number or rows and columns),
and they are assumed to be non-jagged (same number of elements in each row).

In each case, the result of the multiplication is stored in the original 2D slice.
*/
func Mul(m [][]float64, val interface{}) {
	switch v := val.(type) {
	case float64:
		for i := range m {
			for j := range m[i] {
				m[i][j] *= v
			}
		}
	case []float64:
		for i := range m {
			if len(v) != len(m[i]) {
				fmt.Println("\nNumgo/mat error.")
				s := "In mat.%v, in row %d, the number of the columns of the first\n"
				s += "slice is %d, but the length of the vector is %d. They must\n"
				s += "match.\n"
				s = fmt.Sprintf(s, "Mul", i, len(m[i]), len(v))
				panic(s)
			}
		}
		for i := range m {
			for j := range v {
				m[i][j] *= v[j]
			}
		}
	case [][]float64:
		if len(m) != len(v) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, the number of the rows of the first slice is %d\n"
			s += "but the number of rows of the second slice is %d. They must\n"
			s += "match.\n"
			s = fmt.Sprintf(s, "Mul", len(m), len(v))
			panic(s)
		}
		for i := range m {
			if len(m[i]) != len(v[i]) {
				fmt.Println("\nNumgo/mat error.")
				s := "In mat.%v, column number %d of the first matrix has length %d,\n"
				s += "while column number %d of the second 2D slice has length %d.\n"
				s += "The length of each column must match.\n"
				s = fmt.Sprintf(s, "Mul", i, len(m[i]), i, len(v[i]))
				panic(s)
			}
			for j := range m[i] {
				m[i][j] *= v[i][j]
			}
		}
	default:
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, expected float64, []float64, or [][]float64 for the second\n"
		s += "argument, but received argument of type: %v."
		s = fmt.Sprintf(s, "Mul", reflect.TypeOf(v))
		panic(s)
	}
}

/*
AddAll Adds the passed value to each element of the passed 2D slice.
*/
func AddAll(m [][]float64, val float64) {
	for i := range m {
		for j := range m[i] {
			m[i][j] += val
		}
	}
}

/*
SubAll subtracts the passed value from each element of the passed 2D slice.
*/
func SubAll(m [][]float64, val float64) {
	for i := range m {
		for j := range m[i] {
			m[i][j] -= val
		}
	}
}

/*
DivAll Diives each element of the 2D slice by the passed value. The passed
value must not be 0.0, to avoid division by zero.
*/
func DivAll(m [][]float64, val float64) {
	if val == 0.0 {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the passed value cannot be 0.0\n"
		s = fmt.Sprintf(s, "DivAll")
		panic(s)
	}
	for i := range m {
		for j := range m[i] {
			m[i][j] /= val
		}
	}
}

/*
Rand sets the values of a 2D slice, m, to random numbers. The range from which
the random numbers are selected is determined based on the arguments passed.

For no additional arguments, such as
	mat.Rand(m)
the range is [0, 1)

For 1 argument, such as
	mat.Rand(m, arg)
the range is [0, arg) for arg > 0, or (arg, 0] is arg < 0.

For 2 arguments, such as
	mat.Rand(m, arg1, arg2)
the range is [arg1, arg2). For this case, arg1 must be less than arg2, or
the function will panic.
*/
func Rand(m [][]float64, args ...float64) {
	switch len(args) {
	case 0:
		for i := range m {
			for j := range m[i] {
				m[i][j] = rand.Float64()
			}
		}
	case 1:
		to := args[0]
		for i := range m {
			for j := range m[i] {
				m[i][j] = rand.Float64() * to
			}
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
			panic(s)
		}
		for i := range m {
			for j := range m[i] {
				m[i][j] = rand.Float64()*(to-from) + from
			}
		}
	default:
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s expected 0 to 2 arguments, but recieved %d."
		s = fmt.Sprintf(s, "Rand", len(args))
		panic(s)
	}
}

/*
Col returns a column from a 2D slice of float64. For example:

	fmt.Println(m) // [[1.0, 2.3], [3.4, 1.7]]
	mat.Col(0, m) // [1.0, 3.4]

Col also accepts negative indices. For example:

	mat.Col(-1, m) // [2.3, 1.7]
*/
func Col(x int, m [][]float64) []float64 {
	if (x >= len(m[0])) || (x < -len(m[0])) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "Col", x, len(m[0]), len(m[0]))
		panic(s)
	}
	v := make([]float64, len(m))
	if x >= 0 {
		for i := range m {
			v[i] = m[i][x]
		}
	} else {
		for i := range m {
			v[i] = m[i][len(m[0])+x]
		}
	}
	return v
}

/*
Row returns a row from a 2D slice of float64. For example:

	fmt.Println(m) // [[1.0, 2.3], [3.4, 1.7]]
	mat.Row(0, m) // [1.0, 2.3]

Row also accepts negative indices. For example:

	mat.Row(-1, m) // [3.4, 1.7]
*/
func Row(x int, m [][]float64) []float64 {
	if (x >= len(m)) || (x < -len(m)) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested row %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "Row", x, len(m), len(m))
		panic(s)
	}
	v := make([]float64, len(m[0]))
	if x >= 0 {
		copy(v, m[x])
	} else {
		copy(v, m[len(m)+x])
	}
	return v
}

/*
Equal checks to see if two 2D slices are equal. That mean that the two slices
have the same number of rows, same number of columns, and have the same float64
in each entry at a given set of indices.
*/
func Equal(m, n [][]float64) bool {
	if len(m) != len(n) {
		return false
	}
	for i := range m {
		if len(m[i]) != len(n[i]) {
			return false
		}
	}
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j] {
				return false
			}
		}
	}
	return true
}

/*
Copy returns a duplicate of a 2D slice. The returned copy is "deep", meaning
that the object can be manipulated without effecting the original.
*/
func Copy(m [][]float64) [][]float64 {
	n := New(len(m), len(m[0]))
	for i := range m {
		copy(n[i], m[i])
	}
	return n
}

/*
T returns the transpose of the original 2D slice. The transpose of a 2D slice
is defined in the usual manner, where every value at row x, and column y is
placed at row y, and column x. The number of rows and column of the transpose
of a slice are equal to the number of columns and rows of the original slice,
respectively. This method creates a new 2D slice, and the original is
left intact. The passed 2d slice is assumed to be non-jagged.
*/
func T(m [][]float64) [][]float64 {
	n := New(len(m[0]), len(m))
	for i := range m {
		for j := range m[i] {
			n[j][i] = m[i][j]
		}
	}
	return n
}

/*
All checks if a supplied function is true for all elements of a mat object.
The supplied function is expected to have the signature of a BooleanFunc, which
takes a float64, returning a bool.
For instance, consider

	positive := func(i *float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}

Then calling

	mat.All(positive, m)

will return true if and only if all elements in m are positive.
*/
func All(f BooleanFunc, m [][]float64) bool {
	for i := range m {
		for j := range m[i] {
			if !f(m[i][j]) {
				return false
			}
		}
	}
	return true
}

/*
Any checks if a supplied function is true for at least one elements of
a 2D slice of float64s. The supplied function must have the signature of
a BooleanFunc, meaning that it takes a float64, and returns a bool.
For instance,

	positive := func(i *float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}

Then calling

	mat.Any(positive, m)

would be true if at least one element of the m is positive.
*/
func Any(f BooleanFunc, m [][]float64) bool {
	for i := range m {
		for j := range m[i] {
			if f(m[i][j]) {
				return true
			}
		}
	}
	return false
}

/*
Add takes each element of the first 2D slice passed to it, and adds that
element to the corresponding element in the second 2D slice passed to this
function, storing the result in the first 2D slice.

The shape of the 2D slices must be the same (same number or rows and columns),
and they are assumed to be non-jagged (same number of elements in each row).
*/
func Add(m, n [][]float64) {
	if len(m) != len(n) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first slice is %d\n"
		s += "but the number of rows of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Add", len(m), len(n))
		panic(s)
	}
	for i := range m {
		if len(m[i]) != len(n[i]) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, column number %d of the first matrix has length %d,\n"
			s += "while column number %d of the second 2D slice has length %d.\n"
			s += "The length of each column must match.\n"
			s = fmt.Sprintf(s, "Add", i, len(m[i]), i, len(n[i]))
			panic(s)
		}
		for j := range m[i] {
			m[i][j] += n[i][j]
		}
	}
}

/*
Sub takes each element of the first 2D slice passed to it, and subtracts from
that element the corresponding element in the second 2D slice passed to this
function, storing the result in the first 2D slice.

The shape of the 2D slices must be the same (same number or rows and columns),
and they are assumed to be non-jagged (same number of elements in each row).
*/
func Sub(m, n [][]float64) {
	if len(m) != len(n) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first slice is %d\n"
		s += "but the number of rows of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Sub", len(m), len(n))
		panic(s)
	}
	for i := range m {
		if len(m[i]) != len(n[i]) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, column number %d of the first matrix has length %d,\n"
			s += "while column number %d of the second 2D slice has length %d.\n"
			s += "The length of each column must match.\n"
			s = fmt.Sprintf(s, "Sub", i, len(m[i]), i, len(n[i]))
			panic(s)
		}
		for j := range m[i] {
			m[i][j] -= n[i][j]
		}
	}
}

/*
Div takes each element of the first 2D slice passed to it, and divides that
element by the corresponding element in the second 2D slice passed to this
function, storing the result in the first 2D slice.

The shape of the 2D slices must be the same (same number or rows and columns),
and they are assumed to be non-jagged (same number of elements in each row).
The second slice must not contain any zero-valued elements.
*/
func Div(m, n [][]float64) {
	if len(m) != len(n) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first slice is %d\n"
		s += "but the number of rows of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Div", len(m), len(n))
		panic(s)
	}
	zero := func(i float64) bool {
		if i == 0.0 {
			return true
		}
		return false
	}
	if Any(zero, n) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the second slice contains a zero-valued element.\n"
		s = fmt.Sprintf(s, "Div")
		panic(s)
	}
	for i := range m {
		if len(m[i]) != len(n[i]) {
			fmt.Println("\nNumgo/mat error.")
			s := "In mat.%v, column number %d of the first matrix has length %d,\n"
			s += "while column number %d of the second 2D slice has length %d.\n"
			s += "The length of each column must match.\n"
			s = fmt.Sprintf(s, "Div", i, len(m[i]), i, len(n[i]))
			panic(s)
		}
		for j := range m[i] {
			m[i][j] /= n[i][j]
		}
	}
}

/*
Sum returns the sum of all elements in a 2D slice of float64.
*/
func Sum(m [][]float64) float64 {
	sum := 0.0
	for i := range m {
		for j := range m[i] {
			sum += m[i][j]
		}
	}
	return sum
}

/*
SumCol returns the sum of the elements of a slice along a specific column.
For example:

	mat.SumCol(2, m)

will return the sum of the 3rd column of m. Negative indices are also
allowed. For example:

	mat.SumCol(-1, m)

will return the sum of the last column of m.
*/
func SumCol(x int, m [][]float64) float64 {
	if (x >= len(m[0])) || (x < -len(m[0])) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "SumCol", x, len(m[0]), len(m[0]))
		panic(s)
	}
	var sum float64
	if x >= 0 {
		for i := range m {
			sum += m[i][x]
		}
	} else {
		for i := range m {
			sum += m[i][len(m[0])+x]
		}
	}
	return sum
}

/*
SumRow returns the sum of the elements of a slice along a specific row.
For example:

	mat.SumRow(2, m)

will return the sum of the 3rd row of m. Negative indices are also
allowed. For example:

	mat.SumRow(-1, m)

will return the sum of the last row of m.
*/
func SumRow(x int, m [][]float64) float64 {
	if (x >= len(m)) || (x < -len(m)) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "SumRow", x, len(m), len(m))
		panic(s)
	}
	var sum float64
	if x >= 0 {
		for i := range m[x] {
			sum += m[x][i]
		}
	} else {
		for i := range m[len(m)+x] {
			sum += m[len(m)+x][i]
		}
	}
	return sum
}

/*
Avg returns the average value of all the elements in a 2D slice.
*/
func Avg(m [][]float64) float64 {
	avg := 0.0
	numItems := 0
	for i := range m {
		for j := range m[i] {
			avg += m[i][j]
			numItems++
		}
	}
	avg /= float64(numItems)
	return avg
}

/*
AvgRow returns the average of the elements of a slice along a specific row.
For example:

	mat.AvgRow(2, m)

will return the average of the 3rd row of m. Negative indices are also
allowed. For example:

	mat.AvgRow(-1, m)

will return the average of the last row of m.
*/
func AvgRow(x int, m [][]float64) float64 {
	if (x >= len(m)) || (x < -len(m)) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "AvgRow", x, len(m), len(m))
		panic(s)
	}
	var sum float64
	if x >= 0 {
		for i := range m[x] {
			sum += m[x][i]
		}
		sum /= float64(len(m[x]))
	} else {
		for i := range m[len(m)+x] {
			sum += m[len(m)+x][i]
		}
		sum /= float64(len(m[len(m)+x]))
	}
	return sum
}

/*
AvgCol returns the average of the elements of a slice along a specific column.
For example:

	mat.AvgCol(2, m)

will return the average of the 3rd column of m. Negative indices are also
allowed. For example:

	mat.AvgCol(-1, m)

will return the Avrage of the last column of m.
*/
func AvgCol(x int, m [][]float64) float64 {
	if (x >= len(m[0])) || (x < -len(m[0])) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "AvgCol", x, len(m[0]), len(m[0]))
		panic(s)
	}
	var sum float64
	if x >= 0 {
		for i := range m {
			sum += m[i][x]
		}
	} else {
		for i := range m {
			sum += m[i][len(m[0])+x]
		}
	}
	sum /= float64(len(m))
	return sum
}

func Dot(m, n [][]float64) [][]float64 {
	//for i := range m {
	//	if len(m) != len(n[i]) {
	//		fmt.Println("\nNumgo/mat error.")
	//		s := "In mat.%s, Column %d of the 2nd argument has %d elements,\n"
	//		s += "while the 1st argument has %d rows. They must match.\n"
	//		s += fmt.Sprintf(s, "Dot", i, len(n[i]), len(m))
	//		panic(s)
	//	}
	//}
	//for i := range n {
	//	if len(n) != len(m[i]) {
	//		fmt.Println("\nNumgo/mat error.")
	//		s := "In mat.%s, Column %d of the 1st argument has %d elements,\n"
	//		s += "while the 2nd argument has %d rows. They must match.\n"
	//		s += fmt.Sprintf(s, "Dot", i, len(m[i]), len(n))
	//		panic(s)
	//	}
	//}
	res := New(len(m), len(n[0]))
	for i := range m {
		for j := range n[0] {
			for k := range m[i] {
				res[i][j] += m[i][k] * n[k][j]
			}
		}
	}
	return res
}

/*
DotC is the concurrent version of Dot(). This function spawns a goroutine
for each row of the first 2D slice which multiplies that row by each
column of 2nd 2D slice.

For sufficiently large slices, the performance of this function is very
close to that of Dot(). The previous statement is intentionally ambiguous,
and the clients of this library are encouraged to experiment for their
particular hardware and slice sizes.
*/
func DotC(m, n [][]float64) [][]float64 {
	// TODO: Add length checking.
	res := New(len(m), len(n[0]))
	var wg sync.WaitGroup
	for i := range m {
		wg.Add(1)
		go func(i int) {
			for j := range n[0] {
				for k := range m[i] {
					res[i][j] += m[i][k] * n[k][j]
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return res
}
