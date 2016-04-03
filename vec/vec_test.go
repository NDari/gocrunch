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
