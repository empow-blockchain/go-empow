package genesis

import (
	"fmt"
	"os"
	"testing"

	"github.com/empow-blockchain/go-empow/account"
	"github.com/empow-blockchain/go-empow/crypto"

	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/db"
	"github.com/empow-blockchain/go-empow/ilog"
)

func randWitness(idx int) *common.Witness {
	k := account.EncodePubkey(crypto.Ed25519.GetPubkey(crypto.Ed25519.GenSeckey()))
	return &common.Witness{Address: account.PubkeyToAddress(common.Base58Decode(k)), Owner: k, Active: k, SignatureBlock: k, Balance: 3 * 1e8}
}

func TestGenGenesis(t *testing.T) {
	ilog.Stop()

	d, err := db.NewMVCCDB("mvcc")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		d.Close()
		os.RemoveAll("mvcc")
	}()
	k := account.EncodePubkey(crypto.Ed25519.GetPubkey(crypto.Ed25519.GenSeckey()))
	blk, err := GenGenesis(d, &common.GenesisConfig{
		WitnessInfo: []*common.Witness{
			randWitness(1),
			randWitness(2),
			randWitness(3),
			randWitness(4),
			randWitness(5),
			randWitness(6),
			randWitness(7),
		},
		TokenInfo: &common.TokenInfo{
			FoundationAccount: account.PubkeyToAddress(common.Base58Decode(k)),
			EMTotalSupply:   90000000000,
			EMDecimal:       8,
		},
		InitialTimestamp: "2006-01-02T15:04:05Z",
		ContractPath:     os.Getenv("GOPATH") + "/src/github.com/empow-blockchain/go-empow/config/genesis/contract/",
		AdminInfo:        randWitness(8),
		FoundationInfo:   &common.Witness{Address: account.PubkeyToAddress(common.Base58Decode(k)), Owner: k, Active: k, Balance: 0},
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(blk)
	return
}
