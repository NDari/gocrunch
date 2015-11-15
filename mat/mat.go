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
	"log"
	"math"
	"numgo/vec"
	"os"
	"runtime"
	"strconv"
)

/*
New returns a 2D slice of `float64` with the given number of row and columns.
This function should be used as a convenience tool, and it is exactly
equivalent to the normal method of allocating a uniform (non-jagged)
2D slice of `float64`.

If it is anticipated that the 2D slice will grow, use the "NewExpand"
function below. For full details, read that function's documentation.
*/
func New(r, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := 0; i < r; i++ {
		arr[i] = make([]float64, c)
	}
	return arr
}

/*
NewExpand returns a 2D slice of `float64`, with the given number of rows
and columns. The difference between this function and the "New" function
above is that the inner slices are allocated with double the capacity,
and hence can grow without the need for reallocation up to column * 2.

Note that this extended capacity will waste memory, so the NewExtend
should be used with care in situations where the performance gained by
avoiding reallocation justifies the extra cost in memory.
*/
func NewExpand(r, c int) [][]float64 {
	arr := make([][]float64, r)
	for i := 0; i < r; i++ {
		arr[i] = make([]float64, c*2)
	}
	return arr
}

/*
I returns an r by r 2D slice for a given r, where the elements along
the diagonal (where the first and the second index are equal) is set
to `1.0`, and all other elements are set to `0.0`.
*/
func I(r int) [][]float64 {
	identity := New(r, r)
	for i := 0; i < r; i++ {
		identity[i][i] = 1.0
	}
	return identity
}

/*
Ones returns a new 2D slice where all the elements are equal to `1.0`.
*/
func Ones(r, c int) [][]float64 {
	f := func(i float64) float64 {
		return 1.0
	}
	return Map(f, New(r, c))
}

/*
Inc returns a 2D slice, where element `[0][0] == 0.0`, and each
subsequent element is incremented by `1.0`.

For example:

```
m := Inc(3, 2)

mat.Print(m) // 1.0, 2.0
             // 3.0, 4.0
			 // 5.0, 6.0
```
*/
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

/*
Col returns a column of a 2D slice of `float64`. Col uses a 1-based index,
hence the first column of a 2D slice, m,  is `Col(1, m)`.

This function also allows for negative indexing. For example, `Col(-1, m)`
is the last column of the 2D slice m, and `Col(-2, m)` is the second to
last column of m, and so on.

Requesting the 0th column is fatal.
*/
func Col(c int, m [][]float64) []float64 {
	if math.Abs(float64(c)) > float64(len(m[0])) {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "The requested column, %v, is outside the range of %v to %v.\n"
		p, f, l, _ := runtime.Caller(1)
		lenm := len(m[0])
		log.Fatalf(msg, "Col", f, runtime.FuncForPC(p).Name(), l, c, -lenm, lenm)
	}
	vec := make([]float64, len(m))
	if c > 0 {
		for r := 0; r < len(m); r++ {
			vec[r] = m[r][c-1]
		}
	} else if c < 0 {
		lenColM := len(m[0])
		for r := 0; r < len(m); r++ {
			vec[r] = m[r][lenColM+c]
		}
	} else {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "The mat.%v function uses a 1-based index. Hence, the 0th column\n"
		msg += "is not defined.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Col", f, runtime.FuncForPC(p).Name(), l, "col")
	}
	return vec
}

/*
Row returns a row of a 2D slice of `float64`. Row uses a 1-based index, hence
the first row of a 2D slice, m, is Row(1, m).

This function also allows for negative indexing. For example, Row(-1, m) is
the last row of m.

Requesting the 0th row is fatal.
*/
func Row(r int, m [][]float64) []float64 {
	if math.Abs(float64(r)) > float64(len(m)) {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "The requested row, %v, is outside the range of %v to %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Col", f, runtime.FuncForPC(p).Name(), l, r, -len(m), len(m))
	}
	if r > 0 {
		return m[r]
	} else if r < 0 {
		return m[len(m)+r]
	} else {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "The mat.%v function uses a 1-based index. Hence, the 0th row\n"
		msg += "is not defined.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Row", f, runtime.FuncForPC(p).Name(), l, "Row")
		return make([]float64, 0) // This is ugly....
	}
}

/*
T returns a copy of a given 2D slice with the elements of the 2D slice
mirrored across the diagonal. For example, the element `[i][j]` becomes the
element `[j][i]` of the returned 2D slice. This function leaves the
original matrix intact.
*/
func T(m [][]float64) [][]float64 {
	transpose := New(len(m[0]), len(m))
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			transpose[j][i] = m[i][j]
		}
	}
	return transpose
}

/*
Equal checks if two 2D slices have the same shape and the same entries in
each row and column. If either the shape or the entries of the arguments
are different, `false` is returned. Otherwise, the return value is `true`.
*/
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

/*
Mul returns a new 2D slice that is the result of element-wise multiplication
of two 2D slices.
*/
func Mul(m, n [][]float64) [][]float64 {
	if len(m) != len(n) {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Number of rows of the first 2D slice is %v, while the number\n"
		msg += "of rows of the second 2D slice is %v. They must match.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Mul", f, runtime.FuncForPC(p).Name(), l, len(m), len(n))
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n[i]) {
			msg := "mat.%v Error: in %v [%v line %v].\n"
			msg += "In column %v, number of elements for the first 2D slice is %v,\n"
			msg += "while the number of elements of the second 2D slice is %v.\n"
			msg += "They must match.\n"
			p, f, l, _ := runtime.Caller(1)
			log.Fatalf(msg, "Mul", f, runtime.FuncForPC(p).Name(), l, i, len(m[i]), len(n[i]))
		}
		o[i] = vec.Mul(m[i], n[i])
	}
	return o
}

/*
Map calls a given elemental function on each Element of a 2D slice, returning
it afterwards. This function modifies the original 2D slice.
*/
func Map(f func(float64) float64, m [][]float64) [][]float64 {
	for i := 0; i < len(m); i++ {
		vec.MapInPlace(f, m[i])
	}
	return m
}

/*
Dot is the matrix multiplication of two 2D slices of `float64`.
*/
func Dot(m, n [][]float64) [][]float64 {
	lenm := len(m)
	// make sure that the length of the row of m matches the length of
	// each column in n.
	for i := 0; i < len(n); i++ {
		if lenm != len(n[i]) {
			msg := "mat.%v Error: in %v [%v line %v].\n"
			msg += "Length of column %v on the second matrix\n"
			msg += "is %v, which does not match the length of the row of the\n"
			msg += "first matrix, which is %v.\n"
			p, f, l, _ := runtime.Caller(1)
			log.Fatalf(msg, "Dot", f, runtime.FuncForPC(p).Name(), l, i, len(n[i]), lenm)
		}
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		if len(m[i]) != len(n) {
			msg := "mat.%v Error: in %v [%v line %v].\n"
			msg += "Length of column %v of the first matrix\n"
			msg += "is %v, which does not match the length of the row of the \n"
			msg += "second matrix, which is %v.\n"
			p, f, l, _ := runtime.Caller(1)
			log.Fatalf(msg, "Dot", f, runtime.FuncForPC(p).Name(), l, i, len(m[i]), len(n))
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

/*
Reset sets the values of all entries in a 2D slice of `float64` to `0.0`.
*/
func Reset(m [][]float64) [][]float64 {
	f := func(i float64) float64 {
		return 0.0
	}
	return Map(f, m)
}

/*
ToString converts a `[][]float64` to `[][]string`.
*/
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

/*
Dump prints the content of a `[][]float64` slice to a file, using comma as the
delimiter between the elements of a row, and a new line between rows.
*/
func Dump(m [][]float64, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Cannot open %v: %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Dump", f, runtime.FuncForPC(p).Name(), l, fileName, err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(ToString(m))
	if err = w.Error(); err != nil {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Error in CSV writer for file %v: %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Dump", f, runtime.FuncForPC(p).Name(), l, fileName, err)
	}
}

/*
FromString converts a `[][]string` to `[][]float64`.
*/
func FromString(str [][]string) [][]float64 {
	var err error
	m := make([][]float64, len(str))
	for i := 0; i < len(str); i++ {
		m[i] = make([]float64, len(str[i]))
		for j := 0; j < len(str[i]); j++ {
			m[i][j], err = strconv.ParseFloat(str[i][j], 64)
			if err != nil {
				msg := "mat.%v Error: in %v [%v line %v].\n"
				msg += "Died on string to float conversion at entry [%v][%v]: %v.\n"
				p, f, l, _ := runtime.Caller(1)
				log.Fatalf(msg, "FromString", f, runtime.FuncForPC(p).Name(), l, i, j, err)
			}
		}
	}
	return m
}

/*
Load generates a 2D slice of floats from a CSV file.
*/
func Load(fileName string) [][]float64 {
	f, err := os.Open(fileName)
	if err != nil {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Cannot open %v: %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Load", f, runtime.FuncForPC(p).Name(), l, fileName, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	str, err := r.ReadAll()
	if err != nil {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Error in CSV reader for file %v: %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Load", f, runtime.FuncForPC(p).Name(), l, fileName, err)
	}
	return FromString(str)
}

/*
Copy copies the content of a 2D slice of float64 into another with
the same shape. This is a deep copy, unlike the built in copy function
that is shallow for 2D slices.
*/
func Copy(m [][]float64) [][]float64 {
	n := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		n[i] = make([]float64, len(m[i]))
		copy(n[i], m[i])
	}
	return n
}

/*
AppendCol appends a column to the right side of a 2D slice of float64s.
*/
func AppendCol(m [][]float64, v []float64) [][]float64 {
	if len(m) != len(v) {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Number of rows of the first 2D slice is %v, while the number\n"
		msg += "of rows of the second 2D slice is %v. They must match.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "AppendCol", f, runtime.FuncForPC(p).Name(), l, len(m), len(v))
	}
	for i := 0; i < len(v); i++ {
		m[i] = append(m[i], v[i])
	}
	return m
}

/*
Concat concatenates the inner slices of two `[][]float64` arguments..

For example, if we have:

```
m := [[1.0, 2.0], [3.0, 4.0]]`
n := [[5.0, 6.0], [7.0, 8.0]]`
o := mat.Concat(m, n)`

mat.Print(o) // 1.0, 2.0, 5.0, 6.0
             // 3.0, 4.0, 7.0, 8.0
```
*/
func Concat(m, n [][]float64) [][]float64 {
	if len(m) != len(n) {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Number of rows of the first 2D slice is %v, while the number\n"
		msg += "of rows of the second 2D slice is %v. They must match.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Concat", f, runtime.FuncForPC(p).Name(), l, len(m), len(n))
	}
	o := make([][]float64, len(m))
	for i := 0; i < len(m); i++ {
		o[i] = make([]float64, len(m[i])+len(n[i]))
		o[i] = append(m[i], n[i]...)
	}
	return o
}

/*
Print prints a `[][]float64` to the stdout.
*/
func Print(m [][]float64) {
	w := csv.NewWriter(os.Stdout)
	w.Comma = rune(' ')
	w.WriteAll(ToString(m))
	if err := w.Error(); err != nil {
		msg := "mat.%v Error: in %v [%v line %v].\n"
		msg += "Error in CSV writer to stdout: %v.\n"
		p, f, l, _ := runtime.Caller(1)
		log.Fatalf(msg, "Print", f, runtime.FuncForPC(p).Name(), l, err)
	}
}
