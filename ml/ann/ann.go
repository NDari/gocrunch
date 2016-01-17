/*
Package ann ....
*/
package ann

import (
	"fmt"
	"math"
	"os"
	"runtime/debug"

	"github.com/NDari/numgo/mat"
)

/*
Net is the main type of this package. It represents a fully connected artificial
neural network.
*/
type Net struct {
	input     *mat.Mat
	hidden    []*mat.Mat
	output    *mat.Mat
	bias      *mat.Mat
	weights   []*mat.Mat
	numLayers int
}

/*
New is the main contructor of this package.
*/
func New(dims ...int) *Net {
	net := &Net{}
	numLayers := len(dims)
	switch numLayers {
	case 0, 1, 2:
		fmt.Println("\nNumgo/ann error.")
		s := "In ann.%s, the number of inputs must be at least 3, but\n"
		s += "recieved %d. A New network can be contructed from 1 input\n"
		s += "layer, 1 (or more) hidden layer(s), and one output layer.\n"
		s = fmt.Sprintf(s, "New", len(dims))
		fmt.Println(s)
		fmt.Println("Stack trace for this error:")
		debug.PrintStack()
		os.Exit(1)
	default:
		if dims[0] < 1 {
			fmt.Println("\nNumgo/ann error.")
			s := "In ann.%s, the number of nodes of the input layer must be\n"
			s += "set to one or more. However, %d nodes were requested.\n"
			s = fmt.Sprintf(s, "New", dims[0])
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		if dims[numLayers-1] < 1 {
			fmt.Println("\nNumgo/ann error.")
			s := "In ann.%s, the number of nodes of the output layer must be\n"
			s += "set to one or more. However, %d nodes were requested.\n"
			s = fmt.Sprintf(s, "New", dims[numLayers-1])
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		inp := mat.New(1, dims[0])
		out := mat.New(1, dims[numLayers-1])
		var hid []*mat.Mat
		// exclude first and last int passed to this function, as they are the
		// input and output layers.
		for i := 1; i < numLayers-1; i++ {
			if dims[i] < 1 {
				fmt.Println("\nNumgo/ann error.")
				s := "In ann.%s, the number of nodes of hidden layer %d\n"
				s += "layer must be set to one or more.\n"
				s += "However, %d nodes were requested for this layer.\n"
				s = fmt.Sprintf(s, "New", i, dims[numLayers-1])
				fmt.Println(s)
				fmt.Println("Stack trace for this error:")
				debug.PrintStack()
				os.Exit(1)
			}
			hid = append(hid, mat.New(1, dims[i]))
		}
		// one bias per hidden layer
		bias := mat.New(1, numLayers-2)

		// set and nitialize the weights. We use http://arxiv.org/abs/1206.5533
		// for setting the random range.
		var weights []*mat.Mat
		for i := 1; i < numLayers; i++ {
			val := 4.0 * math.Sqrt(6.0/float64(dims[i-1]+dims[i]))
			w := mat.New(dims[i-1], dims[i]).Rand(-val, val)
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

/*
Print writes the weights of the hidden layers to a file.
TODO: Change this to actually print to a file, not stdout.
*/
func (n *Net) Print() {
	var str string
	for i := 0; i < len(n.weights); i++ {
		s := "Hidden layer %d weights:\n"
		s = fmt.Sprintf(s, i)
		str += s
		str += n.weights[i].ToString()
	}
	str += "The bias weights:\n"
	str += n.bias.ToString()
	fmt.Println(str)
}
