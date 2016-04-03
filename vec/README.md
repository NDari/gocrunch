# vec
--
    import "github.com/NDari/gocrunch/vec"

Package vec implements a large set of functions which act on one and one
dimensional slices of float64. For functions that act on [][]float64,
look at the `gocrunch/mat` package.

Many of the functions in this library expect either a float64 or a []float64,
and do "the right thing" based on what is passed. For example, consider the
function:

```go
vec.Mul(m, n)
```
In this function, m is a []float64, where as n could be a float64, or a
[]float64. This allows the same function to be called
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

Create a []float64, then set all elements to 11.0:

```go
v := make([]float64, 10)
vec.Set(v, 11.0)
```
Check if all elements of v are set to 11.0:

```go
isEleven := func(i float64) float64 {
  if i == 11.0 {
    return true
  }
  return false
}

if !vec.All(m, isEleven) {
  log.Fatal("Non-11 values found!")
}
```
make a copy and check if they are equal:

```go
v2 := make([]float64, 10)
copy(v1, v) // builtin copy function
if vec.Equal(v2, v) {
  fmt.Println("yay")
}
```
Add 2.0 to all elements of v:

```go
vec.Add(v, 2.0)
```
subtract v1 from v, which will make all the elements of v to be equal to 2.0:

```go
vec.Sub(v, v1)
```

## Documentation

Full documentation is at godoc.org [![GoDoc](https://godoc.org/github.com/NDari/gocrunch/vec?status.svg)](https://godoc.org/github.com/NDari/gocrunch/vec)

## Badges

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)

