package txpool

import (
	"github.com/empow-blockchain/go-empow/core/block"
	"github.com/empow-blockchain/go-empow/core/blockcache"
	"github.com/empow-blockchain/go-empow/core/tx"
)

//go:generate mockgen -destination mock/mock_txpool.go -package txpool_mock github.com/empow-blockchain/go-empow/core/txpool TxPool

// TxPool defines all the API of txpool package.
type TxPool interface {
	Close()
	AddTx(tx *tx.Tx, from string) error
	DelTx(hash []byte) error
	GetFromPending(hash []byte) (*tx.Tx, error)
	PendingTx() (*SortedTxMap, *blockcache.BlockCacheNode)

	// TODO: The following interfaces need to be moved from txpool to chainbase.
	AddLinkedNode(linkedNode *blockcache.BlockCacheNode) error
	ExistTxs(hash []byte, chainBlock *block.Block) bool
	GetFromChain(hash []byte) (*tx.Tx, *tx.TxReceipt, error)
}
