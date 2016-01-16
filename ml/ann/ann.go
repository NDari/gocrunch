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
	input  *mat.Mat
	hidden []*mat.Mat
	output *mat.Mat
	bias   *mat.Mat
}

/*
New is the main contructor of this package.
*/
func New(dims ...int) *Net {
	net := &Net{}
	switch len(dims) {
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
		if dims[len(dims)-1] < 1 {
			fmt.Println("\nNumgo/ann error.")
			s := "In ann.%s, the number of nodes of the output layer must be\n"
			s += "set to one or more. However, %d nodes were requested.\n"
			s = fmt.Sprintf(s, "New", dims[len(dims)-1])
			fmt.Println(s)
			fmt.Println("Stack trace for this error:")
			debug.PrintStack()
			os.Exit(1)
		}
		inp := mat.New(1, dims[0])
		out := mat.New(1, dims[len(dims)-1])
		var hid []*mat.Mat
		// exclude first and last int passed to this function, as they are the
		// input and output layers.
		for i := 1; i < len(dims)-1; i++ {
			if dims[i] < 1 {
				fmt.Println("\nNumgo/ann error.")
				s := "In ann.%s, the number of nodes of hidden layer %d\n"
				s += "layer must be set to one or more.\n"
				s += "However, %d nodes were requested for this layer.\n"
				s = fmt.Sprintf(s, "New", i, dims[len(dims)-1])
				fmt.Println(s)
				fmt.Println("Stack trace for this error:")
				debug.PrintStack()
				os.Exit(1)
			}
			hid = append(hid, mat.New(1, dims[i]))
		}
		// one bias per hidden layer
		bias := mat.New(1, len(dims)-2)

		// Initialize the weights. We use http://arxiv.org/abs/1206.5533
		// for setting the random range.
		for i := 0; i < len(dims)-2; i++ {
			if i == 0 {
				val := 4.0 * math.Sqrt(6.0/float64(dims[0]+dims[2]))
				hid[0].Rand(-val, val)
				continue
			}
			val := 4.0 * math.Sqrt(6.0/float64(dims[i-1]+dims[i+1]))
			hid[i].Rand(-val, val)
		}
		net = &Net{
			inp,
			hid,
			out,
			bias,
		}
	}
	return net
}

/*
Print writes the weights of the hidden layers to a file.
TODO: Change this to actually print to a file, not stdout.
*/
func (n *Net) Print() {
	for i := 0; i < len(n.hidden); i++ {
		fmt.Printf("Hidden layer %d weights:\n", i)
		n.hidden[i].Print()
	}
	fmt.Println("The bias weights:")
	n.bias.Print()
}
