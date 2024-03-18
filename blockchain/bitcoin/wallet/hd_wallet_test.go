package wallet

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

func TestNewHDWallet(t *testing.T) {
	m, err := NewMnemonic()
	require.NoError(t, err)

	NewHDWallet(m, "123", 0, 0)
}

func TestNewHDSegWitWallet(t *testing.T) {
	m, err := NewMnemonic()
	require.NoError(t, err)

	// btc mainnet
	h, err := NewHDWallet(m, "123", 3652501241, 0)
	require.NoError(t, err)

	w, err := h.NewSegWitWallet(0, 0, 0)
	require.NoError(t, err)

	fmt.Println("address: ", w.DeriveAddress())
	fmt.Println("private key: ", w.DerivePrivateKey())
	fmt.Println("public key: ", w.DerivePublicKey())
	fmt.Println("symbol: ", w.Symbol())
	fmt.Println("chain id: ", w.ChainId())
}

func TestNewHDNativeSegWitWallet(t *testing.T) {
	m, err := NewMnemonic()
	require.NoError(t, err)

	// btc mainnet
	h, err := NewHDWallet(m, "123", 3652501241, 0)
	require.NoError(t, err)

	w, err := h.NewNativeSegWitWallet(0, 0, 0)
	require.NoError(t, err)

	fmt.Println("address: ", w.DeriveAddress())
	fmt.Println("private key: ", w.DerivePrivateKey())
	fmt.Println("public key: ", w.DerivePublicKey())
	fmt.Println("symbol: ", w.Symbol())
	fmt.Println("chain id: ", w.ChainId())
}

func TestGetChainId(t *testing.T) {
	fmt.Println(int(wire.MainNet))
}
