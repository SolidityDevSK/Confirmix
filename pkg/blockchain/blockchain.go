package blockchain

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

// Blockchain represents the entire chain of blocks
type Blockchain struct {
	Blocks     []*Block
	Validators map[string]*validator.Authority // ValidatorAddress -> Authority
}

// AddValidator adds a new validator to the blockchain
func (bc *Blockchain) AddValidator(authority *validator.Authority) {
	if bc.Validators == nil {
		bc.Validators = make(map[string]*validator.Authority)
	}
	bc.Validators[authority.Address] = authority
}

// AddBlock zincire yeni bir blok ekler
func (bc *Blockchain) AddBlock(data string, v *validator.Authority) error {
	// Validator'ın yetkili olup olmadığını kontrol et
	if _, exists := bc.Validators[v.Address]; !exists {
		return errors.New("unauthorized validator")
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
	fmt.Printf("Yeni blok eklendi! Validator: %s\n", v.Address[:10])
	return nil
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
	blockchain := &Blockchain{
		Blocks:     []*Block{},
		Validators: make(map[string]*validator.Authority),
	}
	
	// Genesis validator'ı ekle
	blockchain.AddValidator(genesisValidator)
	
	// Genesis bloğunu oluştur
	genesisBlock, err := NewBlock("Genesis Block", []byte{}, genesisValidator)
	if err != nil {
		return nil, err
	}
	
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
	return blockchain, nil
} 