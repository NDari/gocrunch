/*
Package mat implements a large set of functions which act on 2-dimensional slices
of float64.

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
	"runtime/debug"
	"strconv"
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
by appending all rows tip to tail. The passed 2D slice is assumed to be
non-jagged.
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
MulAll multiples all elements of a 2D slice by the passed value.
*/
func MulAll(m [][]float64, val float64) {
	for i := range m {
		for j := range m[i] {
			m[i][j] *= val
		}
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
	Rand(m)
the range is [0, 1)

For 1 argument, such as
	Rand(m, arg)
the range is [0, arg) for arg > 0, or (arg, 0] is arg < 0.

For 2 arguments, such as
	Rand(m, arg1, arg2)
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
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	}
}

/*
Col returns a new mat object whole values are equal to a column of the original
mat object. The number of Rows of the returned mat object is equal to the
number of rows of the original mat, and the number of columns is equal to 1.
*/
func Col(x int, m [][]float64) []float64 {
	if (x >= len(m[0])) || (x < -len(m[0])) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%s the requested column %d is outside of bounds [-%d, %d)\n"
		s = fmt.Sprintf(s, "Col", x, len(m[0]))
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
left intact.
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

//
///*
//Filter applies a function to each element of a mat object, creating a new
//mat object from all elements for which the function returned true. For
//example consider the following function:
//
//	positive := func(i *float64) bool {
//		if i > 0.0 {
//			return true
//		}
//		return false
//	}
//
//then calling
//
//	m.Filter(positive)
//
//will create a new mat object which contains the positive elements of the
//original matrix. If no elements return true for a given function, nil is
//returned. It is expected that the caller of this method checks the
//returned value to ensure that it is not nil.
//*/
//func (m *Mat) Filter(f BooleanFunc) *Mat {
//	var res []float64
//	for i := 0; i < m.r*m.c; i++ {
//		if f(&m.vals[i]) {
//			res = append(res, m.vals[i])
//		}
//	}
//	if len(res) == 0 {
//		return nil
//	}
//	return From1DSlice(res)
//}
//

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

	Any(positive, m)

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
Mul takes each element of the first 2D slice passed to it, and multiples that
element by the corresponding element in the second 2D slice passed to this
function, storing the result in the first 2D slice.

The shape of the 2D slices must be the same (same number or rows and columns),
and they are assumed to be non-jagged (same number of elements in each row).
*/
func Mul(m, n [][]float64) {
	if len(m) != len(n) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the rows of the first slice is %d\n"
		s += "but the number of rows of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Mul", len(m), len(n))
		panic(s)
	}
	if len(m[0]) != len(n[0]) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first slice is %d\n"
		s += "but the number of columns of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Mul", len(m[0]), len(n[0]))
		panic(s)
	}
	for i := range m {
		for j := range m[i] {
			m[i][j] *= n[i][j]
		}
	}
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
	if len(m[0]) != len(n[0]) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first slice is %d\n"
		s += "but the number of columns of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Add", len(m[0]), len(n[0]))
		panic(s)
	}
	for i := range m {
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
	if len(m[0]) != len(n[0]) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first slice is %d\n"
		s += "but the number of columns of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Sub", len(m[0]), len(n[0]))
		panic(s)
	}
	for i := range m {
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
	if len(m[0]) != len(n[0]) {
		fmt.Println("\nNumgo/mat error.")
		s := "In mat.%v, the number of the columns of the first slice is %d\n"
		s += "but the number of columns of the second slice is %d. They must\n"
		s += "match.\n"
		s = fmt.Sprintf(s, "Div", len(m[0]), len(n[0]))
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
		for j := range m[i] {
			m[i][j] /= n[i][j]
		}
	}
}

///*
//Sum returns the sum of the elements along a specific row or specific column.
//The first argument selects the row or column (0 or 1), and the second argument
//selects which row or column for which we want to calculate the sum. For
//example:
//
//	m.Sum(0, 2)
//
//will return the sum of the 3rd row of mat m.
//*/
//func (m *Mat) Sum(axis, slice int) float64 {
//	if axis != 0 && axis != 1 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the first argument must be 0 or 1, however %d "
//		s += "was recieved.\n"
//		s = fmt.Sprintf(s, "Sum", axis)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if axis == 0 {
//		if (slice >= m.r) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Sum", slice, m.r)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	} else if axis == 1 {
//		if (slice >= m.c) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Sum", slice, m.c)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	}
//	return m.precheckedSum(axis, slice)
//}
//
//func (m *Mat) precheckedSum(axis, slice int) float64 {
//	x := 0.0
//	if axis == 0 {
//		for i := 0; i < m.c; i++ {
//			x += m.vals[slice*m.c+i]
//		}
//	} else if axis == 1 {
//		for i := 0; i < m.r; i++ {
//			x += m.vals[i*m.c+slice]
//		}
//	}
//	return x
//}
//
///*
//Average returns the average of the elements along a specific row or specific
//column.
//The first argument selects the row or column (0 or 1), and the second argument
//selects which row or column for which we want to calculate the average. For
//example:
//
//	m.Average(0, 2)
//
//will return the average of the 3rd row of mat m.
//*/
//func (m *Mat) Average(axis, slice int) float64 {
//	if axis != 0 && axis != 1 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the first argument must be 0 or 1, however %d "
//		s += "was recieved.\n"
//		s = fmt.Sprintf(s, "Average", axis)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if axis == 0 {
//		if (slice >= m.r) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Average", slice, m.r)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	} else if axis == 1 {
//		if (slice >= m.c) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Average", slice, m.c)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	}
//	return m.precheckedAverage(axis, slice)
//}
//
//func (m *Mat) precheckedAverage(axis, slice int) float64 {
//	sum := m.precheckedSum(axis, slice)
//	if axis == 0 {
//		return sum / float64(m.c)
//	}
//	return sum / float64(m.r)
//}
//
///*
//Prod returns the product of the elements along a specific row or specific
//column.
//The first argument selects the row or column (0 or 1), and the second argument
//selects which row or column for which we want to calculate the product. For
//example:
//
//	m.Prod(1, 2)
//
//will return the product of the 3rd column of mat m.
//*/
//func (m *Mat) Prod(axis, slice int) float64 {
//	if axis != 0 && axis != 1 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the first argument must be 0 or 1, however %d "
//		s += "was recieved.\n"
//		s = fmt.Sprintf(s, "Prod", axis)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if axis == 0 {
//		if (slice >= m.r) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the requested row %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Prod", slice, m.r)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	} else if axis == 1 {
//		if (slice >= m.c) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Prod", slice, m.c)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	}
//	x := 1.0
//	if axis == 0 {
//		for i := 0; i < m.c; i++ {
//			x *= m.vals[slice*m.c+i]
//		}
//	} else if axis == 1 {
//		for i := 0; i < m.r; i++ {
//			x *= m.vals[i*m.c+slice]
//		}
//	}
//	return x
//}
//
///*
//Std returns the standard deviation of the elements along a specific row
//or specific column. The standard deviation is defined as the square root of
//the mean distance of each element from the mean. Look at:
//http://mathworld.wolfram.com/StandardDeviation.html
//
//For example:
//
//	m.Std(1, 0)
//
//will return the standard deviation of the first column of mat m.
//*/
//func (m *Mat) Std(axis, slice int) float64 {
//	if axis != 0 && axis != 1 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the first argument must be 0 or 1, however %d "
//		s += "was recieved.\n"
//		s = fmt.Sprintf(s, "Std", axis)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if axis == 0 {
//		if (slice >= m.r) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the row %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Std", slice, m.r)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	} else if axis == 1 {
//		if (slice >= m.c) || (slice < 0) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the column %d is outside of bounds [0, %d)\n"
//			s = fmt.Sprintf(s, "Std", slice, m.c)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//	}
//	avg := m.precheckedAverage(axis, slice)
//	var s []float64
//	if axis == 0 {
//		s = make([]float64, m.c)
//		for i := 0; i < m.c; i++ {
//			s[i] = avg - m.vals[slice*m.c+i]
//			s[i] *= s[i]
//		}
//	} else {
//		s = make([]float64, m.r)
//		for i := 0; i < m.r; i++ {
//			s[i] = avg - m.vals[i*m.c+slice]
//			s[i] *= s[i]
//		}
//	}
//	sum := 0.0
//	for i := range s {
//		sum += s[i]
//	}
//	return math.Sqrt(sum / float64(len(s)))
//}
//
///*
//Dot is the matrix multiplication of two mat objects. Consider the following two
//mats:
//
//	m := New(5, 6)
//	n := New(6, 10)
//
//then
//
//	o := m.Dot(n)
//
//is a 5 by 10 mat whose element at row i and column j is given by:
//
//	Sum(m.Row(i).Mul(n.col(j))
//*/
//func (m *Mat) Dot(n *Mat) *Mat {
//	if m.c != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of columns of the first mat is %d\n"
//		s += "which is not equal to the number of rows of the second mat,\n"
//		s += "which is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "Dot", m.c, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	o := New(m.r, n.c)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < n.c; j++ {
//			for k := 0; k < m.c; k++ {
//				o.vals[i*o.c+j] += (m.vals[i*m.c+k] * n.vals[k*n.c+j])
//			}
//		}
//	}
//	return o
//}
//
///*
//ToString returns the string representation of a mat. This is done by putting
//every row into a line, and separating the entries of that row by a space. note
//that the last line does not contain a newline.
//*/
//func (m *Mat) ToString() string {
//	var str string
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			str += strconv.FormatFloat(m.vals[i*m.c+j], 'f', 14, 64)
//			str += " "
//		}
//		if i+1 <= m.r {
//			str += "\n"
//		}
//	}
//	return str
//}
//
///*
//AppendCol appends a column to the right side of a Mat.
//*/
//func (m *Mat) AppendCol(v []float64) *Mat {
//	if m.r != len(v) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of rows of the reciever is %d, while\n"
//		s += "the number of rows of the vector is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "AppendCol", m.r, len(v))
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	// TODO: redo this by hand, instead of taking this shortcut... or check if
//	// this is a huge bottleneck
//	q := m.ToSlice()
//	for i := range q {
//		q[i] = append(q[i], v[i])
//	}
//	m.c++
//	m.vals = append(m.vals, v...)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			m.vals[i*m.c+j] = q[i][j]
//		}
//	}
//	return m
//}
//
///*
//AppendRow appends a row to the bottom of a Mat.
//*/
//func (m *Mat) AppendRow(v []float64) *Mat {
//	if m.c != len(v) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of cols of the reciever is %d, while\n"
//		s += "the number of rows of the vector is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "AppendCol", m.c, len(v))
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	m.vals = append(m.vals, v...)
//	m.r++
//	return m
//}
//
///*
//Concat concatenates the inner slices of two `[][]float64` arguments..
//
//For example, if we have:
//
//	m := [[1.0, 2.0], [3.0, 4.0]]
//	n := [[5.0, 6.0], [7.0, 8.0]]
//	o := mat.Concat(m, n).Print // 1.0, 2.0, 5.0, 6.0
//															// 3.0, 4.0, 7.0, 8.0
//
//*/
//func (m *Mat) Concat(n *Mat) *Mat {
//	if m.r != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of rows of the receiver is %d, while\n"
//		s += "the number of rows of the second Mat is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "Concat", m.r, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	q := m.ToSlice()
//	t := n.Vals()
//	r := n.ToSlice()
//	m.vals = append(m.vals, t...)
//	for i := range q {
//		q[i] = append(q[i], r[i]...)
//	}
//	m.c += n.c
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			m.vals[i*m.c+j] = q[i][j]
//		}
//	}
//	return m
//}
//
///*
//Print displays the content of a Mat to the screen.
//*/
//func (m *Mat) Print() {
//	fmt.Println(m.ToString())
//}
//
///*
//Set sets the value of a mat at a given row and column to the given value
//*/
//func (m *Mat) Set(r, c int, val float64) *Mat {
//	if (r >= m.r) || (r < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested row %d is outside of bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "Set", r, m.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if (c >= m.c) || (c < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested column %d is outside of bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "Set", r, m.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	m.vals[r*m.c+c] = val
//	return m
//}
//
///*
//CombineWithRows combines a slice elementally with each of the rows in a Mat. The allowed
//combinations are ["add", "sub", "mul", "div"]. Consider:
//	v := make([]float64, 5)
//	for i := range v {
//		v[i] = float64(i) // v is now [0.0, 1.0, 2.0, 3.0, 4.0]
//	}
//	m := mat.New(2, 5).Inc() // note that m's number of columns is equal to len(v)
//	m.Row(0).Print() // 0.0 1.0 2.0 3.0 4.0
//	m.Row(1).Print() // 5.0 6.0 7.0 8.0 9.0
//	m.CombineWithRows("add", v)
//	m.Row(0).Print() // 0.0 2.0 4.0 6.0 8.0
//	m.Row(1).Print() // 5.0 7.0 9.0 11.0 13.0
//In other words, each element of v is added to the corresponding element of each row of
//m.
//
//Note that for the combination method of "Div", all elements of the passed slice must be
//non-zero to avoid division by zero.
//*/
//func (m *Mat) CombineWithRows(how string, v []float64) *Mat {
//	if m.c != len(v) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of cols of the reciever is %d, while\n"
//		s += "the number of rows of the vector is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "AddToRows", m.c, len(v))
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	switch how {
//	case "add":
//		for i := 0; i < m.r; i++ {
//			for j := 0; j < m.c; j++ {
//				m.vals[i*m.c+j] += v[j]
//			}
//		}
//	case "sub":
//		for i := 0; i < m.r; i++ {
//			for j := 0; j < m.c; j++ {
//				m.vals[i*m.c+j] -= v[j]
//			}
//		}
//	case "mul":
//		for i := 0; i < m.r; i++ {
//			for j := 0; j < m.c; j++ {
//				m.vals[i*m.c+j] *= v[j]
//			}
//		}
//	case "div":
//		for i := range v {
//			if v[i] == 0.0 {
//				fmt.Println("\nNumgo/mat error.")
//				s := "In mat.%s a zero-valued element was found in v at index %d.\n"
//				s += "Division by zero is not allowed.\n"
//				s = fmt.Sprintf(s, "CombineWithRows", i)
//				fmt.Println(s)
//				fmt.Println("Stack trace for this error:")
//				debug.PrintStack()
//				os.Exit(1)
//			}
//		}
//		for i := 0; i < m.r; i++ {
//			for j := 0; j < m.c; j++ {
//				m.vals[i*m.c+j] /= v[j]
//			}
//		}
//	default:
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the allowed combination methods are ['add', 'sub', 'mul', 'div'].\n"
//		s += "However, %s was recieved.\n"
//		s = fmt.Sprintf(s, "CombineWithRows", how)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	return m
//}
//
///*
//CombineWithCols combines a slice elementally with each of the columns in a Mat. The allowed
//combinations are ["add", "sub", "mul", "div"]. Consider:
//	v := make([]float64, 5)
//	for i := range v {
//		v[i] = float64(i) // v is now [0.0, 1.0, 2.0, 3.0, 4.0]
//	}
//	m := mat.New(5, 2).Inc() // note that m's number of rows is equal to len(v)
//	m.Col(0).Print() // 0.0
//			 // 1.0
//			 // 2.0
//			 // 3.0
//			 // 4.0
//	m.Col(1).Print() // 5.0
//			 // 6.0
//			 // 7.0
//			 // 8.0
//			 // 9.0
//	m.CombineWithCols("add", v)
//	m.Col(0).Print() // 0.0
//	                 // 2.0
//	                 // 4.0
//	                 // 6.0
//	                 // 8.0
//	m.Col(1).Print() // 5.0
//	                 // 7.0
//	                 // 9.0
//	                 // 11.0
//	                 // 13.0
//In other words, each element of v is added to the corresponding element of each column of
//m.
//
//Note that for the combination method of "Div", all elements of the passed slice must be
//non-zero to avoid division by zero.
//*/
//func (m *Mat) CombineWithCols(how string, v []float64) *Mat {
//	if m.r != len(v) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the number of columns of the reciever is %d, while\n"
//		s += "the number of rows of the vector is %d. They must be equal.\n"
//		s = fmt.Sprintf(s, "AddToRows", m.r, len(v))
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	switch how {
//	case "add":
//		for i := 0; i < m.c; i++ {
//			for j := 0; j < m.r; j++ {
//				m.vals[j*m.c+i] += v[j]
//			}
//		}
//	case "sub":
//		for i := 0; i < m.c; i++ {
//			for j := 0; j < m.r; j++ {
//				m.vals[j*m.c+i] -= v[j]
//			}
//		}
//	case "mul":
//		for i := 0; i < m.c; i++ {
//			for j := 0; j < m.r; j++ {
//				m.vals[j*m.c+i] *= v[j]
//			}
//		}
//	case "div":
//		for i := range v {
//			if v[i] == 0.0 {
//				fmt.Println("\nNumgo/mat error.")
//				s := "In mat.%s a zero-valued element was found in v at index %d.\n"
//				s += "Division by zero is not allowed.\n"
//				s = fmt.Sprintf(s, "CombineWithRows", i)
//				fmt.Println(s)
//				fmt.Println("Stack trace for this error:")
//				debug.PrintStack()
//				os.Exit(1)
//			}
//		}
//		for i := 0; i < m.c; i++ {
//			for j := 0; j < m.r; j++ {
//				m.vals[j*m.c+i] /= v[j]
//			}
//		}
//	default:
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the allowed combination methods are ['add', 'sub', 'mul', 'div'].\n"
//		s += "However, %s was recieved.\n"
//		s = fmt.Sprintf(s, "CombineWithRows", how)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	return m
//}
