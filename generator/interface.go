package generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

// ErrGenOver shows the generation is over
var ErrGenOver = fmt.Errorf("the last point reached")

// ErrNewCounter means that the value and the deviation are meaningless
var ErrNewCounter = fmt.Errorf("value and deviation are meaningless")

// ErrEmptyGens represents that there are no Generators
var ErrEmptyGens = fmt.Errorf("no Generators")

// ErrNotImplemented represents the generator is not implemented yet
var ErrNotImplemented = fmt.Errorf("generator type is not implemented")

// Generator represents
type Generator interface {
	// Next calculates next value of generator and returns ErrGenOver when the latest point is reached
	Next() error
	// Name return metric name
	Name() string
	// Point returns the metric in carbon format, e.g. 'metric.name 123.33 1234567890\n'
	Point() []byte
	// SetStop sets the stop field for the Generator
	SetStop(stop uint)
	// Stop returns the current value for the stop field
	Stop() uint
	// WriteTo writes the point []byte representation to a given io.Writer
	WriteTo(io.Writer) (int64, error)
	// Append append the point []byte representation to a given []byte slice
	Append([]byte) []byte
}

// Generators is a slice of Generator. Next() and Point() works accordingly
type Generators struct {
	name       string
	typeName   string
	step       uint
	randomized bool
	gens       []Generator
}

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
		return NewCounter(name, start, stop, step, randomizeStart, value, deviation)
	case RandomType:
		return NewRandom(name, start, stop, step, randomizeStart, value, deviation), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrNotImplemented, typeName)
}

// NewExpand expands name as shell expansion
// (e.g. metric.name{1..3} will produce 3 metrics metric.name1, metric.name2 and metric.name3)
// and creates slice of Generator with names.
func NewExpand(typeName, expandableName string, variables map[string]string,
	start, stop, step uint, randomizeStart bool, value, deviation float64) (Generators, error) {

	names := expandString(expandableName, variables)
	if len(names) == 0 {
		return Generators{}, ErrEmptyGens
	}
	gg := Generators{
		name:       expandableName,
		typeName:   typeName,
		step:       step,
		randomized: randomizeStart,
		gens:       make([]Generator, len(names)),
	}
	for i, name := range names {
		g, err := New(typeName, name, start, stop, step, randomizeStart, value, deviation)
		if err != nil {
			return Generators{}, err
		}
		gg.gens[i] = g
	}
	return gg, nil
}

// List returns the list of []Generator
func (gg *Generators) List() []Generator {
	return gg.gens
}

// Next iterates over each element and calls Next. If any of calls returns an error, it breaks
func (gg *Generators) Next() error {
	if len(gg.gens) == 0 {
		return ErrEmptyGens
	}
	for _, g := range gg.gens {
		if err := g.Next(); err != nil {
			return err
		}
	}
	return nil
}

// Point returns []byte representation of all generator Point() calls
func (gg *Generators) Point() []byte {
	buf := new(bytes.Buffer)
	for _, g := range gg.gens {
		g.WriteTo(buf)
	}
	return buf.Bytes()
}

// Randomized shows if expanded generators have randomized start timestamp
func (gg *Generators) Randomized() bool {
	return gg.randomized
}

// SetList replaces existing []Generator with a given one
func (gg *Generators) SetList(gens []Generator) {
	gg.gens = gens
}

// SetStop invokes the according method for each Generagor
func (gg *Generators) SetStop(stop uint) {
	for _, g := range gg.gens {
		g.SetStop(stop)
	}
}

// Step returns the common step for Generators
func (gg *Generators) Step() uint {
	return gg.step
}

// Append append point's []byte representation to []byte slice
func (gg *Generators) Append(buf []byte) []byte {
	for _, g := range gg.gens {
		buf = g.Append(buf)
	}
	return buf
}

// WriteTo writes point's []byte representation to io.Writer
func (gg *Generators) WriteTo(w io.Writer, buf *[]byte) (n int64, err error) {
	var add int
	for _, g := range gg.gens {
		*buf = g.Append((*buf)[:0])
		add, err = w.Write(*buf)
		n += int64(add)
		if err != nil {
			return
		}
	}
	return n, nil
}

// WriteAllTo writes all points for Generators to io.Writer
func (gg *Generators) WriteAllTo(w io.Writer, buf *[]byte) (int64, error) {
	var n int64
	var err error
	for ; err == nil; err = gg.Next() {
		add, err := gg.WriteTo(w, buf)
		n += add
		if err != nil {
			return n, err
		}
	}
	if !errors.Is(err, ErrGenOver) {
		return n, err
	}
	return n, nil
}

// WriteAllToWithContext writes all points, but may be stopped by the passing a struct to a stop channel
func (gg *Generators) WriteAllToWithContext(ctx context.Context, w io.Writer, buf *[]byte) (int64, error) {
	var (
		n int64
	)
	wr := func() error {
		add, err := gg.WriteTo(w, buf)
		n += add
		return err
	}
	var err error
	for ; err == nil; err = gg.Next() {
		select {
		case <-ctx.Done():
			return n, ctx.Err()
		default:
		}
		if err := wr(); err != nil {
			return n, err
		}
	}
	if !errors.Is(err, ErrGenOver) {
		return n, err
	}
	return n, nil
}
