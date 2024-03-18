package wallet

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	w := NewLegacyWallet()
	addr := w.DeriveAddress()
	fmt.Println(addr)
}
