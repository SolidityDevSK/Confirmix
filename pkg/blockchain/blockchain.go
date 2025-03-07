package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/consensus"
)

// Blockchain represents the entire chain of blocks
type Blockchain struct {
	Blocks     []*Block
	Validators map[string]*validator.Authority // ValidatorAddress -> Authority
	consensus  *consensus.RoundRobin
	lastBlockTime time.Time
}

// AddValidator adds a new validator to the blockchain
func (bc *Blockchain) AddValidator(authority *validator.Authority) {
	if bc.Validators == nil {
		bc.Validators = make(map[string]*validator.Authority)
	}
	bc.Validators[authority.Address] = authority
	bc.consensus.AddValidator(authority)
}

// RemoveValidator removes a validator from the blockchain
func (bc *Blockchain) RemoveValidator(address string) {
	delete(bc.Validators, address)
	bc.consensus.RemoveValidator(address)
}

// AddBlock zincire yeni bir blok ekler
func (bc *Blockchain) AddBlock(data string, v *validator.Authority) error {
	// Validator'ın yetkili olup olmadığını kontrol et
	if _, exists := bc.Validators[v.Address]; !exists {
		return errors.New("unauthorized validator")
	}

	// Validator sırası kontrolü
	if !bc.consensus.IsValidatorTurn(v.Address) {
		return fmt.Errorf("not validator's turn: %s", v.Address[:10])
	}

	// Bloklar arası minimum süre kontrolü
	if time.Since(bc.lastBlockTime) < bc.consensus.GetBlockInterval() {
		return fmt.Errorf("minimum block interval not reached, wait %v", bc.consensus.GetBlockInterval()-time.Since(bc.lastBlockTime))
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
	bc.lastBlockTime = time.Now()

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
	// Round-Robin konsensüs mekanizmasını oluştur (5 saniyelik blok aralığı)
	roundRobin := consensus.NewRoundRobin(5 * time.Second)
	
	blockchain := &Blockchain{
		Blocks:     []*Block{},
		Validators: make(map[string]*validator.Authority),
		consensus:  roundRobin,
	}
	
	// Genesis validator'ı ekle
	blockchain.AddValidator(genesisValidator)
	
	// Genesis bloğunu oluştur
	genesisBlock, err := NewBlock("Genesis Block", []byte{}, genesisValidator)
	if err != nil {
		return nil, err
	}
	
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
	blockchain.lastBlockTime = time.Now()
	return blockchain, nil
} 