package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstNew(t *testing.T) {
	var c, _ = NewConst("metric.name", 12, 15, 1, false, 30, 0, 100)
	expected := &Const{}
	expected.base = base{
		name:          "metric.name",
		generatorType: ConstType,
		start:         12,
		stop:          15,
		step:          1,
		time:          12,
		value:         30,
		deviation:     0,
	}
	expected.constant = 30
	assert.Equal(t, expected, c)
	randomized := false
	for i := 0; i < 100; i++ {
		c, _ = NewConst("metric.name", 12, 15, 100, true, 30, 0, 100)
		if c.Time() != 12 {
			randomized = true
			break
		}
	}
	assert.True(t, randomized)
}

func TestConstNext(t *testing.T) {
	// Check error
	c := &Const{}
	c.time = 12
	c.stop = 11
	c.step = 2
	assert.ErrorIs(t, c.Next(), ErrGenOver)

	// check normal next
	c.stop = 12
	c.value = 0
	c.constant = 12
	err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, uint(14), c.time)
	assert.Equal(t, float64(12), c.value)

	randomized := false
	c.time = 11
	c.step = 1
	c.stop = 111
	c.deviation = 5
	for {
		err := c.Next()
		assert.NoError(t, err)
		if err != nil {
			break
		}
		if c.Value() != c.constant {
			randomized = true
			break
		}
	}
	assert.True(t, randomized)
}
