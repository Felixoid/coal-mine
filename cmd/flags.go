package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func commonFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.StringVarP(&cfgFile, "config", "c", "", "config file")
	f.String("carbon", viper.GetString("carbon"), "carbon-server address or '-' for STDOUT, should be set as '-', 'tcp://server:port' or 'udp://server:port'")
	f.StringSlice("const", []string{}, "constant generators")
	f.StringSlice("counter", []string{}, "counter generators")
	f.StringSlice("random", []string{}, "random generators")
	f.Bool("randomize", viper.GetBool("randomize"), "toggle if starting point of generators should be randomized")
	f.Float64("value", viper.GetFloat64("value"), "starting value for generators")
	f.Float64("deviation", viper.GetFloat64("deviation"), "deviation for the next point in generator")
	f.Uint("step", viper.GetUint("step"), "generators interval in seconds")
}

func bindCommonFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	viper.BindPFlag("carbon", f.Lookup("carbon"))
	viper.BindPFlag("const", f.Lookup("const"))
	viper.BindPFlag("counter", f.Lookup("counter"))
	viper.BindPFlag("random", f.Lookup("random"))
	viper.BindPFlag("randomize", f.Lookup("randomize"))
	viper.BindPFlag("value", f.Lookup("value"))
	viper.BindPFlag("deviation", f.Lookup("deviation"))
	viper.BindPFlag("step", f.Lookup("step"))
}
