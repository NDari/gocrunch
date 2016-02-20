package mat

import (
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	rows := 13
	cols := 7
	m := New(rows)
	if len(m) != rows {
		t.Errorf("Expected %d, got %d", rows, len(m))
	}
	for i := range m {
		if len(m[i]) != rows {
			t.Errorf("at index %d, expected %d, got %d", i, rows, len(m[i]))
		}
	}
	m = New(rows, cols)
	if len(m) != rows {
		t.Errorf("Expected %d, got %d", rows, len(m))
	}
	for i := range m {
		if len(m[i]) != cols {
			t.Errorf("at index %d, expected %d, got %d", i, cols, len(m[i]))
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
	if len(m) != rows {
		t.Errorf("expected length m to be %d, but got %d", rows, len(m))
	}
	for i := range m {
		if len(m[i]) != cols {
			t.Errorf("expected %d, got %d", cols, len(m[i]))
		}
	}
	if m[0][0] != 1.0 || m[0][1] != 1.0 {
		t.Errorf("The first two entries are not 1.0: %f, %f", m[0][0], m[0][1])
	}
	os.Remove(filename)
}

func TestFlatten(t *testing.T) {
	row, col := 5, 3
	m := New(row, col)
	n := Flatten(m)
	if len(n) != row*col {
		t.Errorf("expected %d, got %d", row*col, len(n))
	}
}

func TestToCSV(t *testing.T) {
	m := New(23, 17)
	filename := "tocsv_test.csv"
	ToCSV(m, filename)
	n := FromCSV(filename)
	if len(n) != len(m) {
		t.Errorf("expected %d, got %d", len(m), len(n))
	}
	if len(n[0]) != len(m[0]) {
		t.Errorf("expected %d, got %d", len(m[0]), len(n[0]))
	}
	os.Remove(filename)
}

func TestMap(t *testing.T) {
	rows := 132
	cols := 24
	f := func(i float64) float64 {
		return 1.0
	}
	m := New(rows, cols)
	Map(f, m)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m[i][j] != 1.0 {
				t.Errorf("expected 1.0, got %f", m[i][j])
			}
		}
	}
}

func BenchmarkMap(b *testing.B) {
	m := New(300, 1000)
	f := func(i float64) float64 {
		return 10.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(f, m)
	}
}

func TestSetAll(t *testing.T) {
	row := 3
	col := 4
	val := 11.0
	m := New(row, col)
	SetAll(m, val)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != val {
				t.Errorf("expected %f, got %f", val, m[i][j])
			}
		}
	}
}

func BenchmarkSetAll(b *testing.B) {
	m := New(300, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetAll(m, 10.0)
	}
}

func TestRand(t *testing.T) {
	row := 31
	col := 42
	m := New(row, col)
	Rand(m)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < 0.0 || m[i][j] >= 1.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
	Rand(m, 100.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < 0.0 || m[i][j] >= 100.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
	Rand(m, -12.0, 2.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < -12.0 || m[i][j] >= 2.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
}

func TestCol(t *testing.T) {
	row := 3
	col := 4
	m := New(row, col)
	for i := 0; i < col; i++ {
		got := Col(i, m)
		if len(got) != row {
			t.Errorf("expected %d, got %d", row, len(got))
		}
		for j := 0; j < row; j++ {
			if got[j] != m[j][i] {
				t.Errorf("expected %f, got %f", m[j][i], got[j])
			}
		}
	}
}

//
//func BenchmarkCol(b *testing.B) {
//	m := New(1721, 311).Inc()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = m.Col(211)
//	}
//}
//
//func TestRow(t *testing.T) {
//	row := 3
//	col := 4
//	m := New(row, col).Inc()
//	idx := 0
//	for i := 0; i < row; i++ {
//		got := m.Row(i)
//		for j := 0; j < col; j++ {
//			if got.vals[j] != m.vals[idx] {
//				t.Errorf("At index %v Col(%v), got %f, want %f", j, i,
//					got.vals[j], m.vals[j*m.r+i])
//			}
//			idx++
//		}
//	}
//}
//
//func BenchmarkRow(b *testing.B) {
//	m := New(1721, 311).Inc()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = m.Row(211)
//	}
//}
//
//func TestEquals(t *testing.T) {
//	m := New(13, 12)
//	if !m.Equals(m) {
//		t.Errorf("m is not equal iteself")
//	}
//}
//
//func TestCopy(t *testing.T) {
//	m := New(13, 13).Inc()
//	n := m.Copy()
//	for i := 0; i < 13*13; i++ {
//		if m.vals[i] != n.vals[i] {
//			t.Errorf("at %d, expected %f, got %f", i, m.vals[i], n.vals[i])
//		}
//	}
//}
//
//func TestT(t *testing.T) {
//	m := New(12, 3).Inc()
//	n := m.T()
//	p := m.ToSlice()
//	q := n.ToSlice()
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			if p[i][j] != q[j][i] {
//				t.Errorf("at %d, %d, expected %f, got %f", i, j, p[i][j], q[j][i])
//			}
//		}
//	}
//}
//
//func BenchmarkT(b *testing.B) {
//	m := New(1000, 251).Inc()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = m.T()
//	}
//}
//
//func TestFilter(t *testing.T) {
//	m := New(100, 21).Inc()
//	neg := func(i *float64) bool {
//		if *i < 0.0 {
//			return true
//		}
//		return false
//	}
//	r := m.Filter(neg)
//	if r != nil {
//		t.Errorf("Found negative values in Inc()")
//	}
//	m.Ones()
//	one := func(i *float64) bool {
//		if *i == 1.0 {
//			return true
//		}
//		return false
//	}
//	r = m.Filter(one)
//	if len(m.vals) != len(r.vals) {
//		t.Errorf("not all values of Ones came through the filter")
//	}
//}
//
//func TestAll(t *testing.T) {
//	m := New(100, 21).Inc()
//	pos := func(i *float64) bool {
//		if *i >= 0.0 {
//			return true
//		}
//		return false
//	}
//	if !m.All(pos) {
//		t.Errorf("All(pos) is false for Inc()")
//	}
//	notOne := func(i *float64) bool {
//		if *i != 1.0 {
//			return true
//		}
//		return false
//	}
//	m.Ones()
//	if m.All(notOne) {
//		t.Errorf("m.Ones() has non-one values in it")
//	}
//}
//
//func TestAny(t *testing.T) {
//	m := New(100, 21).Inc()
//	neg := func(i *float64) bool {
//		if *i < 0.0 {
//			return true
//		}
//		return false
//	}
//	if m.Any(neg) {
//		t.Errorf("Any(neg) is true for Inc()")
//	}
//	one := func(i *float64) bool {
//		if *i == 1.0 {
//			return true
//		}
//		return false
//	}
//	m.Ones()
//	if !m.Any(one) {
//		t.Errorf("m.Ones() has no values equal to 1.0 in it")
//	}
//}
//
//func TestMul(t *testing.T) {
//	m := New(10, 11).Inc()
//	n := m.Copy()
//	m.Mul(m)
//	for i := 0; i < 110; i++ {
//		if m.vals[i] != n.vals[i]*n.vals[i] {
//			t.Errorf("expected %f, got %f", n.vals[i]*n.vals[i], m.vals[i])
//		}
//	}
//}
//
//func BenchmarkMul(b *testing.B) {
//	n := New(1000, 1000).Inc()
//	m := New(1000, 1000).Inc()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		m.Mul(n)
//	}
//}
//
//func TestAdd(t *testing.T) {
//	m := New(10, 11).Inc()
//	n := m.Copy()
//	m.Add(m)
//	for i := 0; i < 110; i++ {
//		if m.vals[i] != 2.0*n.vals[i] {
//			t.Errorf("expected %f, got %f", 2.0*n.vals[i], m.vals[i])
//		}
//	}
//}
//
//func TestSub(t *testing.T) {
//	m := New(10, 11).Inc()
//	m.Sub(m)
//	for i := 0; i < 110; i++ {
//		if m.vals[i] != 0.0 {
//			t.Errorf("expected 0.0, got %f", m.vals[i])
//		}
//	}
//}
//
//func TestDiv(t *testing.T) {
//	m := New(10, 11).Inc()
//	m.vals[0] = 1.0
//	m.Div(m)
//	for i := 0; i < 110; i++ {
//		if m.vals[i] != 1.0 {
//			t.Errorf("expected 1.0, got %f", m.vals[i])
//		}
//	}
//}
//
//func TestScale(t *testing.T) {
//	m := New(12, 13).Inc()
//	n := m.Copy()
//	m.Scale(1.7)
//	for i := 0; i < m.r*m.c; i++ {
//		if m.vals[i] != 1.7*n.vals[i] {
//			t.Errorf("At %d, expected %f, got %f", i, 1.7*n.vals[i], m.vals[i])
//		}
//	}
//}
//
//func TestSum(t *testing.T) {
//	row := 12
//	col := 17
//	m := New(row, col).Ones()
//	for i := 0; i < row; i++ {
//		q := m.Sum(0, i)
//		if q != float64(col) {
//			t.Errorf("at row %d expected sum to be %d, got %f", i, col, q)
//		}
//	}
//	for i := 0; i < col; i++ {
//		q := m.Sum(1, i)
//		if q != float64(row) {
//			t.Errorf("at col %d expected sum to be %d, got %f", i, row, q)
//		}
//	}
//}
//
//func TestAverage(t *testing.T) {
//	row := 12
//	col := 17
//	m := New(row, col).Ones()
//	for i := 0; i < row; i++ {
//		q := m.Average(0, i)
//		if q != 1.0 {
//			t.Errorf("at row %d expected average to be 1.0, got %f", i, q)
//		}
//	}
//	for i := 0; i < col; i++ {
//		q := m.Average(1, i)
//		if q != 1.0 {
//			t.Errorf("at col %d expected average to be 1.0, got %f", i, q)
//		}
//	}
//}
//
//func TestProd(t *testing.T) {
//	row := 12
//	col := 17
//	m := New(row, col).Ones()
//	for i := 0; i < row; i++ {
//		q := m.Prod(0, i)
//		if q != 1.0 {
//			t.Errorf("at row %d expected product to be 1.0, got %f", i, q)
//		}
//	}
//	for i := 0; i < col; i++ {
//		q := m.Prod(1, i)
//		if q != 1.0 {
//			t.Errorf("at col %d expected product to be 1.0, got %f", i, q)
//		}
//	}
//}
//
//func TestStd(t *testing.T) {
//	row := 12
//	col := 17
//	m := New(row, col).Ones()
//	for i := 0; i < row; i++ {
//		q := m.Std(0, i)
//		if q != 0.0 {
//			t.Errorf("at row %d expected std-div to be 0.0, got %f", i, q)
//		}
//	}
//	for i := 0; i < col; i++ {
//		q := m.Std(1, i)
//		if q != 0.0 {
//			t.Errorf("at col %d expected product to be 0.0, got %f", i, q)
//		}
//	}
//}
//
//func TestDot(t *testing.T) {
//	var (
//		row = 10
//		col = 4
//	)
//	m := New(row, col).Inc()
//	n := New(col, row).Inc()
//	o := m.Dot(n)
//	if o.r != row {
//		t.Errorf("o.r: expected %d, got %d", row, o.r)
//	}
//	if o.c != row {
//		t.Errorf("o.c: expected %d, got %d", row, o.c)
//	}
//	p := New(row, row)
//	q := o.Dot(p)
//	for i := 0; i < row*row; i++ {
//		if q.vals[i] != 0.0 {
//			t.Errorf("at index %d expected 0.0 got %f", i, q.vals[i])
//		}
//	}
//}
//
//func BenchmarkDot(b *testing.B) {
//	m := New(150, 130).Inc()
//	n := New(130, 150).Inc()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_ = m.Dot(n)
//	}
//}
//
//func TestToString(t *testing.T) {
//	var (
//		row = 10
//		col = 4
//	)
//	m := New(row, col).Inc()
//	fmt.Println(m.ToString())
//}
//
//func TestAppendCol(t *testing.T) {
//	var (
//		row = 10
//		col = 4
//	)
//	m := New(row, col).Inc()
//	v := make([]float64, row)
//	m.AppendCol(v)
//	if m.c != col+1 {
//		t.Errorf("Expected number of columns to be %d, but got %d", col+1, m.c)
//	}
//}
//
//func TestAppendRow(t *testing.T) {
//	var (
//		row = 10
//		col = 4
//	)
//	m := New(row, col).Inc()
//	v := make([]float64, col)
//	m.AppendRow(v)
//	if m.r != row+1 {
//		t.Errorf("Expected number of rows to be %d, but got %d", row+1, m.r)
//	}
//}
//
//func TestConcat(t *testing.T) {
//	var (
//		row = 10
//		col = 4
//	)
//	m := New(row, col).Inc()
//	n := New(row, row).Inc()
//	m.Concat(n)
//	if m.c != row+col {
//		t.Errorf("Expected number of cols to be %d, but got %d", row+col, m.c)
//	}
//	idx1 := 0
//	idx2 := 0
//	for i := 0; i < row; i++ {
//		for j := 0; j < col+row; j++ {
//			if j < col {
//				if m.vals[i*m.c+j] != float64(idx1) {
//					t.Errorf("At row %d, column %d, expected %f got %f", i, j,
//						float64(idx1), m.vals[i*m.c+j])
//				}
//				idx1++
//				continue
//			}
//			if m.vals[i*m.c+j] != float64(idx2) {
//				t.Errorf("At row %d, column %d, expected %f got %f", i, j,
//					float64(idx2), m.vals[i*m.c+j])
//			}
//			idx2++
//		}
//	}
//}
//
//func TestSet(t *testing.T) {
//	m := New(10, 12).Inc()
//	m.Set(0, 10, 12)
//	if m.vals[10] != 12.0 {
//		t.Errorf("Expected 12.0, got %f", m.vals[10])
//	}
//}
//
//func TestCombineWithRows(t *testing.T) {
//	v := make([]float64, 5)
//	for i := range v {
//		v[i] = float64(i)
//	}
//	v[0] = -1.0
//	m := New(2, 5).Inc()
//	n := m.Copy()
//	n.CombineWithRows("add", v)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			if n.vals[i*n.c+j] != m.vals[i*m.c+j]+v[j] {
//				t.Errorf("expected %f, got %f", m.vals[i*m.c+j]+v[j], n.vals[i*n.c+j])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithRows("sub", v)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			if n.vals[i*n.c+j] != m.vals[i*m.c+j]-v[j] {
//				t.Errorf("expected %f, got %f", m.vals[i*m.c+j]-v[j], n.vals[i*n.c+j])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithRows("mul", v)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			if n.vals[i*n.c+j] != m.vals[i*m.c+j]*v[j] {
//				t.Errorf("expected %f, got %f", m.vals[i*m.c+j]*v[j], n.vals[i*n.c+j])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithRows("div", v)
//	for i := 0; i < m.r; i++ {
//		for j := 0; j < m.c; j++ {
//			if n.vals[i*n.c+j] != m.vals[i*m.c+j]/v[j] {
//				t.Errorf("expected %f, got %f", m.vals[i*m.c+j]/v[j], n.vals[i*n.c+j])
//			}
//		}
//	}
//
//}
//
//func TestCombineWithCols(t *testing.T) {
//	v := make([]float64, 5)
//	for i := range v {
//		v[i] = float64(i)
//	}
//	v[0] = -1.0
//	m := New(5, 2).Inc()
//	n := m.Copy()
//	n.CombineWithCols("add", v)
//	for i := 0; i < m.c; i++ {
//		for j := 0; j < m.r; j++ {
//			if n.vals[j*n.c+i] != m.vals[j*m.c+i]+v[j] {
//				t.Errorf("expected %f, got %f", m.vals[j*m.c+i]+v[j], n.vals[j*n.c+i])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithCols("sub", v)
//	for i := 0; i < m.c; i++ {
//		for j := 0; j < m.r; j++ {
//			if n.vals[j*n.c+i] != m.vals[j*m.c+i]-v[j] {
//				t.Errorf("expected %f, got %f", m.vals[j*m.c+i]-v[j], n.vals[j*n.c+i])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithCols("mul", v)
//	for i := 0; i < m.c; i++ {
//		for j := 0; j < m.r; j++ {
//			if n.vals[j*n.c+i] != m.vals[j*m.c+i]*v[j] {
//				t.Errorf("expected %f, got %f", m.vals[j*m.c+i]*v[j], n.vals[j*n.c+i])
//			}
//		}
//	}
//	n = m.Copy()
//	n.CombineWithCols("div", v)
//	for i := 0; i < m.c; i++ {
//		for j := 0; j < m.r; j++ {
//			if n.vals[j*n.c+i] != m.vals[j*m.c+i]/v[j] {
//				t.Errorf("expected %f, got %f", m.vals[j*m.c+i]/v[j], n.vals[j*n.c+i])
//			}
//		}
//	}
//
//}
