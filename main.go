package main

import (
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/ssv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configPath string

func init() {
	runCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config/config.yaml", "Path to configuration file")
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ssv-notify",
		Short: "ssv-notify",
		Long:  `ssv operator monitoring notifications.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(runCmd)
	_ = rootCmd.Execute()
}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run monitor",
	Example: "./ssv-notify run",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("monitor service start")
		config.Init(configPath)

		var shutdown = make(chan struct{})
		go ssv.StartMonitor(shutdown)

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		close(shutdown)

		log.Println("monitor service shutting down")
	},
}
