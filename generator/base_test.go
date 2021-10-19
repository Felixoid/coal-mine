package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	assert.Equal(t, []string{"undefined", "const", "counter", "random"}, types)

	// Check logic for predefined types
	backupMap := map[string]Type{}
	for k, v := range typesMap {
		backupMap[k] = v
	}
	typesMap["newTestType"] = UndefinedType
	assert.Panics(t, initTypes)
	typesMap["newTestType"] = endType
	assert.Panics(t, initTypes)

	// restore types and typesMap
	typesMap = backupMap
	initTypes()
}

func TestBaseWith(t *testing.T) {
	b := new(base).WithName("metric.name").WithStep(12).WithDeviation(123.12)
	assert.Equal(t, &base{name: "metric.name", step: 12, deviation: 123.12}, b)
	assert.Equal(t, "metric.name", b.Name())
	assert.Equal(t, uint(12), b.Step())
	assert.Equal(t, 123.12, b.Deviation())
}

func TestBaseGetters(t *testing.T) {
	b := base{time: 1234432112, value: 333.333}
	assert.Equal(t, uint(1234432112), b.Time())
	assert.Equal(t, 333.333, b.Value())
}

func TestBaseType(t *testing.T) {
	b := new(base)
	for _, gt := range types {
		genType, err := GetType(gt)
		assert.NoError(t, err)
		b.generatorType = genType
		assert.Equal(t, genType, b.Type())
		assert.Equal(t, gt, b.TypeName())
	}
}

func TestGetType(t *testing.T) {
	for i, gt := range types {
		genType, err := GetType(gt)
		assert.NoError(t, err)
		assert.Equal(t, genType, typesMap[gt])
		assert.Equal(t, Type(i), typesMap[gt])
	}
	genType, err := GetType("incorrect")
	assert.ErrorIs(t, err, ErrWrongType)
	assert.Equal(t, Type(0), genType)
}

func TestBaseRandomizeStart(t *testing.T) {
	success := false
	b := base{start: 60, step: 30}
	for i := 0; i < 100; i++ {
		b.RandomizeStart(true)
		if b.time != b.start {
			success = true
			break
		}
	}
	assert.True(t, success)
}

func TestBasePoint(t *testing.T) {
	b := base{}
	b.name = "metric.name"
	b.value = 123.333
	b.time = 13
	b.probability.current = 100
	assert.Equal(t, []byte("metric.name 123.333 13\n"), b.Point())
}

func TestBaseStop(t *testing.T) {
	b := base{}
	for _, stop := range []uint{1, 2, 34, 5, 6, 7} {
		b.SetStop(stop)
		assert.Equal(t, stop, b.stop)
	}
}

func TestBaseCheckProbability(t *testing.T) {
	b := base{}
	for _, iter := range []uint8{10, 20, 4, 5, 36, 49, 99} {
		b.probability.current = 0
		b.probability.start = iter
		assert.False(t, b.checkProbability())
	}
	for _, iter := range []uint8{99, 50, 56, 89, 65, 67} {
		b.probability.current = iter
		b.probability.start = iter
		assert.True(t, b.checkProbability())
		assert.Equal(t, b.probability.current, 2*iter-100)
	}
}
