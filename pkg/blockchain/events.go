package blockchain

import (
    "time"
    "github.com/ethereum/go-ethereum/common"
)

// EventType represents different types of blockchain events
type EventType string

const (
    // Contract related events
    EventContractDeployStarted  EventType = "CONTRACT_DEPLOY_STARTED"
    EventContractDeploySuccess  EventType = "CONTRACT_DEPLOY_SUCCESS"
    EventContractDeployFailed   EventType = "CONTRACT_DEPLOY_FAILED"
    EventContractVerified       EventType = "CONTRACT_VERIFIED"
)

// Event represents a blockchain event
type Event struct {
    Type      EventType           `json:"type"`
    Timestamp int64              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}

// ContractEvent represents a smart contract event
type ContractEvent struct {
    Address     common.Address    `json:"address"`
    Name        string           `json:"name"`
    Args        map[string]interface{} `json:"args"`
    BlockNumber uint64           `json:"blockNumber"`
    TxHash      common.Hash      `json:"transactionHash"`
    Timestamp   int64           `json:"timestamp"`
}

// EventEmitter handles event emission and subscription
type EventEmitter struct {
    subscribers map[EventType][]chan Event
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
    return &EventEmitter{
        subscribers: make(map[EventType][]chan Event),
    }
}

// Subscribe to specific event type
func (e *EventEmitter) Subscribe(eventType EventType) chan Event {
    ch := make(chan Event, 100)
    e.subscribers[eventType] = append(e.subscribers[eventType], ch)
    return ch
}

// Unsubscribe from specific event type
func (e *EventEmitter) Unsubscribe(eventType EventType, ch chan Event) {
    if subs, ok := e.subscribers[eventType]; ok {
        for i, sub := range subs {
            if sub == ch {
                e.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
                close(ch)
                break
            }
        }
    }
}

// Emit an event
func (e *EventEmitter) Emit(eventType EventType, data map[string]interface{}) {
    event := Event{
        Type:      eventType,
        Timestamp: time.Now().Unix(),
        Data:      data,
    }

    if subs, ok := e.subscribers[eventType]; ok {
        for _, ch := range subs {
            select {
            case ch <- event:
            default:
                // Channel is full, skip
            }
        }
    }
}

// EmitContractEvent emits a contract-specific event
func (e *EventEmitter) EmitContractEvent(event ContractEvent) {
    data := map[string]interface{}{
        "address":     event.Address,
        "name":        event.Name,
        "args":        event.Args,
        "blockNumber": event.BlockNumber,
        "txHash":      event.TxHash,
        "timestamp":   event.Timestamp,
    }

    e.Emit(EventType(event.Name), data)
} 