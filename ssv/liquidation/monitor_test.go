package liquidation

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"testing"
)

func TestMonitor(t *testing.T) {
	_ = logging.SetLogLevel("*", "INFO")
	config.Init("../../config/config.yaml")
	notifyClient, err := notify.NewNotify()
	if err != nil {
		t.Fatal(err)
	}
	monitor(notifyClient, 0)
}
