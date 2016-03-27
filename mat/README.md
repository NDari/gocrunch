# mat
--
    import "github.com/NDari/gocrunch/mat"

Package mat implements a large set of functions which act on one and two
dimensional slices of float64.

Many of the functions in this library expect either a float64, a []float64,
or [][]float64, and do "the right thing" based on what is passed. For example,
consider the function

```go
mat.Mul(m, n)
```
In this function, m can be a [][]float64, where as n could be
a float64, []float64, or [][]float64. This allows the same function to be called
for wide range of situations. This trades compile time safety for runtime errors.
We believe that Go's fast compile time, along with the verbose errors in this
package make up for that, however.

All errors encountered in this package, such as attempting to access an
element out of bounds are treated as critical error, and thus, the code
immediately panics. In such cases, the function in which the error was
encountered is printed to the screen along with the reason for the panic,
in addition to the full stack trace, in order to help fix any issues
rapidly.

As mentioned, all the functions in this library act on Go primitive types,
which allows the code to be easily modified to serve in different situations.

## Usage

Create a square matrix, and a non-square matrix:

```go
m, n := mat.New(10), mat.New(10, 5)
```
Create an Identity matrix:

```go
q := mat.I(10)
```
Set all values of m to 11.0:

```go
mat.Set(m, 11.0)
```
Check if all elements of m are set to 11.0:

```go
isEleven := func(i float64) float64 {
  if i == 11.0 {
    return true
  }
  return false
}

if !mat.All(m, isEleven) {
  log.Fatal("Non-11 values found!")
}
```
Calculate the  dot product of m with the identity matrix:

```go
m2 := mat.Dot(m, q)
```
Check if m and m2 are equal:

 ```go
if !mat.Equal(m, m2) {
  log.Fatal("We have a problem...")
}

 ```
## Documentation

Full documentation is at below

[![GoDoc](https://godoc.org/github.com/NDari/gocrunch/mat?status.svg)](https://godoc.org/github.com/NDari/gocrunch/mat)

## Badges

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
