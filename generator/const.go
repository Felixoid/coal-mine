package generator

import (
	"math/rand"
)

// Const represents generator for constant values. When deviation is set, each value is calculated around the first value
type Const struct {
	base
	constant float64
}

// NewConst returns new generator for constant points. Possibly it can randomize values around constant.
func NewConst(name string, start, stop, step uint, randomizeStart bool, value, deviation float64) *Const {
	c := &Const{
		base: base{
			name:          name,
			generatorType: ConstType,
			start:         start,
			stop:          stop,
			step:          step,
			value:         value,
			deviation:     deviation,
		},
		constant: value,
	}
	c.RandomizeStart(randomizeStart)
	return c
}

// Next sets value and time for the next point
func (c *Const) Next() error {
	err := c.nextTime()
	if err != nil {
		return err
	}
	c.value = c.constant
	if c.Deviation() != 0 {
		c.value = c.constant + c.Deviation()*(1-rand.Float64()*2)
	}
	return nil
}
