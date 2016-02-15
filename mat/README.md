# mat
--
    import "github.com/NDari/numgo/mat"

Package mat implements a "mat" object, which behaves like a 2-dimensional array
or list in other programming languages. Under the hood, the Mat is a flat
slice, which provides for optimal performance in Go, while the methods and
constructors provide for a higher level of performance and abstraction when
compared to the "2D" slices of go (slices of slices).

All errors encountered in this package, such as attempting to access an element
out of bounds are treated as critical error, and thus, the code immediately
exits with signal 1. In such cases, the function/method in which the error was
encountered is printed to the screen, in addition to the full stack trace, in
order to help fix the issue rapidly.

## Usage


## Documentation

Full documentation is at [this page.](https://godoc.org/github.com/NDari/numgo/mat)
