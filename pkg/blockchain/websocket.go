package blockchain

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // In production, implement proper origin checking
    },
}

// WebSocketServer handles WebSocket connections
type WebSocketServer struct {
    blockchain  *Blockchain
    clients     map[*websocket.Conn]bool
    clientsLock sync.RWMutex
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(bc *Blockchain) *WebSocketServer {
    return &WebSocketServer{
        blockchain: bc,
        clients:    make(map[*websocket.Conn]bool),
    }
}

// HandleWebSocket handles WebSocket connections
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Printf("Failed to upgrade connection: %v\n", err)
        return
    }

    // Register client
    s.clientsLock.Lock()
    s.clients[conn] = true
    s.clientsLock.Unlock()

    // Subscribe to blockchain events
    eventChan := s.blockchain.EventEmitter.Subscribe(EventContractDeployStarted)
    successChan := s.blockchain.EventEmitter.Subscribe(EventContractDeploySuccess)
    failedChan := s.blockchain.EventEmitter.Subscribe(EventContractDeployFailed)
    verifiedChan := s.blockchain.EventEmitter.Subscribe(EventContractVerified)

    // Handle client messages
    go func() {
        defer func() {
            conn.Close()
            s.clientsLock.Lock()
            delete(s.clients, conn)
            s.clientsLock.Unlock()

            // Unsubscribe from events
            s.blockchain.EventEmitter.Unsubscribe(EventContractDeployStarted, eventChan)
            s.blockchain.EventEmitter.Unsubscribe(EventContractDeploySuccess, successChan)
            s.blockchain.EventEmitter.Unsubscribe(EventContractDeployFailed, failedChan)
            s.blockchain.EventEmitter.Unsubscribe(EventContractVerified, verifiedChan)
        }()

        for {
            select {
            case event := <-eventChan:
                s.broadcastEvent(conn, event)
            case event := <-successChan:
                s.broadcastEvent(conn, event)
            case event := <-failedChan:
                s.broadcastEvent(conn, event)
            case event := <-verifiedChan:
                s.broadcastEvent(conn, event)
            }
        }
    }()
}

// broadcastEvent sends an event to a specific client
func (s *WebSocketServer) broadcastEvent(conn *websocket.Conn, event Event) {
    message := map[string]interface{}{
        "type": event.Type,
        "data": event.Data,
        "timestamp": event.Timestamp,
    }

    jsonBytes, err := json.Marshal(message)
    if err != nil {
        fmt.Printf("Error marshaling message: %v\n", err)
        return
    }

    err = conn.WriteMessage(websocket.TextMessage, jsonBytes)
    if err != nil {
        fmt.Printf("Error sending message to client: %v\n", err)
        conn.Close()
        s.clientsLock.Lock()
        delete(s.clients, conn)
        s.clientsLock.Unlock()
    }
}

// BroadcastToAll sends an event to all connected clients
func (s *WebSocketServer) BroadcastToAll(event Event) {
    s.clientsLock.RLock()
    defer s.clientsLock.RUnlock()

    for conn := range s.clients {
        s.broadcastEvent(conn, event)
    }
} 