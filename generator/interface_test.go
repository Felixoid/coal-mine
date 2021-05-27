package generator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type bufWithLimit struct {
	buf    []byte
	remain int
}

func newBufWithLimit(l int) *bufWithLimit {
	return &bufWithLimit{make([]byte, 0, l), l}
}

func (b *bufWithLimit) Write(p []byte) (int, error) {
	l := len(p)
	if b.remain < l {
		b.buf = append(b.buf, p[:b.remain]...)
		remain := b.remain
		b.remain = 0
		return remain, fmt.Errorf("the buffer size excited")
	}
	b.buf = append(b.buf, p...)
	b.remain -= l
	return l, nil
}

func TestNewGenerator(t *testing.T) {
	g, err := New("invalid", "", 0, 0, 0, false, 0, 0)
	assert.Nil(t, g)
	assert.ErrorIs(t, err, ErrWrongType)
	g, err = New("const", "", 0, 0, 0, false, 0, 0)
	assert.IsType(t, &Const{}, g)
	assert.NoError(t, err)
	g, err = New("counter", "", 0, 0, 0, false, 0, 0)
	assert.IsType(t, &Counter{}, g)
	assert.NoError(t, err)
	g, err = New("random", "", 0, 0, 0, false, 0, 0)
	assert.IsType(t, &Random{}, g)
	assert.NoError(t, err)
}

func TestNewExpand(t *testing.T) {
	gg, err := NewExpand("const", "{metric,name}{1..2}{a..b}", 0, 0, 0, false, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, gg, 8)

	gg, err = NewExpand("const", "", 0, 0, 0, false, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, gg, 0)

	gg, err = NewExpand("const", "metric.name", 0, 0, 0, false, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, gg, 1)

	gg, err = NewExpand("invalid", "metric.name", 0, 0, 0, false, 0, 0)
	assert.Error(t, err)
	assert.Len(t, gg, 0)
}

func TestGeneratorsNext(t *testing.T) {
	gg, err := NewExpand("const", "name{1..3}", 1, 2, 1, false, 0, 0)
	assert.NoError(t, err)
	c, ok := gg[1].(*Const)
	assert.True(t, ok)
	c.time = 3
	// this Next success for [0], fails for [1] and skips [2] element
	assert.ErrorIs(t, gg.Next(), ErrGenOver)
	out := []byte(`name1 0 2
name2 0 3
name3 0 1
`)
	assert.Equal(t, out, gg.Point())
}

func TestGeneratorsWriteTo(t *testing.T) {
	buf := newBufWithLimit(200)
	gg, err := NewExpand("const", "root.level{01..05}.node{01..10}", 0, 2, 1, false, 0, 0)
	assert.NoError(t, err)
	n, err := gg.WriteTo(buf)
	assert.Error(t, err)
	assert.Equal(t, int64(200), n)
}
