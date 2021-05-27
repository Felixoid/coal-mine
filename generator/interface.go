package generator

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Felixoid/braxpansion"
)

// ErrGenOver shows the generation is over
var ErrGenOver = fmt.Errorf("the last point reached")

// ErrNotImplemented represents the generator is not implemented yet
var ErrNotImplemented = fmt.Errorf("generator type is not implemented")

// Generator represents
type Generator interface {
	// Next calculates next value of generator and returns ErrGenOver
	Next() error
	// Point returns the metric in carbon format, e.g. 'metric.name 123.33 1234567890\n'
	Point() []byte
	// WriteTo writes the point []byte representation to a given io.Writer
	WriteTo(io.Writer) (int64, error)
}

// Generators is a slice of Generator. Next() and Point() works accordingly
type Generators []Generator

// New returns new Generator for given parameters
func New(typeName, name string, start, stop, step uint, randomizeStart bool, value, deviation float64) (Generator, error) {
	gt, err := GetType(typeName)
	if err != nil {
		return nil, err
	}
	switch gt {
	case ConstType:
		return NewConst(name, start, stop, step, randomizeStart, value, deviation), nil
	case CounterType:
		return NewCounter(name, start, stop, step, randomizeStart, value, deviation), nil
	case RandomType:
		return NewRandom(name, start, stop, step, randomizeStart, value, deviation), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrNotImplemented, typeName)
}

// NewExpand expands name as shell expansion
// (e.g. metric.name{1..3} will produce 3 metrics metric.name1, metric.name2 and metric.name3)
// and creates slice of Generator with names.
func NewExpand(typeName, expandableName string, start, stop, step uint, randomizeStart bool, value, deviation float64) (Generators, error) {
	names := braxpansion.ExpandString(expandableName)
	gens := make(Generators, 0, len(names))
	for _, name := range names {
		g, err := New(typeName, name, start, stop, step, randomizeStart, value, deviation)
		if err != nil {
			return make(Generators, 0), err
		}
		gens = append(gens, g)
	}
	return gens, nil
}

// Next iterates over each element and calls Next. If any of calls returns an error, it breaks
func (gg Generators) Next() error {
	for _, g := range gg {
		if err := g.Next(); err != nil {
			return err
		}
	}
	return nil
}

// Point returns []byte representation of all generator Point() calls
func (gg Generators) Point() []byte {
	buf := new(bytes.Buffer)
	for _, g := range gg {
		g.WriteTo(buf)
	}
	return buf.Bytes()
}

// WriteTo writes point's []byte representation to io.Writer
func (gg Generators) WriteTo(w io.Writer) (n int64, err error) {
	var add int64
	for _, g := range gg {
		add, err = g.WriteTo(w)
		n += add
		if err != nil {
			return
		}
	}
	return n, nil
}
