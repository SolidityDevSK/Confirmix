package blockchain

import (
	"container/heap"
	"fmt"
	"math/big"
	"sync"
)

// TransactionPriority işlemlerin öncelik sırasını belirler
type TransactionPriority struct {
	tx       *Transaction
	priority *big.Int
	index    int
}

// PriorityQueue işlem öncelik kuyruğu için heap implementasyonu
type PriorityQueue []*TransactionPriority

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority.Cmp(pq[j].priority) > 0
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TransactionPriority)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Mempool işlem havuzunu yönetir
type Mempool struct {
	mu            sync.RWMutex
	transactions  map[string]*Transaction     // Hash -> Transaction
	byNonce      map[string]map[uint64]*Transaction // Address -> Nonce -> Transaction
	priorityQueue PriorityQueue
	maxSize      uint64
	currentSize  uint64
}

// NewMempool yeni bir işlem havuzu oluşturur
func NewMempool(maxSize uint64) *Mempool {
	mp := &Mempool{
		transactions:  make(map[string]*Transaction),
		byNonce:      make(map[string]map[uint64]*Transaction),
		priorityQueue: make(PriorityQueue, 0),
		maxSize:      maxSize,
		currentSize:  0,
	}
	heap.Init(&mp.priorityQueue)
	return mp
}

// AddTransaction işlemi havuza ekler
func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// İşlem zaten varsa
	if _, exists := mp.transactions[string(tx.Hash)]; exists {
		return fmt.Errorf("transaction already exists in mempool")
	}

	// Havuz boyut kontrolü
	if mp.currentSize+tx.GetSize() > mp.maxSize {
		return fmt.Errorf("mempool size limit exceeded")
	}

	// Nonce kontrolü
	if _, exists := mp.byNonce[tx.From]; !exists {
		mp.byNonce[tx.From] = make(map[uint64]*Transaction)
	}
	if _, exists := mp.byNonce[tx.From][tx.Nonce]; exists {
		return fmt.Errorf("transaction with same nonce already exists")
	}

	// İşlemi havuza ekle
	mp.transactions[string(tx.Hash)] = tx
	mp.byNonce[tx.From][tx.Nonce] = tx
	mp.currentSize += tx.GetSize()

	// Öncelik kuyruğuna ekle
	priority := calculatePriority(tx)
	item := &TransactionPriority{
		tx:       tx,
		priority: priority,
	}
	heap.Push(&mp.priorityQueue, item)

	return nil
}

// RemoveTransaction işlemi havuzdan kaldırır
func (mp *Mempool) RemoveTransaction(hash []byte) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	tx, exists := mp.transactions[string(hash)]
	if !exists {
		return
	}

	// İşlemi haritalardan kaldır
	delete(mp.transactions, string(hash))
	delete(mp.byNonce[tx.From], tx.Nonce)
	if len(mp.byNonce[tx.From]) == 0 {
		delete(mp.byNonce, tx.From)
	}

	// Öncelik kuyruğundan kaldır
	for i, item := range mp.priorityQueue {
		if string(item.tx.Hash) == string(hash) {
			heap.Remove(&mp.priorityQueue, i)
			break
		}
	}

	mp.currentSize -= tx.GetSize()
}

// GetTransaction belirtilen hash'e sahip işlemi döndürür
func (mp *Mempool) GetTransaction(hash []byte) *Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.transactions[string(hash)]
}

// GetPendingNonce bir adres için bekleyen en yüksek nonce değerini döndürür
func (mp *Mempool) GetPendingNonce(address string) uint64 {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if nonces, exists := mp.byNonce[address]; exists {
		maxNonce := uint64(0)
		for nonce := range nonces {
			if nonce > maxNonce {
				maxNonce = nonce
			}
		}
		return maxNonce + 1
	}
	return 0
}

// GetBestTransactions en yüksek öncelikli işlemleri döndürür
func (mp *Mempool) GetBestTransactions(gasLimit uint64) []*Transaction {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if len(mp.transactions) == 0 {
		return nil
	}

	// Mevcut öncelik kuyruğunu kopyala
	tempQueue := make(PriorityQueue, len(mp.priorityQueue))
	copy(tempQueue, mp.priorityQueue)
	heap.Init(&tempQueue)

	result := make([]*Transaction, 0)
	totalGas := uint64(0)

	// Gas limitini aşmadan en yüksek öncelikli işlemleri seç
	for len(tempQueue) > 0 && totalGas < gasLimit {
		item := heap.Pop(&tempQueue).(*TransactionPriority)
		if totalGas+item.tx.GasLimit <= gasLimit {
			result = append(result, item.tx)
			totalGas += item.tx.GasLimit
		}
	}

	return result
}

// Clear havuzu temizler
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.transactions = make(map[string]*Transaction)
	mp.byNonce = make(map[string]map[uint64]*Transaction)
	mp.priorityQueue = make(PriorityQueue, 0)
	mp.currentSize = 0
}

// GetSize havuzun mevcut boyutunu döndürür
func (mp *Mempool) GetSize() uint64 {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return mp.currentSize
}

// calculatePriority işlem önceliğini hesaplar
func calculatePriority(tx *Transaction) *big.Int {
	// Öncelik sadece gasPrice'a göre belirlenir
	return new(big.Int).SetUint64(tx.GasPrice)
} 