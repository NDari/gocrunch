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
	"os"
	"strconv"
)

type elementFunc func(*float64)
type booleanFunc func(*float64) bool
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
assume that the number of rows of the resultant mat object is equal to the
number of lines, and the number of columns is equal to the number of entries
in each line. As before, we make sure that each line contains the same number
of elements.

The file to be read is assumed to be very large, and hence it is read one line
at a time. This results in some major inefficiencies, and it is recommended
that this function be used sparingly, and not as a major component of your
library/executable.

Unlike other mat creation methods in this package, the mat object created here,
the capacity of the mat object created here is the same as its length since we
assume the mat to be very large.
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
			s = fmt.Sprintf(s, "Load", filename, err)
			panic(s)
		}
	}
	return m
}

///*
//Flatten...
//*/
//func Flatten(m [][]float64) []float64 {
//	// things...
//}
//
///*
//Reshape changes the row and the columns of the mat object as long as the total
//number of values contained in the mat object remains constant. The order and
//the values of the mat does not change with this function.
//*/
//func (m *Mat) Reshape(rows, cols int) *Mat {
//	if rows <= 0 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s, the number of rows must be greater than '0', but\n"
//		s += "recieved %d. "
//		s = fmt.Sprintf(s, "Reshape", rows)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if cols <= 0 {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s, the number of columns must be greater than '0', but\n"
//		s += "recieved %d. "
//		s = fmt.Sprintf(s, "Reshape", cols)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if rows*cols != m.r*m.c {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s, The total number of entries of the old and new shape\n"
//		s += "must match.\n"
//		s = fmt.Sprintf(s, "Reshape")
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	} else {
//		m.r = rows
//		m.c = cols
//	}
//	return m
//}
//
///*
//Dims returns the number of rows and columns of a mat object.
//*/
//func (m *Mat) Dims() (int, int) {
//	return m.r, m.c
//}
//
///*
//Vals returns the values contained in a mat object as a 1D slice of float64s.
//*/
//func (m *Mat) Vals() []float64 {
//	s := make([]float64, m.r*m.c)
//	copy(s, m.vals)
//	return s
//}
//
///*
//ToSlice returns the values of a mat object as a 2D slice of float64s.
//*/
//func (m *Mat) ToSlice() [][]float64 {
//	s := make([][]float64, m.r)
//	for i := 0; i < m.r; i++ {
//		s[i] = make([]float64, m.c)
//		for j := 0; j < m.c; j++ {
//			s[i][j] = m.vals[i*m.c+j]
//		}
//	}
//	return s
//}
//
///*
//ToCSV creates a file with the passed name, and writes the content of a mat
//object to it, by putting each row in a single comma separated line. The
//number of entries in each line is equal to the columns of the mat object.
//*/
//func (m *Mat) ToCSV(fileName string) {
//	f, err := os.Create(fileName)
//	if err != nil {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, cannot open %s due to error: %v.\n"
//		s = fmt.Sprintf(s, "ToCSV", fileName, err)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	defer f.Close()
//	str := ""
//	idx := 0
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			str += strconv.FormatFloat(m.vals[idx], 'e', 14, 64)
//			if j+1 != m.c {
//				str += ","
//			}
//			idx++
//		}
//		if i+1 != m.r {
//			str += "\n"
//		}
//	}
//	_, err = f.Write([]byte(str))
//	if err != nil {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, cannot write to %s due to error: %v.\n"
//		s = fmt.Sprintf(s, "ToCSV", fileName, err)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//}
//
///*
//At returns a pointer to the float64 stored in the given row and column.
//*/
//func (m *Mat) At(r, c int) float64 {
//	if (r >= m.r) || (r < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested row %d is outside of bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "At", r, m.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if (c >= m.c) || (c < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested column %d is outside of bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "At", r, m.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	return m.vals[r*m.c+c]
//}
//
///*
//Map applies a given function to each element of a mat object. The given
//function must take a pointer to a float64, and return nothing.
//*/
//func (m *Mat) Map(f elementFunc) *Mat {
//	for i := 0; i < m.r*m.c; i++ {
//		f(&m.vals[i])
//	}
//	return m
//}
//
///*
//Ones sets all values of a mat to be equal to 1.0
//*/
//func (m *Mat) Ones() *Mat {
//	f := func(i *float64) {
//		*i = 1.0
//		return
//	}
//	return m.Map(f)
//}
//
///*
//Inc takes each element of a mat object, and starting from 0.0, sets their
//value to be the value of the previous entry plus 1.0. In other words, the
//first few values of a mat object after this operation would be 0.0, 1.0,
//2.0, ...
//*/
//func (m *Mat) Inc() *Mat {
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] = float64(i)
//	}
//	return m
//}
//
///*
//Reset sets all values of a mat object to 0.0
//*/
//func (m *Mat) Reset() *Mat {
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] = 0.0
//	}
//	return m
//}
//
///*
//SetAllTo sets all values of a mat to the passed float64 value.
//*/
//func (m *Mat) SetAllTo(val float64) *Mat {
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] = val
//	}
//	return m
//}
//
///*
//Rand sets the values of a mat to random numbers. The range from which the random
//numbers are selected is determined based on the arguments passed.
//
//For no arguments, such as
//	m.Rand()
//the range is [0, 1)
//
//For 1 argument, such as
//	m.Rand(arg)
//the range is [0, arg) for arg > 0, or (arg, 0] is arg < 0.
//
//For 2 arguments, such as
//	m.Rand(arg1, arg2)
//the range is [arg1, arg2).
//*/
//func (m *Mat) Rand(args ...float64) *Mat {
//	switch len(args) {
//	case 0:
//		for i := 0; i < m.r*m.c; i++ {
//			m.vals[i] = rand.Float64()
//		}
//	case 1:
//		to := args[0]
//		for i := 0; i < m.r*m.c; i++ {
//			m.vals[i] = rand.Float64() * to
//		}
//	case 2:
//		from := args[0]
//		to := args[1]
//		if !(from < to) {
//			fmt.Println("\nNumgo/mat error.")
//			s := "In mat.%s the first argument, %f, is not less than the\n"
//			s += "second argument, %f. The first argument must be strictly\n"
//			s += "less than the second.\n"
//			s = fmt.Sprintf(s, "Rand", from, to)
//			fmt.Println(s)
//			fmt.Println("Stack trace for this error:")
//			debug.PrintStack()
//			os.Exit(1)
//		}
//		for i := 0; i < m.r*m.c; i++ {
//			m.vals[i] = rand.Float64()*(to-from) + from
//		}
//	default:
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s expected 0 to 2 arguments, but recieved %d."
//		s = fmt.Sprintf(s, "Rand", len(args))
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	return m
//}
//
///*
//Col returns a new mat object whole values are equal to a column of the original
//mat object. The number of Rows of the returned mat object is equal to the
//number of rows of the original mat, and the number of columns is equal to 1.
//*/
//func (m *Mat) Col(x int) *Mat {
//	if (x >= m.c) || (x < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested column %d is outside of bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "Col", x, m.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	v := New(m.r, 1)
//	for r := 0; r < m.r; r++ {
//		v.vals[r] = m.vals[r*m.c+x]
//	}
//	return v
//}
//
///*
//Cols returns a generator which, upon invocation, returns the next column of
//a Mat in form of a Mat with 1 column, and the same number of rows of the
//method receiver. Consider the following:
//
//	m := mat.New(3, 2).Inc()
//	for col := range m.Cols() {
//		col.Print()
//	}
//
//The col.Print() above prints:
//	// 0.0
//	// 2.0
//	// 4.0
//
//and then
//	// 1.0
//	// 3.0
//	// 5.0
//*/
//func (m *Mat) Cols() <-chan *Mat {
//	res := make(chan *Mat, m.c)
//	go func() {
//		for i := 0; i < m.c; i++ {
//			res <- m.Col(i)
//		}
//		close(res)
//	}()
//	return res
//}
//
///*
//Row returns a new mat object whose values are equal to a row of the original
//mat object. The number of Rows of the returned mat object is equal to 1, and
//the number of columns is equal to the number of columns of the original mat.
//*/
//func (m *Mat) Row(x int) *Mat {
//	if (x >= m.r) || (x < 0) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%s the requested row %d is outside of the bounds [0, %d)\n"
//		s = fmt.Sprintf(s, "Row", x, m.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	v := New(1, m.c)
//	for r := 0; r < m.c; r++ {
//		v.vals[r] = m.vals[x*m.c+r]
//	}
//	return v
//}
//
///*
//Rows returns a generator which, upon invocation, returns the next row of
//a Mat in form of a Mat with 1 row, and the same number of columns of the
//method receiver. Consider the following:
//
//	m := mat.New(3, 2).Inc()
//	m.Print() // 0.0 1.0
//	          // 2.0 3.0
//			  // 4.0 5.0
//
//	for row := range m.Rows() {
//		row.Print()
//	}
//
//The col.Print() above prints:
//	// 0.0 1.0
//
//and then
//	// 2.0 3.0
//
//and finally
//	// 4.0 5.0
//*/
//func (m *Mat) Rows() <-chan *Mat {
//	res := make(chan *Mat, m.r)
//	go func() {
//		for i := 0; i < m.r; i++ {
//			res <- m.Row(i)
//		}
//		close(res)
//	}()
//	return res
//}
//
///*
//Equals checks to see if two mat objects are equal. That mean that the two mats
//have the same number of rows, same number of columns, and have the same float64
//in each entry at a given index.
//*/
//func (m *Mat) Equals(n *Mat) bool {
//	if m.r != n.r {
//		return false
//	}
//	if m.c != n.c {
//		return false
//	}
//	for i := 0; i < m.r*m.c; i++ {
//		if m.vals[i] != n.vals[i] {
//			return false
//		}
//	}
//	return true
//}
//
///*
//Copy returns a duplicate of a mat object. The returned copy is "deep", meaning
//that the object can be manipulated without effecting the original mat object.
//*/
//func (m *Mat) Copy() *Mat {
//	n := New(m.r, m.c)
//	copy(n.vals, m.vals)
//	return n
//}
//
///*
//T returns the transpose of the original matrix. The transpose of a mat object
//is defined in the usual manner, where every value at row x, and column y is
//placed at row y, and column x. The number of rows and column of the transposed
//mat are equal to the number of columns and rows of the original matrix,
//respectively. This method creates a new mat object, and the original is
//left intact.
//*/
//func (m *Mat) T() *Mat {
//	n := New(m.c, m.r)
//	idx := 0
//	for i := 0; i < m.c; i++ {
//		for j := 0; j < m.r; j++ {
//			n.vals[idx] = m.vals[j*m.c+i]
//			idx++
//		}
//	}
//	return n
//}
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
//func (m *Mat) Filter(f booleanFunc) *Mat {
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
///*
//All checks if a supplied function is true for all elements of a mat object.
//For instance, consider
//
//	positive := func(i *float64) bool {
//		if i > 0.0 {
//			return true
//		}
//		return false
//	}
//
//Then calling
//
//	m.All(positive)
//
//will return true if and only if all elements in m are positive.
//*/
//func (m *Mat) All(f booleanFunc) bool {
//	for i := 0; i < m.r*m.c; i++ {
//		if !f(&m.vals[i]) {
//			return false
//		}
//	}
//	return true
//}
//
///*
//Any checks if a supplied function is true for one elements of a mat object.
//For instance,
//
//	positive := func(i *float64) bool {
//		if i > 0.0 {
//			return true
//		}
//		return false
//	}
//
//Then calling
//
//	m.Any(positive)
//
//would be true if at least one element of the mat object is positive.
//*/
//func (m *Mat) Any(f booleanFunc) bool {
//	for i := 0; i < m.r*m.c; i++ {
//		if f(&m.vals[i]) {
//			return true
//		}
//	}
//	return false
//}
//
///*
//Mul is the element-wise multiplication of a mat object by another which is
//passed to this method.
//
//The shape of the mat objects must be the same (same number or rows and columns)
//and the results of the element-wise multiplication is stored in the original
//mat on which the method was invoked.
//*/
//func (m *Mat) Mul(n *Mat) *Mat {
//	if m.r != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the rows of the first mat is %d\n"
//		s += "but the number of rows of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Mul", m.r, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if m.c != n.c {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the columns of the first mat is %d\n"
//		s += "but the number of columns of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Mul", m.c, n.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] *= n.vals[i]
//	}
//	return m
//}
//
///*
//Add is the element-wise addition of a mat object with another which is passed
//to this method.
//
//The shape of the mat objects must be the same (same number or rows and columns)
//and the results of the element-wise addition is stored in the original
//mat on which the method was invoked.
//*/
//func (m *Mat) Add(n *Mat) *Mat {
//	if m.r != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the rows of the first mat is %d\n"
//		s += "but the number of rows of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Add", m.r, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if m.c != n.c {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the columns of the first mat is %d\n"
//		s += "but the number of columns of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Add", m.c, n.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] += n.vals[i]
//	}
//	return m
//}
//
///*
//Sub is the element-wise subtraction of a mat object which is passed
//to this method from the original mat which called the method.
//
//The shape of the mat objects must be the same (same number or rows and columns)
//and the results of the element-wise subtraction is stored in the original
//mat on which the method was invoked.
//*/
//func (m *Mat) Sub(n *Mat) *Mat {
//	if m.r != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the rows of the first mat is %d\n"
//		s += "but the number of rows of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Sub", m.r, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if m.c != n.c {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the columns of the first mat is %d\n"
//		s += "but the number of columns of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Sub", m.c, n.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] -= n.vals[i]
//	}
//	return m
//}
//
///*
//Div is the element-wise dicition of a mat object by another which is passed
//to this method.
//
//The shape of the mat objects must be the same (same number or rows and columns)
//and the results of the element-wise division is stored in the original
//mat on which the method was invoked. The dividing mat object (passed to this
//method) must not contain any elements which are equal to 0.0.
//*/
//func (m *Mat) Div(n *Mat) *Mat {
//	if m.r != n.r {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the rows of the first mat is %d\n"
//		s += "but the number of rows of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Div", m.r, n.r)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	if m.c != n.c {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, the number of the columns of the first mat is %d\n"
//		s += "but the number of columns of the second mat is %d. They must\n"
//		s += "match.\n"
//		s = fmt.Sprintf(s, "Div", m.c, n.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	zero := func(i *float64) bool {
//		if *i == 0.0 {
//			return true
//		}
//		return false
//	}
//	if n.Any(zero) {
//		fmt.Println("\nNumgo/mat error.")
//		s := "In mat.%v, one or more elements of the second matrix are 0.0\n"
//		s += "Division by zero is not allowed.\n"
//		s = fmt.Sprintf(s, "Div", m.c, n.c)
//		fmt.Println(s)
//		fmt.Println("Stack trace for this error:")
//		debug.PrintStack()
//		os.Exit(1)
//	}
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] /= n.vals[i]
//	}
//	return m
//}
//
///*
//Scale is the element-wise multiplication of a mat object by a scalar.
//
//The results of the element-wise multiplication is stored in the original
//mat on which the method was invoked.
//*/
//func (m *Mat) Scale(f float64) *Mat {
//	for i := 0; i < m.r*m.c; i++ {
//		m.vals[i] *= f
//	}
//	return m
//}
//
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
