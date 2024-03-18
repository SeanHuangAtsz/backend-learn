package wallet

import (
	"encoding/hex"
	"errors"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/accounts"
)

var ErrAddressNotMatch = errors.New("address not match")

type BtcWallet struct {
	symbol      string
	segWitType  SegWitType
	chainParams *chaincfg.Params
	privateKey  *btcec.PrivateKey
	publicKey   *btcec.PublicKey
}

func NewBtcWallet(privateKey string, chainId int, segWitType SegWitType) (*BtcWallet, error) {
	chainParams, err := GetBtcChainParams(chainId)
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(chainParams) {
		return nil, errors.New("key network doesn't match")
	}

	return &BtcWallet{symbol: SymbolBtc,
		chainParams: chainParams, segWitType: segWitType,
		privateKey: wif.PrivKey,
		publicKey:  wif.PrivKey.PubKey()}, nil
}

func NewBtcWalletByPath(path string, seed []byte, chainId int, segWitType SegWitType) (*BtcWallet, error) {
	chainParams, err := GetBtcChainParams(chainId)
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, chainParams)
	if err != nil {
		return nil, err
	}

	privateKey, err := DerivePrivateKeyByPath(masterKey, path, IsFixIssue172)
	if err != nil {
		return nil, err
	}

	return &BtcWallet{symbol: SymbolBtc,
		chainParams: chainParams, segWitType: segWitType,
		privateKey: privateKey,
		publicKey:  privateKey.PubKey()}, nil
}

func (w *BtcWallet) ChainId() int {
	return int(w.chainParams.Net)
}

func (w *BtcWallet) ChainParams() *chaincfg.Params {
	return w.chainParams
}

func (w *BtcWallet) Symbol() string {
	return w.symbol
}

func (w *BtcWallet) DeriveAddress() string {
	addr := w.DeriveNativeAddress()
	if addr != nil {
		return addr.EncodeAddress()
	}
	return ""
}

func (w *BtcWallet) DerivePublicKey() string {
	return hex.EncodeToString(w.publicKey.SerializeCompressed())
}

func (w *BtcWallet) DerivePrivateKey() string {
	wif, err := btcutil.NewWIF(w.privateKey, w.chainParams, true)
	if err != nil {
		log.Println("DerivePrivateKey error:", err)
		return ""
	}
	return wif.String()
}

func (w *BtcWallet) DeriveNativeAddress() btcutil.Address {
	switch w.segWitType {
	case SegWitNone:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		p2pkhAddr, err := btcutil.NewAddressPubKeyHash(keyHash, w.chainParams)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return nil
		}
		return p2pkhAddr
	case SegWitScript:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(keyHash).Script()
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return nil
		}
		addr, err := btcutil.NewAddressScriptHash(scriptSig, w.chainParams)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return nil
		}
		return addr
	case SegWitNative:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(keyHash, w.chainParams)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return nil
		}
		return p2wpkh
	}
	return nil
}

func (w *BtcWallet) DeriveNativePrivateKey() *btcec.PrivateKey {
	return w.privateKey
}

func DerivePrivateKeyByPath(masterKey *hdkeychain.ExtendedKey, path string, fixIssue172 bool) (*btcec.PrivateKey, error) {
	dpath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range dpath {
		if fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// txauthor.SecretsSource
func (w *BtcWallet) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	if w.DeriveAddress() == addr.EncodeAddress() {
		return w.privateKey, true, nil
	}
	return nil, false, ErrAddressNotMatch
}

func (w *BtcWallet) GetScript(addr btcutil.Address) ([]byte, error) {
	return nil, errors.New("GetScript not supported")
}
