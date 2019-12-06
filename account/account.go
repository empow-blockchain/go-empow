package account

import (
	"fmt"
	"strings"

	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/crypto"
)

// KeyPair account of the ios
type KeyPair struct {
	Algorithm crypto.Algorithm
	Pubkey    []byte
	Seckey    []byte
}

// AddressKeyPair address with keypair
type AddressKeyPair struct {
	Algorithm crypto.Algorithm
	Address   string
	Pubkey    []byte
	Seckey    []byte
}

// AddressToPubkey convert address to public key
func AddressToPubkey(address string) ([]byte, error) {
	if len(address) != 49 || !strings.HasPrefix(address, "EM") {
		return nil, fmt.Errorf("wrong address")
	}

	removePrefix := address[2:len(address)]
	pubkeyWithPrefix := common.Base58Decode(removePrefix)

	if len(pubkeyWithPrefix) < 2 {
		return nil, fmt.Errorf("wrong address")
	}

	pubkey := pubkeyWithPrefix[2:]

	return pubkey, nil
}

// PubkeyToAddress Convert PublicKey to Address
func PubkeyToAddress(pubkey []byte) string {
	emSlice := []byte("EM")

	addressSlice := append(emSlice, pubkey...)
	addressNoPrefix := common.Base58Encode(addressSlice)
	address := fmt.Sprintf("%v%v", "EM", addressNoPrefix)

	return address
}

// NewAddress to create an address
func NewAddress(seckey []byte, algo crypto.Algorithm) (*AddressKeyPair, error) {
	if seckey == nil {
		seckey = algo.GenSeckey()
	}

	err := algo.CheckSeckey(seckey)
	if err != nil {
		return nil, err
	}

	pubkey := algo.GetPubkey(seckey)

	address := PubkeyToAddress(pubkey)

	account := &AddressKeyPair{
		Algorithm: algo,
		Address:   address,
		Pubkey:    pubkey,
		Seckey:    seckey,
	}
	return account, nil
}

// NewKeyPair create an account
func NewKeyPair(seckey []byte, algo crypto.Algorithm) (*KeyPair, error) {
	if seckey == nil {
		seckey = algo.GenSeckey()
	}

	err := algo.CheckSeckey(seckey)
	if err != nil {
		return nil, err
	}

	pubkey := algo.GetPubkey(seckey)

	account := &KeyPair{
		Algorithm: algo,
		Pubkey:    pubkey,
		Seckey:    seckey,
	}
	return account, nil
}

// Sign sign a tx
func (a *KeyPair) Sign(info []byte) *crypto.Signature {
	return crypto.NewSignature(a.Algorithm, info, a.Seckey)
}

// ReadablePubkey ...
func (a *KeyPair) ReadablePubkey() string {
	return EncodePubkey(a.Pubkey)
}

// EncodePubkey ...
func EncodePubkey(pubkey []byte) string {
	return common.Base58Encode(pubkey)
}

// DecodePubkey ...
func DecodePubkey(readablePubKey string) []byte {
	return common.Base58Decode(readablePubKey)
}
