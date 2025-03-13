export interface ContractDetail {
  address: string;
  name: string;
  type: 'token' | 'nft' | 'defi' | 'other';
  status: 'active' | 'inactive';
  owner: string;
  transactions: number;
  deployedAt: number;
  balance: string;
  abi: Array<{
    name: string;
    type: string;
    inputs: Array<{
      name: string;
      type: string;
    }>;
    outputs?: Array<{
      type: string;
    }>;
  }>;
  bytecode: string;
  sourceCode: string;
} 