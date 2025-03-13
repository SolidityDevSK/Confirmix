package consensus

import (
	"testing"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
)

func TestRoundRobin(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second * 3,
		MinActiveValidators:  2,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)

	// Test adding validators
	v1, err := validator.NewAuthority()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}
	
	v2, err := validator.NewAuthority()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	if err := rr.AddValidator(v1); err != nil {
		t.Errorf("Failed to add validator: %v", err)
	}

	if err := rr.AddValidator(v2); err != nil {
		t.Errorf("Failed to add validator: %v", err)
	}

	// Test duplicate validator
	if err := rr.AddValidator(v1); err == nil {
		t.Error("Expected error when adding duplicate validator")
	}

	// Test validator rotation
	current, err := rr.GetCurrentValidator()
	if err != nil {
		t.Errorf("Failed to get current validator: %v", err)
	}
	if current.Address != v1.Address {
		t.Error("Expected first validator to be current")
	}

	// Wait for block interval
	time.Sleep(config.BlockInterval)

	next, err := rr.NextValidator()
	if err != nil {
		t.Errorf("Failed to get next validator: %v", err)
	}
	if next.Address != v2.Address {
		t.Error("Expected second validator to be next")
	}

	// Test validator status
	v2.UpdateStatus(validator.StatusInactive)
	_, err = rr.GetCurrentValidator()
	if err != ErrValidatorInactive {
		t.Error("Expected validator inactive error")
	}

	// Test block validation
	v1.UpdateStatus(validator.StatusActive)
	rr.currentIndex = 0 // Reset to first validator
	
	// Wait for block interval
	time.Sleep(config.BlockInterval)
	
	if err := rr.ValidateBlock(v1.Address, time.Now()); err != nil {
		t.Errorf("Failed to validate block: %v", err)
	}

	if err := rr.ValidateBlock(v2.Address, time.Now()); err != ErrNotValidatorTurn {
		t.Error("Expected not validator's turn error")
	}

	// Test removing validator
	rr.RemoveValidator(v2.Address)
	validators := rr.GetValidators()
	if len(validators) != 1 {
		t.Error("Expected one validator after removal")
	}

	// Test active validator count
	count := rr.GetActiveValidatorCount()
	if count != 1 {
		t.Errorf("Expected 1 active validator, got %d", count)
	}

	// Test block production recording
	blockTime := time.Now()
	rr.RecordBlockProduction(blockTime)
	if rr.lastBlockTime != blockTime {
		t.Error("Block time not updated correctly")
	}

	// Test validator timeout
	time.Sleep(config.ValidatorTimeout)
	v1.UpdateStatus(validator.StatusInactive)
	_, err = rr.NextValidator()
	if err != ErrNoValidators {
		t.Error("Expected no validators error after timeout")
	}
}

func TestConsecutiveMisses(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  1,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	v1, _ := validator.NewAuthority()
	rr.AddValidator(v1)

	// Test consecutive misses
	for i := 0; i < int(config.MaxConsecutiveMisses); i++ {
		time.Sleep(config.ValidatorTimeout)
		_, _ = rr.NextValidator()
	}

	if v1.Status != validator.StatusPenalized {
		t.Error("Expected validator to be penalized after max consecutive misses")
	}
}

func TestValidatorPerformance(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  1,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	v1, _ := validator.NewAuthority()
	rr.AddValidator(v1)

	// Record some block productions
	for i := 0; i < 5; i++ {
		time.Sleep(config.BlockInterval)
		rr.RecordBlockProduction(time.Now())
	}

	metrics := v1.GetPerformanceMetrics()
	if metrics["producedBlocks"].(uint64) != 5 {
		t.Errorf("Expected 5 produced blocks, got %d", metrics["producedBlocks"])
	}
}

func TestValidatorRotationWithMultiple(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  3,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	
	// Add multiple validators
	validators := make([]*validator.Authority, 3)
	for i := 0; i < 3; i++ {
		v, _ := validator.NewAuthority()
		validators[i] = v
		rr.AddValidator(v)
	}

	// Test rotation through all validators
	for i := 0; i < len(validators); i++ {
		current, err := rr.GetCurrentValidator()
		if err != nil {
			t.Errorf("Failed to get validator at step %d: %v", i, err)
		}
		if current.Address != validators[i].Address {
			t.Errorf("Wrong validator at step %d", i)
		}
		time.Sleep(config.BlockInterval)
		_, _ = rr.NextValidator()
	}
}

func TestBlockIntervalEnforcement(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 200,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  1,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	v1, _ := validator.NewAuthority()
	rr.AddValidator(v1)

	// Try to validate block before interval
	rr.RecordBlockProduction(time.Now())
	err := rr.ValidateBlock(v1.Address, time.Now())
	if err == nil {
		t.Error("Expected error when validating block before interval")
	}

	// Wait for interval and try again
	time.Sleep(config.BlockInterval)
	err = rr.ValidateBlock(v1.Address, time.Now())
	if err != nil {
		t.Errorf("Unexpected error after waiting for interval: %v", err)
	}
}

func TestValidatorRecovery(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  1,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	v1, _ := validator.NewAuthority()
	rr.AddValidator(v1)

	// Cause validator to be penalized
	for i := 0; i < int(config.MaxConsecutiveMisses); i++ {
		time.Sleep(config.ValidatorTimeout)
		_, _ = rr.NextValidator()
	}

	if v1.Status != validator.StatusPenalized {
		t.Error("Expected validator to be penalized")
	}

	// Simulate recovery by producing blocks
	v1.UpdateStatus(validator.StatusActive)
	for i := 0; i < 5; i++ {
		time.Sleep(config.BlockInterval)
		rr.RecordBlockProduction(time.Now())
	}

	metrics := v1.GetPerformanceMetrics()
	if metrics["consecutiveMisses"].(uint64) != 0 {
		t.Error("Expected consecutive misses to be reset after recovery")
	}
}

func TestMinValidatorRequirement(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  3,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	
	// Add only 2 validators when minimum is 3
	v1, _ := validator.NewAuthority()
	v2, _ := validator.NewAuthority()
	rr.AddValidator(v1)
	rr.AddValidator(v2)

	activeCount := rr.GetActiveValidatorCount()
	if activeCount >= config.MinActiveValidators {
		t.Errorf("Expected active validators (%d) to be less than minimum required (%d)", 
			activeCount, config.MinActiveValidators)
	}

	// Add third validator to meet minimum requirement
	v3, _ := validator.NewAuthority()
	rr.AddValidator(v3)

	activeCount = rr.GetActiveValidatorCount()
	if activeCount < config.MinActiveValidators {
		t.Errorf("Expected active validators (%d) to meet minimum requirement (%d)", 
			activeCount, config.MinActiveValidators)
	}
}

func TestValidatorReactivation(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  2,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	
	// Add two validators
	v1, _ := validator.NewAuthority()
	v2, _ := validator.NewAuthority()
	rr.AddValidator(v1)
	rr.AddValidator(v2)

	// Deactivate first validator
	v1.UpdateStatus(validator.StatusInactive)
	
	// Verify second validator is current
	current, err := rr.GetCurrentValidator()
	if err != nil {
		t.Errorf("Failed to get current validator: %v", err)
	}
	if current.Address != v2.Address {
		t.Error("Expected second validator to be current after first was deactivated")
	}

	// Reactivate first validator
	v1.UpdateStatus(validator.StatusActive)
	time.Sleep(config.BlockInterval)
	
	// Verify rotation continues correctly
	next, err := rr.NextValidator()
	if err != nil {
		t.Errorf("Failed to get next validator: %v", err)
	}
	if next.Address != v1.Address {
		t.Error("Expected first validator to be next after reactivation")
	}
}

func TestConcurrentValidatorUpdates(t *testing.T) {
	config := ConsensusConfig{
		BlockInterval:        time.Millisecond * 100,
		ValidatorTimeout:     time.Second,
		MinActiveValidators:  1,
		MaxConsecutiveMisses: 3,
	}
	
	rr := NewRoundRobin(config)
	
	// Add initial validator
	v1, _ := validator.NewAuthority()
	rr.AddValidator(v1)

	// Start goroutine to simulate concurrent validator updates
	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			v, _ := validator.NewAuthority()
			rr.AddValidator(v)
			time.Sleep(config.BlockInterval / 2)
		}
		done <- true
	}()

	// Simultaneously perform validator operations
	for i := 0; i < 5; i++ {
		current, err := rr.GetCurrentValidator()
		if err != nil {
			t.Errorf("Failed to get current validator during concurrent updates: %v", err)
		}
		if current == nil {
			t.Error("Got nil validator during concurrent updates")
		}
		time.Sleep(config.BlockInterval / 2)
		_, _ = rr.NextValidator()
	}

	<-done // Wait for concurrent operations to complete

	// Verify final validator count
	validators := rr.GetValidators()
	if len(validators) != 6 { // Initial + 5 added
		t.Errorf("Expected 6 validators after concurrent updates, got %d", len(validators))
	}
} 