package blockchain

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
)

const testMempoolSize = uint64(1000000)

// createMempoolTestTransaction mempool testleri için işlem oluşturur
func createMempoolTestTransaction(nonce uint64, from string, value uint64, gasPrice uint64, gasLimit uint64) *Transaction {
	return &Transaction{
		Hash:      make([]byte, 32),
		From:      from,
		To:        "0x0987654321",
		Value:     big.NewInt(int64(value)),
		GasPrice:  gasPrice,
		GasLimit:  gasLimit,
		GasUsed:   0,
		Nonce:     nonce,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
}

// TestNewMempool mempool oluşturmayı test eder
func TestNewMempool(t *testing.T) {
	mp := NewMempool(testMempoolSize)
	if mp == nil {
		t.Error("Mempool oluşturulamadı")
	}
	if len(mp.transactions) != 0 {
		t.Error("Yeni mempool boş olmalı")
	}
}

// TestMempoolAddTransaction işlem eklemeyi test eder
func TestMempoolAddTransaction(t *testing.T) {
	mp := NewMempool(testMempoolSize)

	// İlk işlemi ekle
	tx1 := &Transaction{
		Hash:      []byte("tx1"),
		From:      "0x1234567890",
		To:        "0x0987654321",
		Value:     big.NewInt(1000),
		Data:      []byte("test"),
		GasPrice:  21000,
		GasLimit:  21000,
		GasUsed:   0,
		Nonce:     1,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
	if err := mp.AddTransaction(tx1); err != nil {
		t.Errorf("İlk işlem eklenemedi: %v", err)
	}

	// Aynı nonce ile ikinci işlemi eklemeyi dene
	tx2 := &Transaction{
		Hash:      []byte("tx2"),
		From:      "0x1234567890",
		To:        "0x0987654321",
		Value:     big.NewInt(2000),
		Data:      []byte("test"),
		GasPrice:  21000,
		GasLimit:  21000,
		GasUsed:   0,
		Nonce:     1,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
	if err := mp.AddTransaction(tx2); err == nil {
		t.Error("Aynı nonce ile ikinci işlem eklendi")
	} else {
		t.Logf("İkinci işlem eklenemedi: %v", err)
	}

	// Farklı nonce ile işlem ekle
	tx3 := &Transaction{
		Hash:      []byte("tx3"),
		From:      "0x1234567890",
		To:        "0x0987654321",
		Value:     big.NewInt(1500),
		Data:      []byte("test"),
		GasPrice:  21000,
		GasLimit:  21000,
		GasUsed:   0,
		Nonce:     2,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
	if err := mp.AddTransaction(tx3); err != nil {
		t.Errorf("Farklı nonce ile işlem eklenemedi: %v", err)
	}

	// İşlem sayısını kontrol et
	if len(mp.transactions) != 2 {
		t.Errorf("Beklenen işlem sayısı 2, alınan: %d", len(mp.transactions))
	}
}

// TestRemoveTransaction işlem silmeyi test eder
func TestRemoveTransaction(t *testing.T) {
	mp := NewMempool(testMempoolSize)

	// İşlem ekle
	tx := &Transaction{
		Hash:      []byte("tx1"),
		From:      "0x1234567890",
		To:        "0x0987654321",
		Value:     big.NewInt(1000),
		Data:      []byte("test"),
		GasPrice:  21000,
		GasLimit:  21000,
		GasUsed:   0,
		Nonce:     1,
		Signature: make([]byte, 64),
		Status:    TxPending,
	}
	if err := mp.AddTransaction(tx); err != nil {
		t.Fatalf("İşlem eklenemedi: %v", err)
	}

	// İşlemi sil
	mp.RemoveTransaction(tx.Hash)

	// İşlemin silindiğini kontrol et
	if len(mp.transactions) != 0 {
		t.Error("İşlem silinmedi")
	}
}

// TestGetPendingNonce bekleyen nonce değerini test eder
func TestGetPendingNonce(t *testing.T) {
	mp := NewMempool(testMempoolSize)

	// İşlemleri ekle
	for i := uint64(1); i <= 5; i++ {
		tx := &Transaction{
			Hash:      []byte(fmt.Sprintf("tx%d", i)),
			From:      "0x1234567890",
			To:        "0x0987654321",
			Value:     big.NewInt(1000),
			Data:      []byte("test"),
			GasPrice:  21000,
			GasLimit:  21000,
			GasUsed:   0,
			Nonce:     i,
			Signature: make([]byte, 64),
			Status:    TxPending,
		}
		if err := mp.AddTransaction(tx); err != nil {
			t.Fatalf("İşlem eklenemedi: %v", err)
		}
	}

	// Nonce değerini kontrol et
	nonce := mp.GetPendingNonce("0x1234567890")
	if nonce != 6 {
		t.Errorf("Beklenen nonce 6, alınan: %d", nonce)
	}
}

// TestGetBestTransactions en iyi işlemleri seçmeyi test eder
func TestGetBestTransactions(t *testing.T) {
	mp := NewMempool(testMempoolSize)

	// İşlemleri ekle
	for i := uint64(1); i <= 3; i++ {
		tx := &Transaction{
			Hash:      []byte(fmt.Sprintf("tx%d", i)),
			From:      "0x1234567890",
			To:        "0x0987654321",
			Value:     big.NewInt(1000),
			Data:      []byte("test"),
			GasPrice:  21000 + i*1000, // Farklı gas price'lar
			GasLimit:  21000,
			GasUsed:   0,
			Nonce:     i,
			Signature: make([]byte, 64),
			Status:    TxPending,
		}
		if err := mp.AddTransaction(tx); err != nil {
			t.Fatalf("İşlem eklenemedi: %v", err)
		}
	}

	// En iyi işlemleri al
	txs := mp.GetBestTransactions(100000)
	if len(txs) != 3 {
		t.Errorf("Beklenen işlem sayısı 3, alınan: %d", len(txs))
	}

	// İşlem sırasını kontrol et (gas price'a göre sıralı olmalı)
	for i := 0; i < len(txs)-1; i++ {
		if txs[i].GasPrice < txs[i+1].GasPrice {
			t.Errorf("İşlemler gas price'a göre sıralı değil: %d < %d", txs[i].GasPrice, txs[i+1].GasPrice)
		}
	}
}

// TestConcurrentAccess eşzamanlı erişimi test eder
func TestConcurrentAccess(t *testing.T) {
	mp := NewMempool(1000000)
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// İşlem ekle
			tx := createMempoolTestTransaction(uint64(id), "0x1234", 1000, 21000, 21000)
			_ = mp.AddTransaction(tx)
			// İşlemi al
			_ = mp.GetTransaction(tx.Hash)
			// İşlemi sil
			mp.RemoveTransaction(tx.Hash)
		}(i)
	}

	wg.Wait()
} 