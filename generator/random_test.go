package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomNew(t *testing.T) {
	c, _ := NewRandom("metric.name", 12, 15, 1, false, 30, 0, 100)
	expected := &Random{}
	expected.base = base{
		name:          "metric.name",
		generatorType: RandomType,
		start:         12,
		stop:          15,
		step:          1,
		time:          12,
		value:         30,
		deviation:     0,
		probability:   c.probability,
	}
	assert.Equal(t, expected, c)
	randomized := false
	for i := 0; i < 100; i++ {
		c, _ = NewRandom("metric.name", 12, 15, 100, true, 30, 0, 100)
		if c.Time() != 12 {
			randomized = true
			break
		}
	}
	assert.True(t, randomized)
}

func TestRandomNext(t *testing.T) {
	// Check error
	r := &Random{}
	r.time = 12
	r.stop = 11
	r.step = 2
	assert.ErrorIs(t, r.Next(), ErrGenOver)

	// check normal next
	r.stop = 12
	r.value = 0
	err := r.Next()
	assert.NoError(t, err)
	assert.Equal(t, uint(14), r.time)
	assert.Equal(t, float64(0), r.value)

	// check normal next
	r.stop = 21
	r.value = 0
	r.deviation = 4
	r.Next()
	assert.NoError(t, err)
	assert.Equal(t, uint(16), r.time)
	assert.NotEqual(t, float64(0), r.value)
}
