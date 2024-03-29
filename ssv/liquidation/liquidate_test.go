package liquidation

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"math/big"
	"testing"
)

func TestCalcLiquidation(t *testing.T) {
	_ = logging.SetLogLevel("*", "INFO")
	config.Init("../../config/config.yaml")
	conf := config.GetConfig()

	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()
	curBlock, err := eth1Client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	cluster := Cluster{
		Owner:       common.HexToAddress("0x7ddee51c8375399127b18a0333F80292C3fb6486"),
		OperatorIds: []uint64{1, 2, 3, 4},
		Cluster:     ISSVNetworkCoreCluster{1, 39164045808, 386439710860, true, big.NewInt(1000174745760000000)},
	}
	balances, err := GetClustersBalance([]Cluster{cluster})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("balance", balances[0].String())
	feeInfo, err := GetSSVFeeInfo(cluster.OperatorIds)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("NetworkFee", feeInfo.NetworkFee)
	t.Log("LiquidationThresholdPeriod", feeInfo.LiquidationThresholdPeriod)
	t.Log("MinimumLiquidationCollateral", feeInfo.MinimumLiquidationCollateral)
	t.Log("OperatorsFee", feeInfo.OperatorsFee)

	activeBlock, activeDay := CalcLiquidation(feeInfo, cluster, balances[0], curBlock)

	t.Log("activeBlock", activeBlock-curBlock)
	t.Log("activeDay", activeDay)
}
