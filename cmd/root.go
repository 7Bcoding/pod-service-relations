package cmd

import (
	"github.com/spf13/cobra"

	"pod-service-relations/config"
	"pod-service-relations/database"
	"pod-service-relations/logging"
	"pod-service-relations/server"
)

func init() {
	cobra.OnInitialize(config.InitConfigs)
	rootCmd.Flags().BoolP("server", "s", false, "server mode")
}

func runServer() {
	logging.Init()
	database.Init()
	server.Init()
}

var rootCmd = &cobra.Command{
	Use:   "pod_service",
	Short: "pod_service is a pod-services relations system for xunlei k8s cluster",
	Long:  "pod_service -s|-w",
	Run: func(cmd *cobra.Command, args []string) {
		serverMode, err := cmd.Flags().GetBool("server")
		if err != nil {
			panic(err)
		}
		if serverMode {
			runServer()
			return
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
