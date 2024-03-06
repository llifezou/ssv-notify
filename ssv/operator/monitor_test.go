package operator

import (
	"fmt"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"os"
	"testing"
)

func TestMonitor(t *testing.T) {
	config.Init("../../config/config.yaml")
	conf := config.GetConfig()

	notifyClient, err := notify.NewNotify()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monitor(notifyClient, "", "holesky", conf.OperatorMonitor.ClusterOwner)
}
