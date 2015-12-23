# mat
--
    import "numgo/mat"

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
func New(r, c int) *Mat
```
New is the primary constructor for the "mat" object. The "r" and "c" params are
expected to be greater than zero, and the values of the mat object are
initialized to 0.0, which is the default behavior of Go for slices of float64s.

#### func (*Mat) Add

```go
func (m *Mat) Add(n *Mat) *Mat
```

#### func (*Mat) All

```go
func (m *Mat) All(f booleanFunc) bool
```
All checks if a supplied function is true for all elements of a mat object. For
instance, if a supplied function returns true for negative values and false
otherwise, then All would be true if and only if all elements of the mat object
are negative.

#### func (*Mat) Any

```go
func (m *Mat) Any(f booleanFunc) bool
```
Any checks if a supplied function is true for one elements of a mat object. For
instance, if a supplied function returns true for negative values and false
otherwise, then Any would be true if at least one element of the mat object is
negative.

#### func (*Mat) Average

```go
func (m *Mat) Average(axis, slice int) float64
```

#### func (*Mat) Col

```go
func (m *Mat) Col(x int) *Mat
```
Col returns a new mat object whole values are equal to a column of the original
mat object. The number of Rows of the returned mat object is equal to the number
of rows of the original mat, and the number of columns is equal to 1.

#### func (*Mat) CombineWith

```go
func (m *Mat) CombineWith(n *Mat, how reducerFunc) *Mat
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
object from all elements for which the function returned true. For example,
given a function that takes a float64, returning true for numbers greater than
zero and false otherwise, This method created a new mat object from all the
positive elements of the original matrix. If no elements return true for a given
function, nil is returned. It is expected that the caller of this method checks
the returned value to ensure that it is not nil.

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

#### func (*Mat) Ones

```go
func (m *Mat) Ones() *Mat
```
Ones sets all values of a mat to be equal to 1.0

#### func (*Mat) Prod

```go
func (m *Mat) Prod(axis, slice int) float64
```

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

#### func (*Mat) Std

```go
func (m *Mat) Std(axis, slice int) float64
```

#### func (*Mat) Sub

```go
func (m *Mat) Sub(n *Mat) *Mat
```

#### func (*Mat) Sum

```go
func (m *Mat) Sum(axis, slice int) float64
```

#### func (*Mat) T

```go
func (m *Mat) T() *Mat
```
T returns the transpose of the original matrix. The transpose of a mat object is
defined in the usual manner, where every value at row x, and column y is placed
at row y, and column x. The number of rows and column of the transposed mat are
equal to the number of columns and rows of the original matrix, respectively.

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

#### func (*Mat) Vals

```go
func (m *Mat) Vals() []float64
```
Vals returns the values contained in a mat object as a 1D slice of float64s.
