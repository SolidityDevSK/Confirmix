package validator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Status represents the validator status
type Status int

const (
	StatusInactive Status = iota
	StatusActive
	StatusPenalized
)

// Authority represents a validator in the PoA system
type Authority struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Status     Status
	LastActive time.Time
	BlockCount uint64
	mu         sync.RWMutex

	// Performance metrics
	MissedBlocks     uint64
	ProducedBlocks   uint64
	LastBlockTime    time.Time
	ConsecutiveMisses uint64
}

// NewAuthority creates a new authority with a keypair
func NewAuthority(privateKey *ecdsa.PrivateKey) (*Authority, error) {
	if privateKey == nil {
		var err error
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
	}

	publicKey := &privateKey.PublicKey
	address := fmt.Sprintf("%x", sha256.Sum256(elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)))

	return &Authority{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Status:     StatusActive,
		LastActive: time.Now(),
	}, nil
}

// Sign signs a message with the authority's private key
func (a *Authority) Sign(message []byte) ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	fmt.Printf("Signing message with validator %s\n", a.Address)
	fmt.Printf("Message hash: %x\n", message)

	r, s, err := ecdsa.Sign(rand.Reader, a.PrivateKey, message)
	if err != nil {
		fmt.Printf("Failed to sign message: %v\n", err)
		return nil, err
	}

	signature := make([]byte, 64)
	r.FillBytes(signature[:32])
	s.FillBytes(signature[32:])

	fmt.Printf("Generated signature: %x\n", signature)
	return signature, nil
}

// Verify verifies a signature with the authority's public key
func (a *Authority) Verify(message, signature []byte) bool {
	if len(signature) != 64 {
		fmt.Printf("Invalid signature length: %d\n", len(signature))
		return false
	}

	fmt.Printf("Verifying signature for validator %s\n", a.Address)
	fmt.Printf("Message hash: %x\n", message)
	fmt.Printf("Signature: %x\n", signature)

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	valid := ecdsa.Verify(a.PublicKey, message, r, s)
	if !valid {
		fmt.Printf("Signature verification failed for validator %s\n", a.Address)
	} else {
		fmt.Printf("Signature verified successfully for validator %s\n", a.Address)
	}

	return valid
}

// RecordBlockProduction records a successful block production
func (a *Authority) RecordBlockProduction() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ProducedBlocks++
	a.LastBlockTime = time.Now()
	a.ConsecutiveMisses = 0
}

// RecordMissedBlock records a missed block opportunity
func (a *Authority) RecordMissedBlock() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.MissedBlocks++
	a.ConsecutiveMisses++
}

// GetPerformanceMetrics returns the validator's performance metrics
func (a *Authority) GetPerformanceMetrics() (uint64, uint64, uint64) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.ProducedBlocks, a.MissedBlocks, a.ConsecutiveMisses
}

// GetStatus returns the validator's current status
func (a *Authority) GetStatus() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.Status
}

// SetStatus sets the validator's status
func (a *Authority) SetStatus(status Status) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Status = status
	if status == StatusActive {
		a.LastActive = time.Now()
	}
} 