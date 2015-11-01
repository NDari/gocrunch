package numgo

import (
	"log"
	"numgo/mat"
	"os"
	"testing"
)

func TestMatI(t *testing.T) {
	var (
		row = 4
	)
	m := mat.I(row)
	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			if i == j {
				if m[i][j] != 1.0 {
					t.Errorf("mat.I[%d,%d] == %f, want 1.0", i, j, m[i][j])
				}
			} else {
				if m[i][j] != 0.0 {
					t.Errorf("mat.I[%d,%d] == %f, want 0.0", i, j, m[i][j])
				}
			}
		}
	}
}

func TestMatOnes(t *testing.T) {
	var (
		row = 12
		col = 7
	)
	o := mat.Ones(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if o[i][j] != 1.0 {
				t.Errorf("mat.Ones[%d][%d]: expected 1.0, got %f", i, j, o[i][j])
			}
		}
	}
}

func TestMatCol(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := mat.Inc(row, col)
	got := mat.Col(2, m)
	if len(got) != row {
		t.Errorf("mat.Col: got.NumRows == %f, want %f", len(got), row)
	}
	want := []float64{2.0, 6.0, 10.0}
	for i := 0; i < row; i++ {
		if want[i] != got[i] {
			t.Errorf("mat.Col: m[%d][2] == %f, want %f", i, got[i], want[i])
		}
	}
	a1 := mat.Col(-1, m)
	a2 := mat.Col(3, m)
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			t.Errorf("mat.Col: at index %d, Col(-1, m) is %f, expected %f", i, a1[i], a2[i])
		}
	}
}

func BenchmarkMatCol(b *testing.B) {
	n := mat.Inc(1721, 311)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mat.Col(211, n)
	}
}

func TestMatRow(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := mat.Inc(row, col)
	got := mat.Row(2, m)
	want := []float64{8.0, 9.0, 10.0, 11.0}
	for i := 0; i < col; i++ {
		if got[i] != want[i] {
			t.Errorf("Mat.Row at index %d: want %f, got %f", i, want[i], got[i])
		}
	}
	a1 := mat.Row(-1, m)
	a2 := mat.Row(2, m)
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			t.Errorf("mat.Col: at index %d, Col(-1, m) is %f, expected %f", i, a1[i], a2[i])
		}
	}
}

func BenchmarkMatRow(b *testing.B) {
	n := mat.Inc(1721, 311)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mat.Row(511, n)
	}
}

func TestMatT(t *testing.T) {
	var (
		row = 5
		col = 7
	)
	m := mat.Inc(row, col)
	n := mat.T(m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n[j][i] != m[i][j] {
				t.Errorf("mat.T: [%d][%d] is %f, want %f", i, j, n[j][i], m[i][j])
			}
		}
	}
	s := mat.T(mat.T(m))
	if !mat.Equal(m, s) {
		t.Errorf("mat.T.T != mat")
		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				if s[i][j] != m[i][j] {
					t.Errorf("At [%d][%d]: want %f, got %f", i, j, m[i][j], s[i][j])
				}
			}
		}
	}
}

func BenchmarkMatT(b *testing.B) {
	m := mat.Inc(1000, 251)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mat.T(m)
	}
}

func TestMatMul(t *testing.T) {
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
	m = mat.Map(func(i float64) float64 { return i * i }, m)
	if !mat.Equal(o, m) {
		t.Errorf("m Mul m != m.Map( i * i for each element i in m)")
	}
}

func BenchmarkMatMul(b *testing.B) {
	n := mat.Inc(1000, 1000)
	m := mat.Inc(1000, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mat.Mul(m, n)
	}
}

func TestMatMap(t *testing.T) {
	var (
		row = 4
		col = 4
	)
	m := mat.New(row, col)
	n := mat.Map(func(i float64) float64 { return i + 1.0 }, m)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			if n[i][j] != 1.0 {
				t.Errorf("expected 1.0, got %v", n[i][j])
			}
		}
	}
}

func TestMatDot(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := mat.New(row, col)
	n := mat.New(col, row)
	o := mat.Dot(m, n)
	if len(o) != row {
		t.Errorf("mat.Dot(m, n)'s numRows expected %v, got %v", row, len(o))
	}
	if len(o[0]) != row {
		t.Errorf("mat.Dot(m, n)'s numCols expected %v, got %v", row, len(o[0]))
	}
	o = mat.Inc(row, row)
	p := mat.Dot(mat.I(row), o)
	if !mat.Equal(p, o) {
		t.Errorf("o x I != o...")
	}
}

func BenchmarkMatDot(b *testing.B) {
	m := mat.Inc(150, 130)
	n := mat.Inc(130, 150)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mat.Dot(m, n)
	}
}

func TestMatReset(t *testing.T) {
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

func TestMatDump(t *testing.T) {
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

func TestMatCopy(t *testing.T) {
	m := mat.Inc(3, 4)
	n := mat.Copy(m)
	if !mat.Equal(m, n) {
		t.Errorf("m != its copy")
	}
}

func TestMatAppendCol(t *testing.T) {
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

func TestMatConcat(t *testing.T) {
	var (
		row = 3
		col = 7
	)
	m := mat.Inc(row, col)
	n := mat.I(row)
	o := mat.Concat(m, n)
	if len(o) != row {
		t.Errorf("len of concatinated 2Dslice is %v, want %v", len(o), row)
	}
	for i := 0; i < len(o); i++ {
		if len(o[i]) != (row + col) {
			t.Errorf("len(o[%v]) is %v, want %v", i, len(o[i]), row+col)
		}
	}
	for i := 0; i < row; i++ {
		for j := 0; j < row+col; j++ {
			if j < col {
				if o[i][j] != m[i][j] {
					t.Errorf("[%v][%v] got %v, want %v", i, j, o[i][j], m[i][j])
				}
			} else {
				if o[i][j] != n[i][j-col] {
					t.Errorf("[%v][%v] got %v, want %v", i, j, o[i][j], n[i][j-col])
				}
			}
		}
	}
}
