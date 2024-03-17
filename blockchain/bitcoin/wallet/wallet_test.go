package wallet

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	w := NewWallet()
	addr := w.GetAddress()
	fmt.Println(addr)
}
