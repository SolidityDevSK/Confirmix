package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
)

func main() {
	// Validator'ları oluştur
	validator1, err := validator.NewAuthority()
	if err != nil {
		log.Fatal("Validator 1 oluşturulamadı:", err)
	}

	validator2, err := validator.NewAuthority()
	if err != nil {
		log.Fatal("Validator 2 oluşturulamadı:", err)
	}

	// Yeni bir blockchain oluştur (validator1 genesis validator olacak)
	bc, err := blockchain.NewBlockchain(validator1)
	if err != nil {
		log.Fatal("Blockchain oluşturulamadı:", err)
	}

	// İkinci validator'ı ekle
	bc.AddValidator(validator2)

	fmt.Println("Blockchain oluşturuldu!")
	fmt.Println("Genesis bloğu hash'i:", bc.Blocks[0].GetHashString())
	fmt.Printf("Genesis Validator: %s\n", validator1.Address[:10])

	// Birkaç blok ekleyelim
	fmt.Println("\nYeni bloklar ekleniyor...")
	
	err = bc.AddBlock("İlk işlem: Alice'den Bob'a 50 coin transfer", validator1)
	if err != nil {
		log.Fatal("Blok eklenemedi:", err)
	}
	time.Sleep(1 * time.Second)
	
	err = bc.AddBlock("İkinci işlem: Bob'dan Charlie'ye 30 coin transfer", validator2)
	if err != nil {
		log.Fatal("Blok eklenemedi:", err)
	}
	time.Sleep(1 * time.Second)
	
	err = bc.AddBlock("Üçüncü işlem: Charlie'den David'e 20 coin transfer", validator1)
	if err != nil {
		log.Fatal("Blok eklenemedi:", err)
	}

	// Blockchain'i görüntüle
	fmt.Println("\nBlockchain'deki tüm bloklar:")
	for i, block := range bc.Blocks {
		fmt.Printf("\nBlok %d:\n", i)
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Data: %s\n", string(block.Data))
		fmt.Printf("Hash: %s\n", block.GetHashString())
		fmt.Printf("Validator: %s\n", block.ValidatorAddress[:10])
		if i > 0 {
			fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		}
	}

	// Blockchain'in geçerliliğini kontrol et
	fmt.Printf("\nBlockchain geçerli mi? %v\n", bc.IsValid())
} 