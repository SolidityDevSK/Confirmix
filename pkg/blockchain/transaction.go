package blockchain

import (
	"math/big"
)

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash      []byte
	From      string
	To        string
	Value     *big.Int
	Data      []byte
	GasPrice  uint64
	GasLimit  uint64
	GasUsed   uint64
	Nonce     uint64
	Signature []byte
	Status    TxStatus
}

// TxStatus represents the status of a transaction
type TxStatus int

const (
	TxPending TxStatus = iota
	TxSuccess
	TxFailed
)

// GetSize işlemin yaklaşık boyutunu hesaplar
func (tx *Transaction) GetSize() uint64 {
	size := uint64(0)
	
	// Temel alanların boyutları
	size += uint64(len(tx.Hash))       // Hash boyutu
	size += uint64(len(tx.From))       // From adresi
	size += uint64(len(tx.To))         // To adresi
	size += 32                         // Value (big.Int)
	size += 8                          // GasPrice
	size += 8                          // GasLimit
	size += 8                          // GasUsed
	size += 8                          // Nonce
	size += uint64(len(tx.Signature))  // İmza boyutu
	size += uint64(len(tx.Data))       // Data boyutu
	
	return size
} 