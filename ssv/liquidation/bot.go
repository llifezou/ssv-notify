package liquidation

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"math/big"
	"os"
	"sync"
	"time"
)

const LiquidationMonitoringThreshold = uint64(100)

func StartLiquidationBot(shutdown <-chan struct{}) {
	notifyClient, err := notify.NewNotify()
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	if config.GetConfig().Key == "" {
		log.Warn("config.yam:key is empty")
		os.Exit(1)
	}

	endBlock, err := ScanSSVCluster(0)
	if err != nil {
		log.Warnw("GetSSVCluster", "err", err)
		os.Exit(1)
	}

	var clusterCh = make(chan Cluster, 20)

	go scanEventLoop(endBlock)
	go liquidateLoop(notifyClient, clusterCh)

	if err = scanLiquidationCluster(clusterCh); err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(6 * time.Minute) // 6 min = 30 block
	for {
		select {
		case <-ticker.C:
			log.Info("Liquidation Bot Monitoring trigger")
			err = scanLiquidationCluster(clusterCh)
			if err != nil {
				log.Warn(err)
			}
		case <-shutdown:
			return
		}
	}
}

func liquidateLoop(notify *notify.Notify, clusterCh <-chan Cluster) {
	for {
		select {
		case cluster := <-clusterCh:
			log.Infow("liquidate monitor", "owner", cluster.Owner, "operator", operatorIdsToString(cluster.OperatorIds), "validatorCount", cluster.Cluster.ValidatorCount)
			go liquidate(notify, cluster)
		}
	}
}

func scanEventLoop(endBlock uint64) {
	ticker := time.NewTicker(12 * time.Second)
	var err error
	var startBlock = endBlock
	for {
		select {
		case <-ticker.C:
			startBlock, err = ScanSSVCluster(startBlock)
			if err != nil {
				log.Errorw("GetSSVCluster", "err", err)
			}
		}
	}
}

var (
	liquidateClusterMap = map[string]struct{}{}
	liquidateLk         sync.Mutex
)

func liquidate(notify *notify.Notify, cluster Cluster) {
	key := cluster.Owner.String() + ":" + operatorIdsToString(cluster.OperatorIds)
	liquidateLk.Lock()
	if _, ok := liquidateClusterMap[key]; ok {
		log.Infow("the cluster added to liquidation monitoring", "cluster", key)
		liquidateLk.Unlock()
		return
	}
	liquidateClusterMap[key] = struct{}{}
	liquidateLk.Unlock()

	done := func() {
		liquidateLk.Lock()
		delete(liquidateClusterMap, key)
		liquidateLk.Unlock()
	}
	defer done()

	conf := config.GetConfig()
	ssvNetworkAddr, ssvNetworkViewAddr, err := getSSVNetworkAddrs(conf.Network)
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		log.Warn(err)
		return
	}
	defer cancel()

	ticker := time.NewTicker(6 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Warnw("liquidation monitoring", "cluster", key)

			cluster = GetClusterOfOwnerAndOperators(cluster.Owner.String(), operatorIdsToString(cluster.OperatorIds))
			if !cluster.Cluster.Active {
				log.Warnw("cluster is liquidated", "owner", key)
				return
			}

			ssvView, err := NewLiquidation(ssvNetworkViewAddr, eth1Client)
			if err != nil {
				log.Warnw("NewLiquidation", "err", err)
				continue
			}

			isLiquidatable, err := ssvView.IsLiquidatable(nil, cluster.Owner, cluster.OperatorIds, cluster.Cluster)
			if err != nil {
				log.Warnw("IsLiquidatable", "err", err)
				continue
			}

			log.Infow("Liquidation will be initiated", "owner", cluster.Owner, "operator", operatorIdsToString(cluster.OperatorIds), "isLiquidatable", isLiquidatable)

			if isLiquidatable {
				ssv, err := NewLiquidation(ssvNetworkAddr, eth1Client)
				if err != nil {
					log.Warnw("NewLiquidation", "err", err)
					continue
				}

				gasPrice, err := eth1Client.SuggestGasPrice(context.Background())
				if err != nil {
					gasPrice = big.NewInt(200000000000) // 200 gwei
				} else {
					gasPrice = big.NewInt(0).Add(gasPrice, gasPrice) // suggestGasPrice*2
				}

				if gasPrice.Uint64() < 20000000000 {
					gasPrice = big.NewInt(20000000000) // min 20 gwei
				}

				opts := makeOpts(conf.Key, gasPrice, 200000, getChainId(conf.Network))
				tx, err := ssv.Liquidate(opts, cluster.Owner, cluster.OperatorIds, cluster.Cluster)
				if err != nil {
					log.Warnw("liquidate Tx", "err", err)
					return
				}

				log.Infow("waiting tx...", "tx", tx.Hash().String())
				txTicker := time.NewTicker(time.Second * 6)
				for {
					hash := tx.Hash()
					<-txTicker.C
					status := getTxStatus(eth1Client, hash)
					if Success == status {
						log.Info("successfully executed liquidate tx: %s", hash.String())
						notify.Send(fmt.Sprintf("Liquidate tx: %s", hash.String()))
						return
					} else if Failed == status {
						log.Warnf("failed to execute the tx, please check: %s", hash.String())
						return
					}
				}
			}
		}
	}
}

func scanLiquidationCluster(clusterCh chan<- Cluster) error {
	conf := config.GetConfig()
	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		return err
	}
	defer cancel()

	curBlock, err := eth1Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	clusters := GetClusters()
	if len(clusters) == 0 {
		log.Warn("no cluster information was scanned")
		return fmt.Errorf("no cluster information was scanned")
	}

	activeClusters := []Cluster{}
	for _, cluster := range clusters {
		c := cluster
		if cluster.Cluster.Active {
			activeClusters = append(activeClusters, c)
		}
	}

	balances, err := GetClustersBalance(activeClusters)
	if err != nil {
		log.Warnw("GetClustersBalance", "err", err)
		return err
	}

	for i, c := range activeClusters {
		feeInfo, err := GetSSVFeeInfo(c.OperatorIds)
		if err != nil {
			log.Warnw("GetSSVFeeInfo", "err", err)
			continue
		}

		activeBlock, activeDay := CalcLiquidation(feeInfo, c, balances[i], curBlock)
		if activeBlock-curBlock < LiquidationMonitoringThreshold {
			msg := fmt.Sprintf("Liquidation Bot Monitoring: clusterOwner: %s; operators: %s; Liquidation blockHeight: %d; Operational Runway: %d Day",
				c.Owner, operatorIdsToString(c.OperatorIds), activeBlock, activeDay)

			log.Warn(msg)
			clusterCh <- c
			continue
		}

		log.Infow("Liquidation Bot Monitoring", "clusterOwner", c.Owner.String(), "OperatorIds", operatorIdsToString(c.OperatorIds), "LiquidateBlock", activeBlock-curBlock)
	}

	return nil
}

func makeOpts(key string, gasPrice *big.Int, gasLimit uint64, chainID int64) *bind.TransactOpts {
	from := crypto.PubkeyToAddress(str2pri(key).PublicKey)

	txOpts := &bind.TransactOpts{
		From: from,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			signedTx, err := SignTxFromPriKey(tx, key, chainID)
			if err != nil {
				return nil, err
			}
			return signedTx, nil
		},
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		Context:  context.Background(),
	}
	return txOpts
}

func SignTxFromPriKey(tx *types.Transaction, key string, chainID int64) (*types.Transaction, error) {
	signer := types.NewLondonSigner(big.NewInt(chainID))
	h := signer.Hash(tx)
	sign, err := crypto.Sign(h[:], str2pri(key))
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, sign)
}

func str2pri(pkStr string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(pkStr)
	if err != nil {
		return nil
	}
	return privateKey
}

const (
	Pending = iota
	Success
	Failed
)

func getTxStatus(client *ethclient.Client, hash common.Hash) int {
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return Pending
	}
	if receipt.Status == 1 {
		return Success
	} else {
		return Failed
	}
}
