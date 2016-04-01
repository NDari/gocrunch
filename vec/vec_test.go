package vec

import (
	"fmt"
	"testing"
)

func TestPop(t *testing.T) {
	v := make([]float64, 2)
	x, v := Pop(v)
	if x != 0 {
		t.Errorf("expected 0, got %f", x)
	}
	if len(v) != 1 {
		t.Errorf("expected length of 1, got %d", len(v))
	}
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
	fmt.Println(v)
	v = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	v = Cut(v, 2, 4)
	fmt.Println(v)
}
