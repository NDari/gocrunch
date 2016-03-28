package vec

import "testing"

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
