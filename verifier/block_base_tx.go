package verifier

import (
	"fmt"
	"github.com/empow-blockchain/go-empow/common"

	"github.com/empow-blockchain/go-empow/core/block"
	"github.com/empow-blockchain/go-empow/core/blockcache"
	"github.com/empow-blockchain/go-empow/core/tx"
)

// NewBaseTx is new baseTx
func NewBaseTx(blk, parent *block.Block, witnessList *blockcache.WitnessList) (*tx.Tx, error) {
	acts := []*tx.Action{}
	if blk.Head.Number > 0 {
		txData, err := baseTxData(blk, parent, witnessList)
		if err != nil {
			return nil, err
		}
		act := tx.NewAction("base.empow", "exec", txData)
		acts = append(acts, act)
	}
	tx := &tx.Tx{
		Publisher: "base.empow",
		GasLimit:  100000000,
		GasRatio:  100,
		Actions:   acts,
		Time:      blk.Head.Time,
		ChainID:   tx.ChainID,
	}
	return tx, nil
}

func baseTxData(b, pb *block.Block, witnessList *blockcache.WitnessList) (string, error) {
	if pb != nil {
		witnessChanged := false
		if witnessList != nil && !common.StringSliceEqual(witnessList.Active(), witnessList.Pending()) {
			witnessChanged = true
		}
		return fmt.Sprintf(`[{"parent":["%v", "%v", %v]}]`, pb.Head.Witness, pb.CalculateGasUsage(), witnessChanged), nil
	}
	if b.Head.Number != 0 {
		return "", fmt.Errorf("block dit not have parent")
	}
	return `[{"parent":["", "0", false]}]`, nil
}
