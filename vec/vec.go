/*
Package vec implements a large set of functions which act on one dimensional
slices of float64.

Many of the functions in this library expect either a float64 or a []float64,
and do "the right thing" based on what is passed. For example,consider the
function

	vec.Mul(m, n)

In this function, m is a []float64, where as n could be a float64 or a
[]float64. This allows the same function to be called for wide range of
situations. This trades compile time safety for runtime errors.

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
*/
package vec

import "fmt"

/*
Pop takes a []float64, and "pops" the last entry, returning it along with the
modified []float64. The other elements of the []float64 remain intact.
*/
func Pop(v []float64) (float64, []float64) {
	x, v := v[len(v)-1], v[:len(v)-1]
	return x, v
}

/*
Push appends a float64 to the end of a []float64, returning the modified
[]float64. The returned []float64 is one element longer, and the other
elements remain intact.
*/
func Push(v []float64, x float64) []float64 {
	v = append(v, x)
	return v
}

/*
Shift removes the first element of a []float64, returning it along with the
modified []float64. All the other elements in the []float64 remain intact,
however their order is changed (the second element is now the first, etc).
*/
func Shift(v []float64) (float64, []float64) {
	x, v := v[0], v[1:]
	return x, v
}

/*
Unshift appends a float64 to the begening of a []float64, returning the
modified []float64. The elements in the original []float64 remain intact,
however their order is now changed (the first element is now the second, etc.)
*/
func Unshift(v []float64, x float64) []float64 {
	v = append([]float64{x}, v...)
	return v
}

/*
Cut removes a range of enteries from a []float64. This function can be used
in one of two ways. First is providing a single int, such as:
	vec.Cut(v, x)
which means that all elements of v whose index is larger than x will be
dropped. The second method of using this function is as:
	vec.Cut(v, 2, 4)
which means that the second and 3rd elements of v are dropped.
*/
func Cut(v []float64, args ...int) []float64 {
	switch len(args) {
	case 1:
		v = v[:args[0]]
	case 2:
		v = append(v[:args[0]], v[args[1]:]...)
	default:
		fmt.Println("\ngocrunch/vec error.")
		s := "In mat.%v, 1 or 2 ints must be passed in addition to the slice.\n"
		s += "However, %d ints were passed.\n"
		s = fmt.Sprintf(s, "Cut()", len(args))
		panic(s)
	}
	return v
}

/*
Equal
*/
func Equal(v, w []float64) bool {
	if len(v) != len(w) {
		return false
	}
	for i := range v {
		if v[i] != w[i] {
			return false
		}
	}
	return true
}

/*
Set
*/
func Set(v []float64, val float64) {
	for i := range v {
		v[i] = val
	}
}

/*
Foreach
*/
func Foreach(v []float64, f func(float64) float64) {
	for i := range v {
		v[i] = f(v[i])
	}
}

/*
All
*/
func All(v []float64, f func(float64) bool) bool {
	for i := range v {
		if !f(v[i]) {
			return false
		}
	}
	return true
}

/*
Any
*/
func Any(v []float64, f func(float64) bool) bool {
	for i := range v {
		if f(v[i]) {
			return true
		}
	}
	return false
}

/*
Sum
*/
func Sum(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum
}

/*
Avg
*/
func Avg(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum / float64(len(v))
}
