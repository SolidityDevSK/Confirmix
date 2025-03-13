import { ethers } from 'ethers';
import type { ContractEvent, BlockchainContract } from '@/types/blockchain';

type EventCallback = (event: ContractEvent) => void;
type ContractStatusCallback = (status: ContractStatus) => void;

type ContractStatus = {
  type: 'CONTRACT_DEPLOY_STARTED' | 'CONTRACT_DEPLOY_SUCCESS' | 'CONTRACT_DEPLOY_FAILED' | 'CONTRACT_VERIFIED' | 'CONTRACT_EVENT';
  address?: string;
  timestamp: number;
  data: any;
};

type PendingSubscription = {
  channel: string;
  abi: any[];
  callback: EventCallback;
};

export class WebSocketService {
  private ws: WebSocket | null = null;
  private url: string;
  private pendingSubscriptions: PendingSubscription[] = [];
  private statusSubscriptions: Map<string, ContractStatusCallback[]> = new Map();
  private isConnecting: boolean = false;
  private reconnectTimeout: NodeJS.Timeout | null = null;
  private maxReconnectAttempts = 5;
  private reconnectAttempts = 0;
  private subscriptions: Map<string, Set<(data: any) => void>> = new Map();

  constructor() {
    const wsUrl = process.env.NEXT_PUBLIC_WEBSOCKET_URL;
    if (!wsUrl) {
      console.warn('WebSocket URL not configured. Using default: ws://localhost:8080/ws');
    }
    this.url = wsUrl || 'ws://localhost:8080/ws';
    this.connect();
  }

  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  public async connect(): Promise<void> {
    if (this.isConnecting) {
      console.log('WebSocket connection already in progress');
      return;
    }

    if (this.isConnected()) {
      console.log('WebSocket already connected');
      return;
    }

    this.isConnecting = true;

    try {
      const wsUrl = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost:8080/ws';
      console.log('Connecting to WebSocket:', wsUrl);
      
      this.ws = new WebSocket(wsUrl);

      await new Promise<void>((resolve, reject) => {
        if (!this.ws) return reject(new Error('WebSocket not initialized'));

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.isConnecting = false;
          this.reconnectAttempts = 0;
          resolve();
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          this.isConnecting = false;
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('WebSocket closed');
          this.isConnecting = false;
          this.ws = null;
          
          if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Reconnecting (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
            setTimeout(() => this.connect(), 1000 * this.reconnectAttempts);
          }
        };

        this.ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            console.log('WebSocket message received:', message);
            
            // Handle contract status messages
            if (message.type && message.type.startsWith('CONTRACT_')) {
              const contractAddress = message.data?.address;
              const callbacks = this.statusSubscriptions.get(contractAddress || 'global') || [];
              
              // Create status object matching the expected format
              const status: ContractStatus = {
                type: message.type,
                timestamp: message.timestamp,
                data: message.data,
                address: contractAddress
              };

              console.log('Notifying subscribers for status:', status);
              callbacks.forEach(callback => callback(status));

              // Also notify global subscribers
              if (contractAddress) {
                const globalCallbacks = this.statusSubscriptions.get('global') || [];
                globalCallbacks.forEach(callback => callback(status));
              }
            }
          } catch (err) {
            console.error('Error processing WebSocket message:', err);
          }
        };
      });
    } catch (err) {
      console.error('Error connecting to WebSocket:', err);
      this.isConnecting = false;
      throw err;
    }
  }

  public subscribeToContractStatus(
    address: string | 'global',
    callback: ContractStatusCallback
  ): () => void {
    const callbacks = this.statusSubscriptions.get(address) || [];
    callbacks.push(callback);
    this.statusSubscriptions.set(address, callbacks);

    // Return unsubscribe function
    return () => {
      const callbacks = this.statusSubscriptions.get(address) || [];
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
        if (callbacks.length === 0) {
          this.statusSubscriptions.delete(address);
        } else {
          this.statusSubscriptions.set(address, callbacks);
        }
      }
    };
  }

  public subscribeToContractEvents(
    address: string,
    abi: any[],
    callback: EventCallback
  ) {
    console.log('Subscribing to contract events:', address);
    const channel = `contract:${address}`;
    this.sendSubscription(channel, abi, callback);
  }

  private sendSubscription(channel: string, abi: any[], callback: EventCallback) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.log('WebSocket not ready, queueing subscription for:', channel);
      this.pendingSubscriptions.push({ channel, abi, callback });
      return;
    }

    try {
      const message = JSON.stringify({
        type: 'subscribe',
        channel,
        abi
      });
      
      console.log('Sending subscription:', message);
      this.ws.send(message);
    } catch (error) {
      console.error('Error sending subscription:', error);
      this.pendingSubscriptions.push({ channel, abi, callback });
    }
  }

  public unsubscribeFromContractEvents(address: string) {
    console.log('Unsubscribing from contract events:', address);
    
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      this.pendingSubscriptions = this.pendingSubscriptions.filter(
        sub => sub.channel !== `contract:${address}`
      );
      return;
    }

    try {
      this.ws.send(JSON.stringify({
        type: 'unsubscribe',
        channel: `contract:${address}`
      }));
    } catch (error) {
      console.error('Error unsubscribing from events:', error);
    }
  }

  public disconnect() {
    console.log('Disconnecting WebSocket service');
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
    }
    this.pendingSubscriptions = [];
    this.statusSubscriptions.clear();
    this.reconnectAttempts = 0;
  }
}

export const websocketService = new WebSocketService(); 