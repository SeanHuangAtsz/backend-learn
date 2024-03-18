package tx

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/wallet"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	// generate mnemonic
	mnemonic, err := wallet.NewMnemonic()
	require.NoError(t, err)
	fmt.Println("mnemonic: ", mnemonic)

	// create hd wallet for btc testnet and eth mainnet
	btcChainId := wallet.BtcChainRegtest
	hdw, err := wallet.NewHDWallet(mnemonic, "", btcChainId, wallet.ChainMainNet)
	require.NoError(t, err)

	// generate btc hd wallets
	w0, err := hdw.NewNativeSegWitWallet(0, 0, 0)
	require.NoError(t, err)
	w1, err := hdw.NewWallet(wallet.SymbolBtc, 0, 0, 0)
	require.NoError(t, err)
	w2, err := hdw.NewSegWitWallet(0, 0, 0)
	require.NoError(t, err)
	w3, err := hdw.NewSegWitWallet(0, 0, 1)
	require.NoError(t, err)

	// get address
	chainParams, _ := wallet.GetBtcChainParams(btcChainId)
	a0 := w0.DeriveAddress()
	a1 := w1.DeriveAddress()
	a2 := w2.DeriveAddress()
	a3 := w3.DeriveAddress()
	addrA0, _ := wallet.DecodeAddress(a0, chainParams)
	addrA1, _ := wallet.DecodeAddress(a1, chainParams)
	addrA2, _ := wallet.DecodeAddress(a2, chainParams)
	addrA3, _ := wallet.DecodeAddress(a3, chainParams)
	fmt.Println("address0: ", addrA0)
	fmt.Println("address1: ", addrA1)
	fmt.Println("address2: ", addrA2)
	fmt.Println("address3: ", addrA3)

	// list unspent
	var utxo = btcjson.ListUnspentResult{}

	transferAmount := 2.2
	var tx *BtcTransaction

	// build tx
	{
		//feePerKb, err := cli.EstimateFeePerKb()
		feePerKb := int64(80 * 1000)

		unspent := BtcUnspent{TxID: utxo.TxID, Vout: utxo.Vout,
			ScriptPubKey: utxo.ScriptPubKey, RedeemScript: utxo.RedeemScript,
			Amount: utxo.Amount}

		out1 := BtcOutput{Address: addrA1, Amount: BtcToSatoshi(transferAmount)}
		out2 := BtcOutput{Address: addrA2, Amount: BtcToSatoshi(transferAmount)}
		out3 := BtcOutput{Address: addrA3, Amount: BtcToSatoshi(transferAmount)}

		tx, err = NewBtcTransaction([]BtcUnspent{unspent}, []BtcOutput{out1, out2, out3},
			addrA0, feePerKb, chainParams)
		require.NoError(t, err)
	}

	{ // fee
		fee := tx.GetFee()
		fmt.Println("fee:", fee)
	}

	{ // sign
		err = tx.Sign(w0.(*wallet.BtcWallet))
		require.NoError(t, err)
	}

	{ // decode
		ret := tx.Decode()
		b, _ := json.MarshalIndent(ret, "", " ")
		fmt.Println("decoded tx:", string(b))
	}
}
