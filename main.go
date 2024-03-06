package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"github.com/llifezou/ssv-notify/ssv/liquidation"
	"github.com/llifezou/ssv-notify/ssv/operator"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var log = logging.Logger("main")

var (
	configPath string
	dir        string
)

var rootCmd = &cobra.Command{
	Use:   "ssv-notify",
	Short: "ssv-notify",
	Long:  `ssv monitoring notifications.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config/config.yaml", "Path to configuration file")
	scanClusterCmd.PersistentFlags().StringVarP(&dir, "dir", "d", "", "Save scan results")
}

func main() {
	_ = logging.SetLogLevel("*", "INFO")
	config.Init(configPath)

	ssvToolsCmd.AddCommand(scanClusterCmd)
	rootCmd.AddCommand(operatorMonitorCmd, liquidationMonitorCmd, ssvToolsCmd, AlarmTestCmd)

	_ = rootCmd.Execute()
}

var liquidationMonitorCmd = &cobra.Command{
	Use:     "liquidation-monitor",
	Short:   "liquidation monitor",
	Example: "./ssv-notify liquidation-monitor",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("liquidation monitor service start")

		var shutdown = make(chan struct{})
		go liquidation.StartLiquidationMonitor(shutdown)

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		close(shutdown)

		log.Info("liquidation monitor service shutting down")
	},
}

var operatorMonitorCmd = &cobra.Command{
	Use:     "operator-monitor",
	Short:   "operator monitor",
	Example: "./ssv-notify operator-monitor",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("operator monitor service start")

		var shutdown = make(chan struct{})
		go operator.StartOperatorMonitor(shutdown)

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		close(shutdown)

		log.Info("operator monitor service shutting down")
	},
}

var ssvToolsCmd = &cobra.Command{
	Use:   "ssv-tools",
	Short: "ssv tools",
	Long:  `./ssv-notify ssv-tools -h`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var scanClusterCmd = &cobra.Command{
	Use:     "scan-cluster",
	Short:   "Scan cluster and calculate operational runway",
	Example: "./ssv-notify ssv-tools scan-cluster",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Start scanning")
		err := liquidation.ScanAllSSVCluster(dir)
		if err != nil {
			log.Warnw("ScanAllSSVCluster failed", "err", err)
		}
	},
}

var AlarmTestCmd = &cobra.Command{
	Use:     "alarm-test",
	Short:   "Test alarms can be used",
	Example: "./ssv-notify alarm-test",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Start alarm test")
		notifyClient, err := notify.NewNotify()
		if err != nil {
			log.Warn(err)
		}
		notifyClient.Send("alarm test!")
	},
}
