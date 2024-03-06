package liquidation

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/llifezou/ssv-notify/config"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	mainnetStartBlock = 17507487
	goerliStartBlock  = 9203578
	holeskyStartBlock = 181612
)

func getStartBlock(network string) (uint64, error) {
	switch strings.ToLower(network) {
	case "mainnet":
		return mainnetStartBlock, nil
	case "goerli":
		return goerliStartBlock, nil
	case "holesky":
		return holeskyStartBlock, nil
	default:
		return 0, errors.New("the network does not support")
	}
}

func ScanAllSSVCluster(dir string) error {
	conf := config.GetConfig()
	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return err
	}

	curBlock, err := eth1Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	cancel()

	clusters, _, err := ScanClusterInfo(0)
	if err != nil {
		return err
	}

	activeClusters, inactiveClusters := []Cluster{}, []Cluster{}
	inactiveClusterBalances := []*big.Int{}
	for _, cluster := range clusters {
		c := cluster
		if cluster.Cluster.Active {
			activeClusters = append(activeClusters, c)
		} else {
			inactiveClusters = append(inactiveClusters, c)
			inactiveClusterBalances = append(inactiveClusterBalances, big.NewInt(0))
		}
	}

	balances, err := GetClustersBalance(activeClusters)
	if err != nil {
		return err
	}

	if len(inactiveClusters) != 0 {
		clusters = append(activeClusters, inactiveClusters...)
		balances = append(balances, inactiveClusterBalances...)
	}

	activeBlocks, activeDays := []uint64{}, []uint64{}

	isPrint := dir == ""
	var w *tabwriter.Writer

	if isPrint {
		w = tabwriter.NewWriter(os.Stdout, 8, 4, 2, ' ', 0)
		fmt.Fprintf(w, "ClusterOwner\tOperatorId\tValidatorCount\tNetworkFeeIndex\tIndex\tActive\tBalance\tLiquidationBlock\tRunway\n")
	}

	for i, cluster := range clusters {
		feeInfo, err := GetSSVFeeInfo(cluster.OperatorIds)
		if err != nil {
			log.Warnw("GetSSVFeeInfo", "err", err)
			continue
		}

		activeBlock, activeDay := CalcLiquidation(feeInfo, cluster, balances[i], curBlock)
		activeBlocks = append(activeBlocks, activeBlock)
		activeDays = append(activeDays, activeDay)

		if isPrint {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%v\t%t\t%s\t%d\t%d\t\n", cluster.Owner, operatorIdsToString(cluster.OperatorIds),
				cluster.Cluster.ValidatorCount, cluster.Cluster.NetworkFeeIndex, cluster.Cluster.Index, cluster.Cluster.Active, cluster.Cluster.Balance.String(), activeBlock, activeDay)
		}
	}

	if !isPrint {
		return writeToCsv(clusters, activeDays, activeBlocks, dir, conf.Network)
	}

	if isPrint {
		if err = w.Flush(); err != nil {
			return fmt.Errorf("flushing output: %+v", err)
		}
	}

	return nil
}

func ScanClusterInfo(startBlock uint64) ([]Cluster, uint64, error) {
	conf := config.GetConfig()
	if startBlock == 0 {
		var err error
		startBlock, err = getStartBlock(conf.Network)
		if err != nil {
			return nil, 0, err
		}
	}

	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return nil, 0, err
	}
	defer cancel()

	curBlock, err := eth1Client.BlockNumber(context.Background())
	if err != nil {
		return nil, 0, err
	}

	ssvNetworkAddr, _, err := getSSVNetworkAddrs(conf.Network)
	if err != nil {
		return nil, 0, err
	}

	// key = owner:operatorIds
	clusterInfoMap := make(map[string]Cluster)

	for fromBlock := startBlock; fromBlock < curBlock; {
		nextBlock := fromBlock + 20000
		if nextBlock >= curBlock {
			nextBlock = curBlock
		}

		filter := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(nextBlock)),
			Addresses: []common.Address{ssvNetworkAddr},
			Topics:    [][]common.Hash{{ValidatorAddedTopic, ValidatorRemovedTopic, ClusterDepositedTopic, ClusterWithdrawnTopic, ClusterLiquidatedTopic, ClusterReactivatedTopic}},
		}

		fromBlock = nextBlock

		log.Infow("scan block", "fromBlock", fromBlock, "nextBlock", nextBlock)
		addLogs, err := eth1Client.FilterLogs(context.Background(), filter)
		if err != nil {
			continue
		}

		for _, l := range addLogs {
			var addr common.Address
			copy(addr[:], l.Topics[1][12:])
			event := TopicToEvent[l.Topics[0]]
			data, err := ssvNetworkABI.Events[event].Inputs.Unpack(l.Data)
			if err != nil {
				log.Warnw("ValidatorAdded events unpack failed", "err", err)
				continue
			}

			operatorIds := data[0].([]uint64)
			clusterInfo, err := json.Marshal(data[len(data)-1])
			if err != nil {
				log.Warnw("Marshal", "err", err)
				continue
			}

			log.Infow("event analysis", "owner", addr, "event", event, "clusterInfo", string(clusterInfo))
			err = updateClusterInfo(clusterInfoMap, addr, operatorIds, clusterInfo)
			if err != nil {
				log.Warnw("updateClusterInfo", "err", err)
				continue
			}
		}
	}

	var clusterInfoSlice = make([]Cluster, 0)
	for _, clusterInfo := range clusterInfoMap {
		if clusterInfo.Cluster.ValidatorCount == 0 {
			log.Infow("the cluster has no validators", "owner", clusterInfo.Owner, "operators", operatorIdsToString(clusterInfo.OperatorIds))
			continue
		}
		clusterInfoSlice = append(clusterInfoSlice, clusterInfo)
	}

	sort.Slice(clusterInfoSlice, func(i, j int) bool {
		return clusterInfoSlice[i].Cluster.ValidatorCount > clusterInfoSlice[j].Cluster.ValidatorCount
	})

	return clusterInfoSlice, curBlock, nil
}

func updateClusterInfo(clusterInfoMap map[string]Cluster, owner common.Address, operatorIds []uint64, clusterInfo []byte) error {
	cluster := ISSVNetworkCoreCluster{}
	err := json.Unmarshal(clusterInfo, &cluster)
	if err != nil {
		return err
	}

	key1 := owner.String()
	key2 := operatorIdsToString(operatorIds)
	key := fmt.Sprintf("%s:%s", key1, key2)

	if _, ok := clusterInfoMap[key]; !ok {
		clusterInfoMap[key] = Cluster{
			Owner:       owner,
			OperatorIds: operatorIds,
			Cluster:     cluster,
		}
	} else {
		clusterInfoMap[key] = Cluster{
			Owner:       owner,
			OperatorIds: operatorIds,
			Cluster:     cluster,
		}
	}

	return nil
}

func writeToCsv(clusterInfos []Cluster, runways, activeBlocks []uint64, dir, network string) error {
	if len(clusterInfos) == 0 {
		return nil
	}

	t := time.Now().Format("2006-01-02T15:04:05")
	path := filepath.Join(dir, network+"-cluster-"+t+".csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"clusterOwner", "operatorId", "validatorCount", "networkFeeIndex", "index", "active", "balance", "liquidationBlock", "runway"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	for i, c := range clusterInfos {
		row := []string{
			c.Owner.String(),
			operatorIdsToString(c.OperatorIds),
			fmt.Sprintf("%d", c.Cluster.ValidatorCount),
			fmt.Sprintf("%d", c.Cluster.NetworkFeeIndex),
			fmt.Sprintf("%d", c.Cluster.Index),
			fmt.Sprintf("%t", c.Cluster.Active),
			c.Cluster.Balance.String(),
			fmt.Sprintf("%d", activeBlocks[i]),
			fmt.Sprintf("%d", runways[i]),
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	log.Infow("ScanCluster successful", "file", path)
	return nil
}

func operatorIdsToString(operatorIds []uint64) string {
	key := ""
	for i := 0; i < len(operatorIds); i++ {
		if key == "" {
			key = fmt.Sprintf("%d", operatorIds[i])
			continue
		}
		key = fmt.Sprintf("%s,%d", key, operatorIds[i])
	}
	return key
}
