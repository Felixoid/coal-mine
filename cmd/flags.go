package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type stringMap map[string]string

func (u stringMap) Set(value string) error {
	k, v, ok := strings.Cut(value, "=")
	if !ok {
		return errors.New("invalid variable '" + value + "', must be with = delimiter")
	}
	if _, ok := u[v]; ok {
		return errors.New("duplicate variable '" + k + "'")
	}
	u[k] = v
	return nil
}

func (u stringMap) String() string {
	if len(u) == 0 {
		return "[]"
	}
	var buf strings.Builder
	first := true
	buf.WriteByte('[')
	for k, v := range u {
		buf.WriteByte('\'')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte('\'')
		if first {
			first = false
		} else {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

func (u stringMap) Type() string {
	return "map[string]string"
}

var variables = make(stringMap)

func commonFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	metricsFlags(cmd)
	f.String("carbon", viper.GetString("carbon"), "carbon-server address or '-' for STDOUT, should be set as '-', 'tcp://server:port' or 'udp://server:port'")
	f.Float64("value", viper.GetFloat64("value"), "starting value for generators")
	f.Float64("deviation", viper.GetFloat64("deviation"), "deviation for the next point in generator")
	f.Uint("step", viper.GetUint("step"), "generators interval in seconds")
}

func bindCommonFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	viper.BindPFlag("carbon", f.Lookup("carbon"))
	bindMetricsFlags(cmd)
	viper.BindPFlag("value", f.Lookup("value"))
	viper.BindPFlag("deviation", f.Lookup("deviation"))
	viper.BindPFlag("step", f.Lookup("step"))
}

func metricsFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&cfgFile, "config", "c", "", "config file")
	f.StringArray("const", []string{}, "constant generators")
	f.StringArray("counter", []string{}, "counter generators")
	f.StringArray("random", []string{}, "random generators")
	f.Bool("randomize", viper.GetBool("randomize"), "toggle if starting point of generators should be randomized")
	f.Var(variables, "var", "variables for expand")
}

func bindMetricsFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	viper.BindPFlag("const", f.Lookup("const"))
	viper.BindPFlag("counter", f.Lookup("counter"))
	viper.BindPFlag("random", f.Lookup("random"))
	viper.BindPFlag("randomize", f.Lookup("randomize"))
	viper.BindPFlag("var", f.Lookup("var"))
}
