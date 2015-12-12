package mat

import "testing"

func TestNew(t *testing.T) {
	rows := 13
	cols := 7
	m := New(rows, cols)
	if m.r != rows {
		t.Errorf("New mat.r is %d, expected %d", m.r, rows)
	}
	if m.c != cols {
		t.Errorf("New mat.c is %d, expected %d", m.c, cols)
	}
	if m.vals == nil {
		t.Errorf("New mat.vals not initilized")
	}
	if m.work == nil {
		t.Errorf("New mat.work not initilized")
	}
	if len(m.vals) != rows*cols {
		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
	}
	if len(m.work) != rows*cols {
		t.Errorf("len(mat.work) is %d, expected %d", len(m.work), rows*cols)
	}
}

func TestBare(t *testing.T) {
	rows := 13
	cols := 7
	m := Bare(rows, cols)
	if m.r != rows {
		t.Errorf("New mat.r is %d, expected %d", m.r, rows)
	}
	if m.c != cols {
		t.Errorf("New mat.c is %d, expected %d", m.c, cols)
	}
	if m.vals == nil {
		t.Errorf("Bare mat.vals not initilized")
	}
	if m.work != nil {
		t.Errorf("Bare mat.work is initilized")
	}
	if len(m.vals) != rows*cols {
		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
	}
	if len(m.work) != 0 {
		t.Errorf("len(mat.work) is %d, expected 0", len(m.work))
	}
}
