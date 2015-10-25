# vec
--
    import "numgo/vec"

Package vec implements functions that create or act upon 1D slices of `float64`.

## Usage

#### func  Add

```go
func Add(v1, v2 []float64) []float64
```
Add returns a new 1D slice that is the result of element-wise addition of two 1D
slices.

#### func  Apply

```go
func Apply(f ElementalFn, v []float64) []float64
```
Apply created a new 1D slice which is populated throw applying the given
function to the corresponding entries of a given 1D slice. This function does
not modify its arguments, instead allocating a new 1D slice to contain the
result. This is a performance hit. If you are OK with mutating the original
vector, then use the "ApllyInPlace" function instead.

#### func  ApplyInPlace

```go
func ApplyInPlace(f ElementalFn, v []float64)
```
ApplyInPlace calls a given elemental function on each Element of a 1D slice,
returning it afterwards. This function modifies the original 1D slice. If a
non-mutating operation is desired, use the "Apply" function instead.

#### func  Div

```go
func Div(v1, v2 []float64) []float64
```
Div returns a new 1D slice that is the result of element-wise division of two 1D
slices. If any elements in the 2nd 1D slice are 0, then this function call
aborts.

#### func  Dot

```go
func Dot(v1, v2 []float64) float64
```
Dot is the inner product of two 1D slices of `float64`.

#### func  Equal

```go
func Equal(v1, v2 []float64) bool
```
Equals checks if two 1D slices have the same length, and contain the same
entries at each slot.

#### func  Inc

```go
func Inc(l int) []float64
```
Inc returns a 1D slice, where element `[0] == 0.0`, and each subsequent element
is incremented by `1.0`.

For example, `m := Inc(3)` is

`[1.0, 2.0 3.0]`.

#### func  Mul

```go
func Mul(v1, v2 []float64) []float64
```
Mul returns a new 1D slice that is the result of element-wise multiplication of
two 1D slices.

#### func  Norm

```go
func Norm(v []float64) float64
```
Norm calculated the norm of a given 1D slice. This is the Euclidean length of
the slice.

#### func  Ones

```go
func Ones(l int) []float64
```
Ones returns a new 1D slice where all the elements are equal to `1.0`.

#### func  Reset

```go
func Reset(v []float64)
```
Reset sets the values of all entries in a 2D slice of `float64` to `0.0`.

#### func  Sub

```go
func Sub(v1, v2 []float64) []float64
```
Sub returns a new 1D slice that is the result of element-wise subtraction of two
1D slices.

#### func  Sum

```go
func Sum(v []float64) float64
```
Sum returns the sum of the elements of a 1D slice of `float64`.

#### type ElementalFn

```go
type ElementalFn func(float64) float64
```

ElementalFn is a function that takes a float64 and returns a `float64`. This
function can therefore be applied to each element of a 2D `float64` slice, and
can be used to construct a new one.
