package rpc

import (
	"encoding/json"

	"github.com/empow-blockchain/go-empow/account"
	"github.com/empow-blockchain/go-empow/common"
	"github.com/empow-blockchain/go-empow/core/block"
	"github.com/empow-blockchain/go-empow/core/contract"
	"github.com/empow-blockchain/go-empow/core/tx"
	"github.com/empow-blockchain/go-empow/crypto"
	rpcpb "github.com/empow-blockchain/go-empow/rpc/pb"
	"github.com/empow-blockchain/go-empow/verifier"
)

func toPbAction(a *tx.Action) *rpcpb.Action {
	return &rpcpb.Action{
		Contract:   a.Contract,
		ActionName: a.ActionName,
		Data:       a.Data,
	}
}

func toPbTxReceipt(tr *tx.TxReceipt) *rpcpb.TxReceipt {
	if tr == nil {
		return nil
	}
	ret := &rpcpb.TxReceipt{
		TxHash:     common.Base58Encode(tr.TxHash),
		GasUsage:   float64(tr.GasUsage) / 100,
		RamUsage:   tr.RAMUsage,
		StatusCode: rpcpb.TxReceipt_StatusCode(tr.Status.Code),
		Message:    tr.Status.Message,
		Returns:    tr.Returns,
	}
	for _, r := range tr.Receipts {
		ret.Receipts = append(ret.Receipts, &rpcpb.TxReceipt_Receipt{
			FuncName: r.FuncName,
			Content:  r.Content,
		})
	}
	return ret
}

func toPbAmountLimit(a *contract.Amount) *rpcpb.AmountLimit {
	return &rpcpb.AmountLimit{
		Token: a.Token,
		Value: a.Val,
	}
}

func toPbTx(t *tx.Tx, tr *tx.TxReceipt) *rpcpb.Transaction {
	ret := &rpcpb.Transaction{
		Hash:       common.Base58Encode(t.Hash()),
		Time:       t.Time,
		Expiration: t.Expiration,
		GasRatio:   float64(t.GasRatio) / 100,
		GasLimit:   float64(t.GasLimit) / 100,
		Delay:      t.Delay,
		ChainId:    t.ChainID,
		Signers:    t.Signers,
		Publisher:  t.Publisher,
		ReferredTx: common.Base58Encode(t.ReferredTx),
		TxReceipt:  toPbTxReceipt(tr),
	}
	for _, a := range t.Actions {
		ret.Actions = append(ret.Actions, toPbAction(a))
	}
	for _, a := range t.AmountLimit {
		ret.AmountLimit = append(ret.AmountLimit, toPbAmountLimit(a))
	}
	return ret
}

func toPbBlock(blk *block.Block, complete bool) *rpcpb.Block {
	ret := &rpcpb.Block{
		Hash:                common.Base58Encode(blk.HeadHash()),
		Version:             blk.Head.Version,
		ParentHash:          common.Base58Encode(blk.Head.ParentHash),
		TxMerkleHash:        common.Base58Encode(blk.Head.TxMerkleHash),
		TxReceiptMerkleHash: common.Base58Encode(blk.Head.TxReceiptMerkleHash),
		Number:              blk.Head.Number,
		Witness:             blk.Head.Witness,
		Time:                blk.Head.Time,
		GasUsage:            float64(blk.Head.GasUsage) / 100,
		TxCount:             int64(len(blk.Txs)),
	}
	var info verifier.Info
	json.Unmarshal(blk.Head.Info, &info)
	ret.Info = &rpcpb.Block_Info{
		Mode:   int32(info.Mode),
		Thread: int32(info.Thread),
	}
	for _, i := range info.Batch {
		ret.Info.BatchIndex = append(ret.Info.BatchIndex, int32(i))
	}
	if complete {
		for i, t := range blk.Txs {
			ret.Transactions = append(ret.Transactions, toPbTx(t, blk.Receipts[i]))
		}
	}
	return ret
}

func toPbAccount(acc *account.Account) *rpcpb.Account {
	ret := &rpcpb.Account{
		Address: acc.Address,
	}
	return ret
}

func toPbContract(c *contract.Contract) *rpcpb.Contract {
	ret := &rpcpb.Contract{
		Id:       c.ID,
		Code:     c.Code,
		Language: c.Info.Lang,
		Version:  c.Info.Version,
	}
	for _, abi := range c.Info.Abi {
		pbABI := &rpcpb.Contract_ABI{
			Name: abi.Name,
			Args: abi.Args,
		}
		for _, al := range abi.AmountLimit {
			pbABI.AmountLimit = append(pbABI.AmountLimit, toPbAmountLimit(al))
		}
		ret.Abis = append(ret.Abis, pbABI)
	}
	return ret
}

func toCoreTx(t *rpcpb.TransactionRequest) *tx.Tx {
	ret := &tx.Tx{
		Time:       t.Time,
		Expiration: t.Expiration,
		GasRatio:   int64(t.GasRatio * 100),
		GasLimit:   int64(t.GasLimit * 100),
		Delay:      t.Delay,
		ChainID:    t.ChainId,
		Signers:    t.Signers,
		Publisher:  t.Publisher,
	}
	for _, a := range t.Actions {
		ret.Actions = append(ret.Actions, &tx.Action{
			Contract:   a.Contract,
			ActionName: a.ActionName,
			Data:       a.Data,
		})
	}
	for _, a := range t.AmountLimit {
		ret.AmountLimit = append(ret.AmountLimit, &contract.Amount{
			Token: a.Token,
			Val:   a.Value,
		})
	}
	for _, s := range t.Signatures {
		ret.Signs = append(ret.Signs, &crypto.Signature{
			Algorithm: crypto.Algorithm(s.Algorithm),
			Pubkey:    s.PublicKey,
			Sig:       s.Signature,
		})
	}
	for _, s := range t.PublisherSigs {
		ret.PublishSigns = append(ret.PublishSigns, &crypto.Signature{
			Algorithm: crypto.Algorithm(s.Algorithm),
			Pubkey:    s.PublicKey,
			Sig:       s.Signature,
		})
	}
	return ret
}
