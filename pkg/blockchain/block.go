package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

// Block represents each block in the blockchain
type Block struct {
	Timestamp        int64  // bloğun oluşturulma zamanı
	Data            []byte // blok içindeki veriler
	PrevHash        []byte // önceki bloğun hash'i
	Hash            []byte // mevcut bloğun hash'i
	Signature       []byte // validator imzası
	ValidatorAddress string // blok oluşturan validator'ın adresi
}

// CalculateHash bloğun hash'ini hesaplar
func (b *Block) CalculateHash() []byte {
	data := append([]byte{}, b.PrevHash...)
	data = append(data, b.Data...)
	timestamp := []byte(string(b.Timestamp))
	data = append(data, timestamp...)
	
	hash := sha256.Sum256(data)
	return hash[:]
}

// NewBlock yeni bir blok oluşturur
func NewBlock(data string, prevHash []byte, v *validator.Authority) (*Block, error) {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Data:     []byte(data),
		PrevHash: prevHash,
		ValidatorAddress: v.Address,
	}
	
	block.Hash = block.CalculateHash()
	
	// Bloğu validator ile imzala
	signature, err := v.Sign(block.Hash)
	if err != nil {
		return nil, err
	}
	block.Signature = signature
	
	return block, nil
}

// GetHashString bloğun hash'ini string olarak döndürür
func (b *Block) GetHashString() string {
	return hex.EncodeToString(b.Hash)
} 