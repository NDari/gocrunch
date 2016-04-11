/*
Package ann ....
*/
package ann

import (
	"fmt"
	"math"
	"math/rand"
)

/*
Net is the main type of this package. It represents a fully connected artificial
neural network.
*/
type Net struct {
	input     []float64
	hidden    [][]float64
	output    []float64
	bias      []float64
	weights   [][][]float64
	numLayers int
}

var (
	errStrings = []string{
		"\ngocrunch/ann error. \nIn %s, the number of layers must be 3 or more, but %d inputs received.\n",
		"\ngocrunch/ann error. \nIn %s, the number of inputs must 1 or more, but %d inputs received.\n",
		"\ngocrunch/ann error. \nIn %s, the number of outputs must 1 or more, but %d inputs received.\n",
		"\ngocrunch/ann error. \nIn %s, the number of nodes in hidden layer %d must 1 or more, but %d nodes requested.\n",
	}
)

/*
New is the main contructor of this package.
*/
func New(dims ...int) *Net {
	net := &Net{}
	numLayers := len(dims)
	switch numLayers {
	case 0, 1, 2:
		panic(fmt.Sprintf(errStrings[0], "New()", len(dims)))
	default:
		if dims[0] < 1 {
			panic(fmt.Sprintf(errStrings[1], "New()", dims[0]))
		}
		if dims[numLayers-1] < 1 {
			panic(fmt.Sprintf(errStrings[2], "New()", dims[0]))
		}
		inp := make([]float64, dims[0])
		out := make([]float64, dims[numLayers-1])
		var hid [][]float64
		// exclude first and last int passed to this function, as they are the
		// input and output layers.
		for i := 1; i < numLayers-1; i++ {
			if dims[i] < 1 {
				panic(fmt.Sprintf(errStrings[3], "New()", i, dims[0]))
			}
			hid = append(hid, make([]float64, dims[i]))
		}
		// one bias per hidden layer
		bias := make([]float64, numLayers-2)

		// set and initialize the weights. We use http://arxiv.org/abs/1206.5533
		// for setting the random range.
		var weights [][][]float64
		for i := 1; i < numLayers; i++ {
			val := 4.0 * math.Sqrt(6.0/float64(dims[i-1]+dims[i]))
			w := make([][]float64, dims[i-1])
			for i := range w {
				w[i] = make([]float64, dims[i])
			}
			for i := range w {
				for j := range w[i] {
					w[i][j] = rand.Float64()*(-2*val) + val
				}
			}
			weights = append(weights, w)
		}
		net = &Net{
			inp,
			hid,
			out,
			bias,
			weights,
			numLayers,
		}
	}
	return net
}
