package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

// Header represents the block header
type Header struct {
	Version         uint32    // Protocol version
	Timestamp       time.Time // Block creation time
	PrevHash        []byte    // Previous block hash
	Height          uint64    // Block height
	StateRoot       []byte    // State trie root hash
	TransactionRoot []byte    // Transaction trie root hash
	ReceiptRoot     []byte    // Receipt trie root hash
	GasLimit        uint64    // Block gas limit
	GasUsed         uint64    // Total gas used by transactions
	ValidatorAddress string    // Block producer address
}

// Block represents a block in the blockchain
type Block struct {
	Header       *Header
	Transactions []*Transaction
	Signature    []byte // Validator's signature of the block header

	// Cached values
	hash []byte
}

// NewBlock creates a new block
func NewBlock(
	height uint64,
	prevHash []byte,
	stateRoot []byte,
	gasLimit uint64,
	v *validator.Authority,
) (*Block, error) {
	header := &Header{
		Version:          1,
		Timestamp:        time.Now().UTC(),
		PrevHash:        prevHash,
		Height:          height,
		StateRoot:       stateRoot,
		GasLimit:        gasLimit,
		ValidatorAddress: v.Address,
	}

	block := &Block{
		Header:       header,
		Transactions: make([]*Transaction, 0),
	}

	// Calculate initial transaction root
	txRoot, err := block.calculateTransactionRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate transaction root: %v", err)
	}
	header.TransactionRoot = txRoot

	// Sign the block
	if err := block.Sign(v); err != nil {
		return nil, fmt.Errorf("failed to sign block: %v", err)
	}

	return block, nil
}

// AddTransaction adds a transaction to the block
func (b *Block) AddTransaction(tx *Transaction) error {
	// Gas limit kontrolü
	if tx.GasLimit > b.Header.GasLimit {
		return fmt.Errorf("transaction gas limit %d exceeds block gas limit %d", tx.GasLimit, b.Header.GasLimit)
	}

	// Toplam gas kullanımı kontrolü
	if b.Header.GasUsed+tx.GasLimit > b.Header.GasLimit {
		return fmt.Errorf("total gas used would exceed block gas limit")
	}

	// Negatif değer kontrolü
	if tx.Value == nil {
		return fmt.Errorf("transaction value cannot be nil")
	}
	if tx.Value.Sign() < 0 {
		return fmt.Errorf("transaction value cannot be negative")
	}

	// Aynı nonce kontrolü
	for _, transaction := range b.Transactions {
		if transaction.From == tx.From && transaction.Nonce == tx.Nonce {
			return fmt.Errorf("transaction with same nonce already exists")
		}
	}

	// Gas kullanımını güncelle
	tx.GasUsed = tx.GasLimit
	b.Header.GasUsed += tx.GasUsed

	// İşlemi ekle
	b.Transactions = append(b.Transactions, tx)

	// Transaction root'u güncelle
	txRoot, err := b.calculateTransactionRoot()
	if err != nil {
		return fmt.Errorf("failed to calculate transaction root: %v", err)
	}
	b.Header.TransactionRoot = txRoot

	// Hash'i sıfırla (yeni transaction eklendi)
	b.hash = nil

	return nil
}

// GetHash returns the block hash, calculating it if necessary
func (b *Block) GetHash() []byte {
	if b.hash == nil {
		b.hash = b.calculateHash()
	}
	return b.hash
}

// GetHashString returns the block hash as a hex string
func (b *Block) GetHashString() string {
	return hex.EncodeToString(b.GetHash())
}

// calculateHash calculates the block hash
func (b *Block) calculateHash() []byte {
	headerHash := b.calculateHeaderHash()
	blockHash := sha256.Sum256(append(headerHash, b.Signature...))
	return blockHash[:]
}

// calculateHeaderHash calculates the hash of the block header
func (b *Block) calculateHeaderHash() []byte {
	// Ensure header is not nil
	if b.Header == nil {
		// Return empty hash if header is nil
		empty := sha256.Sum256(nil)
		return empty[:]
	}

	// Ensure all byte slices are initialized
	if b.Header.PrevHash == nil {
		b.Header.PrevHash = make([]byte, 0)
	}
	if b.Header.StateRoot == nil {
		b.Header.StateRoot = make([]byte, 0)
	}
	if b.Header.TransactionRoot == nil {
		b.Header.TransactionRoot = make([]byte, 0)
	}
	if b.Header.ReceiptRoot == nil {
		b.Header.ReceiptRoot = make([]byte, 0)
	}

	headerData, err := json.Marshal(b.Header)
	if err != nil {
		// Log error but return empty hash instead of panicking
		fmt.Printf("Error marshaling header: %v\n", err)
		empty := sha256.Sum256(nil)
		return empty[:]
	}
	
	headerHash := sha256.Sum256(headerData)
	return headerHash[:]
}

// calculateTransactionRoot calculates the merkle root of transactions
func (b *Block) calculateTransactionRoot() ([]byte, error) {
	if len(b.Transactions) == 0 {
		empty := sha256.Sum256(nil)
		return empty[:], nil
	}

	var txHashes [][]byte
	for _, tx := range b.Transactions {
		// Transaction'ın tüm alanlarını hash'e dahil et
		txData, err := json.Marshal(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal transaction: %v", err)
		}
		txHash := sha256.Sum256(txData)
		txHashes = append(txHashes, txHash[:])
	}

	// Merkle ağacı oluştur
	for len(txHashes) > 1 {
		if len(txHashes)%2 == 1 {
			txHashes = append(txHashes, txHashes[len(txHashes)-1])
		}
		var temp [][]byte
		for i := 0; i < len(txHashes); i += 2 {
			combined := append(txHashes[i], txHashes[i+1]...)
			hash := sha256.Sum256(combined)
			temp = append(temp, hash[:])
		}
		txHashes = temp
	}

	return txHashes[0], nil
}

// Verify verifies the block's signature
func (b *Block) Verify(v *validator.Authority) bool {
	if v == nil {
		fmt.Printf("Validator is nil for block %d\n", b.Header.Height)
		return false
	}

	// Verify validator address
	if b.Header.ValidatorAddress != v.Address {
		fmt.Printf("Validator address mismatch for block %d. Expected: %s, Got: %s\n", 
			b.Header.Height, v.Address, b.Header.ValidatorAddress)
		return false
	}

	// Create a minimal header for verification
	verifyHeader := &Header{
		Height:          b.Header.Height,
		ValidatorAddress: b.Header.ValidatorAddress,
	}
	headerData, err := json.Marshal(verifyHeader)
	if err != nil {
		fmt.Printf("Failed to marshal verification header for block %d: %v\n", b.Header.Height, err)
		return false
	}
	
	headerHash := sha256.Sum256(headerData)
	fmt.Printf("Verifying signature for block %d with header hash: %x\n", b.Header.Height, headerHash)

	// Verify signature
	if !v.Verify(headerHash[:], b.Signature) {
		fmt.Printf("Signature verification failed for block %d\n", b.Header.Height)
		return false
	}

	fmt.Printf("Block signature verified for block %d\n", b.Header.Height)
	return true
}

// Sign signs the block with the given validator
func (b *Block) Sign(v *validator.Authority) error {
	if v == nil {
		return errors.New("validator is nil")
	}

	fmt.Printf("Signing block %d with validator %s\n", b.Header.Height, v.Address)

	// Set validator address
	b.Header.ValidatorAddress = v.Address

	// Create a minimal header for signing
	signHeader := &Header{
		Height:          b.Header.Height,
		ValidatorAddress: b.Header.ValidatorAddress,
	}
	headerData, err := json.Marshal(signHeader)
	if err != nil {
		return fmt.Errorf("failed to marshal signing header: %v", err)
	}
	
	headerHash := sha256.Sum256(headerData)
	fmt.Printf("Calculated header hash for block %d: %x\n", b.Header.Height, headerHash)

	// Sign header hash
	signature, err := v.Sign(headerHash[:])
	if err != nil {
		fmt.Printf("Failed to sign block %d: %v\n", b.Header.Height, err)
		return fmt.Errorf("failed to sign block: %v", err)
	}

	b.Signature = signature
	fmt.Printf("Block %d signed successfully\n", b.Header.Height)
	return nil
}

// GetTransactionByHash returns a transaction by its hash
func (b *Block) GetTransactionByHash(hash []byte) *Transaction {
	for _, tx := range b.Transactions {
		if hex.EncodeToString(tx.Hash) == hex.EncodeToString(hash) {
			return tx
		}
	}
	return nil
}

// GetBlockSize returns the approximate size of the block in bytes
func (b *Block) GetBlockSize() uint64 {
	size := uint64(0)
	
	// Header size
	headerData, _ := json.Marshal(b.Header)
	size += uint64(len(headerData))
	
	// Transactions size
	for _, tx := range b.Transactions {
		txData, _ := json.Marshal(tx)
		size += uint64(len(txData))
	}
	
	// Signature size
	size += uint64(len(b.Signature))
	
	return size
} 