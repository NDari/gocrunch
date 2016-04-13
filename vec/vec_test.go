package vec

import (
	"fmt"
	"sync"
	"testing"
)

func TestPop(t *testing.T) {
	v := make([]float64, 1)
	x, v := Pop(v)
	if x != 0.0 {
		t.Errorf("expected 0, got %f", x)
	}
	if len(v) != 0 {
		t.Errorf("expected length of 0, got %d", len(v))
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[0], "Pop()", "Pop()")
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		x, v = Pop(v)
	}()
	wg.Wait()
}

func TestPush(t *testing.T) {
	val := 3.0
	length := 3
	v := make([]float64, length)
	v = Push(v, val)
	if len(v) != 4 {
		t.Errorf("expected length of 4, got %d", len(v))
	}
	if v[length] != val {
		t.Errorf("expected %f, got %f", val, v[length])
	}
}

func TestShift(t *testing.T) {
	v := []float64{1.0, 2.0, 3.0, 4.0}
	x, v := Shift(v)
	if x != 1.0 {
		t.Errorf("expected 1.0, got %f", x)
	}
	if len(v) != 3 {
		t.Errorf("expected length of 3, got %d", len(v))
	}
	if v[0] != 2.0 {
		t.Errorf("expected first element to be 2.0, got %f", v[0])
	}
	x, v = Shift(v)
	x, v = Shift(v)
	x, v = Shift(v)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[0], "Shift()", "Shift()")
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		x, v = Shift(v)
	}()
	wg.Wait()
}

func TestSUnshift(t *testing.T) {
	v := []float64{1.0, 2.0, 3.0, 4.0}
	v = Unshift(v, 0.0)
	if len(v) != 5 {
		t.Errorf("expected length of 5, got %d", len(v))
	}
	if v[0] != 0.0 {
		t.Errorf("expected first element to be 0.0, got %f", v[0])
	}
}

func TestCut(t *testing.T) {
	v := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	v = Cut(v, 2)
	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	v = Cut(v, 2, 4)
	var wg sync.WaitGroup
	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[1], "Cut()", -1, len(v))
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, -1)
	}()

	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[1], "Cut()", len(v), len(v))
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, len(v))
	}()

	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[1], "Cut()", -1, len(v))
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, -1, 1)
	}()

	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[1], "Cut()", len(v), len(v))
			if r != expectedErr {
				t.Errorf("Expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, len(v), 1)
	}()

	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[2], "Cut()", len(v)+1, 1, len(v))
			if r != expectedErr {
				t.Errorf("expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, 1, len(v)+1)
	}()
	wg.Wait()

	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			expectedErr := fmt.Sprintf(errStrings[3], "Cut()", 1, 3)
			if r != expectedErr {
				t.Errorf("expected %s, got %v", expectedErr, r)
			}
			wg.Done()
		}()
		v = Cut(v, 3, 1)
	}()
	wg.Wait()
}

func TestEqual(t *testing.T) {
	if !Equal([]float64{1.0}, []float64{1.0}) {
		t.Errorf("expected equal, got not equal")
	}
	if Equal([]float64{1.0}, []float64{1.0, 2.0}) {
		t.Errorf("expected not equal, got equal")
	}
	if Equal([]float64{1.0}, []float64{2.0}) {
		t.Errorf("expected not equal, got equal")
	}
}

func TestSet(t *testing.T) {
	w := make([]float64, 14)
	w = Set(w, 10.0)
	for i := range w {
		if w[i] != 10.0 {
			t.Errorf("at index %d, expected 10.0, got %f", i, w[i])
		}
	}
}

func TestForeach(t *testing.T) {
	m := make([]float64, 10)
	m = Set(m, 2.0)
	double := func(i float64) float64 {
		return i * i
	}
	m = Foreach(m, double)
	for i := range m {
		if m[i] != 4.0 {
			t.Errorf("at index %d expected 4.0, got %f", i, m[i])
		}
	}
}

func TestAll(t *testing.T) {
	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	v := make([]float64, 10)
	v = Set(v, -12.0)
	if !All(v, negative) {
		t.Errorf("Expected all to be negative, got otherwise")
	}
}

func TestAny(t *testing.T) {
	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	v := make([]float64, 10)
	v = Set(v, 12.0)
	if Any(v, negative) {
		t.Errorf("Expected no negative values, got otherwise")
	}
}

func TestSum(t *testing.T) {
	v := make([]float64, 10)
	s := Sum(v)
	if s != 0.0 {
		t.Errorf("expected sum of zeros to be zero, but got %f", s)
	}
	v = Set(v, 1.0)
	s = Sum(v)
	if s != float64(len(v)) {
		t.Errorf("expected the sum to be %f, got %f", float64(len(v)), s)
	}
	v = Set(v, -1.0)
	s = Sum(v)
	if s != -float64(len(v)) {
		t.Errorf("expected the sum to be %f, got %f", -float64(len(v)), s)
	}
}

func TestProd(t *testing.T) {
	v := make([]float64, 10)
	s := Prod(v)
	if s != 0.0 {
		t.Errorf("expected product of zeros to be zero, but got %f", s)
	}
	v = Set(v, 1.0)
	s = Prod(v)
	if s != 1.0 {
		t.Errorf("expected the prod to be 1.0, got %f", s)
	}
	v = Set(v, -1.0)
	s = Prod(v)
	if s != 1.0 {
		t.Errorf("expected the prod to be -1.0, got %f", s)
	}
}

func TestAvg(t *testing.T) {
	v := make([]float64, 10)
	s := Avg(v)
	if s != 0.0 {
		t.Errorf("expected average of zeros to be zero, but got %f", s)
	}
	v = Set(v, 1.0)
	s = Avg(v)
	if s != 1.0 {
		t.Errorf("expected the average to be 1.0, got %f", s)
	}
	v = Set(v, -1.0)
	s = Avg(v)
	if s != -1.0 {
		t.Errorf("expected the average to be -1.0, got %f", s)
	}
}

func TestMul(t *testing.T) {
	v := make([]float64, 10)
	v = Set(v, 10.0)
	v = Mul(v, 3.0)
	for i := range v {
		if v[i] != 30.0 {
			t.Errorf("at index %d, expected, 30.0, got %f", i, v[i])
		}
	}
	v1 := make([]float64, 10)
	v2 := make([]float64, 10)
	v2 = Set(v2, 10.0)
	v3 := Mul(v1, v2)
	for i := range v3 {
		if v3[i] != 0.0 {
			t.Errorf("at index %d, expected 0.0, got %f", i, v3[i])
		}
	}
	v1 = Set(v1, 12.0)
	v3 = Mul(v1, v2)
	for i := range v3 {
		if v3[i] != 120.0 {
			t.Errorf("at index %d, expected 120.0, got %f", i, v3[i])
		}
	}
}

func TestAdd(t *testing.T) {
	v := make([]float64, 10)
	v = Set(v, 10.0)
	v = Add(v, 3.0)
	for i := range v {
		if v[i] != 13.0 {
			t.Errorf("at index %d, expected, 13.0, got %f", i, v[i])
		}
	}
	v1 := make([]float64, 10)
	v2 := make([]float64, 10)
	v2 = Set(v2, 10.0)
	v3 := Add(v1, v2)
	for i := range v3 {
		if v3[i] != 10.0 {
			t.Errorf("at index %d, expected 10.0, got %f", i, v3[i])
		}
	}
	v1 = Set(v1, 12.0)
	v3 = Add(v1, v2)
	for i := range v3 {
		if v3[i] != 22.0 {
			t.Errorf("at index %d, expected 22.0, got %f", i, v3[i])
		}
	}
}

func TestSub(t *testing.T) {
	v := make([]float64, 10)
	v = Set(v, 10.0)
	v = Sub(v, 3.0)
	for i := range v {
		if v[i] != 7.0 {
			t.Errorf("at index %d, expected, 7.0, got %f", i, v[i])
		}
	}
	v1 := make([]float64, 10)
	v2 := make([]float64, 10)
	v2 = Set(v2, 10.0)
	v3 := Sub(v1, v2)
	for i := range v3 {
		if v3[i] != -10.0 {
			t.Errorf("at index %d, expected -10..0, got %f", i, v3[i])
		}
	}
	v1 = Set(v1, 12.0)
	v3 = Sub(v1, v2)
	for i := range v3 {
		if v3[i] != 2.0 {
			t.Errorf("at index %d, expected 2.0, got %f", i, v3[i])
		}
	}
}

func TestDiv(t *testing.T) {
	v := make([]float64, 10)
	v = Set(v, 10.0)
	v = Div(v, 2.0)
	for i := range v {
		if v[i] != 5.0 {
			t.Errorf("at index %d, expected, 5.0, got %f", i, v[i])
		}
	}
	v1 := make([]float64, 10)
	v2 := make([]float64, 10)
	v1 = Set(v1, 10.0)
	v2 = Set(v2, 10.0)
	v3 := Div(v1, v2)
	for i := range v3 {
		if v3[i] != 1.0 {
			t.Errorf("at index %d, expected 1.0, got %f", i, v3[i])
		}
	}
	v1 = Set(v1, 10.0)
	v2 = Set(v2, -1.0)
	v3 = Div(v1, v2)
	for i := range v3 {
		if v3[i] != -10.0 {
			t.Errorf("at index %d, expected -10.0, got %f", i, v3[i])
		}
	}
}

func TestDot(t *testing.T) {
	v1 := make([]float64, 13)
	v2 := make([]float64, 13)
	v1 = Set(v1, 1.0)
	v2 = Set(v2, 3.0)
	res := Dot(v1, v2)
	if res != 13.0*3.0 {
		t.Errorf("expected result to be %f, but got %f", 13.0*3.0, res)
	}
}
