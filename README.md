# mat64
--
    import "mat64"

Package mat64 supplies functions that create or act on 2D slices of float64s,
for the Go language.

## Usage

#### func  AppendCol

```go
func AppendCol(m [][]float64, v []float64) [][]float64
```
AppendCol appends a column to the right side of a 2D slice of float64s.

#### func  Apply

```go
func Apply(f ElementalFn, m [][]float64) [][]float64
```
Apply calls a given elemental function on each Element of a 2D slice, returning
it afterwards.

#### func  Col

```go
func Col(c int, m [][]float64) []float64
```
Col returns a column of a 2D slice of float64s.

#### func  Copy

```go
func Copy(m [][]float64) [][]float64
```
Copy copies the content of a 2D slice of float64s into another with the same
shape. This is a deep copy, unlike the builtin copy function that is shallow for
2D slices.

#### func  Dot

```go
func Dot(m, n [][]float64) [][]float64
```
Dot is the matrix multiplication of two 2D slices of float64s

#### func  Dump

```go
func Dump(m [][]float64, fileName string)
```
Dump prints the content of a [][]float64 object to a file, using comma as the
delimeter between the elements of a row, and a new line between rows.

#### func  Equal

```go
func Equal(m, n [][]float64) bool
```
Equals checks if two mat objects have the same shape and the same entries in
each row and column.

#### func  FromString

```go
func FromString(str [][]string) [][]float64
```
FromString converts a 2D slice of strings into a 2D slice of float64s.

#### func  I

```go
func I(r int) [][]float64
```
I returns an r by r identity matrix for a given r.

#### func  Inc

```go
func Inc(r, c int) [][]float64
```
Inc returns a 2D slice, where element [0][0] == 0.0, and each subsequent elemnt
is incrmeneted by 1.0

#### func  Load

```go
func Load(fileName string) [][]float64
```
Load generates a 2D slice of floats from a csv file.

#### func  New

```go
func New(r, c int) [][]float64
```
New returns a 2D slice of float64s with the given row and columns.

#### func  Ones

```go
func Ones(r, c int) [][]float64
```
Ones returns a new 2D slice where all the elements are equal to 1.0

#### func  Print

```go
func Print(m [][]float64)
```
Print prints a 2D slice of float64s to the std out.

#### func  Reset

```go
func Reset(m [][]float64) [][]float64
```
Reset sets the values of all entries in a 2D slice of float64s to 0.0

#### func  Row

```go
func Row(r int, m [][]float64) []float64
```
Row returns a row of a 2D slice of float64s

#### func  T

```go
func T(m [][]float64) [][]float64
```
T returns a copy of a given matrix with the elements mirrored across the
diagonal. for example, the element At(i, j) becomes the element At(j, i). This
function leaves the original matrix intact.

#### func  Times

```go
func Times(m, n [][]float64) [][]float64
```
Times returns a new 2D slice that is the result of element-wise multiplication
of two 2D slices.

#### func  ToString

```go
func ToString(m [][]float64) [][]string
```
ToString converts a 2D slice of float64s into a 2D slice of strings.

#### type ElementalFn

```go
type ElementalFn func(float64) float64
```

ElementalFn is a function that takes a float64 and returns a float64. This
function can therefore be applied to each element of a 2D float64 slice, and can
be used to construct a new one.
