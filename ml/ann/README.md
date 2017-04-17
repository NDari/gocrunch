# ann
--
    import "github.Ndari/numgo/ml/ann"

Package ann ....

## Usage

#### type Net

```go
type Net struct {
}
```

Net is the main type of this package. It represents a fully connected artificial
neural network.

#### func  New

```go
func New(dims ...int) *Net
```
New is the main contructor of this package.

#### func (*Net) Print

```go
func (n *Net) Print()
```
Print writes the weights of the hidden layers to a file. TODO: Change this to
actually print to a file, not stdout.
