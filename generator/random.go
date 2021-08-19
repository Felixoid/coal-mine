package generator

import "math/rand"

// Random works like Const, but each next value is calculated like value±deviation
type Random struct {
	base
}

// NewRandom returns new generator for growing points. Without deviation it behaves like constant.
func NewRandom(name string, start, stop, step uint, randomizeStart bool, value, deviation float64) *Random {
	c := &Random{
		base: base{
			name:          name,
			generatorType: RandomType,
			start:         start,
			stop:          stop,
			step:          step,
			value:         value,
			deviation:     deviation,
		},
	}
	c.RandomizeStart(randomizeStart)
	return c
}

// Next sets value and time for the next point
func (r *Random) Next() error {
	if r.getNextCommon() != nil {
		return ErrGenOver
	}
	if r.Deviation() != 0 {
		r.value += r.Deviation() * (1 - rand.Float64()*2)
	}
	return nil
}
