package generator

import "math/rand"

// Counter represents generator for growing-up metrics
type Counter struct {
	base
	increment float64
}

// NewCounter returns new generator for growing points. Starting value is an increment as well.
// Possibly it can randomize values around increment.
// When deviation is set, it the next value won't be less then previous.
func NewCounter(name string, start, stop, step uint, randomizeStart bool, value, deviation float64) *Counter {
	c := &Counter{
		base: base{
			name:          name,
			generatorType: CounterType,
			start:         start,
			stop:          stop,
			step:          step,
			value:         value,
			deviation:     deviation,
		},
		increment: value,
	}
	c.RandomizeStart(randomizeStart)
	return c
}

// Next sets value and time for the next point
func (c *Counter) Next() error {
	if c.time > c.stop {
		return ErrGenOver
	}
	c.time += c.step
	increment := c.increment
	if c.Deviation() != 0 {
		increment = c.increment + c.Deviation()*(1-rand.Float64()*2)
	}
	if 0 < increment {
		c.value += increment
	}
	return nil
}
