package liquidation

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"strings"
)

var (
	mainnetSSVNetworkAddr     = common.HexToAddress("0xDD9BC35aE942eF0cFa76930954a156B3fF30a4E1")
	mainnetSSVNetworkViewAddr = common.HexToAddress("0xafE830B6Ee262ba11cce5F32fDCd760FFE6a66e4")

	holeskySSVNetworkAddr     = common.HexToAddress("0x38A4794cCEd47d3baf7370CcC43B560D3a1beEFA")
	boleskySSVNetworkViewAddr = common.HexToAddress("0x352A18AEe90cdcd825d1E37d9939dCA86C00e281")
)

var (
	ssvNetworkAbi     = `[{"inputs":[],"name":"ApprovalNotWithinTimeframe","type":"error"},{"inputs":[],"name":"CallerNotOwner","type":"error"},{"inputs":[],"name":"CallerNotWhitelisted","type":"error"},{"inputs":[],"name":"ClusterAlreadyEnabled","type":"error"},{"inputs":[],"name":"ClusterDoesNotExists","type":"error"},{"inputs":[],"name":"ClusterIsLiquidated","type":"error"},{"inputs":[],"name":"ClusterNotLiquidatable","type":"error"},{"inputs":[],"name":"ExceedValidatorLimit","type":"error"},{"inputs":[],"name":"FeeExceedsIncreaseLimit","type":"error"},{"inputs":[],"name":"FeeIncreaseNotAllowed","type":"error"},{"inputs":[],"name":"FeeTooLow","type":"error"},{"inputs":[],"name":"IncorrectClusterState","type":"error"},{"inputs":[],"name":"IncorrectValidatorState","type":"error"},{"inputs":[],"name":"InsufficientBalance","type":"error"},{"inputs":[],"name":"InvalidOperatorIdsLength","type":"error"},{"inputs":[],"name":"InvalidPublicKeyLength","type":"error"},{"inputs":[],"name":"MaxValueExceeded","type":"error"},{"inputs":[],"name":"NewBlockPeriodIsBelowMinimum","type":"error"},{"inputs":[],"name":"NoFeeDeclared","type":"error"},{"inputs":[],"name":"NotAuthorized","type":"error"},{"inputs":[],"name":"OperatorAlreadyExists","type":"error"},{"inputs":[],"name":"OperatorDoesNotExist","type":"error"},{"inputs":[],"name":"OperatorsListNotUnique","type":"error"},{"inputs":[],"name":"SameFeeChangeNotAllowed","type":"error"},{"inputs":[],"name":"TargetModuleDoesNotExist","type":"error"},{"inputs":[],"name":"TokenTransferFailed","type":"error"},{"inputs":[],"name":"UnsortedOperatorsList","type":"error"},{"inputs":[],"name":"ValidatorAlreadyExists","type":"error"},{"inputs":[],"name":"ValidatorDoesNotExist","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterDeposited","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterLiquidated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterReactivated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterWithdrawn","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"shares","type":"bytes"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ValidatorAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ValidatorRemoved","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"deposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"liquidate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"reactivate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"publicKey","type":"bytes"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"bytes","name":"sharesData","type":"bytes"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"registerValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"publicKey","type":"bytes"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"removeValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	ssvNetworkViewAbi = `[{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"getBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getLiquidationThresholdPeriod","outputs":[{"internalType":"uint64","name":"","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getMinimumLiquidationCollateral","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getNetworkFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"getOperatorFee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"isLiquidatable","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`

	ssvNetworkABI     abi.ABI
	ssvNetworkViewABI abi.ABI
)

var (
	ValidatorAddedTopic     = common.HexToHash("0x48a3ea0796746043948f6341d17ff8200937b99262a0b48c2663b951ed7114e5")
	ValidatorRemovedTopic   = common.HexToHash("0xccf4370403e5fbbde0cd3f13426479dcd8a5916b05db424b7a2c04978cf8ce6e")
	ClusterDepositedTopic   = common.HexToHash("0x2bac1912f2481d12f0df08647c06bee174967c62d3a03cbc078eb215dc1bd9a2")
	ClusterWithdrawnTopic   = common.HexToHash("0x39d1320bbda24947e77f3560661323384aa0a1cb9d5e040e617e5cbf50b6dbe0")
	ClusterLiquidatedTopic  = common.HexToHash("0x1fce24c373e07f89214e9187598635036111dbb363e99f4ce498488cdc66e688")
	ClusterReactivatedTopic = common.HexToHash("0xc803f8c01343fcdaf32068f4c283951623ef2b3fa0c547551931356f456b6859")
)

const (
	ValidatorAdded     = "ValidatorAdded"
	ValidatorRemoved   = "ValidatorRemoved"
	ClusterDeposited   = "ClusterDeposited"
	ClusterWithdrawn   = "ClusterWithdrawn"
	ClusterLiquidated  = "ClusterLiquidated"
	ClusterReactivated = "ClusterReactivated"
)

var TopicToEvent = map[common.Hash]string{
	ValidatorAddedTopic:     ValidatorAdded,
	ValidatorRemovedTopic:   ValidatorRemoved,
	ClusterDepositedTopic:   ClusterDeposited,
	ClusterWithdrawnTopic:   ClusterWithdrawn,
	ClusterLiquidatedTopic:  ClusterLiquidated,
	ClusterReactivatedTopic: ClusterReactivated,
}

const (
	getNetworkFee                   = "getNetworkFee"
	getLiquidationThresholdPeriod   = "getLiquidationThresholdPeriod"
	getMinimumLiquidationCollateral = "getMinimumLiquidationCollateral"
	getOperatorFee                  = "getOperatorFee"
	getBalance                      = "getBalance"
)

var log = logging.Logger("liquidation-monitor")

func init() {
	var err error
	ssvNetworkViewABI, err = abi.JSON(strings.NewReader(ssvNetworkViewAbi))
	if err != nil {
		panic(err)
	}
	ssvNetworkABI, err = abi.JSON(strings.NewReader(ssvNetworkAbi))
}

func getSSVNetworkAddrs(network string) (ssvNetwork common.Address, ssvNetworkView common.Address, err error) {
	switch strings.ToLower(network) {
	case "mainnet":
		return mainnetSSVNetworkAddr, mainnetSSVNetworkViewAddr, nil
	case "holesky":
		return holeskySSVNetworkAddr, boleskySSVNetworkViewAddr, nil
	default:
		return common.Address{}, common.Address{}, errors.New("the network does not support")
	}
}

func getChainId(network string) int64 {
	switch strings.ToLower(network) {
	case "mainnet":
		return 1
	case "holesky":
		return 17000
	default:
		return 0
	}
}
