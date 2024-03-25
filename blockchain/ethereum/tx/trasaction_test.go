package tx

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/ethereum/node"
	"github.com/SeanHuangAtsz/backend-learn/blockchain/wallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	// generate mnemonic
	mnemonic, err := wallet.NewMnemonic()
	require.NoError(t, err)

	// get eth chain id and chain params
	ethChainId := wallet.ChainPrivate
	chainParam, err := wallet.GetEthChainParams(ethChainId)
	require.NoError(t, err)

	// generate hd wallet
	hdw, err := wallet.NewHDWallet(mnemonic, "", wallet.BtcChainMainNet, ethChainId)
	require.NoError(t, err)

	// generate eth hd wallet
	w0, err := hdw.NewWallet(wallet.SymbolEth, 0, 0, 0)
	require.NoError(t, err)
	w1, err := hdw.NewWallet(wallet.SymbolEth, 0, 0, 1)
	require.NoError(t, err)

	// derive address
	a0 := w0.DeriveAddress()
	a1 := w1.DeriveAddress()
	fmt.Printf("a0: %s\na1: %s\n", a0, a1)

	addrA0, _ := HexToAddress(a0)
	addrA1, _ := HexToAddress(a1)

	// coonect eth client
	// TODO: find a url
	cli, err := node.NewEthClient("")
	require.NoError(t, err)

	// get balance
	{
		fmt.Println("Get balance --------")

		rightBalances := []struct {
			address common.Address
			balance *big.Int
		}{
			{address: addrA0, balance: big.NewInt(0).Mul(big.NewInt(wallet.WeiPerEther), big.NewInt(100))},
			{address: addrA1, balance: big.NewInt(0)},
		}

		for _, rightBal := range rightBalances {
			bal, err := cli.RpcClient.BalanceAt(context.Background(), rightBal.address, nil)
			require.Nil(t, err, "get balance fail")
			fmt.Println("addr balance:", rightBal.address, bal)
			require.True(t, bal.Cmp(rightBal.balance) == 0, "Wrong balance")
		}
	}

	// transfer 6 eth
	transferAmount := big.NewInt(6 * wallet.WeiPerEther)

	// Transfer ether
	{
		// from and amount
		baseParam := TransactBaseParam{
			From:     addrA0,
			EthValue: transferAmount,
		}

		// get gas info
		err = baseParam.EnsureGasPrice(cli.RpcClient)
		require.NoError(t, err)
		fmt.Printf("baseParam: %+v\n", baseParam)

		// build transaction
		opts, err := MakeTransactOpts(w0.(*wallet.EthWallet), baseParam, -1, -1)
		require.NoError(t, err)
		fmt.Printf("opts: %+v\n", opts)

		// build, sign and broadcast transaction
		tx, err := TransferEther(opts, cli.RpcClient, addrA1)
		require.NoError(t, err)

		fmt.Printf("transfer ether, txid: %s, gas: %d, fee: %d\n", tx.Hash().String(), tx.Gas(),
			CalcEthFee(baseParam.GetGasPrice(), int64(tx.Gas())))

		// marshal
		b, _ := json.MarshalIndent(tx, "", " ")
		fmt.Println("tx:", string(b))

		fmt.Println("get transaction receipt ----------")

		// transaction receipt
		tx2, from, blockNumber, err := cli.TransactionByHash(context.Background(), tx.Hash())
		require.NoError(t, err)
		require.True(t, from.String() == a0)

		receipt, err := cli.RpcClient.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		fmt.Println("tx block number:", receipt.BlockNumber)
		require.True(t, receipt.BlockNumber.Cmp(blockNumber) == 0)
		block, err := cli.RpcClient.BlockByNumber(context.Background(), blockNumber)
		require.NoError(t, err)
		sig := types.MakeSigner(chainParam, receipt.BlockNumber, block.Time())
		from2, err := types.Sender(sig, tx2)
		require.NoError(t, err)
		fmt.Println("tx from address:", from2)
		require.True(t, from2.String() == a0)
	}
}
