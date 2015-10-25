/*
Package numgo supplies functions that are used to create and act on
vectors (1D slices) and matrices (2D slices). This package is created
in the spirit of the Numpy library for Python, while adhering to the
principles of good Go code.

The 2D slices acted on or created by the functions below are assumed to
be non-jagged. This means that for a given [][]float64, the inner slices
are assumed to be of the same length.

All errors encountered in this package, such as shape mismatch between
two arrays to be multiplied together in an element-wise fashion, are
considered fatal. In other words, when encountering an error, the
functions in this package will die, instead of returning an error to be
handled by the caller.
*/
package numgo
