package liquidation

import (
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/ssv/utils"
	"math/big"
	"strings"
)

var clustersMap = make(map[string]map[string]Cluster, 0)

func GetClusterOfOwner(owner string) []Cluster {
	var clusters = []Cluster{}
	for _, cluster := range clustersMap[owner] {
		c := cluster
		clusters = append(clusters, c)
	}

	return clusters
}

func GetSSVCluster(startBlock uint64) (uint64, error) {
	allClusters, endBlock, err := ScanClusterInfo(startBlock)
	if err != nil {
		return 0, err
	}

	for _, cluster := range allClusters {
		owner := strings.ToLower(cluster.Owner.String())
		if _, ok := clustersMap[owner]; !ok {
			clustersMap[owner] = map[string]Cluster{}
			c := cluster
			clustersMap[owner][operatorIdsToString(c.OperatorIds)] = c
		} else {
			c := cluster
			clustersMap[owner][operatorIdsToString(c.OperatorIds)] = c
		}

	}

	log.Infow("GetSSVCluster", "startBlock", startBlock, "endBlock", endBlock)

	return endBlock, nil
}

func GetSSVFeeInfo(operatorIds []uint64) (*FeeInfo, error) {
	conf := config.GetConfig()
	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return nil, err
	}
	defer cancel()

	_, ssvNetworkViewAddr, err := getSSVNetworkAddrs(conf.Network)
	if err != nil {
		return nil, err
	}

	var callStructs = make([]utils.Struct0, 0)
	callStructs = append(callStructs, utils.Struct0{
		Target:   ssvNetworkViewAddr,
		CallData: ssvNetworkViewABI.Methods[getNetworkFee].ID,
	})
	callStructs = append(callStructs, utils.Struct0{
		Target:   ssvNetworkViewAddr,
		CallData: ssvNetworkViewABI.Methods[getLiquidationThresholdPeriod].ID,
	})
	callStructs = append(callStructs, utils.Struct0{
		Target:   ssvNetworkViewAddr,
		CallData: ssvNetworkViewABI.Methods[getMinimumLiquidationCollateral].ID,
	})
	for _, opId := range operatorIds {
		data, err := ssvNetworkViewABI.Methods[getOperatorFee].Inputs.Pack(opId)
		if err != nil {
			return nil, err
		}
		callStructs = append(callStructs, utils.Struct0{
			Target:   ssvNetworkViewAddr,
			CallData: append(ssvNetworkViewABI.Methods[getOperatorFee].ID, data...),
		})
	}

	addr, err := utils.GetMultiCallAddr(conf.Network)
	if err != nil {
		return nil, err
	}

	multiCall, err := utils.NewMulticall(addr, eth1Client)
	if err != nil {
		return nil, err
	}

	log.Info("multicall aggregate: getSSVFeeInfo")
	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return nil, err
	}
	results := make([]*big.Int, 0)
	for _, r := range outs[1].([][]uint8) {
		results = append(results, big.NewInt(0).SetBytes(r))
	}

	var feeInfo FeeInfo
	feeInfo.NetworkFee = results[0].Uint64()
	feeInfo.LiquidationThresholdPeriod = results[1].Uint64()
	feeInfo.MinimumLiquidationCollateral = results[2]
	operatorsFee := make([]uint64, 0)
	for _, r := range results[3:] {
		operatorsFee = append(operatorsFee, r.Uint64())
	}
	feeInfo.OperatorsFee = operatorsFee

	return &feeInfo, nil
}

func GetClustersBalance(clusters []Cluster) ([]*big.Int, error) {
	conf := config.GetConfig()
	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return nil, err
	}
	defer cancel()

	_, ssvNetworkViewAddr, err := getSSVNetworkAddrs(conf.Network)
	if err != nil {
		return nil, err
	}
	addr, err := utils.GetMultiCallAddr(conf.Network)
	if err != nil {
		return nil, err
	}

	multiCall, err := utils.NewMulticall(addr, eth1Client)
	if err != nil {
		return nil, err
	}

	results := make([]*big.Int, 0)

	batchCluster := [][]Cluster{clusters}
	clustersLen := len(clusters)
	if clustersLen > 10 {
		batchCluster = [][]Cluster{}
		batchCount := clustersLen / 10
		if clustersLen%10 != 0 {
			batchCount++
		}
		for i := 0; i < clustersLen; i += 10 {
			end := i + 10
			if end > clustersLen {
				end = clustersLen
			}
			batchCluster = append(batchCluster, clusters[i:end])
		}
	}

	for _, cs := range batchCluster {
		var callStructs = make([]utils.Struct0, 0)
		for _, c := range cs {
			data, err := ssvNetworkViewABI.Methods[getBalance].Inputs.Pack(c.Owner, c.OperatorIds, c.Cluster)
			if err != nil {
				return nil, err
			}
			callStructs = append(callStructs, utils.Struct0{
				Target:   ssvNetworkViewAddr,
				CallData: append(ssvNetworkViewABI.Methods[getBalance].ID, data...),
			})
		}

		log.Infow("multicall aggregate: getBalance", "count", len(cs))
		outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
		if err != nil {
			log.Warnw("multicall failed")
			return nil, err
		}
		for _, r := range outs[1].([][]uint8) {
			results = append(results, big.NewInt(0).SetBytes(r))
		}
	}

	return results, nil
}
