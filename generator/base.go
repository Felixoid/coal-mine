package generator

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strconv"
)

// Type represents the generator type
type Type uint8

// ErrWrongType represents the bad generator type name
var ErrWrongType error = fmt.Errorf("type is not valid")

const (
	// UndefinedType is the default value
	UndefinedType Type = iota
	// ConstType represent metrics with constant values
	ConstType
	// CounterType represents metrics with growing values
	CounterType
	// RandomType represents metrics with random values
	RandomType
	endType
)

var typesMap map[string]Type = map[string]Type{
	"undefined": UndefinedType,
	"const":     ConstType,
	"counter":   CounterType,
	"random":    RandomType,
}

var types []string

func init() {
	initTypes()
}

func initTypes() {
	types = make([]string, endType)
	for k, v := range typesMap {
		if types[v] != "" {
			panic(fmt.Sprintf("check the typesMap, it is not in sync with Type const: %v", typesMap))
		}
		types[v] = k
	}
}

// GetType returns the Type by name or rises ErrWrongType error if name is invalid
func GetType(typeName string) (Type, error) {
	var t Type = 0
	var ok bool
	if t, ok = typesMap[typeName]; !ok {
		return t, fmt.Errorf("%w: %s not in %v", ErrWrongType, typeName, types)
	}
	return t, nil
}

func (b *base) nextTime() error {
	if b.time > b.stop {
		return ErrGenOver
	}
	b.time += b.step
	return nil
}

// Probability returns true if doble b.probability more than 100
func (b *base) checkProbability() bool {
	if b.probability.start == 100 {
		return true
	}
	b.probability.current = b.probability.current + b.probability.start
	if b.probability.current >= 100 {
		b.probability.current -= 100
		return true
	}
	return false
}
func probabilityIsCorrect(probabilityStart uint8) bool {
	if probabilityStart > 100 || probabilityStart < 1 {
		return false
	}
	return true
}

type base struct {
	name          string
	generatorType Type
	start         uint
	stop          uint
	step          uint
	time          uint
	value         float64
	deviation     float64
	probability   Probability
}

type Probability struct {
	start   uint8
	current uint8
}

// Randomize first current
func newProbability(probabilityStart uint8) Probability {

	return Probability{
		start:   probabilityStart,
		current: uint8(rand.Intn(100))}
}

// Point returns the metric in carbon format, e.g. 'metric.name 123.33 1234567890\n'
func (b *base) Point() []byte {
	buf := new(bytes.Buffer)
	b.WriteTo(buf)
	return buf.Bytes()
}

func (b *base) WriteTo(w io.Writer) (int64, error) {
	if !b.checkProbability() {
		return 0, nil
	}
	buf := new(bytes.Buffer)
	buf.WriteString(b.Name())
	buf.WriteString(" ")
	buf.WriteString(strconv.FormatFloat(b.Value(), 'f', -1, 64))
	buf.WriteString(" ")
	buf.WriteString(strconv.Itoa(int(b.Time())))
	buf.WriteString("\n")
	return buf.WriteTo(w)
}

// WithName sets the metric name for generator
func (b *base) WithName(name string) *base {
	b.name = name
	return b
}

// Name returns the metric name
func (b *base) Name() string {
	return b.name
}

// Type returns the generator type
func (b *base) Type() Type {
	return b.generatorType
}

// TypeName returns the name for generator type
func (b *base) TypeName() string {
	return types[b.generatorType]
}

// WithStep sets the step for the generator
func (b *base) WithStep(step uint) *base {
	b.step = step
	return b
}

// Step returns the generator step
func (b *base) Step() uint {
	return b.step
}

// SetStop sets stop to a given value
func (b *base) SetStop(stop uint) {
	b.stop = stop
}

// Stop returns value of stop field for the generator
func (b *base) Stop() uint {
	return b.stop
}

// Time returns the generator current time
func (b *base) Time() uint {
	return b.time
}

// Value returns the generator step
func (b *base) Value() float64 {
	return b.value
}

// WithDeviation deviation for the generator
func (b *base) WithDeviation(deviation float64) *base {
	b.deviation = deviation
	return b
}

// Deviation returns the generator's deviation
func (b *base) Deviation() float64 {
	return b.deviation
}

func (b *base) RandomizeStart(randomizeStart bool) {
	b.time = b.start
	if randomizeStart {
		b.time = b.start + uint(rand.Intn(int(b.step)))
	}
}
