package mat

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	rows := 13
	cols := 7
	m := New()
	if m.r != 0 {
		t.Errorf("expected m.r to be 0, but got %d", m.r)
	}
	if m.c != 0 {
		t.Errorf("expected m.c to be 0, but got %d", m.c)
	}
	if len(m.vals) != 0 {
		t.Errorf("expected len(m.vals) to be 0, but got %d", len(m.vals))
	}
	if cap(m.vals) != 0 {
		t.Errorf("expected cap(m.vals) to be 0, but got %d", cap(m.vals))
	}

	m = New(rows)
	if m.r != rows {
		t.Errorf("expected m.r to be %d, but got %d", rows, m.r)
	}
	if m.c != rows {
		t.Errorf("expected m.c to be %d, but got %d", rows, m.c)
	}
	if m.vals == nil {
		t.Errorf("New mat.vals not initilized")
	}
	if len(m.vals) != rows*rows {
		t.Errorf("expected len(m.vals) to be %d, but got %d", rows*rows, len(m.vals))
	}
	if cap(m.vals) != 2*rows*rows {
		t.Errorf("expected cap(m.vals) to be %d, but got %d", 2*rows*rows, cap(m.vals))
	}
	m = New(rows, cols)
	if m.r != rows {
		t.Errorf("New mat.r is %d, expected %d", m.r, rows)
	}
	if m.c != cols {
		t.Errorf("New mat.c is %d, expected %d", m.c, cols)
	}
	if m.vals == nil {
		t.Errorf("New mat.vals not initilized")
	}
	if len(m.vals) != rows*cols {
		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
	}
	if cap(m.vals) != 2*rows*cols {
		t.Errorf("cap(mat.vals) is %d, expected %d", cap(m.vals), 2*rows*cols)
	}

	m = New(rows, cols, rows*cols)
	if m.r != rows {
		t.Errorf("New mat.r is %d, expected %d", m.r, rows)
	}
	if m.c != cols {
		t.Errorf("New mat.c is %d, expected %d", m.c, cols)
	}
	if m.vals == nil {
		t.Errorf("New mat.vals not initilized")
	}
	if len(m.vals) != rows*cols {
		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
	}
	if cap(m.vals) != rows*cols {
		t.Errorf("cap(mat.vals) is %d, expected %d", cap(m.vals), rows*cols)
	}
}

func TestFromSlice(t *testing.T) {
	rows := 11
	cols := 5
	s := make([][]float64, rows)
	for i := range s {
		s[i] = make([]float64, cols)
	}
	for i := range s {
		for j := range s[i] {
			s[i][j] = float64(i + j)
		}
	}
	m := FromSlice(s)
	if len(m.vals) != rows*cols {
		t.Errorf("expected length of mat to be %d, but got %d", rows*cols, len(m.vals))
	}
	if cap(m.vals) != 2*rows*cols {
		t.Errorf("expected capacity of mat to be %d, but got %d", 2*rows*cols, cap(m.vals))
	}
	idx := 0
	for i := range s {
		for j := range s[i] {
			if s[i][j] != m.vals[idx] {
				t.Errorf("slice[%d][%d]: %f, mat: %f", i, j,
					s[i][j], m.vals[idx])
			}
			idx++
		}
	}

}

func TestIsJagged(t *testing.T) {
	s := make([][]float64, 10)
	for i := range s {
		s[i] = make([]float64, i+1)
	}
	if !isJagged(s) {
		t.Errorf("Jagged 2D slice passed the jagged test...")
	}
	q := make([][]float64, 11)
	for i := range q {
		q[i] = make([]float64, 5)
	}
	if isJagged(q) {
		t.Errorf("Non-jagged 2D slice failed the jagged test...")
	}
}

func TestFrom1DSlice(t *testing.T) {
	rows := 1
	cols := 113
	s := make([]float64, cols)
	for i := 0; i < len(s); i++ {
		s[i] = float64(i * i)
	}
	m := From1DSlice(s)
	if len(m.vals) != rows*cols {
		t.Errorf("expected length of mat to be %d, but got %d", rows*cols, len(m.vals))
	}
	if cap(m.vals) != 2*rows*cols {
		t.Errorf("expected capacity of mat to be %d, but got %d", 2*rows*cols, cap(m.vals))
	}
	for i := 0; i < len(s); i++ {
		if s[i] != m.vals[i] {
			t.Errorf("slice[%d]: %f, mat: %f", i, s[i], m.vals[i])
		}
	}
}

func TestFromCSV(t *testing.T) {
	rows := 4
	cols := 4
	filename := "test.csv"
	str := "1.0,1.0,2.0,3.0\n"
	str += "5.0,8.0,13.0,21.0\n"
	str += "34.0,55.0,89.0,144.0\n"
	str += "233.0,377.0,610.0,987.0\n"
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	m := FromCSV(filename)
	if len(m.vals) != rows*cols {
		t.Errorf("expected length of mat to be %d, but got %d", rows*cols, len(m.vals))
	}
	if cap(m.vals) != rows*cols {
		t.Errorf("expected capacity of mat to be %d, but got %d", rows*cols, cap(m.vals))
	}
	if m.vals[0] != 1.0 || m.vals[1] != 1.0 {
		t.Errorf("The first two entries are not 1.0: %f, %f", m.vals[0], m.vals[1])
	}
	for i := 2; i < m.r*m.c; i++ {
		if m.vals[i] != (m.vals[i-1] + m.vals[i-2]) {
			t.Errorf("expected %f, got %f", m.vals[i-1]+m.vals[i-2], m.vals[i])
		}
	}
	os.Remove(filename)
}

func TestReshape(t *testing.T) {
	s := make([]float64, 120)
	for i := 0; i < len(s); i++ {
		s[i] = float64(i * 3)
	}
	m := From1DSlice(s)
	m.Reshape(10, 12)
	if m.r != 10 {
		t.Errorf("expected rows = 10, got %d", m.r)
	}
	if m.c != 12 {
		t.Errorf("expected cols = 12, got %d", m.c)
	}
	for i := 0; i < len(s); i++ {
		if m.vals[i] != s[i] {
			t.Errorf("at index %d, expected %f, got %f", i, s[i], m.vals[i])
		}
	}
}

func TestDims(t *testing.T) {
	m := New(11, 10)
	r, c := m.Dims()
	if m.r != r {
		t.Errorf("m.r expected 11, got %d", m.r)
	}
	if m.c != c {
		t.Errorf("m.r expected 10, got %d", m.c)
	}
}

func TestVals(t *testing.T) {
	m := New(22, 22).Ones()
	s := m.Vals()
	if len(s) != 22*22 {
		t.Errorf("expected len(s) to be %d, got %d", 22*22, len(s))
	}
	for i := 0; i < len(s); i++ {
		if s[i] != 1.0 {
			t.Errorf("At index %d: expected 1.0, got %f", i, s[i])
		}
	}
}

func TestToSlice(t *testing.T) {
	rows := 13
	cols := 21
	m := New(rows, cols)
	for i := 0; i < m.r*m.c; i++ {
		m.vals[i] = float64(i)
	}
	s := m.ToSlice()
	if m.r != len(s) {
		t.Errorf("mat.r: %d and len(s): %d must match", m.r, len(s))
	}
	if m.c != len(s[0]) {
		t.Errorf("mat.c: %d and len(s[0]): %d must match", m.c, len(s[0]))
	}
	idx := 0
	for i := range s {
		for j := range s[i] {
			if s[i][j] != m.vals[idx] {
				t.Errorf("slice[%d][%d]: %f, mat: %f", i, j,
					s[i][j], m.vals[idx])
			}
			idx++
		}
	}
}

func TestToCSV(t *testing.T) {
	m := New(23, 17).Inc()
	filename := "tocsv_test.csv"
	m.ToCSV(filename)
	n := FromCSV(filename)
	if !n.Equals(m) {
		t.Errorf("m and n are not equal")
	}
	os.Remove(filename)
}

func TestAt(t *testing.T) {
	rows := 17
	cols := 13
	m := New(rows, cols).Inc()
	idx := 0
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m.At(i, j) != m.vals[idx] {
				t.Errorf("at index %d expected %f, got %f", i, m.vals[idx], m.At(i, j))
			}
			idx++
		}
	}
}

func TestMap(t *testing.T) {
	rows := 132
	cols := 24
	f := func(i *float64) {
		*i = 1.0
		return
	}
	m := New(rows, cols)
	m.Map(f)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("At %d, expected 1.0, got %f", i, m.vals[i])
		}
	}
}

func TestOnes(t *testing.T) {
	rows := 13
	cols := 7
	m := New(rows, cols).Ones()
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("At %d, expected 1.0, got %f", i, m.vals[i])
		}
	}
}

func TestInc(t *testing.T) {
	rows := 17
	cols := 3
	m := New(rows, cols).Inc()
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] != float64(i) {
			t.Errorf("At %d, expected %f, got %f", i, float64(i), m.vals[i])
		}
	}
}

func TestReset(t *testing.T) {
	row := 3
	col := 4
	m := New(row, col).Inc()
	m.Reset()
	for i := 0; i < row*col; i++ {
		if m.vals[i] != 0.0 {
			t.Errorf("at index %d, not equal to 0.0", i)
		}
	}
}

func TestSetAllTo(t *testing.T) {
	row := 3
	col := 4
	val := 11.0
	m := New(row, col).SetAllTo(val)
	for i := 0; i < row*col; i++ {
		if m.vals[i] != val {
			t.Errorf("at index %d, not equal to %f", i, val)
		}
	}
}

func TestRand(t *testing.T) {
	row := 31
	col := 42
	m := New(row, col).Rand()
	for i := 0; i < row*col; i++ {
		if m.vals[i] < 0.0 || m.vals[i] >= 1.0 {
			t.Errorf("at index %d, expected [0, 1.0), got %f", i, m.vals[i])
		}
	}
	m.Rand(100.0)
	for i := 0; i < row*col; i++ {
		if m.vals[i] < 0.0 || m.vals[i] >= 100.0 {
			t.Errorf("at index %d, expected [0, 100.0), got %f", i, m.vals[i])
		}
	}
	m.Rand(-12.0, 2.0)
	for i := 0; i < row*col; i++ {
		if m.vals[i] < -12.0 || m.vals[i] >= 2.0 {
			t.Errorf("at index %d, expected [-12.0, 2.0), got %f", i, m.vals[i])
		}
	}
}

func TestCol(t *testing.T) {
	row := 3
	col := 4
	m := New(row, col).Inc()
	for i := 0; i < col; i++ {
		got := m.Col(i)
		for j := 0; j < row; j++ {
			if got.vals[j] != m.vals[j*m.c+i] {
				t.Errorf("At index %v Col(%v), got %f, want %f", j, i,
					got.vals[j], m.vals[j*m.c+i])
			}
		}
	}
}

func BenchmarkCol(b *testing.B) {
	m := New(1721, 311).Inc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Col(211)
	}
}

func TestRow(t *testing.T) {
	row := 3
	col := 4
	m := New(row, col).Inc()
	idx := 0
	for i := 0; i < row; i++ {
		got := m.Row(i)
		for j := 0; j < col; j++ {
			if got.vals[j] != m.vals[idx] {
				t.Errorf("At index %v Col(%v), got %f, want %f", j, i,
					got.vals[j], m.vals[j*m.r+i])
			}
			idx++
		}
	}
}

func BenchmarkRow(b *testing.B) {
	m := New(1721, 311).Inc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Row(211)
	}
}

func TestEquals(t *testing.T) {
	m := New(13, 12)
	if !m.Equals(m) {
		t.Errorf("m is not equal iteself")
	}
}

func TestCopy(t *testing.T) {
	m := New(13, 13).Inc()
	n := m.Copy()
	for i := 0; i < 13*13; i++ {
		if m.vals[i] != n.vals[i] {
			t.Errorf("at %d, expected %f, got %f", i, m.vals[i], n.vals[i])
		}
	}
}

func TestT(t *testing.T) {
	m := New(12, 3).Inc()
	n := m.T()
	p := m.ToSlice()
	q := n.ToSlice()
	for i := 0; i < m.r; i++ {
		for j := 0; j < m.c; j++ {
			if p[i][j] != q[j][i] {
				t.Errorf("at %d, %d, expected %f, got %f", i, j, p[i][j], q[j][i])
			}
		}
	}
}

func BenchmarkT(b *testing.B) {
	m := New(1000, 251).Inc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.T()
	}
}

func TestFilter(t *testing.T) {
	m := New(100, 21).Inc()
	neg := func(i *float64) bool {
		if *i < 0.0 {
			return true
		}
		return false
	}
	r := m.Filter(neg)
	if r != nil {
		t.Errorf("Found negative values in Inc()")
	}
	m.Ones()
	one := func(i *float64) bool {
		if *i == 1.0 {
			return true
		}
		return false
	}
	r = m.Filter(one)
	if len(m.vals) != len(r.vals) {
		t.Errorf("not all values of Ones came through the filter")
	}
}

func TestAll(t *testing.T) {
	m := New(100, 21).Inc()
	pos := func(i *float64) bool {
		if *i >= 0.0 {
			return true
		}
		return false
	}
	if !m.All(pos) {
		t.Errorf("All(pos) is false for Inc()")
	}
	notOne := func(i *float64) bool {
		if *i != 1.0 {
			return true
		}
		return false
	}
	m.Ones()
	if m.All(notOne) {
		t.Errorf("m.Ones() has non-one values in it")
	}
}

func TestAny(t *testing.T) {
	m := New(100, 21).Inc()
	neg := func(i *float64) bool {
		if *i < 0.0 {
			return true
		}
		return false
	}
	if m.Any(neg) {
		t.Errorf("Any(neg) is true for Inc()")
	}
	one := func(i *float64) bool {
		if *i == 1.0 {
			return true
		}
		return false
	}
	m.Ones()
	if !m.Any(one) {
		t.Errorf("m.Ones() has no values equal to 1.0 in it")
	}
}

func TestMul(t *testing.T) {
	m := New(10, 11).Inc()
	n := m.Copy()
	m.Mul(m)
	for i := 0; i < 110; i++ {
		if m.vals[i] != n.vals[i]*n.vals[i] {
			t.Errorf("expected %f, got %f", n.vals[i]*n.vals[i], m.vals[i])
		}
	}
}

func BenchmarkMul(b *testing.B) {
	n := New(1000, 1000).Inc()
	m := New(1000, 1000).Inc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mul(n)
	}
}

func TestAdd(t *testing.T) {
	m := New(10, 11).Inc()
	n := m.Copy()
	m.Add(m)
	for i := 0; i < 110; i++ {
		if m.vals[i] != 2.0*n.vals[i] {
			t.Errorf("expected %f, got %f", 2.0*n.vals[i], m.vals[i])
		}
	}
}

func TestSub(t *testing.T) {
	m := New(10, 11).Inc()
	m.Sub(m)
	for i := 0; i < 110; i++ {
		if m.vals[i] != 0.0 {
			t.Errorf("expected 0.0, got %f", m.vals[i])
		}
	}
}

func TestDiv(t *testing.T) {
	m := New(10, 11).Inc()
	m.vals[0] = 1.0
	m.Div(m)
	for i := 0; i < 110; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("expected 1.0, got %f", m.vals[i])
		}
	}
}

func TestScale(t *testing.T) {
	m := New(12, 13).Inc()
	n := m.Copy()
	m.Scale(1.7)
	for i := 0; i < m.r*m.c; i++ {
		if m.vals[i] != 1.7*n.vals[i] {
			t.Errorf("At %d, expected %f, got %f", i, 1.7*n.vals[i], m.vals[i])
		}
	}
}

func TestSum(t *testing.T) {
	row := 12
	col := 17
	m := New(row, col).Ones()
	for i := 0; i < row; i++ {
		q := m.Sum(0, i)
		if q != float64(col) {
			t.Errorf("at row %d expected sum to be %d, got %f", i, col, q)
		}
	}
	for i := 0; i < col; i++ {
		q := m.Sum(1, i)
		if q != float64(row) {
			t.Errorf("at col %d expected sum to be %d, got %f", i, row, q)
		}
	}
}

func TestAverage(t *testing.T) {
	row := 12
	col := 17
	m := New(row, col).Ones()
	for i := 0; i < row; i++ {
		q := m.Average(0, i)
		if q != 1.0 {
			t.Errorf("at row %d expected average to be 1.0, got %f", i, q)
		}
	}
	for i := 0; i < col; i++ {
		q := m.Average(1, i)
		if q != 1.0 {
			t.Errorf("at col %d expected average to be 1.0, got %f", i, q)
		}
	}
}

func TestProd(t *testing.T) {
	row := 12
	col := 17
	m := New(row, col).Ones()
	for i := 0; i < row; i++ {
		q := m.Prod(0, i)
		if q != 1.0 {
			t.Errorf("at row %d expected product to be 1.0, got %f", i, q)
		}
	}
	for i := 0; i < col; i++ {
		q := m.Prod(1, i)
		if q != 1.0 {
			t.Errorf("at col %d expected product to be 1.0, got %f", i, q)
		}
	}
}

func TestStd(t *testing.T) {
	row := 12
	col := 17
	m := New(row, col).Ones()
	for i := 0; i < row; i++ {
		q := m.Std(0, i)
		if q != 0.0 {
			t.Errorf("at row %d expected std-div to be 0.0, got %f", i, q)
		}
	}
	for i := 0; i < col; i++ {
		q := m.Std(1, i)
		if q != 0.0 {
			t.Errorf("at col %d expected product to be 0.0, got %f", i, q)
		}
	}
}

func TestDot(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col).Inc()
	n := New(col, row).Inc()
	o := m.Dot(n)
	if o.r != row {
		t.Errorf("o.r: expected %d, got %d", row, o.r)
	}
	if o.c != row {
		t.Errorf("o.c: expected %d, got %d", row, o.c)
	}
	p := New(row, row)
	q := o.Dot(p)
	for i := 0; i < row*row; i++ {
		if q.vals[i] != 0.0 {
			t.Errorf("at index %d expected 0.0 got %f", i, q.vals[i])
		}
	}
}

func BenchmarkDot(b *testing.B) {
	m := New(150, 130).Inc()
	n := New(130, 150).Inc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Dot(n)
	}
}

func TestToString(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col).Inc()
	fmt.Println(m.ToString())
}

func TestAppendCol(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col).Inc()
	v := make([]float64, row)
	m.AppendCol(v)
	if m.c != col+1 {
		t.Errorf("Expected number of columns to be %d, but got %d", col+1, m.c)
	}
}

func TestAppendRow(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col).Inc()
	v := make([]float64, col)
	m.AppendRow(v)
	if m.r != row+1 {
		t.Errorf("Expected number of rows to be %d, but got %d", row+1, m.r)
	}
}

func TestConcat(t *testing.T) {
	var (
		row = 10
		col = 4
	)
	m := New(row, col).Inc()
	n := New(row, row).Inc()
	m.Concat(n)
	if m.c != row+col {
		t.Errorf("Expected number of cols to be %d, but got %d", row+col, m.c)
	}
	idx1 := 0
	idx2 := 0
	for i := 0; i < row; i++ {
		for j := 0; j < col+row; j++ {
			if j < col {
				if m.vals[i*m.c+j] != float64(idx1) {
					t.Errorf("At row %d, column %d, expected %f got %f", i, j,
						float64(idx1), m.vals[i*m.c+j])
				}
				idx1++
				continue
			}
			if m.vals[i*m.c+j] != float64(idx2) {
				t.Errorf("At row %d, column %d, expected %f got %f", i, j,
					float64(idx2), m.vals[i*m.c+j])
			}
			idx2++
		}
	}
}
