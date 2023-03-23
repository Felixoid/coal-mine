package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Generate metrics list",
	Long:  `The command generates metrics to stdout`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindCommonFlags(cmd)

		readConfig()
		unmarshalConfig()
	},
	RunE: metricsGeneration,
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	f := metricsCmd.Flags()
	f.SortFlags = false

	commonFlags(metricsCmd)
}

func metricsGeneration(cmd *cobra.Command, args []string) (asyncErr error) {

	ggg, err := config.ToGenerators()
	if err != nil {
		return err
	}
	if len(ggg) == 0 {
		return nil
	}

	for _, gg := range ggg {
		for _, g := range gg.List() {
			fmt.Println(g.Name())
		}
	}

	return
}
