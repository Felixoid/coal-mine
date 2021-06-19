package cmd

/*
Copyright © 2021 Mikhail f. Shiryaev <mr.felixoid@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bytes"
	"fmt"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

// configExampleCmd represents the configExample command
var configExampleCmd = &cobra.Command{
	Use:   "config-example",
	Short: "Prints a configuration example to STDOUT.",
	RunE:  printConfig,
}

func init() {
	rootCmd.AddCommand(configExampleCmd)
}

func printConfig(cmd *cobra.Command, args []string) error {
	buf := new(bytes.Buffer)
	config := &Config{General: General{From: "-2d", Until: "now", Step: 120, Randomize: true, Value: 333, Deviation: 15.15}}
	config.Const = []string{"metric.const.example1", "metric.const.example{2..5}"}
	config.Counter = []string{"metric.counter.example1", "metric.counter.example{2..5}"}
	config.Random = []string{"metric.random.example{1,{2..5},.subdir}"}
	config.Custom = append(config.Custom, Custom{
		Name: "custom.random.generator{1..10}",
		Type: "random",
		General: General{
			From:      "-2d",
			Until:     "1d",
			Step:      300,
			Randomize: false,
			Value:     1000,
			Deviation: 123.456,
		},
	})
	config.Custom = append(config.Custom, Custom{
		Name: "custom.counter.generator{1..10}",
		Type: "counter",
		General: General{
			From:      "-2h",
			Until:     "now",
			Step:      10,
			Randomize: true,
			Value:     100,
			Deviation: 123.456,
		},
	})
	encoder := toml.NewEncoder(buf).CompactComments(true).Indentation(" ").Order(toml.OrderPreserve)
	encoder.Encode(config)
	fmt.Fprint(cmd.OutOrStdout(), buf.String())
	return nil
}
