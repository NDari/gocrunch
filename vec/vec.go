/*
Package vec implements a large set of functions which act on one dimensional
slices of float64.

Many of the functions in this library expect either a float64, a []float64,
or [][]float64, and do "the right thing" based on what is passed. For example,
consider the function

	mat.Mul(m, n)

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
*/
package vec

func Pop(v []float64) (float64, []float64) {
	x, v := v[len(v)-1], v[:len(v)-1]
	return x, v
}

func Push(v []float64, x float64) []float64 {
	v = append(v, x)
	return v
}

func Shift(v []float64) (float64, []float64) {
	x, v := v[0], v[1:]
	return x, v
}

func Unshift(v []float64, x float64) []float64 {
	v = append([]float64{x}, v...)
	return v
}

func Cut(v []float64, from, to int) []float64 {
	v = append(v[:from], v[to:]...)
	return v
}

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

func Set(v []float64, val float64) {
	for i := range v {
		v[i] = val
	}
}

func Foreach(v []float64, f func(float64) float64) {
	for i := range v {
		v[i] = f(v[i])
	}
}

func All(v []float64, f func(float64) bool) bool {
	for i := range v {
		if !f(v[i]) {
			return false
		}
	}
	return true
}

func Any(v []float64, f func(float64) bool) bool {
	for i := range v {
		if f(v[i]) {
			return true
		}
	}
	return false
}

func Sum(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum
}

func Avg(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum / float64(len(v))
}
