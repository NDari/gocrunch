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
	"reflect"
)

/*
Pop takes a []float64, and "pops" the last entry, returning it along with the
modified []float64. The other elements of the []float64 remain intact.
*/
func Pop(v []float64) (float64, []float64) {
	if len(v) == 0 {
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, cannot use Pop() on an empty []float64.\n"
		s = fmt.Sprintf(s, "Pop()")
		panic(s)
	}
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
	if len(v) == 0 {
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, cannot use Shift() on an empty []float64.\n"
		s = fmt.Sprintf(s, "Shift()")
		panic(s)
	}
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
		s := "In vec.%s, 1 or 2 ints must be passed in addition to the slice.\n"
		s += "However, %d ints were passed.\n"
		s = fmt.Sprintf(s, "Cut()", len(args))
		panic(s)
	}
	return v
}

// TODO: Add out of bounds check to Cut()

/*
Equal checks if two []float64s are equal, by checking that they have the same length,
and the same enteries in each index.
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
Set sets all elements in a []float64 to the passed value.
*/
func Set(v []float64, val float64) {
	for i := range v {
		v[i] = val
	}
}

/*
Foreach applies a function to each element of a []float64, storing the result
in the same index.  Consider:

	double := func(i float64) float64 {
		return i * i
	}
	v := []float64{1.0, 2.0, 3.0}
	vec.Foreach(v, double) // v is {1.0, 4.0, 9.0}

Thus the []float64 is modified by Foreach().
*/
func Foreach(v []float64, f func(float64) float64) {
	for i := range v {
		v[i] = f(v[i])
	}
}

/*
All checks to see if all elements of a []float64 return true for a passed
function. If not, All() returns false. consider:

	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	v := make([]float64, 10)
	vec.Set(v, -12.0)
	allNegatives := vev.All(v, negative) // true

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
	vec.Set(v, 12.0)
	anyNegatives := vev.Any(v, negative) // false

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
	vec.Mul(v, val) // v is now {10.0, 20.0, 30.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	vec.Mul(v, w) // v is now {3.0, 4.0, 12.0}

The original []float64 (the first argument) is modified, but the second is not.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Mul(v []float64, val interface{}) {
	switch w := val.(type) {
	case float64:
		for i := range v {
			v[i] *= w
		}
	case []float64:
		if len(v) != len(w) {
			fmt.Println("\ngocrunch/vec error.")
			s := "In vec.%s, the length of the first []float64 is %d, while\n "
			s += "the length of the second []float64 is 5d. They must match."
			s = fmt.Sprintf(s, "Mul()", len(v), len(w))
			panic(s)
		}
		for i := range v {
			v[i] *= w[i]
		}
	default:
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, only accepts a float64 or []float64 as the second argument.\n"
		s += "However, an argument of type %v was recieved."
		s = fmt.Sprintf(s, "Mul()", reflect.TypeOf(v))
		panic(s)
	}
}

/*
Add takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the addition operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 10
	v := []float64{1.0, 2.0, 3.0}
	vec.Mul(v, val) // v is now {11.0, 12.0, 13.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	vec.Mul(v, w) // v is now {4.0, 4.0, 7.0}

The original []float64 (the first argument) is modified, but the second is not.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Add(v []float64, val interface{}) {
	switch w := val.(type) {
	case float64:
		for i := range v {
			v[i] += w
		}
	case []float64:
		if len(v) != len(w) {
			fmt.Println("\ngocrunch/vec error.")
			s := "In vec.%s, the length of the first []float64 is %d, while\n "
			s += "the length of the second []float64 is 5d. They must match."
			s = fmt.Sprintf(s, "Add()", len(v), len(w))
			panic(s)
		}
		for i := range v {
			v[i] += w[i]
		}
	default:
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, only accepts a float64 or []float64 as the second argument.\n"
		s += "However, an argument of type %v was recieved."
		s = fmt.Sprintf(s, "Add()", reflect.TypeOf(v))
		panic(s)
	}
}

/*
Sub takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the subtraction operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 10
	v := []float64{1.0, 2.0, 3.0}
	vec.Sub(v, val) // v is now {-9.0, -8.0, -7.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 5.0}
	w := []float64{3.0, 2.0, 4.0}
	vec.Sub(v, w) // v is now {-2.0, 0.0, 1.0}

The original []float64 (the first argument) is modified, but the second is not.
In the case where the second argument is a []float64, the length of both
arguments must be equal.
*/
func Sub(v []float64, val interface{}) {
	switch w := val.(type) {
	case float64:
		for i := range v {
			v[i] -= w
		}
	case []float64:
		if len(v) != len(w) {
			fmt.Println("\ngocrunch/vec error.")
			s := "In vec.%s, the length of the first []float64 is %d, while\n "
			s += "the length of the second []float64 is 5d. They must match."
			s = fmt.Sprintf(s, "Sub()", len(v), len(w))
			panic(s)
		}
		for i := range v {
			v[i] -= w[i]
		}
	default:
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, only accepts a float64 or []float64 as the second argument.\n"
		s += "However, an argument of type %v was recieved."
		s = fmt.Sprintf(s, "Sub()", reflect.TypeOf(v))
		panic(s)
	}
}

/*
Div takes a []float64, and a second argument, which can be a float64 or a
[]float64, and applies the division operation on each element, storing
the result in the first []float64. for clarification, consider the case where
the second argument is a float64:

	val := 1.0
	v := []float64{1.0, 2.0, 3.0}
	vec.Div(v, val) // v is now {1.0, 2.0, 3.0}

The second argument can also we a []float64, as below:

	v := []float64{1.0, 2.0, 3.0}
	w := []float64{3.0, 2.0, 4.0}
	vec.Div(v, w) // v is now {0.33, 1.0, 0.75}

The original []float64 (the first argument) is modified, but the second is not.
when the first argument is a float64, it cannot be 0.0. as division by zero is
not allowed.

In the case where the second argument is a []float64, the length of both
arguments must be equal. Adiitionally, the second argument must not contain
any elements whos value is 0.0.
*/
func Div(v []float64, val interface{}) {
	switch w := val.(type) {
	case float64:
		if w == 0.0 {
			fmt.Println("\ngocrunch/vec error.")
			s := "In vec.%s, the devisor cannot be zero"
			s = fmt.Sprintf(s, "Div()")
			panic(s)
		}
		for i := range v {
			v[i] /= w
		}
	case []float64:
		if len(v) != len(w) {
			fmt.Println("\ngocrunch/vec error.")
			s := "In vec.%s, the length of the first []float64 is %d, while\n "
			s += "the length of the second []float64 is 5d. They must match."
			s = fmt.Sprintf(s, "Div()", len(v), len(w))
			panic(s)
		}
		for i := range w {
			if w[i] == 0.0 {
				fmt.Println("\ngocrunch/vec error.")
				s := "In vec.%s, the value of the second argument at index %d is zero.\n "
				s += "Division by zero is not allowed."
				s = fmt.Sprintf(s, "Div()", i)
				panic(s)
			}
		}
		for i := range v {
			v[i] /= w[i]
		}
	default:
		fmt.Println("\ngocrunch/vec error.")
		s := "In vec.%s, only accepts a float64 or []float64 as the second argument.\n"
		s += "However, an argument of type %v was recieved."
		s = fmt.Sprintf(s, "Div()", reflect.TypeOf(v))
		panic(s)
	}
}
