package generator

import (
	"fmt"
	"math"
	"math/rand"
)

// Counter represents generator for growing-up metrics
type Counter struct {
	base
	increment float64
}

// NewCounter returns new generator for growing points. Starting value is an increment as well.
// Possibly it can randomize values around increment.
// When deviation is set, it the next value won't be less then previous.
func NewCounter(name string, start, stop, step uint, randomizeStart bool, value, deviation float64, probabilityStart uint8) (*Counter, error) {
	if value < 0 && math.Abs(value) <= math.Abs(deviation) {
		return nil, fmt.Errorf("%w: with negative value deviation must be greater than value", ErrNewCounter)
	}
	if !probabilityIsCorrect(probabilityStart) {
		return nil, ErrProbabilityStart
	}
	c := &Counter{
		base: base{
			name:          name,
			generatorType: CounterType,
			start:         start,
			stop:          stop,
			step:          step,
			value:         value,
			deviation:     deviation,
			probability: Probability{
				start:   probabilityStart,
				current: NewProbability()},
		},
		increment: value,
	}
	c.RandomizeStart(randomizeStart)
	return c, nil
}

// Next sets value and time for the next point
func (c *Counter) Next() error {
	err := c.nextTime()
	if err != nil {
		return err
	}
	increment := c.increment
	if c.Deviation() != 0 {
		increment = c.increment + c.Deviation()*(1-rand.Float64()*2)
	}
	if 0 < increment {
		c.value += increment
	}
	return nil
}
