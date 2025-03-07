package consensus

import (
	"errors"
	"sync"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

// RoundRobin represents the round-robin consensus mechanism
type RoundRobin struct {
	validators     []*validator.Authority
	currentIndex   int
	mu            sync.RWMutex
	blockInterval time.Duration // Her blok arasındaki minimum süre
}

// NewRoundRobin creates a new round-robin consensus instance
func NewRoundRobin(blockInterval time.Duration) *RoundRobin {
	return &RoundRobin{
		validators:     make([]*validator.Authority, 0),
		currentIndex:   0,
		blockInterval: blockInterval,
	}
}

// AddValidator adds a new validator to the round-robin rotation
func (rr *RoundRobin) AddValidator(v *validator.Authority) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	// Validator zaten ekli mi kontrol et
	for _, existing := range rr.validators {
		if existing.Address == v.Address {
			return
		}
	}
	
	rr.validators = append(rr.validators, v)
}

// RemoveValidator removes a validator from the round-robin rotation
func (rr *RoundRobin) RemoveValidator(address string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	for i, v := range rr.validators {
		if v.Address == address {
			rr.validators = append(rr.validators[:i], rr.validators[i+1:]...)
			// Eğer silinen validator current index'ten önceyse, index'i güncelle
			if i < rr.currentIndex {
				rr.currentIndex--
			}
			// Eğer son validator siliniyorsa ve current index son indexse
			if rr.currentIndex >= len(rr.validators) {
				rr.currentIndex = 0
			}
			return
		}
	}
}

// GetCurrentValidator returns the current validator in the rotation
func (rr *RoundRobin) GetCurrentValidator() (*validator.Authority, error) {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	if len(rr.validators) == 0 {
		return nil, errors.New("no validators available")
	}
	
	return rr.validators[rr.currentIndex], nil
}

// NextValidator moves to the next validator in the rotation
func (rr *RoundRobin) NextValidator() (*validator.Authority, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	if len(rr.validators) == 0 {
		return nil, errors.New("no validators available")
	}
	
	rr.currentIndex = (rr.currentIndex + 1) % len(rr.validators)
	return rr.validators[rr.currentIndex], nil
}

// IsValidatorTurn checks if it's the given validator's turn
func (rr *RoundRobin) IsValidatorTurn(address string) bool {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	if len(rr.validators) == 0 {
		return false
	}
	
	return rr.validators[rr.currentIndex].Address == address
}

// GetValidators returns all validators in the rotation
func (rr *RoundRobin) GetValidators() []*validator.Authority {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	validators := make([]*validator.Authority, len(rr.validators))
	copy(validators, rr.validators)
	return validators
}

// GetBlockInterval returns the minimum time between blocks
func (rr *RoundRobin) GetBlockInterval() time.Duration {
	return rr.blockInterval
} 