package cmd

import (
	"github.com/onspaceship/agent/pkg/watcher"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch Deployments for updates",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := ctrl.SetupSignalHandler()

		w := watcher.NewWatcher()
		go w.Start(ctx)

		<-ctx.Done()
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
