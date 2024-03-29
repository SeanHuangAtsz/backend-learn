package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/params"
)

type SegWitType int

const (
	SymbolEth = "ETH"
	SymbolBtc = "BTC"
	SymbolTrx = "TRX"

	BtcChainMainNet  = int(wire.MainNet)
	BtcChainTestNet3 = int(wire.TestNet3)
	BtcChainRegtest  = int(wire.TestNet)
	BtcChainSimNet   = int(wire.SimNet)

	ChainMainNet      = 1 // for ETH
	ChainRopsten      = 3 // for ETH
	ChainRinkeby      = 4 // for ETH
	ChainGoerli       = 5 // for ETH
	ChainHolesky      = 17000
	ChainSepolia      = 11155111
	ChainBsc          = 56    // for Binance Smart Chain Mainnet
	ChainBscTestnet   = 97    // for Binance Smart Chain Testnet
	ChainMatic        = 137   // for Polygon Matic Mainnet
	ChainMaticTestnet = 80001 // for Polygon Matic Testnet
	ChainPrivate      = 1337  // for ETH

	SegWitNone   SegWitType = 0
	SegWitScript SegWitType = 1
	SegWitNative SegWitType = 2

	ChangeTypeExternal = 0
	ChangeTypeInternal = 1 // Usually used for change, not visible to the outside world

	SatoshiPerBitcoin = 1e8
	SunPerTrx         = 1e6
	GweiPerEther      = 1e9
	WeiPerGwei        = 1e9
	WeiPerEther       = 1e18

	EtherTransferGas = 21000

	TokenShowDecimals = 9

	IsFixIssue172 = false
)

const (
	p2pkh = "p2pkh"
	p2sh  = "p2sh"
)

var (
	BscChainConfig          = &params.ChainConfig{}
	BscTestnetChainConfig   = &params.ChainConfig{}
	MaticChainConfig        = &params.ChainConfig{}
	MaticTestnetChainConfig = &params.ChainConfig{}
)

func GetEthChainParams(chainId int) (*params.ChainConfig, error) {
	switch chainId {
	case ChainMainNet:
		return params.MainnetChainConfig, nil
	//case ChainRopsten:
	//	return params.RopstenChainConfig, nil
	//case ChainRinkeby:
	//	return params.RinkebyChainConfig, nil
	case ChainGoerli:
		return params.GoerliChainConfig, nil
	case ChainHolesky:
		return params.HoleskyChainConfig, nil
	case ChainSepolia:
		return params.SepoliaChainConfig, nil
	case ChainMatic:
		return MaticChainConfig, nil
	case ChainMaticTestnet:
		return MaticTestnetChainConfig, nil
	case ChainBsc:
		return BscChainConfig, nil
	case ChainBscTestnet:
		return BscTestnetChainConfig, nil
	case ChainPrivate:
		return params.AllEthashProtocolChanges, nil
	default:
		return nil, fmt.Errorf("unknown eth chainId: %d", chainId)
	}
}
