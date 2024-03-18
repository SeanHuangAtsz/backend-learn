package ethereum

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	// // generate mnemonic
	// mnemonic, err := wallet.NewMnemonic()
	// require.NoError(t, err)

	// // get eth chain id and chain params, then generate hd wallet
	// ethChainId := wallet.ChainPrivate
	// chainParam, err := wallet.GetEthChainParams(ethChainId)
	// require.NoError(t, err)
	// hdw, err := wallet.NewHDWallet(mnemonic, "", wallet.BtcChainMainNet, ethChainId)
	// require.NoError(t, err)

	// // generate eth hd wallet
	// w0, err := hdw.NewWallet(wallet.SymbolEth, 0, 0, 0)
	// require.NoError(t, err)
	// w1, err := hdw.NewWallet(wallet.SymbolEth, 0, 0, 1)
	// require.NoError(t, err)

	// // derive address
	// a0 := w0.DeriveAddress()
	// a1 := w1.DeriveAddress()
	// fmt.Printf("a0: %s\na1: %s\n", a0, a1)

	// addrA0, _ := eth.HexToAddress(a0)
	// addrA1, _ := eth.HexToAddress(a1)

}
