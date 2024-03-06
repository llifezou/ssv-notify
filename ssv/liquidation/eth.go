package liquidation

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetEthClient(rpcHost string) (*ethclient.Client, func(), error) {
	if rpcHost == "" {
		return nil, nil, fmt.Errorf("config.yaml is missing 'ethrpc'")
	}

	client, err := ethclient.Dial(rpcHost)
	if err != nil {
		return nil, nil, err
	}

	return client, func() {
		client.Close()
	}, nil
}
