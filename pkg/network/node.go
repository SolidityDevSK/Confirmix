package network

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"

	"github.com/SolidityDevSK/Confirmix/pkg/blockchain"
)

const (
	ProtocolID          = "/confirmix/1.0.0"
	BlockchainSync      = "/blockchain/sync"
	BlockAnnouncement   = "/block/announcement"
	ValidatorAnnouncement = "/validator/announcement"
)

// Message represents a network message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Node represents a P2P network node
type Node struct {
	host       host.Host
	blockchain *blockchain.Blockchain
	peers      map[peer.ID]struct{}
	mu         sync.RWMutex
}

// NewNode creates a new P2P node
func NewNode(listenPort int, bc *blockchain.Blockchain) (*Node, error) {
	// Listen adresi oluştur
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort))
	if err != nil {
		return nil, err
	}

	// Host oluştur
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
	)
	if err != nil {
		return nil, err
	}

	node := &Node{
		host:       host,
		blockchain: bc,
		peers:      make(map[peer.ID]struct{}),
	}

	// Stream handler'ları ayarla
	node.host.SetStreamHandler(protocol.ID(BlockchainSync), node.handleBlockchainSync)
	node.host.SetStreamHandler(protocol.ID(BlockAnnouncement), node.handleBlockAnnouncement)
	node.host.SetStreamHandler(protocol.ID(ValidatorAnnouncement), node.handleValidatorAnnouncement)

	return node, nil
}

// Connect connects to a peer
func (n *Node) Connect(ctx context.Context, peerAddr string) error {
	// Peer adresini parse et
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return err
	}

	// Peer bilgisini çıkar
	peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	// Peer'a bağlan
	if err := n.host.Connect(ctx, *peerInfo); err != nil {
		return err
	}

	// Peer'ı listeye ekle
	n.mu.Lock()
	n.peers[peerInfo.ID] = struct{}{}
	n.mu.Unlock()

	// Blockchain senkronizasyonunu başlat
	return n.syncBlockchain(ctx, peerInfo.ID)
}

// Broadcast broadcasts a message to all peers
func (n *Node) Broadcast(ctx context.Context, msgType string, payload interface{}) error {
	msg := Message{
		Type:    msgType,
		Payload: payload,
	}

	// Mesajı JSON'a çevir
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Tüm peer'lara gönder
	n.mu.RLock()
	defer n.mu.RUnlock()

	for peerID := range n.peers {
		stream, err := n.host.NewStream(ctx, peerID, protocol.ID(msgType))
		if err != nil {
			fmt.Printf("Failed to create stream to peer %s: %s\n", peerID, err)
			continue
		}
		defer stream.Close()

		_, err = stream.Write(data)
		if err != nil {
			fmt.Printf("Failed to write to stream: %s\n", err)
		}
	}

	return nil
}

// handleBlockchainSync handles blockchain sync requests
func (n *Node) handleBlockchainSync(stream network.Stream) {
	defer stream.Close()

	// Blockchain'i JSON'a çevir
	data, err := json.Marshal(n.blockchain)
	if err != nil {
		fmt.Printf("Failed to marshal blockchain: %s\n", err)
		return
	}

	// Blockchain'i gönder
	_, err = stream.Write(data)
	if err != nil {
		fmt.Printf("Failed to write blockchain: %s\n", err)
	}
}

// handleBlockAnnouncement handles new block announcements
func (n *Node) handleBlockAnnouncement(stream network.Stream) {
	defer stream.Close()

	// Mesajı oku
	var msg Message
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&msg); err != nil {
		fmt.Printf("Failed to decode message: %s\n", err)
		return
	}

	// Blok verisini parse et
	blockData, err := json.Marshal(msg.Payload)
	if err != nil {
		fmt.Printf("Failed to marshal block data: %s\n", err)
		return
	}

	var block blockchain.Block
	if err := json.Unmarshal(blockData, &block); err != nil {
		fmt.Printf("Failed to unmarshal block: %s\n", err)
		return
	}

	// TODO: Bloğu doğrula ve blockchain'e ekle
}

// handleValidatorAnnouncement handles new validator announcements
func (n *Node) handleValidatorAnnouncement(stream network.Stream) {
	defer stream.Close()

	// Mesajı oku
	var msg Message
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&msg); err != nil {
		fmt.Printf("Failed to decode message: %s\n", err)
		return
	}

	// Validator verisini parse et
	validatorData, err := json.Marshal(msg.Payload)
	if err != nil {
		fmt.Printf("Failed to marshal validator data: %s\n", err)
		return
	}

	// TODO: Validator'ı doğrula ve ekle
}

// syncBlockchain syncs the blockchain with a peer
func (n *Node) syncBlockchain(ctx context.Context, peerID peer.ID) error {
	// Blockchain sync stream'i oluştur
	stream, err := n.host.NewStream(ctx, peerID, protocol.ID(BlockchainSync))
	if err != nil {
		return err
	}
	defer stream.Close()

	// Blockchain verisini oku
	var remoteChain blockchain.Blockchain
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&remoteChain); err != nil {
		return err
	}

	// TODO: Blockchain'leri karşılaştır ve senkronize et
	return nil
}

// GetMultiaddr returns the node's multiaddress
func (n *Node) GetMultiaddr() string {
	return fmt.Sprintf("%s/p2p/%s", n.host.Addrs()[0], n.host.ID())
}

// Close closes the node
func (n *Node) Close() error {
	return n.host.Close()
} 