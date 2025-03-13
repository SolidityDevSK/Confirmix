package consensus

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

var (
	ErrNoValidators      = errors.New("no validators available")
	ErrValidatorInactive = errors.New("validator is not active")
	ErrNotValidatorTurn  = errors.New("not validator's turn")
)

// ConsensusConfig represents the configuration for the consensus mechanism
type ConsensusConfig struct {
	BlockInterval        time.Duration
	ValidatorTimeout     time.Duration
	MinActiveValidators  int
	MaxConsecutiveMisses uint64
}

// RoundRobin represents the round-robin consensus mechanism
type RoundRobin struct {
	validators     []*validator.Authority
	currentIndex   int
	mu            sync.RWMutex
	config        ConsensusConfig
	lastBlockTime time.Time
}

// NewRoundRobin creates a new round-robin consensus instance
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		validators:     make([]*validator.Authority, 0),
		currentIndex:   0,
		config: ConsensusConfig{
			BlockInterval:        100 * time.Millisecond,
			ValidatorTimeout:     1 * time.Second,
			MinActiveValidators:  1,
			MaxConsecutiveMisses: 3,
		},
		lastBlockTime: time.Now(),
	}
}

// AddValidator adds a new validator to the round-robin rotation
func (rr *RoundRobin) AddValidator(v *validator.Authority) error {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	// Check if validator already exists
	for _, existing := range rr.validators {
		if existing.Address == v.Address {
			return errors.New("validator already exists")
		}
	}
	
	v.SetStatus(validator.StatusActive)
	rr.validators = append(rr.validators, v)
	return nil
}

// RemoveValidator removes a validator from the round-robin rotation
func (rr *RoundRobin) RemoveValidator(address string) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	for i, v := range rr.validators {
		if v.Address == address {
			// Set validator status to inactive before removal
			v.SetStatus(validator.StatusInactive)
			rr.validators = append(rr.validators[:i], rr.validators[i+1:]...)
			
			if i < rr.currentIndex {
				rr.currentIndex--
			}
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
	
	fmt.Printf("Getting current validator (index: %d, total validators: %d)\n", rr.currentIndex, len(rr.validators))
	
	if len(rr.validators) == 0 {
		return nil, ErrNoValidators
	}
	
	v := rr.validators[rr.currentIndex]
	fmt.Printf("Current validator: %s (status: %v)\n", v.Address, v.GetStatus())
	
	if v.GetStatus() != validator.StatusActive {
		return nil, ErrValidatorInactive
	}
	
	return v, nil
}

// NextValidator moves to the next validator in the rotation
func (rr *RoundRobin) NextValidator() {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	rr.currentIndex = (rr.currentIndex + 1) % len(rr.validators)
}

// IsValidatorTurn checks if it's the given validator's turn
func (rr *RoundRobin) IsValidatorTurn(address string) bool {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	if len(rr.validators) == 0 {
		fmt.Println("No validators available")
		return false
	}
	
	current := rr.validators[rr.currentIndex]
	isTurn := current.Address == address && current.GetStatus() == validator.StatusActive
	fmt.Printf("Checking validator turn - Address: %s, Current: %s, Status: %v, IsTurn: %v\n", 
		address, current.Address, current.GetStatus(), isTurn)
	
	return isTurn
}

// GetBlockInterval returns the configured block interval
func (rr *RoundRobin) GetBlockInterval() time.Duration {
	return rr.config.BlockInterval
}

// RecordBlockProduction records a successful block production
func (rr *RoundRobin) RecordBlockProduction(timestamp time.Time) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	if len(rr.validators) == 0 {
		fmt.Println("No validators available, skipping block production recording")
		return
	}

	oldIndex := rr.currentIndex
	rr.lastBlockTime = timestamp
	rr.currentIndex = (rr.currentIndex + 1) % len(rr.validators)
	
	fmt.Printf("Block production recorded. Previous validator index: %d, New validator index: %d\n", oldIndex, rr.currentIndex)
	if rr.currentIndex < len(rr.validators) {
		fmt.Printf("Next validator: %s\n", rr.validators[rr.currentIndex].Address)
	}
}

// GetValidatorCount returns the number of validators
func (rr *RoundRobin) GetValidatorCount() int {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	return len(rr.validators)
}

// GetActiveValidatorCount returns the number of active validators
func (rr *RoundRobin) GetActiveValidatorCount() int {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	count := 0
	for _, v := range rr.validators {
		if v.GetStatus() == validator.StatusActive {
			count++
		}
	}
	return count
}

// GetValidators returns all validators
func (rr *RoundRobin) GetValidators() []*validator.Authority {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	validators := make([]*validator.Authority, len(rr.validators))
	copy(validators, rr.validators)
	return validators
}

// ValidateBlock validates if a block can be produced by the given validator
func (rr *RoundRobin) ValidateBlock(address string, blockTime time.Time) error {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	fmt.Printf("Validating block from validator %s\n", address)
	
	// Check if enough time has passed since last block
	timeSinceLastBlock := time.Since(rr.lastBlockTime)
	fmt.Printf("Time since last block: %v (minimum interval: %v)\n", timeSinceLastBlock, rr.config.BlockInterval)
	
	if timeSinceLastBlock < rr.config.BlockInterval {
		return fmt.Errorf("minimum block interval not reached (waited %v, need %v)", timeSinceLastBlock, rr.config.BlockInterval)
	}
	
	// Check if it's the validator's turn
	if !rr.IsValidatorTurn(address) {
		currentValidator := rr.validators[rr.currentIndex]
		fmt.Printf("Not validator's turn. Expected: %s, Got: %s\n", currentValidator.Address, address)
		return ErrNotValidatorTurn
	}
	
	// Check if validator is active
	current := rr.validators[rr.currentIndex]
	fmt.Printf("Current validator status: %v\n", current.GetStatus())
	
	if current.GetStatus() != validator.StatusActive {
		return ErrValidatorInactive
	}
	
	fmt.Printf("Block validation successful for validator %s\n", address)
	return nil
} 