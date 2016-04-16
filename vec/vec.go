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

import (
	"fmt"
	"math"
	"math/rand"
)

var (
	errStrings = []string{
		"\ngocrunch/vec error.\nIn vec.%s, cannot use %s on an empty []float64.\n",
		"\ngocrunch/vec error.\nIn vec.%s, %d is outside of range [0, %d).\n",
		"\ngocrunch/vec error.\nIn vec.%s, %d is outside of range (%d, %d).\n",
		"\ngocrunch/vec error.\nIn vec.%s, second arg, %d is not greater than third arg, %d.\n",
		"\ngocrunch/vec error.\nIn vec.%s, incorrect number of arguments received.\n",
		"\ngocrunch/vec error.\nIn vec.%s, the length the passed slices does not match: %d and %d.\n",
		"\ngocrunch/vec error.\nIn vec.%s, second arg must be float64 or []float64, received %v.\n",
		"\ngocrunch/vec error.\nIn vec.%s, the passed float64 cannot be 0.0\n",
		"\ngocrunch/vec error.\nIn vec.%s, in the second []float64, zero value found at index %d.\n",
		"\ngocrunch/vec error.\nIn vec.%s, the length of slice %d is not divisible by the stride %d.\n",
		"\ngocrunch/vec error.\nIn vec.%s, the first argument %f must be less than the second, %f.\n",
		"\ngocrunch/vec error.\nIn vec.%s, expected 0 to 0 float64 arguments, but got %d.\n",
	}
)

/*
Pop takes a []float64, and "pops" the last entry, returning it along with the
modified []float64. The other elements of the []float64 remain intact. For
example:

	v := []float64{1.0, 2.0, 3.0}
	x, v := vec.Pop(v) // x is 3.0, v is [1.0, 2.0]

This function mutates the passed []float64, and panics if the length of the
slice is 0 since it cannot be "popped".
*/
func Pop(v []float64) (float64, []float64) {
	if len(v) == 0 {
		panic(fmt.Sprintf(errStrings[0], "Pop()", "Pop()"))
	}
	x, v := v[len(v)-1], v[:len(v)-1]
	return x, v
}

/*
Push appends a float64 to the end of a []float64, returning the modified
[]float64. The returned []float64 is one element longer, and the other
elements remain intact. For example:

	v := []float64{1.0, 2.0, 3.0}
	v = vec.Push(v, 4.0) // v is [1.0, 2.0, 3.0, 4.0]

This function alters the passed []float64.
*/
func Push(v []float64, x float64) []float64 {
	v = append(v, x)
	return v
}

/*
Shift removes the first element of a []float64, returning it along with the
modified []float64. All the other elements in the []float64 remain intact,
however their order is changed (the second element is now the first, etc).
For example:

	v := []float64{1.0, 2.0, 3.0}
	x, v := vec.Shift(v) // x is 1.0, v is [2.0, 3.0]

This function alters the passed []float64, and panics if the length of the
slice is 0 since it cannot be "shifted".
*/
func Shift(v []float64) (float64, []float64) {
	if len(v) == 0 {
		panic(fmt.Sprintf(errStrings[0], "Shift()", "Shift()"))
	}
	x, v := v[0], v[1:]
	return x, v
}

/*
Unshift appends a float64 to the beginning of a []float64, returning the
modified []float64. The elements in the original []float64 remain intact,
however their order is now changed (the first element is now the second, etc.)
For example:

	v := []float64{1.0, 2.0, 3.0}
	v = vec.Shift(v, 10.0) // v is [10.0, 1.0, 2.0, 3.0]

The passed slice is mutated in this function.
*/
func Unshift(v []float64, x float64) []float64 {
	v = append([]float64{x}, v...)
	return v
}

/*
Cut removes a range of entries from a []float64. This function can be used
in one of two ways. First is providing a single int, such as:

	vec.Cut(v, x)

which means that all elements of v whose index is x or larger will be
dropped. The second method of using this function is as:

	vec.Cut(v, 2, 4)

which means that the second and 3rd elements of v are dropped. The passed
slice is mutated in this function.
*/
func Cut(v []float64, args ...int) []float64 {
	switch len(args) {
	case 1:
		if args[0] < 0 || args[0] >= len(v) {
			panic(fmt.Sprintf(errStrings[1], "Cut()", args[0], len(v)))
		}
		v = v[:args[0]]
	case 2:
		if args[0] < 0 || args[0] >= len(v) {
			panic(fmt.Sprintf(errStrings[1], "Cut()", args[0], len(v)))
		}
		if args[1] >= len(v) {
			panic(fmt.Sprintf(errStrings[2], "Cut()", args[1], args[0], len(v)))
		}
		if args[1] <= args[0] {
			panic(fmt.Sprintf(errStrings[3], "Cut()", args[1], args[0]))
		}
		v = append(v[:args[0]], v[args[1]:]...)
	default:
		panic(fmt.Sprintf(errStrings[4], "Cut()"))
	}
	return v
}

/*
To2D converts a []float64 to a [][]float64, using a passed stride. The values
of the entries in the []float64 are put into the [][]float64 row by row. For
example:

	v := []float64{0.0, 1.0, 2.0, 3.0}
	m := vec.To2D(v, 2) // m is [[0.0, 1.0], [2.0, 3.0]]

The original []float64 is not mutated in this function. The length of the
[]float64 must be exactly divisible by the passed stride, otherwise this
function will panic.
*/
func To2D(v []float64, stride int) [][]float64 {
	if math.Mod(float64(len(v)), float64(stride)) != 0.0 {
		panic(fmt.Sprintf(errStrings[9], "To2D()", len(v), stride))
	}
	m := make([][]float64, len(v)/stride)
	for i := range m {
		m[i] = make([]float64, stride)
	}
	idx := 0
	for i := range m {
		for j := range m[i] {
			m[i][j] = v[idx]
			idx++
		}
	}
	return m
}

/*
Rand creates []float64 of length x with the entries set to random numbers.
The x is and integers passed to this function. The range from which
the random numbers are selected is determined based on the arguments passed.

For no additional arguments, such as

	vec.Rand(x)

the range is [0, 1), which is from and including 0.0, to and excluding 1.0.

For 1 argument, such as

	vec.Rand(x, arg)

the range is [0, arg) for arg > 0, or (arg, 0] if arg < 0.

For 2 arguments, such as

	vec.Rand(x, arg1, arg2)

the range is [arg1, arg2). For this case, arg1 must be less than arg2, or
the function will panic.
*/
func Rand(x int, args ...float64) []float64 {
	v := make([]float64, x)
	switch len(args) {
	case 0:
		for i := range v {
			v[i] = rand.Float64()
		}
	case 1:
		for i := range v {
			v[i] = rand.Float64() * args[0]
		}
	case 2:
		if !(args[2] > args[1]) {
			panic(fmt.Sprintf(errStrings[10], "Rand()", args[0], args[1]))
		}
	default:
		panic(fmt.Sprintf(errStrings[11], "Rand()", len(args)))
	}
	return v
}

/*
Clone replicated the passed []slice. The returned slice is a copy of
original, both in terms of the length and the value of the elements
at each index. The returned copy is "deep", and manupilating it does
not effect the original slice.
*/
func Clone(v []float64) []float64 {
	c := make([]float64, len(v))
	copy(c, v)
	return c
}

/*
Equal checks if two []float64s are equal, by checking that they have the same length,
and the same entries in each index.
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
Set returns a copy of the passed []float64 where all of the elements are set to
the passed float64 in the second argument.

The original []float64 is not mutated in this function.
*/
func Set(v []float64, val float64) []float64 {
	c := Clone(v)
	for i := range c {
		c[i] = val
	}
	return c
}

/*
Foreach applies a function to each element of a []float64, storing the result
in a new []float64 which is returned.  Consider:

	double := func(i float64) float64 {
		return i * i
	}
	v := []float64{1.0, 2.0, 3.0}
	c := vec.Foreach(v, double) // c is {1.0, 4.0, 9.0}

Thus the original []float64 is not modified in this function.
*/
func Foreach(v []float64, f func(float64) float64) []float64 {
	c := Clone(v)
	for i := range v {
		c[i] = f(v[i])
	}
	return c
}

/*
All checks to see if all elements of a []float64 return true for a passed
function. If not, All() returns false. Consider:

	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	v := make([]float64, 10)
	v = vec.Set(v, -12.0)
	allNegatives := vec.All(v, negative) // true

To check if any element of a []float64 pass a certain function, look
at vec.Any().
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
Any checks to see if any element of a []float64 returns true for a passed
function. If no elements return true, then Any() will return false. Consider:

	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	v := make([]float64, 10)
	v = vec.Set(v, 12.0)
	anyNegatives := vec.Any(v, negative) // false

To check if all elements of a []float64 pass a certain function, look
at vec.All().
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
Sum adds all elements in a []float64. Consider:

	v := []float64{ 1.0, 2.0, 3.0 }
	s := vec.Sum(v) // 6.0

This function does not alter the original []float64.
*/
func Sum(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum
}

/*
Prod multiplies all elements in a []float64. Consider

	v := []float64{ 2.0, 2.0, 2.0 }
	s := vec.Prod(v) // 8.0

This function does not alter the original []float64.
*/
func Prod(v []float64) float64 {
	prod := 1.0
	for i := range v {
		prod *= v[i]
	}
	return prod
}

/*
Avg returns the average value of a []float64. Consider

	v := []float64{ 1.0, 2.0, 3.0 }
	s := vec.Avg(v) // 2.0

This function does not alter the original []float64.
*/
func Avg(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i]
	}
	return sum / float64(len(v))
}

/*
Mul takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the multiplication operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 10
	v := []float64{1.0, 2.0, 3.0}
	v := vec.Mul(v, val) // v is now {10.0, 20.0, 30.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	q := vec.Mul(v, w) // q is now {3.0, 4.0, 12.0}

The original arguments are not modified in this function.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Mul(v []float64, val interface{}) []float64 {
	c := Clone(v)
	switch w := val.(type) {
	case float64:
		for i := range c {
			c[i] *= w
		}
	case []float64:
		if len(v) != len(w) {
			panic(fmt.Sprintf(errStrings[5], "Mul()", len(c), len(w)))
		}
		for i := range c {
			c[i] *= w[i]
		}
	default:
		panic(fmt.Sprintf(errStrings[6], "Mul()", w))
	}
	return c
}

/*
Add takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the addition operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 10
	v := []float64{1.0, 2.0, 3.0}
	v := vec.Add(v, val) // v is {11.0, 12.0, 13.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	q := vec.Add(v, w) // q is {4.0, 4.0, 7.0}

The original arguments are not modified in this function.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Add(v []float64, val interface{}) []float64 {
	c := Clone(v)
	switch w := val.(type) {
	case float64:
		for i := range c {
			c[i] += w
		}
	case []float64:
		if len(c) != len(w) {
			panic(fmt.Sprintf(errStrings[5], "Add()", len(c), len(w)))
		}
		for i := range c {
			c[i] += w[i]
		}
	default:
		panic(fmt.Sprintf(errStrings[6], "Mul()", w))
	}
	return c
}

/*
Sub takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the subtraction operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 10
	v := []float64{1.0, 2.0, 3.0}
	v = vec.Sub(v, val) // v is now {-9.0, -8.0, -7.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 5.0}
	w := []float64{3.0, 2.0, 4.0}
	q := vec.Sub(v, w) // q is now {-2.0, 0.0, 1.0}

The original arguments are not modified in this function.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Sub(v []float64, val interface{}) []float64 {
	c := Clone(v)
	switch w := val.(type) {
	case float64:
		for i := range c {
			c[i] -= w
		}
	case []float64:
		if len(c) != len(w) {
			panic(fmt.Sprintf(errStrings[5], "Sub()", len(c), len(w)))
		}
		for i := range c {
			c[i] -= w[i]
		}
	default:
		panic(fmt.Sprintf(errStrings[6], "Mul()", w))
	}
	return c
}

/*
Div takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the division operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 1.0
	v := []float64{1.0, 2.0, 3.0}
	v = vec.Div(v, val) // v is now {1.0, 2.0, 3.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	q := vec.Div(v, w) // q is now {0.33, 1.0, 0.75}

The original arguments are not modified in this function.
when the first argument is a float64, it cannot be 0.0. as division by zero is
not allowed.

In the case where the second argument is a []float64, the length of both
arguments must be equal. Additionally, the second argument must not contain
any elements whose value is 0.0.
*/
func Div(v []float64, val interface{}) []float64 {
	c := Clone(v)
	switch w := val.(type) {
	case float64:
		if w == 0.0 {
			panic(fmt.Sprintf(errStrings[7], "Div()"))
		}
		for i := range c {
			c[i] /= w
		}
	case []float64:
		if len(c) != len(w) {
			panic(fmt.Sprintf(errStrings[5], "Div()", len(c), len(w)))
		}
		for i := range w {
			if w[i] == 0.0 {
				panic(fmt.Sprintf(errStrings[8], "Div()", i))
			}
		}
		for i := range c {
			c[i] /= w[i]
		}
	default:
		panic(fmt.Sprintf(errStrings[6], "Mul()", w))
	}
	return c
}

/*
Dot returns the sum of the element-wise multiplication of two []float64s passed
to it. The passed slices are not altered in this function.
*/
func Dot(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic(fmt.Sprintf(errStrings[5], "Dot()", len(v1), len(v2)))
	}
	result := 0.0
	for i := range v1 {
		result += (v1[i] * v2[i])
	}
	return result
}
