package vec

import (
	"testing"
)

func TestOnes(t *testing.T) {
	o := Ones(13)
	for i := 0; i < 13; i++ {
		if o[i] != 1.0 {
			t.Errorf("Ones[%d]: expected 1.0, got %f", i, o[i])
		}
	}
}

func TestInc(t *testing.T) {
	v := Inc(10)
	for i := 0; i < 10; i++ {
		if v[i] != float64(i) {
			t.Errorf("Inc at index %d: expected %f got %f", i, float64(i), v[i])
		}
	}
}

func TestEqual(t *testing.T) {
	v := Inc(20)
	if !Equal(v, v) {
		t.Errorf("Equal: v != v")
	}
	q := Inc(10)
	if Equal(v, q) {
		t.Errorf("expected false, but got true for %v == %v", v, q)
	}
	if Equal(q, v) {
		t.Errorf("expected false, but got true for %v == %v", q, v)
	}
	g := Ones(10)
	if Equal(g, q) {
		t.Errorf("expected false, but got true for %v == %v", g, q)
	}
	if Equal(q, g) {
		t.Errorf("expected false, but got true for %v == %v", q, g)
	}
}

func TestMul(t *testing.T) {
	v := Ones(12)
	w := make([]float64, 12)
	y := Mul(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != w[i] {
			t.Errorf("Mul at index %d: expected 0.0, got %f", i, y[i])
		}
	}
	w = Ones(12)
	g := Mul(v, w)
	for i := 0; i < 12; i++ {
		if g[i] != (v[i] * w[i]) {
			t.Errorf("Mul at index %d: expected %f, got %f", i, v[i]*w[i], g[i])
		}
	}
	o := Mul(w, w)
	for i := 0; i < 12; i++ {
		if o[i] != 1.0 {
			t.Errorf("Mul at index %d: expected 1.0, got %f", i, o[i])
		}
	}
}

func TestAdd(t *testing.T) {
	v := Ones(12)
	w := make([]float64, 12)
	y := Add(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("Mul at index %d: expected %f, got %f", i, v[i], y[i])
		}
	}
	y = Add(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != (v[i] * 2.0) {
			t.Errorf("Mul at index %d: expected %f, got %f", i, v[i]*2.0, y[i])
		}
	}
}

func TestSub(t *testing.T) {
	v := Ones(12)
	w := make([]float64, 12)
	y := Sub(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("Mul at index %d: expected %f, got %f", i, v[i], y[i])
		}
	}
	y = Sub(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 0.0 {
			t.Errorf("Mul at index %d: expected 0.0, got %f", i, y[i])
		}
	}
}

func TestDiv(t *testing.T) {
	v := Ones(12)
	y := Div(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 1.0 {
			t.Errorf("Mul at index %d: expected 1.0, got %f", i, y[i])
		}
	}
	w := Inc(12)
	y = Div(w, v)
	for i := 0; i < 12; i++ {
		if y[i] != (w[i]) {
			t.Errorf("Mul at index %d: expected %f, got %f", i, w[i], y[i])
		}
	}
}

func TestMap(t *testing.T) {
	v := make([]float64, 17)
	y := Map(func(i float64) float64 { return 1.0 }, v)
	for i := 0; i < 17; i++ {
		if y[i] != 1.0 {
			t.Errorf("Map at index %d: xpected 1.0, got %f", i, y[i])
		}
	}
}

func TestDot(t *testing.T) {
	v := Ones(13)
	w := make([]float64, 13)
	r := Dot(v, w)
	if r != 0.0 {
		t.Errorf("Dot: expected 0.0, got %v", r)
	}
	q := Inc(13)
	if Dot(q, v) != Sum(q) {
		t.Errorf("Dot: Inc * ones is %f, expected %f", Dot(q, v), Sum(q))
	}
	if Dot(q, v) != Dot(v, q) {
		t.Errorf("Dot: v dot q != q dot v")
	}
}

func TestReset(t *testing.T) {
	v := Ones(22)
	Reset(v)
	for i := 0; i < 22; i++ {
		if v[i] != 0.0 {
			t.Errorf("Reset at index %d: expected 0.0, got %f", i, v[i])
		}
	}
}

func TestSum(t *testing.T) {
	v := Ones(22)
	if Sum(v) != 22.0 {
		t.Errorf("Sum expected 22.0, got %f", Sum(v))
	}
	Reset(v)
	if Sum(v) != 0.0 {
		t.Errorf("Sum expected 0.0, got %f", Sum(v))
	}
}

func TestNorm(t *testing.T) {
	m := Ones(4)
	if Norm(m) != 2.0 {
		t.Errorf("Norm, expected 2.0, got %f", Norm(m))
	}
}

func TestAll(t *testing.T) {
	f := func(i float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}
	m := Ones(10)
	v := All(f, m)
	if v != true {
		t.Errorf("All, not all entries of Ones are > 0.0")
	}
}

func TestAny(t *testing.T) {
	f := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	m := Ones(10)
	v := All(f, m)
	if v != false {
		t.Errorf("All, atleast one entry of Ones is < 0.0")
	}
}
