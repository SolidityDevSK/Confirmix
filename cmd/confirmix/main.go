package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SolidityDevSK/Confirmix/internal/validator"
	"github.com/SolidityDevSK/Confirmix/pkg/api"
	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
	"github.com/SolidityDevSK/Confirmix/pkg/network"
)

func main() {
	// Komut satırı parametrelerini tanımla
	apiPort := flag.Int("api-port", 8080, "HTTP API port")
	p2pPort := flag.Int("p2p-port", 9000, "P2P network port")
	bootstrapNode := flag.String("bootstrap", "", "Bootstrap node address")
	flag.Parse()

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

	// P2P node'unu başlat
	node, err := network.NewNode(*p2pPort, bc)
	if err != nil {
		log.Fatal("P2P node oluşturulamadı:", err)
	}
	defer node.Close()

	// Bootstrap node'una bağlan
	if *bootstrapNode != "" {
		ctx := context.Background()
		if err := node.Connect(ctx, *bootstrapNode); err != nil {
			log.Printf("Bootstrap node'una bağlanılamadı: %v", err)
		}
	}

	// HTTP API sunucusunu başlat
	server := api.NewServer(bc)
	go func() {
		log.Printf("Genesis Validator Address: %s", genesisValidator.Address)
		log.Printf("Second Validator Address: %s", validator2.Address)
		log.Printf("P2P Node Address: %s", node.GetMultiaddr())
		log.Printf("HTTP API sunucusu başlatılıyor: http://localhost:%d", *apiPort)
		if err := server.Run(fmt.Sprintf(":%d", *apiPort)); err != nil {
			log.Fatal("Sunucu başlatılamadı:", err)
		}
	}()

	// Graceful shutdown için sinyal bekle
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Uygulama kapatılıyor...")
} 