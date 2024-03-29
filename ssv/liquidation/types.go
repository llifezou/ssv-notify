package liquidation

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type FeeInfo struct {
	NetworkFee                   uint64
	LiquidationThresholdPeriod   uint64
	MinimumLiquidationCollateral *big.Int
	OperatorsFee                 []uint64
}

type Cluster struct {
	Owner       common.Address
	OperatorIds []uint64
	Cluster     ISSVNetworkCoreCluster
}
