package Mat64

import (
	"fmt"
	"testing"
)

func TestIdentity(t *testing.T) {
	var (
		row = 4
	)
	m := Identity(row)
	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			if i == j {
				if m.Vals[i*row+j] != 1.0 {
					t.Errorf("I[%v,%v] == %v, want 1.0", i, j, m.Vals[i*4+j])
				}
			} else {
				if m.Vals[i*row+j] != 0.0 {
					t.Errorf("I[%v,%v] == %v, want 0.0", i, j, m.Vals[i*4+j])
				}
			}
		}
	}
}

func TestCol(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := New(row, col)
	for i := 0; i < row*col; i++ {
		m.Vals[i] = float64(i)
	}
	got := m.Col(2)
	if got.NumRows != row {
		t.Errorf("got.NumRows == %v, want %v", got.NumRows, row)
	}
	if got.NumCols != 1 {
		t.Errorf("got.NumCols == %v, want 1", got.NumCols)
	}

	want := []float64{2.0, 6.0, 10.0}
	for i := 0; i < row; i++ {
		if want[i] != got.Vals[i] {
			t.Errorf("m.Col(2][%v] == %v, want %v", i, got.Vals[i], want[i])
		}
	}
}

func TestRow(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := New(row, col)
	for i := 0; i < row*col; i++ {
		m.Vals[i] = float64(i)
	}
	got := m.Row(1)
	if got.NumCols != col {
		t.Errorf("got.NumCols == %v, want 4", got.NumCols)
	}
	want := []float64{4.0, 5.0, 6.0, 7.0}
	for i := 0; i < col; i++ {
		if want[i] != got.Vals[i] {
			t.Errorf("got.Vals[%v] is %v, want %v", i, got.Vals[i], want[i])
		}
	}
}

func TestAt(t *testing.T) {
	m := New(3, 4)
	for i := 0; i < 12; i++ {
		m.Vals[i] = float64(i)
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
				t.Errorf("transpose.At(%v, %v) is %v, but expecting %v", i, j, n.At(j, i), m.At(i, j))
			}
		}
	}
	if !m.Equals(m.Transpose().Transpose()) {
		t.Errorf("mat.T.T != mat")
	}
}

func TestTimes(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := New(row, row)
	q := Identity(row)
	if !m.Times(q).Equals(m) {
		t.Errorf("A Square matrix times the identity matrix should be equal to itself")
	}
	m = New(row, col)
	n := New(row, col)
	for i := 0; i < row*col; i++ {
		m.Vals[i] = float64(i)
	}
	o := m.Times(n)
	for i := 0; i < row*col; i++ {
		if o.Vals[i] != 0.0 {
			t.Errorf("Times product with 0 matrix, expect 0.0, got %v", o.Vals[i])
		}
	}
	o = m.Times(m)
	p := m.Apply(func(i float64) float64 { return i * i })
	if !o.Equals(p) {
		t.Errorf("m times m != m.Apply( i * i for each element i in m)")
	}
}

func TestTimesInPlace(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := New(row, row)
	q := Identity(row)
	if !m.Times(q).Equals(m.TimesInPlace(q)) {
		t.Errorf("A Square matrix times the identity matrix should be equal to itself")
	}
	m = New(row, col)
	n := New(row, col)
	for i := 0; i < row*col; i++ {
		m.Vals[i] = float64(i)
	}
	m.TimesInPlace(n)
	for i := 0; i < row*col; i++ {
		if m.Vals[i] != 0.0 {
			t.Errorf("TimesInPlace product with 0 matrix, expect 0.0, got %v", m.Vals[i])
		}
	}
	for i := 0; i < row*col; i++ {
		n.Vals[i] = float64(i)
	}
	p := n.Apply(func(i float64) float64 { return i * i })
	n.TimesInPlace(n)
	for i := 0; i < row*col; i++ {
		if n.Vals[i] != p.Vals[i] {
			t.Errorf("Times product matrix with itself, expect %v, got %v", p.Vals[i], n.Vals[i])
		}
	}

}

func TestApply(t *testing.T) {
	m := New(4, 4)
	n := m.Apply(func(i float64) float64 { return i + 1.0 })
	for i := 0; i < 16; i++ {
		if n.Vals[i] != 1.0 {
			t.Errorf("expected 1.0, got %v", n.Vals[i])
		}
	}
}

func TestApplyInPlace(t *testing.T) {
	m := New(4, 4)
	m.ApplyInPlace(func(i float64) float64 { return i + 1.0 })
	for i := 0; i < 16; i++ {
		if m.Vals[i] != 1.0 {
			t.Errorf("expected 1.0, got %v", m.Vals[i])
		}
	}
}

func TestDot(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col)
	n := New(col, row)
	o := m.Dot(n)
	if o.NumRows != row {
		t.Errorf("m.Dot(n).numRows expected 10, got %v", m.NumRows)
	}
	if o.NumCols != row {
		t.Errorf("m.Dot(n).numCols expected 10, got %v", m.NumCols)
	}
	for i := 0; i < row*row; i++ {
		o.Vals[i] = float64(i)
	}
	if !o.Dot(Identity(row)).Equals(o) {
		t.Errorf("mat64 x identity != mat64...")
	}
}
