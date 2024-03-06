package liquidation

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"testing"
)

func TestScanAllSSVCluster(t *testing.T) {
	_ = logging.SetLogLevel("*", "INFO")
	config.Init("../../config/config.yaml")
	err := ScanAllSSVCluster("")
	if err != nil {
		t.Fatal(err)
	}
}
