export type ContractType = 'erc20' | 'erc721' | 'defi' | 'other';

export interface BlockchainContract {
  address: string;
  name: string;
  type: ContractType;
  status: 'completed' | 'pending' | 'failed';
  owner: string;
  transactions: number;
  deployedAt: number;
  balance: string;
  abi: ContractABI[];
  bytecode?: string;
  sourceCode?: string;
  verified?: boolean;
}

export interface ContractTemplate {
  name: string;
  type: ContractType;
  description: string;
  code: string;
}

export interface ContractABI {
  name: string;
  type: string;
  inputs: ABIParameter[];
  outputs?: ABIParameter[];
  stateMutability?: string;
}

export interface ABIParameter {
  name: string;
  type: string;
  components?: ABIParameter[];
}

export interface ContractEvent {
  name: string;
  transactionHash: string;
  blockNumber: number;
  timestamp: number;
  args: Record<string, any>;
}

export interface ContractTransaction {
  hash: string;
  method?: string;
  from: string;
  to: string;
  value: string;
  status: 'success' | 'failed';
  timestamp: number;
  blockNumber: number;
  gasUsed: string;
  gasPrice: string;
}

export interface BlockchainInfo {
  blocks: number;
  active_validators: number;
  validator_count: number;
  current_block: number;
  is_valid: boolean;
  pending_transactions: number;
}

export interface Block {
  height: number;
  hash: string;
  timestamp: number;
  transactions: number;
  validator: string;
}

export interface Transaction {
  hash: string;
  from: string;
  to: string;
  value: string;
  type: 'transfer' | 'contract' | 'validator';
  timestamp: number;
  status: 'success' | 'failed' | 'pending';
}

export interface ContractTransactionsResponse {
  transactions: ContractTransaction[];
  total: number;
} 