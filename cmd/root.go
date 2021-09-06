package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "agent",
	Short: "The Spaceship cluster agent",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if _, err := os.Stat(".env"); err == nil {
		log.Info("Loading environment file")
		err := godotenv.Load()
		if err != nil {
			log.WithError(err).Fatal("Error loading environment file")
		}
	}

	viper.SetEnvPrefix("spaceship")
	viper.AutomaticEnv()

	if _, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST"); inCluster {
		log.SetHandler(json.New(os.Stdout))
		if viper.GetBool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
	} else {
		log.SetLevel(log.DebugLevel)
		log.SetHandler(text.New(os.Stdout))
	}
}
