package numgo

import (
	"log"
	"numgo/mat"
	"numgo/vec"
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
					t.Errorf("mat.I[%v,%v] == %v, want 1.0", i, j, m[i][j])
				}
			} else {
				if m[i][j] != 0.0 {
					t.Errorf("mat.I[%v,%v] == %v, want 0.0", i, j, m[i][j])
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
				t.Errorf("mat.Ones[%v][%v]: expected 1.0, got %v", i, j, o[i][j])
			}
		}
	}
}

func TestVecOnes(t *testing.T) {
	o := vec.Ones(13)
	for i := 0; i < 13; i++ {
		if o[i] != 1.0 {
			t.Errorf("vec.Ones[%v]: expected 1.0, got %v", i, o[i])
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
		t.Errorf("mat.Col: got.NumRows == %v, want %v", len(got), row)
	}
	want := []float64{2.0, 6.0, 10.0}
	for i := 0; i < row; i++ {
		if want[i] != got[i] {
			t.Errorf("mat.Col: m[%v][2] == %v, want %v", i, got[i], want[i])
		}
	}
}

func TestMatRow(t *testing.T) {
	var (
		row = 3
		col = 4
	)
	m := mat.Inc(row, col)
	got := mat.Row(1, m)
	if len(got) != col {
		t.Errorf("mat.Row: len(got) == %v, want %v", len(got), col)
	}
	want := []float64{4.0, 5.0, 6.0, 7.0}
	for i := 0; i < col; i++ {
		if want[i] != got[i] {
			t.Errorf("mat.Row: m[1][%v] is %v, want %v", i, got[i], want[i])
		}
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
				t.Errorf("mat.T: [%v][%v] is %v, want %v", i, j, n[j][i], m[i][j])
			}
		}
	}
	s := mat.T(mat.T(m))
	if !mat.Equal(m, s) {
		t.Errorf("mat.T.T != mat")
		for i := 0; i < row; i++ {
			for j := 0; j < col; j++ {
				if s[i][j] != m[i][j] {
					t.Errorf("At [%v][%v]: want %v, got %v", i, j, m[i][j], s[i][j])
				}
			}
		}
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

func TestVecMul(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Mul(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != w[i] {
			t.Errorf("vec.Mul at index %v: expected 0.0, got %v", i, y[i])
		}
	}
	w = vec.Ones(12)
	y = vec.Mul(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != (v[i] * w[i]) {
			t.Errorf("vec.Mul at index %v: expected %v, got %v", i, v[i]*w[i], y[i])
		}
	}
}

func TestVecAdd(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Add(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("vec.Mul at index %v: expected %v, got %v", v[i], y[i])
		}
	}
	y = vec.Add(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != (v[i] * 2.0) {
			t.Errorf("vec.Mul at index %v: expected %v, got %v", v[i]*2.0, y[i])
		}
	}
}

func TestVecSub(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Sub(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("vec.Mul at index %v: expected %v, got %v", v[i], y[i])
		}
	}
	y = vec.Sub(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 0.0 {
			t.Errorf("vec.Mul at index %v: expected 0.0, got %v", y[i])
		}
	}
}

func TestVecDiv(t *testing.T) {
	v := vec.Ones(12)
	y := vec.Div(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 1.0 {
			t.Errorf("vec.Mul at index %v: expected 1.0, got %v", y[i])
		}
	}
	w := vec.Inc(12)
	y = vec.Div(w, v)
	for i := 0; i < 12; i++ {
		if y[i] != (w[i]) {
			t.Errorf("vec.Mul at index %v: expected %v, got %v", w[i], y[i])
		}
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

func TestVecMap(t *testing.T) {
	v := make([]float64, 17)
	y := vec.Map(func(i float64) float64 { return 1.0 }, v)
	for i := 0; i < 17; i++ {
		if y[i] != 1.0 {
			t.Errorf("vec.Map at index %v: xpected 1.0, got %v", y[i])
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

func TestVecDot(t *testing.T) {
	v := vec.Ones(13)
	w := make([]float64, 13)
	r := vec.Dot(v, w)
	if r != 0.0 {
		t.Errorf("vec.Dot: expected 0.0, got %v", r)
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

func TestVecReset(t *testing.T) {
	v := vec.Ones(22)
	vec.Reset(v)
	for i := 0; i < 22; i++ {
		if v[i] != 0.0 {
			t.Errorf("vec.Reset at index %v: expected 0.0, got %v", v[i])
		}
	}
}

func TestVecSum(t *testing.T) {
	v := vec.Ones(22)
	if vec.Sum(v) != 22.0 {
		t.Errorf("vec.Sum expected 22.0, got %v", vec.Sum(v))
	}
	vec.Reset(v)
	if vec.Sum(v) != 0.0 {
		t.Errorf("vec.Sum expected 0.0, got %v", vec.Sum(v))
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
