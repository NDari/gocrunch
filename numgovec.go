package numgo

import (
	"numgo/vec"
	"testing"
)

func TestVecOnes(t *testing.T) {
	o := vec.Ones(13)
	for i := 0; i < 13; i++ {
		if o[i] != 1.0 {
			t.Errorf("vec.Ones[%d]: expected 1.0, got %f", i, o[i])
		}
	}
}

func TestVecInc(t *testing.T) {
	v := vec.Inc(10)
	for i := 0; i < 10; i++ {
		if v[i] != float64(i) {
			t.Errorf("vec.Inc at index %d: expected %f got %f", i, float64(i), v[i])
		}
	}
}

func TestVecEqual(t *testing.T) {
	v := vec.Inc(20)
	if !vec.Equal(v, v) {
		t.Errorf("vec.Equal: v != v")
	}
	q := vec.Inc(10)
	if vec.Equal(v, q) {
		t.Errorf("expected false, but got true for %v == %v", v, q)
	}
	if vec.Equal(q, v) {
		t.Errorf("expected false, but got true for %v == %v", q, v)
	}
	g := vec.Ones(10)
	if vec.Equal(g, q) {
		t.Errorf("expected false, but got true for %v == %v", g, q)
	}
	if vec.Equal(q, g) {
		t.Errorf("expected false, but got true for %v == %v", q, g)
	}
}

func TestVecMul(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Mul(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != w[i] {
			t.Errorf("vec.Mul at index %d: expected 0.0, got %f", i, y[i])
		}
	}
	w = vec.Ones(12)
	g := vec.Mul(v, w)
	for i := 0; i < 12; i++ {
		if g[i] != (v[i] * w[i]) {
			t.Errorf("vec.Mul at index %d: expected %f, got %f", i, v[i]*w[i], g[i])
		}
	}
	o := vec.Mul(w, w)
	for i := 0; i < 12; i++ {
		if o[i] != 1.0 {
			t.Errorf("vec.Mul at index %d: expected 1.0, got %f", i, o[i])
		}
	}
}

func TestVecAdd(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Add(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("vec.Mul at index %d: expected %f, got %f", i, v[i], y[i])
		}
	}
	y = vec.Add(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != (v[i] * 2.0) {
			t.Errorf("vec.Mul at index %d: expected %f, got %f", i, v[i]*2.0, y[i])
		}
	}
}

func TestVecSub(t *testing.T) {
	v := vec.Ones(12)
	w := make([]float64, 12)
	y := vec.Sub(v, w)
	for i := 0; i < 12; i++ {
		if y[i] != v[i] {
			t.Errorf("vec.Mul at index %d: expected %f, got %f", i, v[i], y[i])
		}
	}
	y = vec.Sub(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 0.0 {
			t.Errorf("vec.Mul at index %d: expected 0.0, got %f", i, y[i])
		}
	}
}

func TestVecDiv(t *testing.T) {
	v := vec.Ones(12)
	y := vec.Div(v, v)
	for i := 0; i < 12; i++ {
		if y[i] != 1.0 {
			t.Errorf("vec.Mul at index %d: expected 1.0, got %f", i, y[i])
		}
	}
	w := vec.Inc(12)
	y = vec.Div(w, v)
	for i := 0; i < 12; i++ {
		if y[i] != (w[i]) {
			t.Errorf("vec.Mul at index %d: expected %f, got %f", i, w[i], y[i])
		}
	}
}

func TestVecMap(t *testing.T) {
	v := make([]float64, 17)
	y := vec.Map(func(i float64) float64 { return 1.0 }, v)
	for i := 0; i < 17; i++ {
		if y[i] != 1.0 {
			t.Errorf("vec.Map at index %d: xpected 1.0, got %f", i, y[i])
		}
	}
}

func TestVecDot(t *testing.T) {
	v := vec.Ones(13)
	w := make([]float64, 13)
	r := vec.Dot(v, w)
	if r != 0.0 {
		t.Errorf("vec.Dot: expected 0.0, got %v", r)
	}
	q := vec.Inc(13)
	if vec.Dot(q, v) != vec.Sum(q) {
		t.Errorf("vec.Dot: Inc * ones is %f, expected %f", vec.Dot(q, v), vec.Sum(q))
	}
	if vec.Dot(q, v) != vec.Dot(v, q) {
		t.Errorf("Vec.Dot: v dot q != q dot v")
	}
}

func TestVecReset(t *testing.T) {
	v := vec.Ones(22)
	vec.Reset(v)
	for i := 0; i < 22; i++ {
		if v[i] != 0.0 {
			t.Errorf("vec.Reset at index %d: expected 0.0, got %f", i, v[i])
		}
	}
}

func TestVecSum(t *testing.T) {
	v := vec.Ones(22)
	if vec.Sum(v) != 22.0 {
		t.Errorf("vec.Sum expected 22.0, got %f", vec.Sum(v))
	}
	vec.Reset(v)
	if vec.Sum(v) != 0.0 {
		t.Errorf("vec.Sum expected 0.0, got %f", vec.Sum(v))
	}
}

func TestVecNorm(t *testing.T) {
	m := vec.Ones(4)
	if vec.Norm(m) != 2.0 {
		t.Errorf("vec.Norm, expected 2.0, got %f", vec.Norm(m))
	}
}

func TestvecAll(t *testing.T) {
	f := func(i float64) bool {
		if i > 0.0 {
			return true
		}
		return false
	}
	m := vec.Ones(10)
	v := vec.All(f, m)
	if v != true {
		t.Errorf("vec.All, not all entries of vec.Ones are > 0.0")
	}
}

func TestvecAny(t *testing.T) {
	f := func(i float64) bool {
		if i < 0.0 {
			return true
		}
		return false
	}
	m := vec.Ones(10)
	v := vec.All(f, m)
	if v != false {
		t.Errorf("vec.All, atleast one entry of vec.Ones is < 0.0")
	}
}
