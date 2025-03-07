package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/consensus"
	"github.com/SolidityDevSK/Confirmix/pkg/contracts"
	"github.com/SolidityDevSK/Confirmix/pkg/common"
)

// Blockchain represents the entire chain of blocks
type Blockchain struct {
	Blocks          []*Block
	PendingBlock    *Block
	Validators      []*validator.Authority
	CurrentIndex    int
	LastBlockTime   time.Time
	ContractManager *contracts.Manager
	mu             sync.RWMutex
	consensus       *consensus.RoundRobin
}

// AddValidator adds a new validator to the blockchain
func (bc *Blockchain) AddValidator(authority *validator.Authority) {
	bc.Validators = append(bc.Validators, authority)
	bc.consensus.AddValidator(authority)
}

// RemoveValidator removes a validator from the blockchain
func (bc *Blockchain) RemoveValidator(address string) {
	for i, validator := range bc.Validators {
		if validator.Address == address {
			bc.Validators = append(bc.Validators[:i], bc.Validators[i+1:]...)
			bc.consensus.RemoveValidator(address)
			break
		}
	}
}

// AddBlock zincire yeni bir blok ekler
func (bc *Blockchain) AddBlock(data string, v *validator.Authority) error {
	// Validator'ın yetkili olup olmadığını kontrol et
	if !bc.consensus.IsValidatorTurn(v.Address) {
		return fmt.Errorf("not validator's turn: %s", v.Address[:10])
	}

	// Bloklar arası minimum süre kontrolü
	if time.Since(bc.LastBlockTime) < bc.consensus.GetBlockInterval() {
		return fmt.Errorf("minimum block interval not reached, wait %v", bc.consensus.GetBlockInterval()-time.Since(bc.LastBlockTime))
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock, err := NewBlock(data, prevBlock.Hash, v)
	if err != nil {
		return err
	}

	// İmzayı doğrula
	if !v.Verify(newBlock.Hash, newBlock.Signature) {
		return errors.New("invalid block signature")
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.LastBlockTime = time.Now()

	// Sıradaki validator'a geç
	_, err = bc.consensus.NextValidator()
	if err != nil {
		return fmt.Errorf("failed to move to next validator: %v", err)
	}

	fmt.Printf("Yeni blok eklendi! Validator: %s\n", v.Address[:10])
	return nil
}

// GetCurrentValidator returns the current validator in the rotation
func (bc *Blockchain) GetCurrentValidator() (*validator.Authority, error) {
	return bc.consensus.GetCurrentValidator()
}

// IsValid zincirin geçerli olup olmadığını kontrol eder
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// Hash'leri kontrol et
		if !bytes.Equal(currentBlock.Hash, currentBlock.CalculateHash()) {
			return false
		}

		// Önceki blok bağlantısını kontrol et
		if !bytes.Equal(currentBlock.PrevHash, previousBlock.Hash) {
			return false
		}

		// Validator'ı kontrol et
		v, exists := bc.Validators[currentBlock.ValidatorAddress]
		if !exists {
			return false
		}

		// İmzayı kontrol et
		if !v.Verify(currentBlock.Hash, currentBlock.Signature) {
			return false
		}
	}
	return true
}

// NewBlockchain yeni bir blockchain oluşturur
func NewBlockchain(genesisValidator *validator.Authority) (*Blockchain, error) {
	// Contract manager'ı oluştur
	contractManager := contracts.NewManager()

	bc := &Blockchain{
		Blocks:          make([]*Block, 0),
		Validators:      []*validator.Authority{genesisValidator},
		CurrentIndex:    0,
		LastBlockTime:   time.Now(),
		ContractManager: contractManager,
	}

	// Round-Robin konsensüs mekanizmasını oluştur (5 saniyelik blok aralığı)
	roundRobin := consensus.NewRoundRobin(5 * time.Second)
	bc.consensus = roundRobin

	// Genesis bloğunu oluştur
	genesisBlock := &Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		Data:      "Genesis Block",
		PrevHash:  "",
		Validator: genesisValidator.Address,
	}

	// Genesis bloğunu imzala ve ekle
	if err := bc.signAndAddBlock(genesisBlock, genesisValidator); err != nil {
		return nil, err
	}

	return bc, nil
}

// DeployContract deploys a new smart contract
func (bc *Blockchain) DeployContract(code []byte, owner common.Address, name, version string) (*contracts.Contract, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	return bc.ContractManager.DeployContract(code, owner, name, version, time.Now().Unix())
}

// ExecuteContract executes a smart contract
func (bc *Blockchain) ExecuteContract(address common.Address, input []byte) ([]byte, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.ContractManager.ExecuteContract(address, input)
}

// GetContract returns a contract by address
func (bc *Blockchain) GetContract(address common.Address) (*contracts.Contract, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.ContractManager.GetContract(address)
}

// ListContracts returns all contracts
func (bc *Blockchain) ListContracts() []*contracts.Contract {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	return bc.ContractManager.ListContracts()
}

// signAndAddBlock is a helper function to sign and add a block to the blockchain
func (bc *Blockchain) signAndAddBlock(block *Block, validator *validator.Authority) error {
	// Sign the block
	if err := block.Sign(validator); err != nil {
		return err
	}

	// Add the block to the blockchain
	bc.Blocks = append(bc.Blocks, block)
	bc.LastBlockTime = time.Now()

	// Sıradaki validator'a geç
	_, err := bc.consensus.NextValidator()
	if err != nil {
		return fmt.Errorf("failed to move to next validator: %v", err)
	}

	fmt.Printf("Yeni blok eklendi! Validator: %s\n", validator.Address[:10])
	return nil
} 