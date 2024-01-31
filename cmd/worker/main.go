/*
Copyright © 2024 Piotr Tobiasz
*/
package cmd

import (
	"log/slog"

	"starling/cmd"
	"starling/internal/worker"

	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts worker",
	Run: func(_ *cobra.Command, _ []string) {
		w := worker.NewWorker()
		w.AddTask("start", func() error {
			slog.Info("Starting")
			return nil
		})
		w.AddTask("stop", func() error {
			slog.Info("Stopping")
			return nil
		})
		w.AddTask("stop", func() error {
			slog.Info("Stopping again")
			return nil
		})
		w.Run()
	},
}

func init() {
	cmd.RootCmd.AddCommand(workerCmd)
}
