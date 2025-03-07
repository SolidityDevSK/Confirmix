package contracts

import (
	"encoding/json"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Contract represents a smart contract
type Contract struct {
	Address    common.Address `json:"address"`
	Code       []byte        `json:"code"`
	Owner      common.Address `json:"owner"`
	Name       string        `json:"name"`
	Version    string        `json:"version"`
	Timestamp  int64         `json:"timestamp"`
	IsEnabled  bool          `json:"is_enabled"`
}

// Manager manages smart contracts
type Manager struct {
	contracts map[common.Address]*Contract
	vm        *VM
	mu        sync.RWMutex
}

// NewManager creates a new contract manager
func NewManager() *Manager {
	return &Manager{
		contracts: make(map[common.Address]*Contract),
		vm:        NewVM(),
	}
}

// DeployContract deploys a new smart contract
func (m *Manager) DeployContract(code []byte, owner common.Address, name, version string, timestamp int64) (*Contract, error) {
	// Kontrat kodunu doğrula
	if err := m.vm.ValidateContract(code); err != nil {
		return nil, err
	}

	// Kontrat adresi oluştur
	address := generateContractAddress(code, owner, timestamp)

	// Kontratın zaten var olup olmadığını kontrol et
	m.mu.RLock()
	if _, exists := m.contracts[address]; exists {
		m.mu.RUnlock()
		return nil, errors.New("contract already exists")
	}
	m.mu.RUnlock()

	// Yeni kontrat oluştur
	contract := &Contract{
		Address:   address,
		Code:      code,
		Owner:     owner,
		Name:      name,
		Version:   version,
		Timestamp: timestamp,
		IsEnabled: true,
	}

	// Kontratı kaydet
	m.mu.Lock()
	m.contracts[address] = contract
	m.mu.Unlock()

	return contract, nil
}

// ExecuteContract executes a smart contract
func (m *Manager) ExecuteContract(address common.Address, input []byte) ([]byte, error) {
	// Kontratı bul
	m.mu.RLock()
	contract, exists := m.contracts[address]
	m.mu.RUnlock()

	if !exists {
		return nil, errors.New("contract not found")
	}

	if !contract.IsEnabled {
		return nil, errors.New("contract is disabled")
	}

	// Kontratı çalıştır
	return m.vm.Execute(contract.Code, input)
}

// GetContract returns a contract by address
func (m *Manager) GetContract(address common.Address) (*Contract, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contract, exists := m.contracts[address]
	if !exists {
		return nil, errors.New("contract not found")
	}

	return contract, nil
}

// DisableContract disables a contract
func (m *Manager) DisableContract(address common.Address, owner common.Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	contract, exists := m.contracts[address]
	if !exists {
		return errors.New("contract not found")
	}

	if contract.Owner != owner {
		return errors.New("not contract owner")
	}

	contract.IsEnabled = false
	return nil
}

// EnableContract enables a contract
func (m *Manager) EnableContract(address common.Address, owner common.Address) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	contract, exists := m.contracts[address]
	if !exists {
		return errors.New("contract not found")
	}

	if contract.Owner != owner {
		return errors.New("not contract owner")
	}

	contract.IsEnabled = true
	return nil
}

// ListContracts returns all contracts
func (m *Manager) ListContracts() []*Contract {
	m.mu.RLock()
	defer m.mu.RUnlock()

	contracts := make([]*Contract, 0, len(m.contracts))
	for _, contract := range m.contracts {
		contracts = append(contracts, contract)
	}

	return contracts
}

// MarshalJSON implements json.Marshaler
func (c *Contract) MarshalJSON() ([]byte, error) {
	type Alias Contract
	return json.Marshal(&struct {
		Address string `json:"address"`
		Owner   string `json:"owner"`
		*Alias
	}{
		Address: c.Address.Hex(),
		Owner:   c.Owner.Hex(),
		Alias:   (*Alias)(c),
	})
}

// generateContractAddress generates a unique contract address
func generateContractAddress(code []byte, owner common.Address, timestamp int64) common.Address {
	// Adres oluşturmak için gerekli verileri birleştir
	data := append(code, owner.Bytes()...)
	data = append(data, common.BigToHash(big.NewInt(timestamp)).Bytes()...)
	
	// Keccak-256 hash'i hesapla
	return common.BytesToAddress(crypto.Keccak256(data))
} 