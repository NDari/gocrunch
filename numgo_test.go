package numgo

import (
	"log"
	"numgo/mat"
	"os"
	"testing"
)

func TestI(t *testing.T) {
	var (
		row = 4
	)
	m := mat.I(row)
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
	o := mat.Ones(row, col)
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
	m := mat.Inc(row, col)
	got := mat.Col(2, m)
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
	m := mat.Inc(row, col)
	got := mat.Row(1, m)
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
	m := mat.Inc(row, col)
	n := mat.T(m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n[j][i] != m[i][j] {
				t.Errorf("T[%v][%v] is %v, but expecting %v", i, j, n[j][i], m[i][j])
			}
		}
	}
	ttm := mat.T(mat.T(m))
	if !mat.Equal(m, ttm) {
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

func TestMul(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := mat.New(row, row)
	q := mat.I(row)
	if !mat.Equal(mat.Mul(m, q), m) {
		t.Errorf("A Square matrix Mul the identity matrix should be equal to itself")
	}
	m = mat.Inc(row, col)
	n := mat.New(row, col)
	o := mat.Mul(m, n)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if o[i][j] != 0.0 {
				t.Errorf("o[%v][%v] == %v, expected 0.0", i, j, o[i][j])
			}
		}
	}
	o = mat.Mul(m, m)
	m = mat.Apply(func(i float64) float64 { return i * i }, m)
	if !mat.Equal(o, m) {
		t.Errorf("m Mul m != m.Apply( i * i for each element i in m)")
	}
}

func TestApply(t *testing.T) {
	var (
		row = 4
		col = 4
	)
	m := mat.New(row, col)
	n := mat.Apply(func(i float64) float64 { return i + 1.0 }, m)
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
	m := mat.New(row, col)
	n := mat.New(col, row)
	o := mat.Dot(m, n)
	if len(o) != row {
		t.Errorf("Dot(m, n)'s numRows expected %v, got %v", row, len(o))
	}
	if len(o[0]) != row {
		t.Errorf("Dot(m, n)'s numCols expected %v, got %v", row, len(o[0]))
	}
	o = mat.Inc(row, row)
	p := mat.Dot(mat.I(row), o)
	if !mat.Equal(p, o) {
		t.Errorf("o x I != o...")
	}
}

func TestReset(t *testing.T) {
	var (
		row = 21
		col = 13
	)
	m := mat.Inc(row, col)
	m = mat.Reset(m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if m[i][j] != 0.0 {
				t.Errorf("Reset(m) at m[%v][%v] is equal %v", i, j, m[i][j])
			}
		}
	}
}

func TestDump(t *testing.T) {
	fileName := "output"
	m := mat.Inc(20, 30)
	mat.Dump(m, fileName)
	n := mat.Load(fileName)
	if !mat.Equal(m, n) {
		t.Errorf("Dumped 2D slice is not Equal Loaded 2D slice")
	}
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCopy(t *testing.T) {
	m := mat.Inc(3, 4)
	n := mat.Copy(m)
	if !mat.Equal(m, n) {
		t.Errorf("m != its copy")
	}
}

func TestAppendCol(t *testing.T) {
	m := mat.Inc(5, 7)
	v := []float64{12, 13, 17, 19, 21}
	m = mat.AppendCol(m, v)
	p := mat.Col(7, m)
	for i := 0; i < len(v); i++ {
		if v[i] != p[i] {
			t.Errorf("In AppendCol, expected %v, got %v", v[i], p[i])
		}
	}
}

func TestConcat(t *testing.T) {
	var (
		row = 3
		col = 7
	)
	m := mat.Inc(row, col)
	n := mat.I(row)
	o := mat.Concat(m, n)
	if len(o) != row {
		t.Errorf("len of concatinated 2Dslice is %v, expected %v", len(o), row)
	}
	for i := 0; i < len(o); i++ {
		if len(o[i]) != (row + col) {
			t.Errorf("length of concatinated slice o[%v] == %v, expected %v", i, len(o[i]), row+col)
		}
	}
	for i := 0; i < row; i++ {
		for j := 0; j < row+col; j++ {
			if j < col {
				if o[i][j] != m[i][j] {
					t.Errorf("concinated array at [%v][%v] is %v, expected %v", i, j, o[i][j], m[i][j])
				}
			} else {
				if o[i][j] != n[i][j-col] {
					t.Errorf("concinated array at [%v][%v] is %v, expected %v", i, j, o[i][j], n[i][j-col])
				}
			}
		}
	}
}
