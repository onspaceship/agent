package cmd

import (
	"github.com/onspaceship/agent/pkg/socket"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Spaceship",
	Run: func(cmd *cobra.Command, args []string) {
		exit := make(chan bool)

		go socket.StartListener(exit)

		<-exit
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
