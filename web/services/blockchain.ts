import type { BlockchainContract, BlockchainInfo, ContractEvent, ContractType } from '@/types/blockchain';
import { ContractTransactionsResponse } from '@/types/blockchain';

class BlockchainService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_BLOCKCHAIN_API_URL || 'http://localhost:8080';
  }

  // Get blockchain information
  async getInfo(): Promise<BlockchainInfo> {
    const response = await fetch(`${this.baseUrl}/info`);
    if (!response.ok) {
      throw new Error('Failed to fetch blockchain info');
    }
    return response.json();
  }

  // Get blocks
  async getBlocks() {
    const response = await fetch(`${this.baseUrl}/blocks`);
    if (!response.ok) {
      throw new Error('Failed to fetch blocks');
    }
    return response.json();
  }

  // Get transactions
  async getTransactions() {
    const response = await fetch(`${this.baseUrl}/transactions`);
    if (!response.ok) {
      throw new Error('Failed to fetch transactions');
    }
    return response.json();
  }

  // Get block by hash
  async getBlock(hash: string) {
    const response = await fetch(`${this.baseUrl}/blocks/${hash}`);
    if (!response.ok) {
      throw new Error('Failed to fetch block');
    }
    return response.json();
  }

  // Get validators
  async getValidators() {
    const response = await fetch(`${this.baseUrl}/validators`);
    if (!response.ok) {
      throw new Error('Failed to fetch validators');
    }
    return response.json();
  }

  // Get current validator
  async getCurrentValidator() {
    const response = await fetch(`${this.baseUrl}/validators/current`);
    if (!response.ok) {
      throw new Error('Failed to fetch current validator');
    }
    return response.json();
  }

  // Add validator
  async addValidator() {
    const response = await fetch(`${this.baseUrl}/validators`, {
      method: 'POST',
    });
    if (!response.ok) {
      throw new Error('Failed to add validator');
    }
    return response.json();
  }

  // Remove validator
  async removeValidator(address: string) {
    const response = await fetch(`${this.baseUrl}/validators/${address}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('Failed to remove validator');
    }
    return response.json();
  }

  // Submit transaction
  async submitTransaction(data: string, validator: string) {
    const response = await fetch(`${this.baseUrl}/transactions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ data, validator }),
    });
    if (!response.ok) {
      throw new Error('Failed to submit transaction');
    }
    return response.json();
  }

  // Compile contract
  async compileContract(sourceCode: string) {
    console.log('Compiling contract source code...');
    
    try {
      const response = await fetch('/api/compile', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sourceCode }),
      });
      
      if (!response.ok) {
        const error = await response.json();
        console.error('Contract compilation failed:', error);
        throw new Error(error.details || 'Failed to compile contract');
      }
      
      const result = await response.json();
      console.log('Contract compiled successfully:', {
        abiLength: result.abi?.length,
        bytecodeLength: result.bytecode?.length
      });
      return result;
    } catch (error) {
      console.error('Contract compilation error:', error);
      throw error;
    }
  }

  // Deploy contract
  async deployContract(data: {
    name: string;
    type: ContractType;
    code: string;
    abi: any[];
    constructorArgs: any[];
    owner: string;
    version: string;
  }): Promise<BlockchainContract> {
    const response = await fetch(`${this.baseUrl}/contracts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: data.name,
        type: data.type,
        code: data.code,
        abi: data.abi,
        constructorArgs: data.constructorArgs,
        owner: data.owner,
        version: data.version
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(`Failed to deploy contract: ${JSON.stringify(errorData)}`);
    }

    return await response.json();
  }

  // List contracts
  async listContracts() {
    const response = await fetch(`${this.baseUrl}/contracts`);
    if (!response.ok) {
      throw new Error('Failed to fetch contracts');
    }
    return response.json();
  }

  // Get contract
  async getContract(address: string) {
    const response = await fetch(`${this.baseUrl}/contracts/${address}`);
    if (!response.ok) {
      throw new Error('Failed to fetch contract');
    }
    return response.json();
  }

  // Execute contract
  async executeContract(address: string, input: string) {
    const response = await fetch(`${this.baseUrl}/contracts/${address}/execute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ input }),
    });
    if (!response.ok) {
      throw new Error('Failed to execute contract');
    }
    return response.json();
  }

  async callContractMethod(
    address: string,
    method: string,
    args: any[] = []
  ): Promise<any> {
    const response = await fetch(`${this.baseUrl}/contracts/${address}/execute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ method, args }),
    });
    if (!response.ok) {
      throw new Error('Failed to call contract method');
    }
    return response.json();
  }

  // Get contract transactions
  async getContractTransactions(
    address: string,
    page: number,
    pageSize: number,
    filters?: {
      fromTimestamp?: number;
      search?: string;
      method?: string;
      status?: 'success' | 'failed';
      sortBy?: 'timestamp' | 'value' | 'gasUsed';
      sortDirection?: 'asc' | 'desc';
    }
  ): Promise<ContractTransactionsResponse> {
    const params = new URLSearchParams();
    params.append('page', page.toString());
    params.append('pageSize', pageSize.toString());
    
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined) {
          params.append(key, value.toString());
        }
      });
    }

    const response = await fetch(`${this.baseUrl}/contracts/${address}/transactions?${params}`);
    if (!response.ok) {
      throw new Error('Failed to fetch contract transactions');
    }
    return response.json();
  }

  async getContractEvents(address: string, fromTimestamp: number): Promise<ContractEvent[]> {
    const params = new URLSearchParams();
    params.append('fromTimestamp', fromTimestamp.toString());

    const response = await fetch(`${this.baseUrl}/contracts/${address}/events?${params}`);
    if (!response.ok) throw new Error('Failed to fetch contract events');
    return response.json();
  }

  async getContracts(filters?: { search?: string; type?: string; status?: string }): Promise<BlockchainContract[]> {
    const params = new URLSearchParams();
    if (filters?.search) params.append('search', filters.search);
    if (filters?.type) params.append('type', filters.type);
    if (filters?.status) params.append('status', filters.status);

    const response = await fetch(`${this.baseUrl}/contracts${params.toString() ? `?${params.toString()}` : ''}`);
    if (!response.ok) throw new Error('Failed to fetch contracts');
    
    const data = await response.json();
    
    // Transform and validate the data
    return data.map((contract: any) => ({
      ...contract,
      type: contract.type === 'token' ? 'erc20' : (contract.type === 'nft' ? 'erc721' : (contract.type || 'other')),
      status: contract.status === 'active' ? 'completed' : (contract.status || 'pending'),
      deployedAt: contract.deployedAt ? Number(contract.deployedAt) : Math.floor(Date.now() / 1000),
      transactions: contract.transactions || 0
    }));
  }

  async verifyContract(address: string, sourceCode: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/contracts/${address}/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ sourceCode }),
    });
    if (!response.ok) {
      throw new Error('Failed to verify contract');
    }
  }

  async estimateContractDeployment(
    bytecode: string,
    abi: any[],
    constructorArgs?: Record<string, string>
  ): Promise<{ gasLimit: string; gasPrice: string; totalCost: string }> {
    console.log('Estimating contract deployment:', { bytecode: bytecode.slice(0, 32) + '...', constructorArgs });
    
    try {
      const validator = await this.getCurrentValidator();
      
      const response = await fetch(`${this.baseUrl}/contracts/estimate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          code: bytecode.slice(2),
          owner: validator.address,
          abi,
          constructorArgs 
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Gas estimation failed:', {
          status: response.status,
          statusText: response.statusText,
          error: errorText
        });
        throw new Error('Failed to estimate contract deployment');
      }

      const result = await response.json();
      console.log('Gas estimation successful:', result);
      return result;
    } catch (error) {
      console.error('Gas estimation error:', error);
      throw error;
    }
  }

  async getTransactionStatus(hash: string) {
    try {
      const response = await fetch(`${this.baseUrl}/transactions/${hash}/status`);
      if (!response.ok) {
        if (response.status === 404) {
          // Transaction henüz işlenmemiş
          return {
            confirmations: 0,
            status: 'pending',
            blockNumber: null,
            gasUsed: null
          };
        }
        throw new Error('Failed to fetch transaction status');
      }
      
      const result = await response.json();
      
      // Eğer transaction başarıyla tamamlandıysa
      if (result.status === 'success' || result.blockNumber) {
        return {
          confirmations: 1, // En az 1 confirmation var
          status: 'success',
          blockNumber: result.blockNumber,
          gasUsed: result.gasUsed
        };
      }
      
      // Transaction hala pending durumunda
      return {
        confirmations: 0,
        status: 'pending',
        blockNumber: null,
        gasUsed: null
      };
    } catch (error) {
      console.error('Error getting transaction status:', error);
      // Hata durumunda varsayılan pending durumu
      return {
        confirmations: 0,
        status: 'pending',
        blockNumber: null,
        gasUsed: null
      };
    }
  }

  async getContractAnalytics(address: string, timeRange: '24h' | '7d' | '30d') {
    const response = await fetch(`${this.baseUrl}/contracts/${address}/analytics?timeRange=${timeRange}`);
    if (!response.ok) {
      throw new Error('Failed to fetch contract analytics');
    }
    return response.json();
  }
}

export const blockchainService = new BlockchainService(); 