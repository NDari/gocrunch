/*
Package vec implements functions that create or act upon 1D slices of
`float64`.
*/
package vec

import (
	"fmt"
	"math"
	"os"
	"runtime/debug"
)

// Ones returns a new 1D slice where all the elements are equal to `1.0`.
func Ones(l int) []float64 {
	o := make([]float64, l)
	f := func(i float64) float64 {
		return 1.0
	}
	MapInPlace(f, o)
	return o
}

// Inc returns a 1D slice, where element `[0] == 0.0`, and each
// subsequent element is incremented by `1.0`.
//
// For example, `m := Inc(3)` is
//
// `[1.0, 2.0 3.0]`.
func Inc(l int) []float64 {
	v := make([]float64, l)
	for i := 0; i < l; i++ {
		v[i] = float64(i)
	}
	return v
}

// Equal checks if two 1D slices have the same length, and contain the same
// entries at each slot.
func Equal(v1, v2 []float64) bool {
	if len(v1) != len(v2) {
		return false
	}
	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

// Mul returns a new 1D slice that is the result of element-wise multiplication
// of two 1D slices.
func Mul(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		fmt.Println("\nnumgo/vec error.")
		s := "In vec.%s the length of the first 1D slice is %d, while\n"
		s += "the length of the second 1D slice is %d. They must be equal\n"
		s = fmt.Sprintf(s, "Mul", len(v1), len(v2))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	o := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		o[i] = v1[i] * v2[i]
	}
	return o
}

// Add returns a new 1D slice that is the result of element-wise addition
// of two 1D slices.
func Add(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		fmt.Println("\nnumgo/vec error.")
		s := "In vec.%s the length of the first 1D slice is %d, while\n"
		s += "the length of the second 1D slice is %d. They must be equal\n"
		s = fmt.Sprintf(s, "Add", len(v1), len(v2))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	o := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		o[i] = v1[i] + v2[i]
	}
	return o
}

// Sub returns a new 1D slice that is the result of element-wise subtraction
// of two 1D slices.
func Sub(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		fmt.Println("\nnumgo/vec error.")
		s := "In vec.%s the length of the first 1D slice is %d, while\n"
		s += "the length of the second 1D slice is %d. They must be equal\n"
		s = fmt.Sprintf(s, "Sub", len(v1), len(v2))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	o := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		o[i] = v1[i] - v2[i]
	}
	return o
}

// Div returns a new 1D slice that is the result of element-wise division
// of two 1D slices. If any elements in the 2nd 1D slice are 0, then this
// function call aborts.
func Div(v1, v2 []float64) []float64 {
	if len(v1) != len(v2) {
		fmt.Println("\nnumgo/vec error.")
		s := "In vec.%s the length of the first 1D slice is %d, while\n"
		s += "the length of the second 1D slice is %d. They must be equal\n"
		s = fmt.Sprintf(s, "Div", len(v1), len(v2))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	o := make([]float64, len(v1))
	for i := 0; i < len(v1); i++ {
		if v2[i] == 0.0 {
			fmt.Println("\nnumgo/vec error.")
			s := "In vec.%s, the entry at the index %d of the second vector is 0.0.\n"
			s += "Division by zero is not defined."
			s = fmt.Sprintf(s, "Div", i)
			fmt.Println(s)
			fmt.Println("Stack trace for this error:\n")
			debug.PrintStack()
			os.Exit(1)
		}
		o[i] = v1[i] / v2[i]
	}
	return o
}

// MapInPlace calls a given elemental function on each Element of a 1D slice,
// returning it afterwards. This function modifies the original 1D slice. If
// a non-mutating operation is desired, use the "Map" function instead.
func MapInPlace(f func(float64) float64, v []float64) {
	for i := 0; i < len(v); i++ {
		v[i] = f(v[i])
	}
}

// Map created a new 1D slice which is populated throw Maping the given
// function to the corresponding entries of a given 1D slice. This function
// does not modify its arguments, instead allocating a new 1D slice to
// contain the result. This is a performance hit. If you are OK with mutating
// the original vector, then use the "ApllyInPlace" function instead.
func Map(f func(float64) float64, v []float64) []float64 {
	o := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		o[i] = f(v[i])
	}
	return o
}

// Dot is the inner product of two 1D slices of `float64`.
func Dot(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		fmt.Println("\nnumgo/vec error.")
		s := "In vec.%s the length of the first 1D slice is %d, while\n"
		s += "the length of the second 1D slice is %d. They must be equal\n"
		s = fmt.Sprintf(s, "Dot", len(v1), len(v2))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:\n")
		debug.PrintStack()
		os.Exit(1)
	}
	var o float64
	for i := 0; i < len(v1); i++ {
		o += v1[i] * v2[i]
	}
	return o
}

// Reset sets the values of all entries in a 2D slice of `float64` to `0.0`.
func Reset(v []float64) {
	f := func(i float64) float64 {
		return 0.0
	}
	MapInPlace(f, v)
	return
}

// Sum returns the sum of the elements of a 1D slice of `float64`.
func Sum(v []float64) float64 {
	var o float64
	for i := 0; i < len(v); i++ {
		o += v[i]
	}
	return o
}

// Norm calculated the norm of a given 1D slice. This is the Euclidean length
// of the slice.
func Norm(v []float64) float64 {
	square := func(i float64) float64 {
		return i * i
	}
	return math.Sqrt(Sum(Map(square, v)))
}

/*
All checks that every item in given 1D slice returns true for a given function.
The passed function must accept a single float64, and return a boolean. For
example, consider:

```
f := func(i float64) bool {
	if i > 0.0 { return true}
	return false
}

m := vec.Ones(10)
t := vec.All(f, m) // t == true
```
*/
func All(f func(float64) bool, v []float64) bool {
	for i := 0; i < len(v); i++ {
		if !f(v[i]) {
			return false
		}
	}
	return true
}

/*
Any checks if there is a single item in a 1D slice returns true for a given
function. The passed function must accept a single float64, and return a
boolean. For example, consider:

```
f := func(i float64) bool {
	if i < 0.0 { return true}
	return false
}

m := vec.Ones(10)
t := vec.All(f, m) // t == false
```
*/
func Any(f func(float64) bool, v []float64) bool {
	for i := 0; i < len(v); i++ {
		if !f(v[i]) {
			return true
		}
	}
	return false
}
