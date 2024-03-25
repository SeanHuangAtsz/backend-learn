package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

type HDWallet struct {
	seed       []byte
	btcChainId int
	ethChainId int
}

func NewHDWallet(mnemonic, password string, btcChainId int, ethChainId int) (*HDWallet, error) {
	mnemonic = strings.ReplaceAll(mnemonic, "\n", "")
	mnemonic = strings.ReplaceAll(mnemonic, "\r", "")

	// generate seed
	seed, err := NewSeedFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, err
	}
	fmt.Println("seed: ", hex.EncodeToString(seed))
	return &HDWallet{seed: seed, btcChainId: btcChainId, ethChainId: ethChainId}, nil
}

func (h *HDWallet) NewWallet(symbol string, accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip44Path(symbol, h.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}

	return h.NewWalletByPath(symbol, path, SegWitNone)
}

func (h *HDWallet) NewSegWitWallet(accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip49Path(SymbolBtc, h.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}
	return h.NewWalletByPath(SymbolBtc, path, SegWitScript)
}

func (h *HDWallet) NewNativeSegWitWallet(accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip84Path(SymbolBtc, h.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}
	return h.NewWalletByPath(SymbolBtc, path, SegWitNative)
}

func (h *HDWallet) NewWalletByPath(symbol string, path string, segWitType SegWitType) (Wallet, error) {
	var w Wallet
	var err error

	switch symbol {
	case SymbolBtc:
		w, err = NewBtcWalletByPath(path, h.seed, h.btcChainId, segWitType)
	case SymbolEth:
		w, err = NewEthWalletByPath(path, h.seed, h.ethChainId)
	// case SymbolTrx:
	// 	w, err = NewTrxWalletByPath(path, h.seed)
	default:
		err = fmt.Errorf("invalid symbol: %s", symbol)
	}

	if err != nil {
		return nil, err
	}
	return w, nil
}

func MakeBip44Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(44, symbol, chainId, accountIndex, changeType, index)
}

func MakeBip49Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(49, symbol, chainId, accountIndex, changeType, index)
}

func MakeBip84Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(84, symbol, chainId, accountIndex, changeType, index)
}

func MakeBipXPath(bipType int, symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	var coinType int
	switch symbol {
	case SymbolEth:
		coinType = 60
	case SymbolBtc:
		chainParams, err := GetBtcChainParams(chainId)
		if err != nil {
			return "", err
		}
		coinType = int(chainParams.HDCoinType)
	case SymbolTrx:
		coinType = 195
	default:
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}

	if accountIndex < 0 || index < 0 {
		return "", errors.New("invalid account index or index")
	}
	if changeType != ChangeTypeExternal && changeType != ChangeTypeInternal {
		return "", errors.New("invalid change type")
	}
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d", bipType, coinType, accountIndex, changeType, index), nil
}

func GetBtcChainParams(chainId int) (*chaincfg.Params, error) {
	switch chainId {
	case BtcChainMainNet:
		return &chaincfg.MainNetParams, nil
	case BtcChainTestNet3:
		return &chaincfg.TestNet3Params, nil
	case BtcChainRegtest:
		return &chaincfg.RegressionNetParams, nil
	case BtcChainSimNet:
		return &chaincfg.SimNetParams, nil
	default:
		return nil, fmt.Errorf("unknown btc chainId: %d", chainId)
	}
}

func NewMnemonic() (string, error) {
	// generate entropy
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	fmt.Println("entropy: ", hex.EncodeToString(entropy))

	// entropy to mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	fmt.Println("mnemonic: ", mnemonic)

	return mnemonic, nil
}

func NewSeedFromMnemonic(mnemonic, password string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}
