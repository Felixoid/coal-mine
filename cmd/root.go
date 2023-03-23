// Package cmd provides cobra commands for github.com/Felixoid/coal-mine application
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
package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Felixoid/coal-mine/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "coal-mine",
	Short: "Metrics generator for graphite-web in carbon format",
	Long: `Metrics generator for graphite-web in carbon plain-text format.

In default mode it generates historical data to mock
a meaningfull graphite-web requests. For example, when
one needs to reproduce an issue without data itself,
he can give a config example to generate points.

The curly braces in  metric names are expanded as it
would be done in zsh. For example, "server{01..10}.soft{1..5}"
will generate 50 metrics with two nodes.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindCommonFlags(cmd)
		f := cmd.Flags()
		viper.BindPFlag("from", f.Lookup("from"))
		viper.BindPFlag("until", f.Lookup("until"))

		readConfig()
		unmarshalConfig()
	},
	RunE:          generation,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
}

// CatchedSignals are catched by program and processed gracefully
var CatchedSignals = make(chan os.Signal, 1)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	setDefaultConfig()

	f := rootCmd.Flags()
	f.SortFlags = false

	pf := rootCmd.PersistentFlags()
	pf.SortFlags = false

	commonFlags(rootCmd)
	f.String("from", viper.GetString("from"), "starting point for generators in graphtie-web format")
	f.String("until", viper.GetString("until"), "final point for generators in graphtie-web format")
}

func generation(cmd *cobra.Command, args []string) error {
	writer, err := config.GetCarbonWriter()
	if err != nil {
		return err
	}

	ggg, err := config.ToGenerators()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-CatchedSignals:
			cancel()
		}
	}()

	errs := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(ggg))

	// preallocate buffers
	buf := make([][]byte, len(ggg))
	for i := 0; i < len(ggg); i++ {
		buf[i] = make([]byte, 0, 256)
	}

	write := func(gg generator.Generators, buf *[]byte) {
		defer wg.Done()
		n, err := gg.WriteAllToWithContext(ctx, writer, buf)
		if err != nil {
			errs <- fmt.Errorf("error while sending metrics, %d bytes sent: %w", n, err)
		}
	}

	wait := make(chan struct{})
	go func() {
		for i, gg := range ggg {
			go write(gg, &buf[i])
		}
		wg.Wait()
		close(wait)
	}()

	select {
	case err := <-errs:
		return err
	case <-wait:
	}

	return nil
}
