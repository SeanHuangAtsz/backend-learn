package tx

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var BigIntEthGWei = big.NewInt(1e9)

type TransactBaseParam struct {
	From      common.Address
	EthValue  *big.Int
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int
	BaseFee   *big.Int
}

func (t *TransactBaseParam) EnsureGasPrice(backend bind.ContractBackend) error {
	// get header
	head, err := backend.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	t.BaseFee = head.BaseFee

	if head.BaseFee == nil {
		// non eip-1559
		if t.GasPrice == nil {
			price, err := backend.SuggestGasPrice(context.Background())
			if err != nil {
				return err
			}
			t.GasPrice = price
		}
	} else {
		// eip-1559
		if t.GasTipCap == nil {
			tip, err := backend.SuggestGasTipCap(context.Background())
			if err != nil {
				return err
			}
			t.GasTipCap = tip
		}
		if t.GasFeeCap == nil {
			gasFeeCap := new(big.Int).Add(
				t.GasTipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			t.GasFeeCap = gasFeeCap
		}
		if t.GasFeeCap.Cmp(t.GasTipCap) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", t.GasFeeCap, t.GasTipCap)
		}
	}
	return nil
}

func (t *TransactBaseParam) GetGasPrice() *big.Int {
	if t.BaseFee != nil {
		return new(big.Int).Add(t.BaseFee, t.GasTipCap)
	} else {
		return t.GasPrice
	}
}

func SignTx(w *wallet.EthWallet, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSigner(w.ChainParams())
	signedTx, err := types.SignTx(tx, signer, w.DeriveNativePrivateKey())
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func MakeTransactOpts(w *wallet.EthWallet, param TransactBaseParam, gasLimit int64, nonce int64) (*bind.TransactOpts, error) {
	var theNonce *big.Int
	if nonce >= 0 {
		theNonce = big.NewInt(nonce)
	}

	if gasLimit < 0 {
		gasLimit = 0
	}

	txOpts := &bind.TransactOpts{
		From:  param.From,
		Nonce: theNonce,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return SignTx(w, tx)
		},
		Value:     param.EthValue,
		GasPrice:  param.GasPrice,
		GasFeeCap: param.GasFeeCap,
		GasTipCap: param.GasTipCap,
		GasLimit:  uint64(gasLimit),
		Context:   context.Background(),
	}
	return txOpts, nil
}

func TransferEther(opts *bind.TransactOpts, backend bind.ContractBackend, addressTo common.Address) (*types.Transaction, error) {
	// nonce
	var nonce uint64
	if opts.Nonce != nil {
		nonce = opts.Nonce.Uint64()
	} else {
		tmp, err := backend.PendingNonceAt(context.Background(), opts.From)
		if err != nil {
			return nil, err
		}
		nonce = tmp
	}

	// gas limit
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		gasLimit = wallet.EtherTransferGas
	}

	param := TransactBaseParam{
		GasPrice:  opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
	}
	// check and set fee
	err := param.EnsureGasPrice(backend)
	if err != nil {
		return nil, err
	}
	opts.GasPrice = param.GasPrice
	opts.GasFeeCap = param.GasFeeCap
	opts.GasTipCap = param.GasTipCap

	var tx *types.Transaction
	var input []byte
	if opts.GasFeeCap == nil {
		baseTx := &types.LegacyTx{
			Nonce:    nonce,
			To:       &addressTo,
			GasPrice: opts.GasPrice,
			Gas:      gasLimit,
			Value:    opts.Value,
			Data:     input,
		}
		tx = types.NewTx(baseTx)
	} else {
		baseTx := &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &addressTo,
			GasFeeCap: opts.GasFeeCap,
			GasTipCap: opts.GasTipCap,
			Gas:       gasLimit,
			Value:     opts.Value,
			Data:      input,
		}
		tx = types.NewTx(baseTx)
	}

	// sign tx
	signedTx, err := opts.Signer(opts.From, tx)
	if err != nil {
		return nil, err
	}

	if opts.NoSend {
		return signedTx, nil
	}

	// send trasaction
	err = backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func HexToAddress(addr string) (common.Address, error) {
	if !common.IsHexAddress(addr) {
		return common.Address{}, errors.New("invalid address")
	}

	return common.HexToAddress(addr), nil
}

func CalcEthFee(gasPrice *big.Int, gas int64) int64 {
	return WeiToGwei(big.NewInt(0).Mul(big.NewInt(gas), gasPrice))
}

func WeiToGwei(v *big.Int) int64 {
	return big.NewInt(0).Div(v, BigIntEthGWei).Int64()
}
