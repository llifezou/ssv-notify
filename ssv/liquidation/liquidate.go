package liquidation

import "math/big"

var (
	maxActiveDay   = uint64(999999)
	maxActiveBlock = maxActiveDay * 7200
)

// CalcLiquidation : Operational Runway
// https://docs.ssv.network/learn/stakers/clusters/cluster-balance
func CalcLiquidation(feeInfo *FeeInfo, cluster Cluster, curBalance *big.Int, curBlock uint64) (activeBlock uint64, operationalRunway uint64) {
	if cluster.Cluster.ValidatorCount == 0 {
		return curBlock + maxActiveBlock, maxActiveDay
	}

	if curBalance.Cmp(feeInfo.MinimumLiquidationCollateral) <= 0 {
		return curBlock, 0
	}

	minimumBlocksBeforeLiquidation := int64(feeInfo.LiquidationThresholdPeriod)

	burnRate := uint64(0)
	for _, opFee := range feeInfo.OperatorsFee {
		burnRate += opFee
	}

	fee := feeInfo.NetworkFee + burnRate

	perLiquidationThreshold := big.NewInt(0).Mul(big.NewInt(minimumBlocksBeforeLiquidation), big.NewInt(int64(fee)))
	liquidationThreshold := big.NewInt(0).Mul(perLiquidationThreshold, big.NewInt(int64(cluster.Cluster.ValidatorCount)))

	if curBalance.Cmp(liquidationThreshold) > 0 {
		activeBalance := big.NewInt(0).Sub(curBalance, feeInfo.MinimumLiquidationCollateral)

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
