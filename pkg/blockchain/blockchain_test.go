package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/SolidityDevSK/Confirmix/pkg/consensus"
	"github.com/SolidityDevSK/Confirmix/pkg/contracts"
)

// createTestKey test için özel anahtar oluşturur
func createTestKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// createTestValidator test için validator oluşturur
func createTestValidator(t *testing.T) (*validator.Authority, error) {
	key, err := createTestKey()
	if err != nil {
		t.Fatalf("Test key oluşturulamadı: %v", err)
		return nil, err
	}
	return validator.NewAuthority(key)
}

func createTestTransaction(nonce uint64, value *big.Int) *Transaction {
	return &Transaction{
		Hash:      make([]byte, 32),
		From:      "0x1234567890",
		To:        "0x0987654321",
		Value:     value,
		Data:      []byte("test"),
		GasPrice:  21000,
		GasLimit:  21000,
		GasUsed:   0,
		Nonce:     nonce,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
}

// createTestBlock test için blok oluşturur
func createTestBlock(prevHash []byte, height uint64, v *validator.Authority) (*Block, error) {
	block, err := NewBlock(height, prevHash, make([]byte, 32), 1000000, v)
	if err != nil {
		return nil, err
	}

	// Sign the block
	if err := block.Sign(v); err != nil {
		return nil, fmt.Errorf("failed to sign block: %v", err)
	}

	return block, nil
}

// TestNewBlockchain yeni blockchain oluşturmayı test eder
func TestNewBlockchain(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	if len(bc.Blocks) != 1 {
		t.Error("Genesis bloğu oluşturulmadı")
	}

	if bc.Blocks[0].Header.Height != 0 {
		t.Errorf("Başlangıç yüksekliği 0 olmalı, alınan: %d", bc.Blocks[0].Header.Height)
	}
}

// TestAddBlock blok eklemeyi test eder
func TestAddBlock(t *testing.T) {
	t.Log("Starting TestAddBlock")
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}
	t.Logf("Test validator created with address: %s", v.Address)

	bc := &Blockchain{
		Blocks:     make([]*Block, 0),
		Validators: make(map[string]*validator.Authority),
		consensus:  consensus.NewRoundRobin(),
	}

	t.Log("Adding genesis validator")
	bc.AddValidator(v)

	t.Log("Creating genesis block")
	genesisBlock, err := NewBlock(0, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Genesis blok oluşturulamadı: %v", err)
	}

	bc.Blocks = append(bc.Blocks, genesisBlock)
	bc.LastBlockTime = time.Now()
	t.Log("Genesis block added to chain")

	t.Log("Creating new block")
	prevHash := bc.GetLatestBlock().GetHash()
	block, err := createTestBlock(prevHash, 1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	t.Log("Adding transaction to block")
	tx := createTestTransaction(1, big.NewInt(1000))
	block.AddTransaction(tx)

	t.Log("Adding block to chain")
	err = bc.AddBlock(block)
	if err != nil {
		t.Errorf("Blok eklenemedi: %v", err)
	}

	t.Logf("Chain length after adding block: %d", len(bc.Blocks))
	if len(bc.Blocks) != 2 {
		t.Errorf("Beklenen blok sayısı 2, alınan: %d", len(bc.Blocks))
	}

	t.Log("Checking validator address")
	if block.Header.ValidatorAddress != v.Address {
		t.Errorf("Beklenen validator adresi %s, alınan: %s", v.Address, block.Header.ValidatorAddress)
	}

	t.Log("Testing invalid block addition")
	invalidBlock, err := createTestBlock(make([]byte, 32), 2, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	t.Log("Attempting to add invalid block")
	err = bc.AddBlock(invalidBlock)
	if err == nil {
		t.Error("Geçersiz önceki hash ile blok eklendi")
	} else {
		t.Logf("Expected error received: %v", err)
	}
}

// TestTransactionValidation işlem doğrulama işlemlerini test eder
func TestTransactionValidation(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Geçerli bir işlem oluştur
	tx := createTestTransaction(1, big.NewInt(1000))

	// İşlemi içeren blok oluştur
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}
	block.AddTransaction(tx)

	// Bloğu zincire ekle
	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// İşlem doğrulamasını kontrol et
	if !bc.IsValid() {
		t.Error("İşlem doğrulaması başarısız")
	}
}

// TestBlockValidation blok doğrulama işlemlerini test eder
func TestBlockValidation(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Geçerli bir blok oluştur
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	// Bloğu zincire ekle
	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// Blok doğrulamasını kontrol et
	if !bc.IsValid() {
		t.Error("Blok doğrulaması başarısız")
	}
}

// TestBlockchainReorganization zincir yeniden düzenleme işlemlerini test eder
func TestBlockchainReorganization(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Ana zincirde blok oluştur
	mainBlock, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Ana zincir bloğu oluşturulamadı: %v", err)
	}

	// Ana zincire bloğu ekle
	err = bc.AddBlock(mainBlock)
	if err != nil {
		t.Fatalf("Ana zincir bloğu eklenemedi: %v", err)
	}

	// Yan zincirde blok oluştur
	sideBlock, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Yan zincir bloğu oluşturulamadı: %v", err)
	}

	// Yan zincire bloğu ekle
	err = bc.AddBlock(sideBlock)
	if err != nil {
		t.Fatalf("Yan zincir bloğu eklenemedi: %v", err)
	}

	// Zincir durumunu kontrol et
	if !bc.IsValid() {
		t.Error("Zincir yeniden düzenleme sonrası geçersiz")
	}
}

// TestBlockchainConsensus konsensüs mekanizmasını test eder
func TestBlockchainConsensus(t *testing.T) {
	t.Log("Starting TestBlockchainConsensus")
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}
	t.Logf("First validator created with address: %s", v.Address)

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	t.Log("Creating second validator")
	v2, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("İkinci validator oluşturulamadı: %v", err)
	}
	t.Logf("Second validator created with address: %s", v2.Address)

	t.Log("Adding second validator to blockchain")
	bc.AddValidator(v2)

	t.Log("Creating first block with first validator")
	block1, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("İlk blok oluşturulamadı: %v", err)
	}

	t.Log("Adding first block to chain")
	err = bc.AddBlock(block1)
	if err != nil {
		t.Fatalf("İlk blok eklenemedi: %v", err)
	}

	t.Log("Creating second block with second validator")
	block2, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v2)
	if err != nil {
		t.Fatalf("İkinci blok oluşturulamadı: %v", err)
	}

	t.Log("Adding second block to chain")
	err = bc.AddBlock(block2)
	if err != nil {
		t.Fatalf("İkinci blok eklenemedi: %v", err)
	}

	t.Log("Checking chain validity")
	if !bc.IsValid() {
		t.Error("Konsensüs sonrası zincir geçersiz")
	}
}

// TestValidatorRotation validator değişimini test eder
func TestValidatorRotation(t *testing.T) {
	v1, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("İlk validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v1)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// İkinci validator oluştur
	v2, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("İkinci validator oluşturulamadı: %v", err)
	}

	bc.AddValidator(v2)

	// İlk blok için v1'i kullan
	block1, err := createTestBlock(bc.GetLatestBlock().GetHash(), 1, v1)
	if err != nil {
		t.Fatalf("İlk blok oluşturulamadı: %v", err)
	}
	err = bc.AddBlock(block1)
	if err != nil {
		t.Fatalf("İlk blok eklenemedi: %v", err)
	}

	// İkinci blok için v2'yi kullan
	block2, err := createTestBlock(bc.GetLatestBlock().GetHash(), 2, v2)
	if err != nil {
		t.Fatalf("İkinci blok oluşturulamadı: %v", err)
	}
	err = bc.AddBlock(block2)
	if err != nil {
		t.Fatalf("İkinci blok eklenemedi: %v", err)
	}

	// Validator sırasını kontrol et
	currentValidator, err := bc.GetCurrentValidator()
	if err != nil {
		t.Fatalf("Mevcut validator alınamadı: %v", err)
	}
	if currentValidator.Address == v2.Address {
		t.Error("Validator rotasyonu hatalı")
	}
}

// TestBlockchainInfo zincir bilgilerini test eder
func TestBlockchainInfo(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc := &Blockchain{
		Blocks:          make([]*Block, 0),
		Validators:      make(map[string]*validator.Authority),
		consensus:       consensus.NewRoundRobin(),
		ContractManager: contracts.NewManager(),
	}

	// Genesis validator'ı ekle
	bc.AddValidator(v)

	// Genesis bloğu oluştur
	genesisBlock, err := NewBlock(0, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Genesis blok oluşturulamadı: %v", err)
	}

	bc.Blocks = append(bc.Blocks, genesisBlock)
	bc.LastBlockTime = time.Now()

	info := bc.GetBlockchainInfo()

	// Temel bilgileri kontrol et
	if info["blockCount"].(int) != 1 {
		t.Error("Blok sayısı hatalı")
	}
	if info["validatorCount"].(int) != 1 {
		t.Error("Validator sayısı hatalı")
	}
	if info["activeValidators"].(int) != 1 {
		t.Error("Aktif validator sayısı hatalı")
	}
	if info["blockInterval"].(float64) <= 0 {
		t.Error("Blok aralığı hatalı")
	}
	if info["contractCount"].(int) != 0 {
		t.Error("Kontrat sayısı hatalı")
	}
}

// TestInvalidValidatorOperations hatalı validator durumlarını test eder
func TestInvalidValidatorOperations(t *testing.T) {
	v1, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("İlk validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v1)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Yetkisiz validator ile blok eklemeyi dene
	unauthorizedValidator, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Yetkisiz validator oluşturulamadı: %v", err)
	}

	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, unauthorizedValidator)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block)
	if err == nil {
		t.Error("Yetkisiz validator ile blok ekleme engellenmedi")
	}

	// Validator'ı devre dışı bırak
	bc.RemoveValidator(v1.Address)

	block2, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v1)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block2)
	if err == nil {
		t.Error("Devre dışı validator ile blok ekleme engellenmedi")
	}

	// Geçersiz adresle validator silmeyi dene
	bc.RemoveValidator("invalid_address")
	if len(bc.Validators) != 0 {
		t.Error("Geçersiz validator silme işlemi hatalı")
	}
}

// TestContractOperations akıllı kontrat işlemlerini test eder
func TestContractOperations(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc := &Blockchain{
		Blocks:          make([]*Block, 0),
		Validators:      make(map[string]*validator.Authority),
		consensus:       consensus.NewRoundRobin(),
		ContractManager: contracts.NewManager(),
	}

	// Genesis validator'ı ekle
	bc.AddValidator(v)

	// Genesis bloğu oluştur
	genesisBlock, err := NewBlock(0, make([]byte, 32), make([]byte, 32), 1000000, v)
	if err != nil {
		t.Fatalf("Genesis blok oluşturulamadı: %v", err)
	}

	bc.Blocks = append(bc.Blocks, genesisBlock)
	bc.LastBlockTime = time.Now()

	// Test kontrat kodu
	code := []byte("test contract code")
	owner := common.HexToAddress("0x1234567890")
	name := "TestContract"
	version := "1.0"

	// Kontrat deploy et
	contract, err := bc.DeployContract(code, owner, name, version)
	if err != nil {
		t.Fatalf("Kontrat deploy edilemedi: %v", err)
	}

	// Kontrat bilgilerini kontrol et
	if contract.Name != name {
		t.Error("Kontrat ismi hatalı")
	}
	if contract.Version != version {
		t.Error("Kontrat versiyonu hatalı")
	}

	// Kontrat listesini kontrol et
	contracts := bc.ListContracts()
	if len(contracts) != 1 {
		t.Error("Kontrat listesi hatalı")
	}

	// Yeni blok oluştur
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	// Kontrat çağrısı için işlem oluştur
	tx := createTestTransaction(1, big.NewInt(0))
	tx.To = contract.Address.Hex() // Address'i string'e çevir
	tx.Data = []byte("test input")
	block.AddTransaction(tx)

	// Bloğu zincire ekle
	err = bc.AddBlock(block)
	if err != nil {
		t.Errorf("Kontrat çağrısı için blok eklenemedi: %v", err)
	}
}

// TestInvalidTransactions geçersiz işlem durumlarını test eder
func TestInvalidTransactions(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Gas limit'i aşan işlem için blok oluştur
	block1, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	invalidTx := createTestTransaction(1, big.NewInt(1000))
	invalidTx.GasLimit = 1000000000 // Çok yüksek gas limit
	err = block1.AddTransaction(invalidTx)
	if err == nil {
		t.Error("Gas limit aşımı engellenmedi")
	}

	// Geçersiz nonce'lu işlem için blok oluştur
	block2, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	duplicateTx := createTestTransaction(1, big.NewInt(1000))
	err = block2.AddTransaction(duplicateTx)
	if err != nil {
		t.Fatalf("İlk işlem eklenemedi: %v", err)
	}
	err = block2.AddTransaction(duplicateTx)
	if err == nil {
		t.Error("Aynı nonce'lu işlemler engellenmedi")
	}

	// Negatif değerli işlem için blok oluştur
	block3, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	negativeTx := createTestTransaction(2, big.NewInt(-1000))
	err = block3.AddTransaction(negativeTx)
	if err == nil {
		t.Error("Negatif değerli işlem engellenmedi")
	}
}

// TestConcurrentOperations eşzamanlı işlemleri test eder
func TestConcurrentOperations(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Eşzamanlı validator ekleme/silme
	const numOperations = 100
	done := make(chan bool)

	go func() {
		for i := 0; i < numOperations; i++ {
			newValidator, err := createTestValidator(t)
			if err != nil {
				t.Errorf("Validator oluşturulamadı: %v", err)
				continue
			}
			bc.AddValidator(newValidator)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < numOperations; i++ {
			newValidator, err := createTestValidator(t)
			if err != nil {
				t.Errorf("Validator oluşturulamadı: %v", err)
				continue
			}
			bc.AddValidator(newValidator)
			bc.RemoveValidator(newValidator.Address)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// İşlemlerin tamamlanmasını bekle
	<-done
	<-done

	// Durumu kontrol et
	if !bc.IsValid() {
		t.Error("Eşzamanlı işlemler sonrası zincir geçersiz")
	}
}

// TestBlockchainPerformance performans metriklerini test eder
func TestBlockchainPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Performans testi atlanıyor")
	}

	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Çok sayıda işlem oluştur
	const numTransactions = 1000
	transactions := make([]*Transaction, numTransactions)
	for i := 0; i < numTransactions; i++ {
		transactions[i] = createTestTransaction(uint64(i), big.NewInt(1000))
	}

	// Blok oluşturma süresini ölç
	start := time.Now()
	batchSize := 100
	for i := 0; i < numTransactions; i += batchSize {
		end := i + batchSize
		if end > numTransactions {
			end = numTransactions
		}

		// Yeni blok oluştur
		block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
		if err != nil {
			t.Fatalf("Blok oluşturulamadı: %v", err)
		}

		// İşlemleri bloğa ekle
		for _, tx := range transactions[i:end] {
			block.AddTransaction(tx)
		}

		// Bloğu zincire ekle
		err = bc.AddBlock(block)
		if err != nil {
			t.Fatalf("Performans testi sırasında blok eklenemedi: %v", err)
		}
	}
	duration := time.Since(start)

	// Performans metriklerini kontrol et
	t.Logf("Toplam süre: %v", duration)
	t.Logf("İşlem başına süre: %v", duration/time.Duration(numTransactions))
	t.Logf("Saniyedeki işlem sayısı: %f", float64(numTransactions)/duration.Seconds())
}

// TestBlockchainRecovery zincir kurtarma işlemlerini test eder
func TestBlockchainRecovery(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Blok ekle
	block1, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block1)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// Zinciri kontrol et
	if !bc.IsValid() {
		t.Error("Zincir geçersiz")
	}
}

// TestBlockchainStress stres testlerini gerçekleştirir
func TestBlockchainStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Stres testi atlanıyor")
	}

	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Çok sayıda blok ekle
	for i := 0; i < 100; i++ {
		block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
		if err != nil {
			t.Fatalf("Blok oluşturulamadı: %v", err)
		}

		// Her bloğa birkaç işlem ekle
		for j := 0; j < 10; j++ {
			tx := createTestTransaction(uint64(j), big.NewInt(1000))
			block.AddTransaction(tx)
		}

		err = bc.AddBlock(block)
		if err != nil {
			t.Fatalf("Blok eklenemedi: %v", err)
		}
	}

	// Zinciri kontrol et
	if !bc.IsValid() {
		t.Error("Stres testi sonrası zincir geçersiz")
	}
}

// TestIsValid zincir geçerliliğini test eder
func TestIsValid(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Geçerli blok ekle
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// Zincir geçerliliğini kontrol et
	if !bc.IsValid() {
		t.Error("Geçerli zincir geçersiz olarak işaretlendi")
	}
}

// TestGetBlock blok alma işlemini test eder
func TestGetBlock(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Blok ekle
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// Bloğu hash ile al
	retrievedBlock := bc.GetBlock(block.GetHash())
	if retrievedBlock == nil {
		t.Error("Blok bulunamadı")
	}
	if !bytes.Equal(retrievedBlock.GetHash(), block.GetHash()) {
		t.Error("Alınan blok hash'i eşleşmiyor")
	}
}

// TestGetBlockByHeight yüksekliğe göre blok alma işlemini test eder
func TestGetBlockByHeight(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Blok ekle
	block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
	if err != nil {
		t.Fatalf("Blok oluşturulamadı: %v", err)
	}

	err = bc.AddBlock(block)
	if err != nil {
		t.Fatalf("Blok eklenemedi: %v", err)
	}

	// Bloğu yükseklik ile al
	retrievedBlock := bc.GetBlockByHeight(block.Header.Height)
	if retrievedBlock == nil {
		t.Error("Blok bulunamadı")
	}
	if retrievedBlock.Header.Height != block.Header.Height {
		t.Error("Alınan blok yüksekliği eşleşmiyor")
	}
}

// TestGetLatestBlock son bloğu alma işlemini test eder
func TestGetLatestBlock(t *testing.T) {
	v, err := createTestValidator(t)
	if err != nil {
		t.Fatalf("Validator oluşturulamadı: %v", err)
	}

	bc, err := NewBlockchain(v)
	if err != nil {
		t.Fatalf("Blockchain oluşturulamadı: %v", err)
	}

	// Birkaç blok ekle
	for i := uint64(1); i <= 3; i++ {
		block, err := createTestBlock(bc.GetLatestBlock().GetHash(), bc.GetLatestBlock().Header.Height+1, v)
		if err != nil {
			t.Fatalf("Blok oluşturulamadı: %v", err)
		}

		err = bc.AddBlock(block)
		if err != nil {
			t.Fatalf("Blok eklenemedi: %v", err)
		}
	}

	// Son bloğu al ve kontrol et
	latestBlock := bc.GetLatestBlock()
	if latestBlock == nil {
		t.Error("Son blok alınamadı")
	}
	if latestBlock.Header.Height != 3 {
		t.Errorf("Son blok yüksekliği hatalı. Beklenen: 3, Alınan: %d", latestBlock.Header.Height)
	}
} 