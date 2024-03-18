package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	// p2pkh
	versionP2pkh = byte(00)
	// // p2sh
	// versionP2sh = byte(05)
	// // p2pkh testnet
	// versionP2pkhTestnet = byte(6f)
	// // p2sh testnet
	// versionP2shTestnet = byte(c4)

	addressChecksumLen = 4
)

type LegacyWallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewLegacyWallet() Wallet {
	// generate private key
	private, public := newKeyPair()
	wallet := LegacyWallet{private, public}

	return &wallet
}

func (w *LegacyWallet) ChainId() int {
	return 0
}

func (w *LegacyWallet) Symbol() string {
	return w.PrivateKey.D.String()
}

func (w *LegacyWallet) DerivePrivateKey() string {
	return w.PrivateKey.D.String()
}

func (w *LegacyWallet) DerivePublicKey() string {
	return hex.EncodeToString(w.PublicKey)
}

func (w *LegacyWallet) DeriveAddress() string {
	var (
		versionedPayload = make([]byte, 0)
		format           string
	)

	fmt.Print("输入地址格式(p2pkh,p2sh):")
	fmt.Scanln(&format)

	switch format {
	case p2pkh:
		// P2PKH address format
		pubKeyHash := hashPubKey(w.PublicKey)
		versionedPayload = append([]byte{versionP2pkh}, pubKeyHash...)
	case p2sh:
	default:
		panic("unsupported address format")
	}

	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return address
}

func hashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, _ = RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// get a curve, 256 here
	curve := elliptic.P256()
	// get the private ket based on the curve
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	// public key
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
