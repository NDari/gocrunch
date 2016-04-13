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

func TestI(t *testing.T) {
	row := 10
	m := I(row)
	for i := range m {
		for j := range m[i] {
			if i == j {
				if m[i][j] != 1.0 {
					t.Errorf("expected 1.0, got %f", m[i][j])
				}
			} else {
				if m[i][j] != 0.0 {
					t.Errorf("expected 0.0, got %f", m[i][j])
				}
			}
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

func TestForeach(t *testing.T) {
	rows := 132
	cols := 24
	f := func(i float64) float64 {
		return 1.0
	}
	m := New(rows, cols)
	m = Foreach(m, f)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m[i][j] != 1.0 {
				t.Errorf("expected 1.0, got %f", m[i][j])
			}
		}
	}
}

func BenchmarkForeach(b *testing.B) {
	m := New(300, 1000)
	f := func(i float64) float64 {
		return 10.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m = Foreach(m, f)
	}
}

func TestSet(t *testing.T) {
	row := 3
	col := 4
	val := 11.0
	m := New(row, col)
	m = Set(m, val)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != val {
				t.Errorf("expected %f, got %f", val, m[i][j])
			}
		}
	}
}

func BenchmarkSet(b *testing.B) {
	m := New(300, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m = Set(m, 10.0)
	}
}

func TestMul(t *testing.T) {
	row, col := 13, 12
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	m = Mul(m, 0.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != 0.0 {
				t.Errorf("expected 0.0, got %f", m[i][j])
			}
		}
	}
	row, col = 11, 13
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	v := make([]float64, col)
	m = Mul(m, v)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != 0.0 {
				t.Errorf("At row %d, col %d, expected 0.0, got %f", i, j, m[i][j])
			}
		}
	}
	row, col = 13, 121
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	n := Copy(m)
	m = Mul(m, m)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j]*n[i][j] {
				t.Errorf("expected %f, got %f", n[i][j]*n[i][j], m[i][j])
			}
		}
	}
}

func BenchmarkMul(b *testing.B) {
	m := New(1000, 1000)
	n := New(1000, 1000)
	q := New(1000, 1000)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*1000 + j)
			n[i][j] = 1.0
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q = Mul(m, n)
	}
}

func TestAdd(t *testing.T) {
	row, col := 13, 12
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	n := Copy(m)
	m = Add(m, 0.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j] {
				t.Errorf("expected %f got %f", n[i][j], m[i][j])
			}
		}
	}
	row, col = 11, 13
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	v := make([]float64, col)
	for i := range v {
		v[i] = 2.0
	}
	n = Copy(m)
	m = Add(m, v)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j]+2.0 {
				t.Errorf("At row %d, col %d, expected %f, got %f", i, j, n[i][j]+2.0, m[i][j])
			}
		}
	}
	row, col = 13, 121
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	n = Copy(m)
	m = Add(m, m)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j]+n[i][j] {
				t.Errorf("expected %f, got %f", n[i][j]*n[i][j], m[i][j])
			}
		}
	}
}

func TestSub(t *testing.T) {
	row, col := 13, 12
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	n := Copy(m)
	m = Sub(m, 0.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j] {
				t.Errorf("expected %f got %f", n[i][j], m[i][j])
			}
		}
	}
	row, col = 11, 13
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	v := make([]float64, col)
	for i := range v {
		v[i] = 2.0
	}
	n = Copy(m)
	m = Sub(m, v)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j]-2.0 {
				t.Errorf("At row %d, col %d, expected %f, got %f", i, j, n[i][j]-2.0, m[i][j])
			}
		}
	}
	row, col = 13, 121
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	m = Sub(m, m)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != 0.0 {
				t.Errorf("expected 0.0, got %f", m[i][j])
			}
		}
	}

}

func TestDiv(t *testing.T) {
	row, col := 13, 12
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	n := Copy(m)
	m = Div(m, 1.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j] {
				t.Errorf("expected %f got %f", n[i][j], m[i][j])
			}
		}
	}
	row, col = 11, 13
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	v := make([]float64, col)
	for i := range v {
		v[i] = 1.0
	}
	n = Copy(m)
	m = Div(m, v)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != n[i][j] {
				t.Errorf("At row %d, col %d, expected %f, got %f", i, j, n[i][j], m[i][j])
			}
		}
	}
	row, col = 13, 121
	m = New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	m[0][0] = 1.0
	m = Div(m, m)
	for i := range m {
		for j := range m[i] {
			if m[i][j] != 1.0 {
				t.Errorf("expected 1.0, got %f", m[i][j])
			}
		}
	}

}

func TestRand(t *testing.T) {
	row := 31
	col := 42
	m := Rand(row, col)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < 0.0 || m[i][j] >= 1.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
	m = Rand(row, col, 100.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < 0.0 || m[i][j] >= 100.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
	m = Rand(row, col, -12.0, 2.0)
	for i := range m {
		for j := range m[i] {
			if m[i][j] < -12.0 || m[i][j] >= 2.0 {
				t.Errorf("at index %d, expected [0, 1.0), got %f", i, m[i][j])
			}
		}
	}
}

func TestCol(t *testing.T) {
	row, col := 3, 5
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	for i := 0; i < col; i++ {
		got := Col(m, i)
		if len(got) != row {
			t.Errorf("expected %d, got %d", row, len(got))
		}
		for j := 0; j < row; j++ {
			if got[j] != m[j][i] {
				t.Errorf("expected %f, got %f", m[j][i], got[j])
			}
		}
	}
	for i := col; i > 0; i-- {
		got := Col(m, -i)
		if len(got) != row {
			t.Errorf("expected %d, got %d", row, len(got))
		}
		for j := 0; j < row; j++ {
			if got[j] != m[j][len(m[0])-i] {
				t.Errorf("expected %f, got %f", m[j][len(m[0])-i], got[j])
			}
		}
	}
}

func BenchmarkCol(b *testing.B) {
	m := New(1721, 311)
	q := make([]float64, 1721)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*1721 + j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q = Col(m, 211)
	}
}

func TestRow(t *testing.T) {
	row, col := 3, 5
	m := New(row, col)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*row + j)
		}
	}
	for i := 0; i < row; i++ {
		got := Row(m, i)
		if len(got) != col {
			t.Errorf("expected %d, got %d", row, len(got))
		}
		for j := 0; j < col; j++ {
			if got[j] != m[i][j] {
				t.Errorf("expected %f, got %f", m[j][i], got[j])
			}
		}
	}
	for i := row; i > 0; i-- {
		got := Row(m, -i)
		if len(got) != col {
			t.Errorf("expected %d, got %d", row, len(got))
		}
		for j := 0; j < col; j++ {
			if got[j] != m[len(m)-i][j] {
				t.Errorf("row %d expected %f, got %f", -i, m[len(m)-i][j], got[j])
			}
		}
	}
}

func BenchmarkRow(b *testing.B) {
	m := New(1721, 311)
	q := make([]float64, 311)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*1721 + j)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q = Row(m, 211)
	}
}

func TestEqual(t *testing.T) {
	m := New(13, 12)
	if !Equal(m, m) {
		t.Errorf("m is not equal iteself")
	}
}

func TestCopy(t *testing.T) {
	m := New(13, 13)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*13 + j)
		}
	}
	n := Copy(m)
	if !Equal(m, n) {
		t.Errorf("not equal to its own copy")
	}
}

func TestT(t *testing.T) {
	m := New(12, 3)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*12 + j)
		}
	}
	n := T(m)
	if len(n) != len(m[0]) {
		t.Errorf("expected %d, got %d", len(m[0]), len(n))
	}
	if len(n[0]) != len(m) {
		t.Errorf("expected %d, got %d", len(m), len(n[0]))
	}
}

func BenchmarkT(b *testing.B) {
	m := New(1000, 251)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = T(m)
	}
}

func TestAll(t *testing.T) {
	m := New(100, 21)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*100 + j)
		}
	}
	positive := func(i float64) bool {
		if i >= 0.0 {
			return true
		}
		return false
	}
	if !All(m, positive) {
		t.Errorf("All(positive) is false, expected true")
	}
	notOne := func(i float64) bool {
		if i != 1.0 {
			return true
		}
		return false
	}
	m = Set(m, 1.0)
	if All(m, notOne) {
		t.Errorf("m has non-one values in it, expected none")
	}
}

func TestAny(t *testing.T) {
	m := New(100, 21)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*100 + j)
		}
	}
	negative := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	if Any(m, negative) {
		t.Errorf("Any(negiative) is true, expected false")
	}
	notOne := func(i float64) bool {
		if i != 1.0 {
			return true
		}
		return false
	}
	m = Set(m, 1.0)
	if Any(m, notOne) {
		t.Errorf("has non-one values in it, expected none")
	}
}

func TestSum(t *testing.T) {
	row, col, val := 131, 12, 2.0
	m := New(row, col)
	m = Set(m, val)
	res := Sum(m)
	if res != float64(row*col)*val {
		t.Errorf("expected %f, got %f", float64(row*col)*val, res)
	}
	row = 12
	col = 17
	m = New(row, col)
	m = Set(m, 1.0)
	for i := 0; i < col; i++ {
		q := Sum(m, 1, i)
		if q != float64(row) {
			t.Errorf("at col %d expected sum to be %f, got %f", i, float64(row), q)
		}
	}
	for i := col; i > 0; i-- {
		q := Sum(m, 1, -i)
		if q != float64(row) {
			t.Errorf("at col %d expected sum to be %f, got %f", i, float64(row), q)
		}
	}
	m = Set(m, 1.0)
	for i := 0; i < row; i++ {
		q := Sum(m, 0, i)
		if q != float64(col) {
			t.Errorf("at col %d expected sum to be %f, got %f", i, float64(col), q)
		}
	}
	for i := row; i > 0; i-- {
		q := Sum(m, 0, -i)
		if q != float64(col) {
			t.Errorf("at col %d expected sum to be %f, got %f", i, float64(col), q)
		}
	}
}
func TestProd(t *testing.T) {
	row, col, val := 3, 2, 2.0
	m := New(row, col)
	m = Set(m, val)
	res := Prod(m)
	if res != 64.0 {
		t.Errorf("expected %f, got %f", 64.0, res)
	}
	row = 12
	col = 17
	m = New(row, col)
	m = Set(m, 1.0)
	for i := 0; i < col; i++ {
		q := Prod(m, 1, i)
		if q != 1.0 {
			t.Errorf("at col %d expected prod to be 1.0, got %f", i, q)
		}
	}
	for i := col; i > 0; i-- {
		q := Prod(m, 1, -i)
		if q != 1.0 {
			t.Errorf("at col %d expected prod to be 1.0, got %f", i, q)
		}
	}
	for i := 0; i < row; i++ {
		q := Prod(m, 0, i)
		if q != 1.0 {
			t.Errorf("at col %d expected Prod to be 1.0, got %f", i, q)
		}
	}
	for i := row; i > 0; i-- {
		q := Prod(m, 0, -i)
		if q != 1.0 {
			t.Errorf("at col %d expected sum to be 1.0, got %f", i, q)
		}
	}
}

func TestAvg(t *testing.T) {
	row, col, val := 7, 6, 3.0
	m := New(row, col)
	m = Set(m, val)
	a := Avg(m)
	if a != val {
		t.Errorf("expected %f, got %f", val, a)
	}
	val = 2.1
	m = Set(m, val)
	a = Avg(m, 1, 0)
	if a != val {
		t.Errorf("expected %f, got %f", val, a)
	}
	val = 1.0
	m = Set(m, val)
	a = Avg(m, 0, 1)
	if a != val {
		t.Errorf("expected %f, got %f", val, a)
	}
}

func TestDot(t *testing.T) {
	m := New(10)
	n := New(10)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*10 + j)
		}
	}
	for i := range n {
		n[i][i] = 1.0
	}
	o := Dot(m, n)
	if !Equal(o, m) {
		t.Errorf("expected equal, got not equal")
	}
}

func BenchmarkDot(b *testing.B) {
	m := New(1000)
	n := New(1000)
	q := New(1000)
	for i := range m {
		for j := range m[i] {
			m[i][j] = float64(i*10 + j)
		}
	}
	for i := range n {
		n[i][i] = 1.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q = Dot(m, n)
	}
}

func TestAppendCol(t *testing.T) {
	v := make([]float64, 10)
	m := New(10, 5)
	m = AppendCol(m, v)
	for i := range m {
		if len(m[i]) != 6 {
			t.Errorf("expected length of 6, got %d", len(m))
		}
	}
}
