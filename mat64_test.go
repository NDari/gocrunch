package mat64

import (
	"testing"
)

func TestI(t *testing.T) {
	m := I(4, 4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				if m.vals[i*4+j] != 1.0 {
					t.Errorf("I[%d,%d] == %f, want 1.0", i, j, m.vals[i*4+j])
				}
			} else {
				if m.vals[i*4+j] != 0.0 {
					t.Errorf("I[%d,%d] == %f, want 0.0", i, j, m.vals[i*4+j])
				}
			}
		}
	}
}

func TestCol(t *testing.T) {
	m := New(3, 4)
	for i := 0; i < 12; i++ {
		m.vals[i] = float64(i)
	}
	got := m.Col(2)
	if len(got) != 3 {
		t.Errorf("len(m.Col(3)) == %d, want 3", len(got))
	}
	want := []float64{2.0, 6.0, 10.0}
	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Errorf("m.Col(2][%v] == %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRow(t *testing.T) {
	m := New(3, 4)
	for i := 0; i < 12; i++ {
		m.vals[i] = float64(i)
	}
	got := m.Row(1)
	if len(got) != 4 {
		t.Errorf("len(m.Col(3)) == %d, want 4", len(got))
	}
	want := []float64{4.0, 5.0, 6.0, 7.0}
	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Errorf("got[%d] is %f, want %f", i, got[i], want[i])
		}
	}
}

func TestAt(t *testing.T) {
	m := New(3, 4)
	for i := 0; i < 12; i++ {
		m.vals[i] = float64(i)
	}
	got := m.At(2, 1)
	if got != 9.0 {
		t.Errorf("got %f, want %f", got, 9.0)
	}
}
