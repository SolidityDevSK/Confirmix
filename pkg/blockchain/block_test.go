package blockchain

import (
	"math/big"
	"testing"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

func TestNewBlock(t *testing.T) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}

	v, err := validator.NewAuthority(key)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	prevHash := make([]byte, 32)
	stateRoot := make([]byte, 32)
	height := uint64(1)
	gasLimit := uint64(1000000)

	block, err := NewBlock(height, prevHash, stateRoot, gasLimit, v)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Check block fields
	if block.Header.Height != height {
		t.Errorf("Expected height %d, got %d", height, block.Header.Height)
	}
	if block.Header.GasLimit != gasLimit {
		t.Errorf("Expected gas limit %d, got %d", gasLimit, block.Header.GasLimit)
	}
	if block.Header.ValidatorAddress != v.Address {
		t.Errorf("Expected validator %s, got %s", v.Address, block.Header.ValidatorAddress)
	}
	if len(block.Transactions) != 0 {
		t.Error("Expected empty transaction list")
	}
	if !block.Verify(v) {
		t.Error("Block signature verification failed")
	}
}

func TestAddTransaction(t *testing.T) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}

	v, err := validator.NewAuthority(key)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	block, err := NewBlock(1, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create test transaction
	tx := &Transaction{
		Hash:      make([]byte, 32),
		From:      "0x1234",
		To:        "0x5678",
		Value:     big.NewInt(1000),
		GasPrice:  1000,
		GasLimit:  21000,
		Nonce:     1,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}

	// Add transaction
	if err := block.AddTransaction(tx); err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	// Check transaction was added
	if len(block.Transactions) != 1 {
		t.Error("Transaction was not added to block")
	}

	// Check gas used was updated
	if block.Header.GasUsed != tx.GasLimit {
		t.Errorf("Expected gas used %d, got %d", tx.GasLimit, block.Header.GasUsed)
	}

	// Try to add transaction that exceeds gas limit
	tx2 := &Transaction{
		GasLimit: block.Header.GasLimit,
	}
	if err := block.AddTransaction(tx2); err == nil {
		t.Error("Expected error when exceeding gas limit")
	}
}

func TestBlockHashing(t *testing.T) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}

	v, err := validator.NewAuthority(key)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	block, err := NewBlock(1, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Get initial hash
	hash1 := block.GetHash()

	// Add transaction and check hash changes
	tx := &Transaction{
		Hash:     make([]byte, 32),
		GasLimit: 21000,
		Value:    big.NewInt(0),
	}
	if err := block.AddTransaction(tx); err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	hash2 := block.GetHash()
	if string(hash1) == string(hash2) {
		t.Error("Block hash did not change after adding transaction")
	}
}

func TestBlockSize(t *testing.T) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}

	v, err := validator.NewAuthority(key)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	block, err := NewBlock(1, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	initialSize := block.GetBlockSize()

	// Add transaction and check size increases
	tx := &Transaction{
		Hash:     make([]byte, 32),
		From:     "0x1234",
		To:       "0x5678",
		Value:    big.NewInt(1000),
		GasLimit: 21000,
	}
	if err := block.AddTransaction(tx); err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	newSize := block.GetBlockSize()
	if newSize <= initialSize {
		t.Error("Block size did not increase after adding transaction")
	}
}

func TestTransactionLookup(t *testing.T) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}

	v, err := validator.NewAuthority(key)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	block, err := NewBlock(1, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Add transaction
	txHash := make([]byte, 32)
	txHash[0] = 1 // Make it unique
	tx := &Transaction{
		Hash:     txHash,
		GasLimit: 21000,
		Value:    big.NewInt(0), // Initialize Value field
	}
	if err := block.AddTransaction(tx); err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}

	// Look up transaction
	found := block.GetTransactionByHash(txHash)
	if found == nil {
		t.Error("Failed to find transaction by hash")
	}

	// Look up non-existent transaction
	notFound := block.GetTransactionByHash(make([]byte, 32))
	if notFound != nil {
		t.Error("Found non-existent transaction")
	}
} 