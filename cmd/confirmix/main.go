package main

import (
	"log"
	"time"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/api"
	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
)

func main() {
	// Genesis validator'ı oluştur
	genesisValidator, err := validator.NewAuthority()
	if err != nil {
		log.Fatal("Genesis validator oluşturulamadı:", err)
	}

	// Blockchain'i başlat
	bc, err := blockchain.NewBlockchain(genesisValidator)
	if err != nil {
		log.Fatal("Blockchain oluşturulamadı:", err)
	}

	// Test için ikinci bir validator ekle
	validator2, err := validator.NewAuthority()
	if err != nil {
		log.Fatal("İkinci validator oluşturulamadı:", err)
	}
	bc.AddValidator(validator2)

	// HTTP API sunucusunu başlat
	server := api.NewServer(bc)
	log.Printf("Genesis Validator Address: %s", genesisValidator.Address)
	log.Printf("Second Validator Address: %s", validator2.Address)
	log.Printf("HTTP API sunucusu başlatılıyor: http://localhost:8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
} 