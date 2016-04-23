# gocrunch

gocrunch is a collection of libraries for numerical libraries in Go. Take a
look at the packages below for more details about the specific libraries. If
you wish to whole-sale install all the packages of gocrunch, then do:

```bash
go get github.com/NDari/gocrunch/...
```

Each library comes pre-packaged with all of its dependencies. Therefore it is
strightforward to pick and choose which libraries you would like to use. For
example, you may with to import and use only the `mat` package. You can get
the package simply by:

```bash
go get github.com/NDari/gocrunch/mat
```

## Directories

- [gocrunch/vec](https://github.com/NDari/gocrunch/tree/master/vec): Package vec
implements functions that act upon one dimentional slices of float64s, `[]float64`.
A one dimentional slice can be thought of as a Vector.
- [gocrunch/mat](https://github.com/NDari/gocrunch/tree/master/mat): Package mat
implements functions that create or act upon two dimentional slices of float64s,
`[][]float64`. A two dimentional slice can be thought of as a Matrix.
- [gocrunch/pso](https://github.com/NDari/gocrunch/tree/master/pso): Package pso
implements the Particle Swarm Optimization method in a fast, flexible, and
extendable way.

## Badges

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
