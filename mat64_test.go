package mat64

import (
	"fmt"
	"testing"
)

func TestI(t *testing.T) {
	m := I(4)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == j {
				if m.vals[i*4+j] != 1.0 {
					t.Errorf("I[%v,%v] == %v, want 1.0", i, j, m.vals[i*4+j])
				}
			} else {
				if m.vals[i*4+j] != 0.0 {
					t.Errorf("I[%v,%v] == %v, want 0.0", i, j, m.vals[i*4+j])
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
		t.Errorf("len(m.Col(3)) == %v, want 3", len(got))
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
		t.Errorf("len(m.Col(3)) == %v, want 4", len(got))
	}
	want := []float64{4.0, 5.0, 6.0, 7.0}
	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Errorf("got[%v] is %v, want %v", i, got[i], want[i])
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
		t.Errorf("got %v, want 9.0", got)
	}
}

func TestTranspose(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := New(row, col)
	n := m.Transpose()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n.At(j, i) != m.At(i, j) {
				fmt.Println(n.At(j, i))
				t.Errorf("transpose.At(%v, %v) is %v, but m.At[%v, %v] is %v", i, j, n.At(j, i), j, i, m.At(i, j))
			}
		}
	}
	o := n.Transpose()
	for i := 0; i < row*col; i++ {
		if o.vals[i] != m.vals[i] {
			t.Errorf("mat.T.T != mat at %v", i)
		}
	}
}

func TestDot(t *testing.T) {
	var (
		row = 5
		col = 7
	)

	m := New(row, row)
	q := I(row)
	if !m.Dot(q).Equals(m) {
		t.Errorf("A Square matrix times the identity matrix should be equal to itself")
	}
	m = New(row, col)
	n := New(row, col)
	for i := 0; i < row*col; i++ {
		m.vals[i] = float64(i)
	}
	o := m.Dot(n)
	for i := 0; i < row*col; i++ {
		if o.vals[i] != 0.0 {
			t.Errorf("Dot product with 0 matrix, expect 0.0, got %v", o.vals[i])
		}
	}
	o = m.Dot(m)
	p := m.Apply(func(i float64) float64 { return i * i })
	for i := 0; i < row*col; i++ {
		if o.vals[i] != p.vals[i] {
			t.Errorf("Dot product matrix with itself, expect %v, got %v", o.vals[i], p.vals[i])
		}
	}

}

func TestDotInPlace(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := New(row, row)
	q := I(row)
	if !m.Dot(q).Equals(m.DotInPlace(q)) {
		t.Errorf("A Square matrix times the identity matrix should be equal to itself")
	}
	m = New(row, col)
	n := New(row, col)
	for i := 0; i < row*col; i++ {
		m.vals[i] = float64(i)
	}
	m.DotInPlace(n)
	for i := 0; i < row*col; i++ {
		if m.vals[i] != 0.0 {
			t.Errorf("DotInPlace product with 0 matrix, expect 0.0, got %v", m.vals[i])
		}
	}
	for i := 0; i < row*col; i++ {
		n.vals[i] = float64(i)
	}
	p := n.Apply(func(i float64) float64 { return i * i })
	n.DotInPlace(n)
	for i := 0; i < row*col; i++ {
		if n.vals[i] != p.vals[i] {
			t.Errorf("Dot product matrix with itself, expect %v, got %v", p.vals[i], n.vals[i])
		}
	}

}

func TestApply(t *testing.T) {
	m := New(4, 4)
	n := m.Apply(func(i float64) float64 { return i + 1.0 })
	for i := 0; i < 16; i++ {
		if n.vals[i] != 1.0 {
			t.Errorf("expected 1.0, got %v", n.vals[i])
		}
	}
}

func TestApplyInPlace(t *testing.T) {
	m := New(4, 4)
	m.ApplyInPlace(func(i float64) float64 { return i + 1.0 })
	for i := 0; i < 16; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("expected 1.0, got %v", m.vals[i])
		}
	}
}
