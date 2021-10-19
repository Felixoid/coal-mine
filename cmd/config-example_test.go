package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigExample(t *testing.T) {
	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"config-example"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	body := `# carbon-server address or '-' for STDOUT, should be set as '-', 'tcp://server:port' or 'udp://server:port'
carbon = ""
# names for constant generators, braces are expanded like in shell
# values are generated with deviation around starting value
const = ["metric.const.example1", "metric.const.example{2..5}"]
# names for counter generators, braces are expanded like in shell
# values are incremented by value with deviation, but not less then the previous value
counter = ["metric.counter.example1", "metric.counter.example{2..5}"]
# names for random generators, braces are expanded like in shell
# values are generated with deviation around the previous value
random = ["metric.random.example{1,{2..5},.subdir}"]
# from in graphite-web format, the local TZ is used
from = "-2d"
# until in graphite-web format, the local TZ is used
until = "now"
# step in seconds
step = 120
# randomize starting time with [0,step)
randomize = true
# first value for all generators
value = 333.0
# deviation of the values, const will be generated around, counter will add [0,value+deviation), random will calculate next value around previous
deviation = 15.15
Probability = 0

[[custom]]
 # names for generator, braces are expanded like in shell
 name = "custom.random.generator{1..10}"
 # type of generator
 type = "random"
 # from in graphite-web format, the local TZ is used
 from = "-2d"
 # until in graphite-web format, the local TZ is used
 until = "1d"
 # step in seconds
 step = 300
 # randomize starting time with [0,step)
 randomize = false
 # first value for all generators
 value = 1000.0
 # deviation of the values, const will be generated around, counter will add [0,value+deviation), random will calculate next value around previous
 deviation = 123.456
 Probability = 0

[[custom]]
 # names for generator, braces are expanded like in shell
 name = "custom.counter.generator{1..10}"
 # type of generator
 type = "counter"
 # from in graphite-web format, the local TZ is used
 from = "-2h"
 # until in graphite-web format, the local TZ is used
 until = "now"
 # step in seconds
 step = 10
 # randomize starting time with [0,step)
 randomize = true
 # first value for all generators
 value = 100.0
 # deviation of the values, const will be generated around, counter will add [0,value+deviation), random will calculate next value around previous
 deviation = 123.456
 Probability = 0
`
	assert.Equal(t, body, buf.String())
}
