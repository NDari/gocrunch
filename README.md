## Note: Alpha package. The API may change at any time! Be warned.

# numgo
--
    import "numgo"

Package numgo supplies functions that are used to create and act on vectors (1D
slices) and matrices (2D slices). This package is created in the spirit of the
Numpy library for Python, while adhering to the principles of good Go code.

All errors encountered in this package, such as shape mismatch between two
arrays to be multiplied together in an element-wise fashion, are considered
fatal. In other words, when encountering an error, the functions in this package
will die, instead of returning an error to be handled by the caller. Before
dying, all errors will report what went wrong, a report where the error 
happened by printing the stack trace, which will help the user in correcting
the error.

## Directories

- [numgo/mat](https://github.com/NDari/numgo/tree/master/mat): Package mat implements function that create or act upon 2D slices of `float64`.
