# mat
--
    import "github.com/NDari/numgo/mat"

Package mat implements a "mat" object, which behaves like a 2-dimensional array
or list in other programming languages. Under the hood, the mat object is a flat
slice, which provides for optimal performance in Go, while the methods and
constructors provide for a higher level of performance and abstraction when
compared to the "2D" slices of go (slices of slices).

All errors encountered in this package, such as attempting to access an element
out of bounds are treated as critical error, and thus, the code immediately
exits with signal 1. In such cases, the function/method in which the error was
encountered is printed to the screen, in addition to the full stack trace, in
order to help fix the issue rapidly.

## Usage

#### type Mat

```go
type Mat struct {
}
```

Mat is the main struct of this library. Mat is a essentially a 1D slice (a
[]float64) that contains two integers, representing rows and columns, which
allow it to behave as if it was a 2D slice. This allows for higher performance
and flexibility for the users of this library, at the expense of some
bookkeeping that is done here.

The fields of this struct are not directly accessible, and they may only change
by the use of the various methods in this library.

#### func  From1DSlice

```go
func From1DSlice(s []float64) *Mat
```
From1DSlice creates a mat object from a slice of float64s. The created mat
object has one row, and the number of columns equal to the length of the 1D
slice from which it was created.

#### func  FromCSV

```go
func FromCSV(filename string) *Mat
```
FromCSV creates a mat object from a CSV (comma separated values) file. Here, we
assume that the number of rows of the resultant mat object is equal to the
number of lines, and the number of columns is equal to the number of entries in
each line. As before, we make sure that each line contains the same number of
elements.

The file to be read is assumed to be very large, and hence it is read one line
at a time. This results in some major inefficiencies, and it is recommended that
this function be used sparingly, and not as a major component of your
library/executable.

Unline other mat creation methods in this package, the mat object created here,
the capacity of the mat opject created here is the same as its length since we
assume the mat to be very large.

#### func  FromSlice

```go
func FromSlice(s [][]float64) *Mat
```
FromSlice generated a mat object from a [][]float64 slice. The slice is checked
for being a non-jagged slice, where each row contains the same number of
elements. The creation of a mat object from jagged 2D slices is not supported as
on yet.

#### func  New

```go
func New(dims ...int) *Mat
```
New is the primary constructor for the "Mat" object. New is a veradic function,
expecting 0 to 3 ints, with differing behavior as follows:

```go
m := New()
```

m is now an empty &Mat{}, where the number of rows, columns and the length and
capacity of the underlying slice are all zero. This is mostly for internal use.

```go
m := New(x)
```

m is a x by x (square) matrix, with the underlying slice of length x, and
capacity 2x.

```go
m := New(x, y)
```

m is an x by y matrix, with the underlying slice of length rc, and capacity of
2rc. This is a good case for when your matrix is going to expand in the future.
There is a negligible hit to performance and a larger memory usage of your code.
But in case expanding matrices, many reallocations are avoided.

```go
m := New(x, y, z)
```

m is a x by u matrix, with the underlying slice of length rc, and capacity z.
This is a good choice for when the size of the matrix is static, or when the
application is memory constrained.

For most cases, we recommend using the New(x) or New(x, y) options, and almost
never the New() option.

#### func (*Mat) Add

```go
func (m *Mat) Add(n *Mat) *Mat
```
Add is the element-wise addition of a mat object with another which is passed to
this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise addition is stored in the original mat on
which the method was invoked.

#### func (*Mat) All

```go
func (m *Mat) All(f booleanFunc) bool
```
All checks if a supplied function is true for all elements of a mat object. For
instance, consider

```go
positive := func(i *float64) bool {
    if i > 0.0 {
    	return true
    } else {
    	return false
    }
}
```


Then calling

```go
m.All(positive)
```

will return true if and only if all elements in m are positive.

#### func (*Mat) Any

```go
func (m *Mat) Any(f booleanFunc) bool
```
Any checks if a supplied function is true for one elements of a mat object. For
instance,

```go
positive := func(i *float64) bool {
    if i > 0.0 {
    	return true
    } else {
    	return false
    }
}
```


Then calling

```go
m.Any(positive)
```

would be true if at least one element of the mat object is positive.

#### func (*Mat) AppendCol

```go
func (m *Mat) AppendCol(v []float64) *Mat
```
AppendCol appends a column to the right side of a Mat. TODO: Fix this... the new
object is not returned.

#### func (*Mat) AppendRow

```go
func (m *Mat) AppendRow(v []float64) *Mat
```
AppendRow appends a row to the bottom of a Mat.

#### func (*Mat) At

```go
func (m *Mat) At(r, c int) float64
```
At returns a pointer to the float64 stored in the given row and column.

#### func (*Mat) Average

```go
func (m *Mat) Average(axis, slice int) float64
```
Average returns the average of the elements along a specific row or specific
column. The first argument selects the row or column (0 or 1), and the second
argument selects which row or column for which we want to calculate the average.
For example:

```go
m.Average(0, 2)
```

will calculate the average of the 3rd row of mat m.

#### func (*Mat) Col

```go
func (m *Mat) Col(x int) *Mat
```
Col returns a new mat object whole values are equal to a column of the original
mat object. The number of Rows of the returned mat object is equal to the number
of rows of the original mat, and the number of columns is equal to 1.

#### func (*Mat) Concat

```go
func (m *Mat) Concat(n *Mat) *Mat
```
Concat concatenates the inner slices of two `[][]float64` arguments..

For example, if we have:

```go
m := [[1.0, 2.0], [3.0, 4.0]] 
n := [[5.0, 6.0], [7.0, 8.0]] 
o := mat.Concat(m, n)

mat.Print(o) // 1.0, 2.0, 5.0, 6.0
             // 3.0, 4.0, 7.0, 8.0

```

#### func (*Mat) Copy

```go
func (m *Mat) Copy() *Mat
```
Copy returns a duplicate of a mat object. The returned copy is "deep", meaning
that the object can be manipulated without effecting the original mat object.

#### func (*Mat) Dims

```go
func (m *Mat) Dims() (int, int)
```
Dims returns the number of rows and columns of a mat object.

#### func (*Mat) Div

```go
func (m *Mat) Div(n *Mat) *Mat
```
Div is the element-wise dicition of a mat object by another which is passed to
this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise divition is stored in the original mat on
which the method was invoked. The dividing mat object (passed to this method)
must not contain any elements which are equal to 0.0.

#### func (*Mat) Dot

```go
func (m *Mat) Dot(n *Mat) *Mat
```
Dot is the matrix multiplication of two mat objects. Consider the following two
mats:

```go
m := New(5, 6) 
n := New(6, 10)
```

then

```go
o := m.Dot(n)
```

is a 5 by 10 mat whose element at row i and column j is given by:

```go
Sum(m.Row(i).Mul(n.col(j))
```

#### func (*Mat) Equals

```go
func (m *Mat) Equals(n *Mat) bool
```
Equals checks to see if two mat objects are equal. That mean that the two mats
have the same number of rows, same number of columns, and have the same float64
in each entry at a given index.

#### func (*Mat) Filter

```go
func (m *Mat) Filter(f booleanFunc) *Mat
```
Filter applies a function to each element of a mat object, creating a new mat
object from all elements for which the function returned true. For example
consider the following function:

```go
f := func(i *float64) bool {
    if i > 0.0 {
    	return true
    } else {
    	return false
    }
}
```

then calling

```go
m.Filter(f)
```

will create a new mat object which contains the positive elements of the
original matrix. If no elements return true for a given function, nil is
returned. It is expected that the caller of this method checks the returned
value to ensure that it is not nil.

#### func (*Mat) Inc

```go
func (m *Mat) Inc() *Mat
```
Inc takes each element of a mat object, and starting from 0.0, sets their value
to be the value of the previous entry plus 1.0. In other words, the first few
values of a mat object after this operation would be 0.0, 1.0, 2.0, ...

#### func (*Mat) Map

```go
func (m *Mat) Map(f elementFunc) *Mat
```
Map applies a given function to each element of a mat object. The given function
must take a pointer to a float64, and return nothing.

#### func (*Mat) Mul

```go
func (m *Mat) Mul(n *Mat) *Mat
```
Mul is the element-wise multiplication of a mat object by another which is
passed to this method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise multiplication is stored in the original mat
on which the method was invoked.

#### func (*Mat) Ones

```go
func (m *Mat) Ones() *Mat
```
Ones sets all values of a mat to be equal to 1.0

#### func (*Mat) Prod

```go
func (m *Mat) Prod(axis, slice int) float64
```
Prod returns the product of the elements along a specific row or specific
column. The first argument selects the row or column (0 or 1), and the second
argument selects which row or column for which we want to calculate the product.
For example:

m.Prod(1, 2)

will calculate the product of the 3rd column of mat m.

#### func (*Mat) Rand

```go
func (m *Mat) Rand(args ...float64) *Mat
```
Rand sets the values of a mat to random numbers. The range from which the random
numbers are selected is determined based on the arguments passed.

- For no arguments, the range is [0, 1) 
- For 1 argument, the range is [0, arg) for arg > 0, or (arg, 0] is arg < 0. 
- For 2 arguments, the range is [arg1, arg2).

#### func (*Mat) Reset

```go
func (m *Mat) Reset() *Mat
```
Reset sets all values of a mat object to 0.0

#### func (*Mat) Reshape

```go
func (m *Mat) Reshape(rows, cols int) *Mat
```
Reshape changes the row and the columns of the mat object as long as the total
number of values contained in the mat object remains constant. The order and the
values of the mat does not change with this function.

#### func (*Mat) Row

```go
func (m *Mat) Row(x int) *Mat
```
Row returns a new mat object whose values are equal to a row of the original mat
object. The number of Rows of the returned mat object is equal to 1, and the
number of columns is equal to the number of columns of the original mat.

#### func (*Mat) Scale

```go
func (m *Mat) Scale(f float64) *Mat
```
Scale is the element-wise multiplication of a mat object by a scalar.

The results of the element-wise multiplication is stored in the original mat on
which the method was invoked.

#### func (*Mat) SetAllTo

```go
func (m *Mat) SetAllTo(val float64) *Mat
```
SetAllTo sets all values of a mat to the passed float64 value.

#### func (*Mat) Std

```go
func (m *Mat) Std(axis, slice int) float64
```
Std returns the standard deviation of the elements along a specific row or
specific column. The standard deviation is defined as the square root of the
mean distance of each element from the mean. Look at:
http://mathworld.wolfram.com/StandardDeviation.html

For example:

```go
m.Std(1, 0)
```

will calculate the standard deviation of the first column of mat m.

#### func (*Mat) Sub

```go
func (m *Mat) Sub(n *Mat) *Mat
```
Sub is the element-wise subtraction of a mat object which is passed to this
method from the original mat which called the method.

The shape of the mat objects must be the same (same number or rows and columns)
and the results of the element-wise subtraction is stored in the original mat on
which the method was invoked.

#### func (*Mat) Sum

```go
func (m *Mat) Sum(axis, slice int) float64
```
Sum returns the sum of the elements along a specific row or specific column. The
first argument selects the row or column (0 or 1), and the second argument
selects which row or column for which we want to calculate the sum. For example:

```go
m.Sum(0, 2)
```

will calculate the sum of the 3rd row of mat m.

#### func (*Mat) T

```go
func (m *Mat) T() *Mat
```
T returns the transpose of the original matrix. The transpose of a mat object is
defined in the usual manner, where every value at row x, and column y is placed
at row y, and column x. The number of rows and column of the transposed mat are
equal to the number of columns and rows of the original matrix, respectively.
This method creates a new mat object, and the original is left intact.

#### func (*Mat) ToCSV

```go
func (m *Mat) ToCSV(fileName string)
```
ToCSV creates a file with the passed name, and writes the content of a mat
object to it, by putting each row in a single comma separated line. The number
of entries in each line is equal to the columns of the mat object.

#### func (*Mat) ToSlice

```go
func (m *Mat) ToSlice() [][]float64
```
ToSlice returns the values of a mat object as a 2D slice of float64s.

#### func (*Mat) ToString

```go
func (m *Mat) ToString() string
```
ToString returns the string representation of a mat. This is done by putting
every row into a line, and separating the entries of that row by a space. note
that the last line does not contain a newline.

#### func (*Mat) Vals

```go
func (m *Mat) Vals() []float64
```
Vals returns the values contained in a mat object as a 1D slice of float64s.
