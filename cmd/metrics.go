package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Generate metrics list",
	Long:  `The command generates metrics to stdout`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindMetricsFlags(cmd)

		readConfig()
		unmarshalConfig()
	},
	RunE: metricsGeneration,
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	f := metricsCmd.Flags()
	f.SortFlags = false

	metricsFlags(metricsCmd)
}

func metricsGeneration(cmd *cobra.Command, args []string) (asyncErr error) {

	ggg, err := config.ToGenerators()
	if err != nil {
		return err
	}
	if len(ggg) == 0 {
		return nil
	}

	fmt.Print("# value=")
	io.WriteString(os.Stdout, strconv.FormatFloat(config.Value, 'f', -1, 64))
	fmt.Print(" dev=")
	io.WriteString(os.Stdout, strconv.FormatFloat(config.Deviation, 'f', -1, 64))
	fmt.Print(" from=")
	fmt.Print(config.From)
	fmt.Print(" until=")
	fmt.Print(config.Until)
	fmt.Print(" step=")
	fmt.Print(config.Step)
	fmt.Println()

	for _, gg := range ggg {
		for _, g := range gg.List() {
			fmt.Println(g.Name())
		}
	}

	return
}
