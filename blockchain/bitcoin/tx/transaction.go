package tx

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/wallet"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
)

type BtcUnspent struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript,omitempty"`
	Amount       float64 `json:"amount"`
}

type BtcOutput struct {
	Address btcutil.Address `json:"address"`
	Amount  int64           `json:"amount"`
}

type BtcTransaction struct {
	txauthor.AuthoredTx
	chainParams *chaincfg.Params
	feePerKb    int64
}

func NewBtcTransaction(unspents []BtcUnspent, outputs []BtcOutput,
	changeAddress btcutil.Address, feePerKb int64, chainCfg *chaincfg.Params) (*BtcTransaction, error) {

	if len(unspents) == 0 || changeAddress == nil || feePerKb <= 0 {
		return nil, errors.New("wrong params")
	}
	if !changeAddress.IsForNet(chainCfg) {
		return nil, errors.New("change address is not the corresponding network address")
	}
	changeBytes, err := txscript.PayToAddrScript(changeAddress)
	if err != nil {
		return nil, err
	}

	feeRatePerKb := btcutil.Amount(feePerKb)

	txOuts, err := makeTxOutputs(outputs, feeRatePerKb, chainCfg)
	if err != nil {
		return nil, err
	}

	changeSource := txauthor.ChangeSource{
		NewScript: func() ([]byte, error) {
			return changeBytes, nil
		},
		ScriptSize: len(changeBytes),
	}

	unsignedTx, err := txauthor.NewUnsignedTransaction(txOuts, feeRatePerKb, makeInputSource(unspents), &changeSource)
	if err != nil {
		return nil, err
	}
	// Randomize change position, if change exists, before signing.  This
	// doesn't affect the serialize size, so the change amount will still
	// be valid.
	if unsignedTx.ChangeIndex >= 0 {
		unsignedTx.RandomizeChangePosition()
	}

	return &BtcTransaction{*unsignedTx, chainCfg, feePerKb}, nil
}

func (t *BtcTransaction) Sign(wallet *wallet.BtcWallet) error {
	return t.SignWithSecretsSource(wallet)
}

func (t *BtcTransaction) SignWithSecretsSource(secretsSource txauthor.SecretsSource) error {
	err := t.AddAllInputScripts(secretsSource)
	if err != nil {
		return err
	}
	err = validateMsgTx(t.Tx, t.PrevScripts, t.PrevInputValues)
	if err != nil {
		return err
	}

	return nil
}

func (t *BtcTransaction) GetFee() int64 {
	fee := t.TotalInput - txauthor.SumOutputValues(t.Tx.TxOut)
	return int64(fee)
}

func (t *BtcTransaction) Decode() *btcjson.TxRawDecodeResult {
	return DecodeMsgTx(t.Tx, t.chainParams)
}

// validateMsgTx verifies transaction input scripts for tx.  All previous output
// scripts from outputs redeemed by the transaction, in the same order they are
// spent, must be passed in the prevScripts slice.
func validateMsgTx(tx *wire.MsgTx, prevScripts [][]byte,
	inputValues []btcutil.Amount) error {

	inputFetcher, err := txauthor.TXPrevOutFetcher(
		tx, prevScripts, inputValues,
	)
	if err != nil {
		return err
	}

	hashCache := txscript.NewTxSigHashes(tx, inputFetcher)
	for i, prevScript := range prevScripts {
		vm, err := txscript.NewEngine(
			prevScript, tx, i, txscript.StandardVerifyFlags, nil,
			hashCache, int64(inputValues[i]), inputFetcher,
		)
		if err != nil {
			return fmt.Errorf("cannot create script engine: %s", err)
		}
		err = vm.Execute()
		if err != nil {
			return fmt.Errorf("cannot validate transaction: %s", err)
		}
	}
	return nil
}

func BtcToSatoshi(v float64) int64 {
	amt, _ := btcutil.NewAmount(v)
	return int64(amt)
}

func makeTxOutputs(outputs []BtcOutput, relayFeePerKb btcutil.Amount, chainCfg *chaincfg.Params) ([]*wire.TxOut, error) {
	outLen := len(outputs)
	if outLen == 0 {
		return nil, errors.New("tx output is empty")
	}

	txOuts := make([]*wire.TxOut, 0, outLen)
	for i := 0; i < outLen; i++ {
		out := &outputs[i]

		if !out.Address.IsForNet(chainCfg) {
			return nil, errors.New("out address is not the corresponding network address")
		}

		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(out.Address)
		if err != nil {
			return nil, err
		}
		txOut := &wire.TxOut{
			Value:    out.Amount,
			PkScript: pkScript,
		}

		// if err = txrules.CheckOutput(txOut, relayFeePerKb); err != nil {
		// 	return nil, err
		// }

		txOuts = append(txOuts, txOut)
	}
	return txOuts, nil
}

func makeInputSource(unspents []BtcUnspent) txauthor.InputSource {
	sz := len(unspents)
	// Current inputs and their total value.  These are closed over by the
	// returned input source and reused across multiple calls.
	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, sz)
	currentInputValues := make([]btcutil.Amount, 0, sz)
	currentScripts := make([][]byte, 0, sz)

	return func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn, []btcutil.Amount, [][]byte, error) {
		for currentTotal < target && len(unspents) != 0 {
			u := unspents[0]
			unspents = unspents[1:]

			hash, _ := chainhash.NewHashFromStr(u.TxID)
			nextInput := wire.NewTxIn(&wire.OutPoint{
				Hash:  *hash,
				Index: u.Vout,
			}, nil, nil)

			amount, _ := btcutil.NewAmount(u.Amount)
			s, _ := hex.DecodeString(u.ScriptPubKey)

			currentTotal += amount
			currentInputs = append(currentInputs, nextInput)
			currentInputValues = append(currentInputValues, amount)
			currentScripts = append(currentScripts, s)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
}
