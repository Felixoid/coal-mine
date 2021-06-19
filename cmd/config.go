package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/Felixoid/coal-mine/generator"
	"github.com/go-graphite/carbonapi/date"
	"github.com/spf13/viper"
)

// General is the general part of configs
type General struct {
	From      string `toml:"from,omitempty" json:"from,omitempty" comment:"from in graphite-web format, the local TZ is used"`
	start     uint
	Until     string `toml:"until,omitempty" json:"until,omitempty" comment:"until in graphite-web format, the local TZ is used"`
	stop      uint
	Step      uint    `toml:"step,omitempty" json:"step,omitempty" comment:"step in seconds"`
	Randomize bool    `toml:"randomize" json:"randomize" comment:"randomize starting time with [0,step)"`
	Value     float64 `toml:"value,omitempty" json:"value,omitempty" comment:"first value for all generators"`
	Deviation float64 `toml:"deviation,omitempty" json:"deviation,omitempty" comment:"deviation of the values, const will be generated around, counter will add [0,value+deviation), random will calculate next value around previous"`
}

// Custom is a config for a generators with special parameters. Is readed only from a config file.
type Custom struct {
	Name    string `toml:"name,omitempty" json:"name,omitempty" comment:"names for generator, braces are expanded like in shell"`
	Type    string `toml:"type,omitempty" json:"type,omitempty" comment:"type of generator"`
	General `mapstructure:",squash"`
}

// ToGenerators returns generator.Generators for a given custom config
func (c *Custom) ToGenerators() (generator.Generators, error) {
	return generator.NewExpand(c.Type, c.Name, c.start, c.stop, c.Step, c.Randomize, c.Value, c.Deviation)
}

// Config is a general application config. Everything besides Generators can be set both from flags and config file.
type Config struct {
	Carbon  string   `toml:"carbon" json:"carbon" comment:"carbon-server address or '-' for STDOUT, should be set as '-', 'tcp://server:port' or 'udp://server:port'"`
	Const   []string `toml:"const,omitempty" json:"const,omitempty" comment:"names for constant generators, braces are expanded like in shell"`
	Counter []string `toml:"counter,omitempty" json:"counter,omitempty" comment:"names for counter generators, braces are expanded like in shell"`
	Random  []string `toml:"random,omitempty" json:"random,omitempty" comment:"names for random generators, braces are expanded like in shell"`
	General `mapstructure:",squash"`
	Custom  []Custom `toml:"custom,omitempty" json:"custom,omitempty" comment:"generators with custom parameters can be specified separately"`
}

var now = time.Now().Unix()

// SetStartStop process graphite-web from/until and sets start and stop fields
func (c *Config) SetStartStop() {
	c.start = uint(date.DateParamToEpoch(c.From, "", now, nil))
	c.stop = uint(date.DateParamToEpoch(c.Until, "", now, nil))
	for _, g := range c.Custom {
		g.start = uint(date.DateParamToEpoch(g.From, "", now, nil))
		g.stop = uint(date.DateParamToEpoch(g.Until, "", now, nil))
	}
}

// ResetStartStop sets all start and stop values of the general and cutom configs to the current timestamp
func (c *Config) ResetStartStop() {
	c.start = uint(now)
	c.stop = uint(now)
	for _, g := range c.Custom {
		g.start = uint(now)
		g.stop = uint(now)
	}
}

// ToGenerators returns slice of generator.Generators for main config and each Config.Custom
func (c *Config) ToGenerators() ([]generator.Generators, error) {
	result := make([]generator.Generators, 0, len(c.Custom)+3)
	for _, n := range c.Const {
		gen, err := generator.NewExpand("const", n, c.start, c.stop, c.Step, c.Randomize, c.Value, c.Deviation)
		if err != nil {
			return nil, fmt.Errorf("unable to create new constant generators: %w", err)
		}
		result = append(result, gen)
	}
	for _, n := range c.Counter {
		gen, err := generator.NewExpand("counter", n, c.start, c.stop, c.Step, c.Randomize, c.Value, c.Deviation)
		if err != nil {
			return nil, fmt.Errorf("unable to create new counter generators: %w", err)
		}
		result = append(result, gen)
	}
	for _, n := range c.Random {
		gen, err := generator.NewExpand("random", n, c.start, c.stop, c.Step, c.Randomize, c.Value, c.Deviation)
		if err != nil {
			return nil, fmt.Errorf("unable to create new random generators: %w", err)
		}
		result = append(result, gen)
	}
	for _, custom := range c.Custom {
		gen, err := custom.ToGenerators()
		if err != nil {
			return nil, fmt.Errorf("unable to create new custom generators for %v: %w", custom, err)
		}
		result = append(result, gen)
	}
	return result, nil
}

// GetCarbonWriter returns net.Conn. If it's unable to parse the Carbon field, an error is not nil.
func (c *Config) GetCarbonWriter() (io.Writer, error) {
	if c.Carbon == "-" {
		return os.Stdout, nil
	}
	u, err := url.Parse(c.Carbon)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL from %s: %w", config.Carbon, err)
	}
	if u.Scheme != "tcp" && u.Scheme != "udp" {
		return nil, fmt.Errorf("scheme %s in %s is not valid", u.Scheme, config.Carbon)
	}

	conn, err := net.Dial(u.Scheme, u.Host)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to carbon: %w", err)
	}
	return conn, nil
}

var (
	cfgFile string
	config  *Config
)

func unmarshalConfig() {
	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Fail to unmarshal config %v", viper.AllSettings())
	}
	config.SetStartStop()
}

func setDefaultConfig() {
	viper.SetDefault("carbon", "-")
	viper.SetDefault("const", []string{})
	viper.SetDefault("counter", []string{})
	viper.SetDefault("random", []string{})
	viper.SetDefault("from", "-24h")
	viper.SetDefault("until", "now")
	viper.SetDefault("step", 60)
	viper.SetDefault("randomize", false)
	viper.SetDefault("value", 10)
	viper.SetDefault("deviation", 5)
	viper.SetDefault("generators", []Custom{})
}

func readConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	if cfgFile == "" {
		return
	}
	viper.SetConfigFile(cfgFile)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Fatalf("Error while reading from %s: %v", viper.ConfigFileUsed(), err)
	}
}
