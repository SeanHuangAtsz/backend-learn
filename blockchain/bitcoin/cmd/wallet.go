package main

import (
	"fmt"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/bitcoin/wallet"
)

func main() {
	w := wallet.NewLegacyWallet()
	fmt.Println(w.DeriveAddress())
}
