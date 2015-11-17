# mat
--
    import "numgo/mat"

Package mat implements function that create or act upon 2D slices of `float64`.
This is in essence the same concept of a matrix in other languages.

The 2D slices acted on or created by the functions below are assumed to be
non-jagged. This means that for a given [][]float64, the inner slices are
assumed to be of the same length.

## Usage

#### func  AppendCol

```go
func AppendCol(m [][]float64, v []float64) [][]float64
```
AppendCol appends a column to the right side of a 2D slice of float64s.

#### func  Col

```go
func Col(c int, m [][]float64) []float64
```
Col returns a column of a 2D slice of `float64`. Col uses a 0-based index, hence
the first column of a 2D slice, m, is `Col(0, m)`.

This function also allows for negative indexing. For example, `Col(-1, m)` is
the last column of the 2D slice m, and `Col(-2, m)` is the second to last column
of m, and so on.

#### func  Concat

```go
func Concat(m, n [][]float64) [][]float64
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

#### func  Copy

```go
func Copy(m [][]float64) [][]float64
```
Copy copies the content of a 2D slice of float64 into another with the same
shape. This is a deep copy, unlike the built in copy function that is shallow
for 2D slices.

#### func  Dot

```go
func Dot(m, n [][]float64) [][]float64
```
Dot is the matrix multiplication of two 2D slices of `float64`.

#### func  Dump

```go
func Dump(m [][]float64, fileName string)
```
Dump prints the content of a `[][]float64` slice to a file, using comma as the
delimiter between the elements of a row, and a new line between rows.

#### func  Equal

```go
func Equal(m, n [][]float64) bool
```
Equal checks if two 2D slices have the same shape and the same entries in each
row and column. If either the shape or the entries of the arguments are
different, `false` is returned. Otherwise, the return value is `true`.

#### func  FromString

```go
func FromString(str [][]string) [][]float64
```
FromString converts a `[][]string` to `[][]float64`.

#### func  I

```go
func I(r int) [][]float64
```
I returns an r by r 2D slice for a given r, where the elements along the
diagonal (where the first and the second index are equal) is set to `1.0`, and
all other elements are set to `0.0`.

#### func  Inc

```go
func Inc(r, c int) [][]float64
```
Inc returns a 2D slice, where element `[0][0] == 0.0`, and each subsequent
element is incremented by `1.0`.

For example:

```go 
m := Inc(3, 2)
mat.Print(m) // 1.0, 2.0
             // 3.0, 4.0
    	     // 5.0, 6.0

```

#### func  Load

```go
func Load(fileName string) [][]float64
```
Load generates a 2D slice of floats from a CSV file.

#### func  Map

```go
func Map(f func(float64) float64, m [][]float64) [][]float64
```
Map calls a given elemental function on each Element of a 2D slice, returning it
afterwards. This function modifies the original 2D slice.

#### func  Mul

```go
func Mul(m, n [][]float64) [][]float64
```
Mul returns a new 2D slice that is the result of element-wise multiplication of
two 2D slices.

#### func  New

```go
func New(r, c int) [][]float64
```
New returns a 2D slice of `float64` with the given number of row and columns.
This function should be used as a convenience tool, and it is exactly equivalent
to the normal method of allocating a uniform (non-jagged) 2D slice of `float64`.

If it is anticipated that the 2D slice will grow, use the "NewExpand" function
below. For full details, read that function's documentation.

#### func  NewExpand

```go
func NewExpand(r, c int) [][]float64
```
NewExpand returns a 2D slice of `float64`, with the given number of rows and
columns. The difference between this function and the "New" function above is
that the inner slices are allocated with double the capacity, and hence can grow
without the need for reallocation up to column * 2.

Note that this extended capacity will waste memory, so the NewExtend should be
used with care in situations where the performance gained by avoiding
reallocation justifies the extra cost in memory.

#### func  Ones

```go
func Ones(r, c int) [][]float64
```
Ones returns a new 2D slice where all the elements are equal to `1.0`.

#### func  Print

```go
func Print(m [][]float64)
```
Print prints a `[][]float64` to the stdout.

#### func  Reset

```go
func Reset(m [][]float64) [][]float64
```
Reset sets the values of all entries in a 2D slice of `float64` to `0.0`.

#### func  Row

```go
func Row(r int, m [][]float64) []float64
```
Row returns a row of a 2D slice of `float64`. Row uses a 0-based index, hence
the first row of a 2D slice, m, is Row(0, m).

This function also allows for negative indexing. For example, Row(-1, m) is the
last row of m.

#### func  T

```go
func T(m [][]float64) [][]float64
```
T returns a copy of a given 2D slice with the elements of the 2D slice mirrored
across the diagonal. For example, the element `[i][j]` becomes the element
`[j][i]` of the returned 2D slice. This function leaves the original matrix
intact.

#### func  ToString

```go
func ToString(m [][]float64) [][]string
```
ToString converts a `[][]float64` to `[][]string`.
