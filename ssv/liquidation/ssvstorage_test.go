package liquidation

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"math/big"
	"testing"
)

//	struct StorageProtocol {
//	   /// @notice The block number when the network fee index was last updated
//	   uint32 networkFeeIndexBlockNumber;
//	   /// @notice The count of validators governed by the DAO
//	   uint32 daoValidatorCount;
//	   /// @notice The block number when the DAO index was last updated
//	   uint32 daoIndexBlockNumber;
//	   /// @notice The maximum limit of validators per operator
//	   uint32 validatorsPerOperatorLimit;
//	   /// @notice The current network fee value
//	   uint64 networkFee;
//	   /// @notice The current network fee index value
//	   uint64 networkFeeIndex;
//	   /// @notice The current balance of the DAO
//	   uint64 daoBalance;
//	   /// @notice The minimum number of blocks before a liquidation event can be triggered
//	   uint64 minimumBlocksBeforeLiquidation;
//	   /// @notice The minimum collateral required for liquidation
//	   uint64 minimumLiquidationCollateral;
//	   /// @notice The period in which an operator can declare a fee change
//	   uint64 declareOperatorFeePeriod;
//	   /// @notice The period in which an operator fee change can be executed
//	   uint64 executeOperatorFeePeriod;
//	   /// @notice The maximum increase in operator fee that is allowed (percentage)
//	   uint64 operatorMaxFeeIncrease;
//	   /// @notice The maximum value in operator fee that is allowed (SSV)
//	   uint64 operatorMaxFee;
//	}
func TestSSVStorageProtocol(t *testing.T) {
	_ = logging.SetLogLevel("*", "INFO")
	config.Init("../../config/config.yaml")
	conf := config.GetConfig()

	eth1Client, cancel, err := GetEthClient(conf.EthRpc)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	genKey := common.BytesToHash(crypto.Keccak256([]byte("ssv.network.storage.protocol")))
	keyInt := big.NewInt(0).SetBytes(genKey.Bytes())
	resultInt := big.NewInt(0).Sub(keyInt, big.NewInt(1))

	key := common.BytesToHash(resultInt.Bytes())
	data, err := eth1Client.StorageAt(context.Background(), holeskySSVNetworkAddr, key, big.NewInt(1192921))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
	t.Log(big.NewInt(0).SetBytes(data[28:]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data[24:28]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data[20:24]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data[16:20]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data[8:16]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data[:8]).Uint64())

	resultInt2 := big.NewInt(0).Sub(keyInt, big.NewInt(0))

	key2 := common.BytesToHash(resultInt2.Bytes())
	data2, err := eth1Client.StorageAt(context.Background(), holeskySSVNetworkAddr, key2, big.NewInt(1192921))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data2)
	t.Log(big.NewInt(0).SetBytes(data2[24:]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data2[16:24]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data2[8:16]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data2[:8]).Uint64())

	resultInt3 := big.NewInt(0).Add(keyInt, big.NewInt(1))

	key3 := common.BytesToHash(resultInt3.Bytes())
	data3, err := eth1Client.StorageAt(context.Background(), holeskySSVNetworkAddr, key3, big.NewInt(1192921))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data3)
	t.Log(big.NewInt(0).SetBytes(data3[24:]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data3[16:24]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data3[8:16]).Uint64())
	t.Log(big.NewInt(0).SetBytes(data3[:8]).Uint64())
}
