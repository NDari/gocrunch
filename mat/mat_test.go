package mat

import (
	"log"
	"os"
	"testing"
)

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
	if len(m.vals) != rows*cols {
		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
	}
}

// func TestBare(t *testing.T) {
// 	rows := 13
// 	cols := 7
// 	m := Bare(rows, cols)
// 	if m.r != rows {
// 		t.Errorf("New mat.r is %d, expected %d", m.r, rows)
// 	}
// 	if m.c != cols {
// 		t.Errorf("New mat.c is %d, expected %d", m.c, cols)
// 	}
// 	if m.vals == nil {
// 		t.Errorf("Bare mat.vals not initilized")
// 	}
// 	if m.work != nil {
// 		t.Errorf("Bare mat.work is initilized")
// 	}
// 	if len(m.vals) != rows*cols {
// 		t.Errorf("len(mat.vals) is %d, expected %d", len(m.vals), rows*cols)
// 	}
// 	if len(m.work) != 0 {
// 		t.Errorf("len(mat.work) is %d, expected 0", len(m.work))
// 	}
// }

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

func TestFromSlice(t *testing.T) {
	s := make([][]float64, 11)
	for i := range s {
		s[i] = make([]float64, 5)
	}
	for i := range s {
		for j := range s[i] {
			s[i][j] = float64(i + j)
		}
	}
	m := FromSlice(s)
	idx := 0
	for i := range s {
		for j := range s[i] {
			if s[i][j] != m.vals[idx] {
				t.Errorf("slice[%d][%d]: %f, mat: %f", i, j,
					s[i][j], m.vals[idx])
			}
			idx += 1
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
			idx += 1
		}
	}
}

func TestMap(t *testing.T) {
	rows := 132
	cols := 24
	f := func(float64) float64 {
		return 1.0
	}
	m := New(rows, cols)
	m.Map(f)
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("At %d, expected 1.0, got %d")
		}
	}
}

func TestFromCSV(t *testing.T) {
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

func TestOnes(t *testing.T) {
	rows := 13
	cols := 7
	m := New(rows, cols).Ones()
	for i := 0; i < rows*cols; i++ {
		if m.vals[i] != 1.0 {
			t.Errorf("At %d, expected 1.0, got %d", i, m.vals[i])
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

func TestCol(t *testing.T) {
	row := 3
	col := 4
	m := New(row, col).Inc()
	for i := 0; i < col; i++ {
		got := m.Col(i)
		for j := 0; j < row; j++ {
			if got.vals[j] != m.vals[j*m.r+i] {
				t.Errorf("At index %v Col(%v), got %f, want %f", j, i,
					got.vals[j], m.vals[j*m.r+i])
			}
		}
	}
}
