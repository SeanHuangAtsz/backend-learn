package main

import (
	"fmt"

	"github.com/SeanHuangAtsz/backend-learn/blockchain/bitcoin/wallet"
)

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.GetAddress())
}
