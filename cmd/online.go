package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Felixoid/coal-mine/generator"
	"github.com/spf13/cobra"
)

// onlineCmd represents the online command
var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "Generate metrics each second to emulate even load",
	Long: `The command behaves like the main one, but doesn't have
--from and --until flags, ignores those parameters from
the config, and generates points for the current second.

It's highly recommended to use it with --randomize
parameter to spread the generation over time.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindCommonFlags(cmd)

		readConfig()
		unmarshalConfig()
	},
	RunE: onlineGeneration,
}

func init() {
	rootCmd.AddCommand(onlineCmd)

	f := onlineCmd.Flags()
	f.SortFlags = false

	commonFlags(onlineCmd)
}

func onlineGeneration(cmd *cobra.Command, args []string) (asyncErr error) {
	config.ResetStartStop()

	writer, err := config.GetCarbonWriter()
	if err != nil {
		return err
	}

	ggg, err := config.ToGenerators()
	if len(ggg) == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, len(ggg))
	_ = ctx

	go func() {
		select {
		case <-CatchedSignals:
			// TODO: use zap logging and log signal
			cancel()
			asyncErr = nil
		case err = <-errs:
			// TODO: use zap logging and log error
			cancel()
			asyncErr = err
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(ggg))

	write := func(gg generator.Generators) {
		defer wg.Done()
		tick := time.Second
		if !gg.Randomized() {
			tick = time.Duration(gg.Step()) * time.Second
		}
		ticker := time.NewTicker(tick)
		for valid := gg; true; valid = getNextGenerators(gg) {
			// TODO: log amount of sent metrics
			select {
			case t := <-ticker.C:
				ctxTimeout, cancel := context.WithTimeout(ctx, tick)
				defer cancel()
				n, err := valid.WriteAllToWithContext(ctxTimeout, writer)
				if err != nil && !errors.Is(err, generator.ErrEmptyGens) {
					errs <- fmt.Errorf("error while sending metrics, %d bytes sent: %w", n, err)
					return
				}
				gg.SetStop(uint(t.Unix()))
			}
		}
	}

	for _, gg := range ggg {
		go write(gg)
	}

	wg.Wait()

	return
}

func getNextGenerators(gg generator.Generators) generator.Generators {
	gens := make([]generator.Generator, 0, len(gg.List()))
	for _, g := range gg.List() {
		if err := g.Next(); err == nil {
			gens = append(gens, g)
		}
	}
	valid := generator.Generators{}
	valid.SetList(gens)
	return valid
}
