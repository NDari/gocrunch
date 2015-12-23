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

#### func  From1DSlice

```go
func From1DSlice(s []float64) *mat
```
From1DSlice creates a mat object from a slice of float64s. The created mat
object has one row, and the number of columns equal to the length of the 1D
slice from which it was created.

#### func  FromCSV

```go
func FromCSV(filename string) *mat
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
func FromSlice(s [][]float64) *mat
```
FromSlice generated a mat object from a [][]float64 slice. The slice is checked
for being a non-jagged slice, where each row contains the same number of
elements. The creation of a mat object from jagged 2D slices is not supported as
on yet.

#### func  New

```go
func New(r, c int) *mat
```
New is the primary constructor for the "mat" object. The "r" and "c" params are
expected to be greater than zero, and the values of the mat object are
initialized to 0.0, which is the default behavior of Go for slices of float64s.
