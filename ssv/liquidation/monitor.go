package liquidation

import (
	"context"
	"fmt"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"math/big"
	"os"
	"strings"
	"time"
)

func StartLiquidationMonitor(shutdown <-chan struct{}) {
	notifyClient, err := notify.NewNotify()
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	var endBlock uint64
	if err, endBlock = monitor(notifyClient, 0); err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(6 * time.Hour)
	for {
		select {
		case <-ticker.C:
			log.Info("Liquidation Monitoring trigger")
			err, endBlock = monitor(notifyClient, endBlock)
			if err != nil {
				log.Warn(err)
			}
		case <-shutdown:
			return
		}
	}
}

func monitor(notify *notify.Notify, startBlock uint64) (error, uint64) {
	conf := config.GetConfig()
	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return err, startBlock
	}
	defer cancel()

	endBlock, err := ScanSSVCluster(startBlock)
	if err != nil {
		log.Warnw("GetSSVCluster", "err", err)
		return err, startBlock
	}

	curBlock, err := eth1Client.BlockNumber(context.Background())
	if err != nil {
		return err, startBlock
	}

	for _, owner := range conf.LiquidationMonitor.ClusterOwner {
		clusters := GetClusterOfOwner(strings.ToLower(owner))
		if len(clusters) == 0 {
			log.Warnw("no cluster information was scanned", "owner", owner)
			continue
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
			log.Warnw("GetClustersBalance", "err", err)
			continue
		}

		if len(inactiveClusters) != 0 {
			clusters = append(activeClusters, inactiveClusters...)
			balances = append(balances, inactiveClusterBalances...)
		}

		for i, c := range clusters {
			feeInfo, err := GetSSVFeeInfo(c.OperatorIds)
			if err != nil {
				log.Warnw("GetSSVFeeInfo", "err", err)
				continue
			}

			activeBlock, activeDay := CalcLiquidation(feeInfo, c, balances[i], curBlock)
			if activeDay <= conf.LiquidationMonitor.Threshold {
				msg := fmt.Sprintf("Liquidation Monitoring: clusterOwner: %s; operators: %s; Liquidation blockHeight: %d; Operational Runway: %d Day",
					owner, operatorIdsToString(c.OperatorIds), activeBlock, activeDay)

				log.Warn(msg)
				notify.Send(msg)
				continue
			}
			log.Infow("Liquidation Monitoring", "clusterOwner", c.Owner.String(), "OperatorIds", operatorIdsToString(c.OperatorIds), "Operational Runway", activeDay)
		}
	}

	return nil, endBlock
}
