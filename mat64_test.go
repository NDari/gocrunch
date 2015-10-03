package mat64

import (
	"testing"
)

func TestI(t *testing.T) {
	var (
		row = 4
	)
	m := I(row)
	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			if i == j {
				if m[i][j] != 1.0 {
					t.Errorf("I[%v,%v] == %v, want 1.0", i, j, m[i][j])
				}
			} else {
				if m[i][j] != 0.0 {
					t.Errorf("I[%v,%v] == %v, want 0.0", i, j, m[i][j])
				}
			}
		}
	}
}

func TestOnes(t *testing.T) {
	var (
		row = 12
		col = 7
	)
	o := Ones(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if o[i][j] != 1.0 {
				t.Errorf("Ones[%v][%v]: expected 1.0, got %v", i, j, o[i][j])
			}
		}
	}
}

func TestCol(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := Inc(row, col)
	got := Col(2, m)
	if len(got) != row {
		t.Errorf("got.NumRows == %v, want %v", len(got), row)
	}
	want := []float64{2.0, 6.0, 10.0}
	for i := 0; i < row; i++ {
		if want[i] != got[i] {
			t.Errorf("m[%v][2] == %v, want %v", i, got[i], want[i])
		}
	}
}

func TestRow(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := Inc(row, col)
	got := Row(1, m)
	if len(got) != col {
		t.Errorf("len(got) == %v, want %v", len(got), col)
	}
	want := []float64{4.0, 5.0, 6.0, 7.0}
	for i := 0; i < col; i++ {
		if want[i] != got[i] {
			t.Errorf("m[1][%v] is %v, want %v", i, got[i], want[i])
		}
	}
}

func TestT(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := Inc(row, col)
	n := T(m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n[j][i] != m[i][j] {
				t.Errorf("T[%v][%v] is %v, but expecting %v", i, j, n[j][i], m[i][j])
			}
		}
	}
	ttm := T(T(m))
	if !Equal(m, ttm) {
		t.Errorf("mat.T.T != mat")
		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				if ttm[i][j] != m[i][j] {
					t.Errorf("At [%v][%v]: expected %v, got %v", i, j, m[i][j], ttm[i][j])
				}
			}
		}
	}
}

func TestTimes(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := New(row, row)
	q := I(row)
	if !Equal(Times(m, q), m) {
		t.Errorf("A Square matrix times the identity matrix should be equal to itself")
	}
	m = Inc(row, col)
	n := New(row, col)
	o := Times(m, n)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if o[i][j] != 0.0 {
				t.Errorf("o[%v][%v] == %v, expected 0.0", i, j, o[i][j])
			}
		}
	}
	o = Times(m, m)
	m = Apply(func(i float64) float64 { return i * i }, m)
	if !Equal(o, m) {
		t.Errorf("m times m != m.Apply( i * i for each element i in m)")
	}
}

func TestApply(t *testing.T) {
	var (
		row = 4
		col = 4
	)
	m := New(row, col)
	n := Apply(func(i float64) float64 { return i + 1.0 }, m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n[i][j] != 1.0 {
				t.Errorf("expected 1.0, got %v", n[i][j])
			}
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
	o := Dot(m, n)
	if len(o) != row {
		t.Errorf("Dot(m, n)'s numRows expected %v, got %v", row, len(o))
	}
	if len(o[0]) != row {
		t.Errorf("Dot(m, n)'s numCols expected %v, got %v", row, len(o[0]))
	}
	o = Inc(row, row)
	p := Dot(I(row), o)
	if !Equal(p, o) {
		t.Errorf("o x I != o...")
	}
}

func TestReset(t *testing.T) {
	var (
		row = 21
		col = 13
	)
	m := Inc(row, col)
	m = Reset(m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if m[i][j] != 0.0 {
				t.Errorf("Reset(m) at m[%v][%v] is equal %v", i, j, m[i][j])
			}
		}
	}
}

func TestCopy(t *testing.T) {
	m := Inc(3, 4)
	n := Copy(m)
	if !Equal(m, n) {
		t.Errorf("m != its copy")
	}
}
