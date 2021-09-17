package cmd

import (
	"github.com/onspaceship/agent/pkg/socket"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Spaceship",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := ctrl.SetupSignalHandler()

		go socket.StartListener(ctx)

		<-ctx.Done()
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
