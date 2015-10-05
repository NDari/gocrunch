package mat64

// Future2DSlice is a channel which is used in async operations
// internally.
type future2DSlice chan [][]float64

// TAsync Runs T() in a goroutine, returning a channel which will
// contain the result when the goroutine is done.
func TAsync(m [][]float64) future2DSlice {
	c := make(future2DSlice)
	go func() { c <- T(m) }()
	return c
}

// TimesAsync runs Times() in a gorroutine, returning a channel
// which will contain the result when the goroutine is done.
func TimesAsync(m, n [][]float64) future2DSlice {
	c := make(future2DSlice)
	go func() { c <- Times(m, n) }()
	return c
}

// DotAsync will apply Dot() in a goroutine, returning a channel that
// with contain the result when the goroutine is done.
func DotAsync(m, n [][]float64) future2DSlice {
	c := make(future2DSlice)
	go func() { c <- Dot(m, n) }()
	return c
}
