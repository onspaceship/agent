package cmd

import (
	"github.com/apex/log"
	"github.com/onspaceship/agent/pkg/update"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Agent to current version",
	Run: func(cmd *cobra.Command, args []string) {
		client := kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie())

		err := update.CheckForUpdate(client)

		if err != nil {
			log.WithError(err).Info("Unable to update Agent")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
