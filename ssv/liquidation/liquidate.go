package liquidation

import "math/big"

var (
	maxActiveDay   = uint64(999999)
	maxActiveBlock = maxActiveDay * 7200
)

// CalcLiquidation : Operational Runway
// https://docs.ssv.network/learn/stakers/clusters/cluster-balance
func CalcLiquidation(feeInfo *FeeInfo, cluster Cluster, curBalance *big.Int, curBlock uint64) (activeBlock uint64, operationalRunway uint64) {
	if curBalance.Cmp(feeInfo.MinimumLiquidationCollateral) <= 0 {
		return curBlock, 0
	}

	period := feeInfo.LiquidationThresholdPeriod
	fee := feeInfo.NetworkFee
	for _, opFee := range feeInfo.OperatorsFee {
		fee += opFee
	}
	if fee == 0 {
		return curBlock + maxActiveBlock, maxActiveDay
	}

	preValidatorFee := big.NewInt(0).Mul(big.NewInt(0).SetUint64(period), big.NewInt(0).SetUint64(fee))
	requireBalance := big.NewInt(0).Mul(preValidatorFee, big.NewInt(int64(cluster.Cluster.ValidatorCount)))
	if curBalance.Cmp(requireBalance) > 0 {
		activeBalance := big.NewInt(0).Sub(curBalance, big.NewInt(0).Mul(big.NewInt(0).Mul(big.NewInt(0).SetUint64(feeInfo.LiquidationThresholdPeriod), big.NewInt(0).SetUint64(fee)), big.NewInt(int64(cluster.Cluster.ValidatorCount))))
		if activeBalance.Uint64() == 0 {
			return curBlock, 0
		}
		if cluster.Cluster.ValidatorCount == 0 {
			return 0, 0
		}

		preValidatorBalance := big.NewInt(0).Div(activeBalance, big.NewInt(int64(cluster.Cluster.ValidatorCount)))
		if preValidatorBalance.Uint64() == 0 {
			return curBlock, 0
		}

		activeBlock = big.NewInt(0).Div(preValidatorBalance, big.NewInt(0).SetUint64(fee)).Uint64()
		if activeBlock == 0 {
			return curBlock, 0
		}

		return curBlock + activeBlock, activeBlock / 7200
	}

	return curBlock, 0
}
