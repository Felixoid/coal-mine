package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterNew(t *testing.T) {
	c, e := NewCounter("metric.name", 12, 15, 1, false, 30, 0, 100)
	expected := &Counter{}
	expected.base = base{
		name:          "metric.name",
		generatorType: CounterType,
		start:         12,
		stop:          15,
		step:          1,
		time:          12,
		value:         30,
		deviation:     0,
		probability:   c.probability,
	}
	expected.increment = 30
	assert.NoError(t, e)
	assert.Equal(t, expected, c)
	randomized := false
	for i := 0; i < 100; i++ {
		c, e = NewCounter("metric.name", 12, 15, 100, true, 30, 0, 100)
		if c.Time() != 12 {
			randomized = true
			break
		}
	}
	assert.NoError(t, e)
	assert.True(t, randomized)

	c, e = NewCounter("metric.name", 12, 15, 1, false, -31, 30, 100)
	assert.Nil(t, c)
	assert.Error(t, e)

	c, e = NewCounter("metric.name", 12, 15, 1, false, -30, 30, 100)
	assert.Nil(t, c)
	assert.Error(t, e)

	c, e = NewCounter("metric.name", 12, 15, 1, false, -30, 31, 100)
	assert.NoError(t, e)
}

func TestCounterNext(t *testing.T) {
	// Check error
	c := &Counter{}
	c.time = 12
	c.stop = 11
	c.step = 2
	assert.ErrorIs(t, c.Next(), ErrGenOver)

	// check normal next
	c.stop = 12
	c.value = 12
	c.increment = 2
	err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, uint(14), c.time)
	assert.Equal(t, float64(14), c.value)

	// no deviation
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.increment = 1
	for i := 0; i < 100; i++ {
		if err = c.Next(); err != nil {
			break
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.Equal(t, float64(56), c.value)

	// negative increment and unsignificant deviation
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.deviation = 0.05
	c.increment = -0.1
	for i := 0; i < 100; i++ {
		if err = c.Next(); err != nil {
			break
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.Equal(t, float64(2), c.value)

	// negative increment and significant deviation
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.deviation = 0.11
	c.increment = -0.1
	for i := 0; i < 100; i++ {
		if err = c.Next(); err != nil {
			break
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.True(t, 2 < c.value && c.value <= 7)

	// negative increment and zero deviation
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.deviation = 0
	c.increment = -0.1
	for i := 0; i < 100; i++ {
		if err = c.Next(); err != nil {
			break
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.Equal(t, float64(2), c.value)

	// Zero increment
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.increment = 0
	for i := 0; i < 100; i++ {
		if err = c.Next(); err != nil {
			break
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.Equal(t, float64(2), c.value)

	// deviation
	c.time = 12
	c.step = 1
	c.stop = 65
	c.value = 2
	c.increment = 1
	c.deviation = 4
	zero := false
	for i := 0; i < 100; i++ {
		cur := c.value
		if err = c.Next(); err != nil {
			break
		}
		if cur == c.value {
			zero = true
		}
	}
	assert.Equal(t, uint(66), c.time)
	assert.NotEqual(t, float64(56), c.value)
	assert.True(t, zero)
}
