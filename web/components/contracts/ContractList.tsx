'use client';

import { Table, Text, Badge, Button, Flex, TextField, Select } from '@radix-ui/themes';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import { blockchainService } from '@/services/blockchain';
import { websocketService } from '@/services/websocket';
import type { BlockchainContract, ContractType } from '@/types/blockchain';

export default function ContractList() {
  const router = useRouter();
  const [contracts, setContracts] = useState<BlockchainContract[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [deployDialogOpen, setDeployDialogOpen] = useState(false);

  useEffect(() => {
    // Initial fetch of contracts
    fetchContracts();

    // Subscribe to WebSocket messages
    console.log('Setting up WebSocket subscription for contract status updates');
    const unsubscribe = websocketService.subscribeToContractStatus('global', (status) => {
      console.log('Received contract status:', status);

      if (status.type === 'CONTRACT_DEPLOY_STARTED') {
        // Add new contract in pending state
        const newContract: BlockchainContract = {
          name: status.data.name,
          type: status.data.type || 'other',
          status: 'pending',
          address: '',
          owner: status.data.owner || '',
          transactions: 0,
          deployedAt: Math.floor(Date.now() / 1000),
          verified: false,
          balance: '0',
          abi: []
        };
        console.log('Adding new pending contract:', newContract);
        setContracts(prev => [...prev, newContract]);
      }
      else if (status.type === 'CONTRACT_DEPLOY_SUCCESS') {
        // Update contract with address and completed status
        setContracts(prev => prev.map(contract => 
          contract.name === status.data.name && contract.status === 'pending'
            ? { 
                ...contract, 
                status: 'completed', 
                address: status.data.address,
                owner: status.data.owner || contract.owner
              }
            : contract
        ));
        console.log('Updated contract status to completed:', status.data.name);
      }
    });

    return () => unsubscribe();
  }, []);

  const fetchContracts = async () => {
    try {
      setLoading(true);
      setError(null);
      const filters = {
        search: searchQuery || undefined,
        type: typeFilter !== 'all' ? typeFilter : undefined,
        status: statusFilter !== 'all' ? statusFilter : undefined,
      };
      const data = await blockchainService.getContracts(filters);
      console.log('Fetched contracts from backend:', {
        contracts: data.map(c => ({
          name: c.name,
          address: c.address,
          status: c.status,
          verified: c.verified
        })),
        timestamp: new Date().toISOString()
      });
      
      // Update contracts list, ensuring verified contracts are marked as completed
      const updatedContracts = data.map(contract => {
        const shouldComplete = contract.verified || contract.status === 'completed';
        if (shouldComplete) {
          console.log('Marking contract as completed:', {
            name: contract.name,
            address: contract.address,
            reason: contract.verified ? 'verified' : 'already completed'
          });
        }
        return {
          ...contract,
          status: shouldComplete ? ('completed' as const) : contract.status
        };
      });
      
      setContracts(updatedContracts);
    } catch (err) {
      setError('Failed to fetch contracts. Please try again later.');
      console.error('Error fetching contracts:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDeploySuccess = () => {
    fetchContracts(); // Refresh the list after successful deployment
  };

  const formatTimeAgo = (timestamp: number) => {
    if (!timestamp) return 'Unknown';
    
    // Convert seconds to milliseconds if needed
    const timestampMs = timestamp * 1000;
    const seconds = Math.floor((Date.now() - timestampMs) / 1000);
    
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes} minutes ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours} hours ago`;
    const days = Math.floor(hours / 24);
    return `${days} days ago`;
  };

  const getTypeBadge = (type: BlockchainContract['type']) => {
    const colors: Record<ContractType, 'blue' | 'purple' | 'green' | 'gray'> = {
      erc20: 'blue',
      erc721: 'purple',
      defi: 'green',
      other: 'gray'
    };

    const displayNames: Record<ContractType, string> = {
      erc20: 'ERC-20',
      erc721: 'ERC-721',
      defi: 'DeFi',
      other: 'Other'
    };

    return (
      <Badge color={colors[type]} size="1">
        {displayNames[type]}
      </Badge>
    );
  };

  const getStatusBadge = (status: BlockchainContract['status']) => {
    const colors: Record<BlockchainContract['status'], 'green' | 'yellow' | 'red'> = {
      completed: 'green' as const,
      pending: 'yellow' as const,
      failed: 'red' as const
    };

    const displayNames: Record<BlockchainContract['status'], string> = {
      completed: 'Completed',
      pending: 'Pending',
      failed: 'Failed'
    };

    return (
      <Badge color={colors[status]} size="1">
        {displayNames[status]}
      </Badge>
    );
  };

  const handleDetailsClick = (address: string) => {
    router.push(`/contracts/${address}`);
  };

  if (error) {
    return (
      <Flex align="center" justify="center" style={{ minHeight: 200 }}>
        <Text color="red" size="2">
          {error}
        </Text>
      </Flex>
    );
  }

  return (
    <div>
      <Flex direction="column" gap="4">
        <Flex gap="4" align="end">
          <div className="flex-1">
            <Text as="div" size="2" mb="1" weight="bold">
              Search
            </Text>
            <div className="relative">
              <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400">
                <MagnifyingGlassIcon height={16} width={16} />
              </span>
              <input
                className="w-full pl-9 pr-3 py-2 bg-gray-900 border border-gray-800 rounded-md text-white"
                placeholder="Search by name, address, or owner"
                value={searchQuery}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
              />
            </div>
          </div>

          <div>
            <Text as="div" size="2" mb="1" weight="bold">
              Type
            </Text>
            <Select.Root value={typeFilter} onValueChange={setTypeFilter}>
              <Select.Trigger />
              <Select.Content>
                <Select.Item value="all">All Types</Select.Item>
                <Select.Item value="erc20">ERC-20</Select.Item>
                <Select.Item value="erc721">ERC-721</Select.Item>
                <Select.Item value="defi">DeFi</Select.Item>
                <Select.Item value="other">Other</Select.Item>
              </Select.Content>
            </Select.Root>
          </div>

          <div>
            <Text as="div" size="2" mb="1" weight="bold">
              Status
            </Text>
            <Select.Root value={statusFilter} onValueChange={setStatusFilter}>
              <Select.Trigger />
              <Select.Content>
                <Select.Item value="all">All Status</Select.Item>
                <Select.Item value="completed">Completed</Select.Item>
                <Select.Item value="pending">Pending</Select.Item>
                <Select.Item value="failed">Failed</Select.Item>
              </Select.Content>
            </Select.Root>
          </div>
        </Flex>

        {loading ? (
          <Flex align="center" justify="center" style={{ minHeight: 200 }}>
            <Text>Loading contracts...</Text>
          </Flex>
        ) : (
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Address</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Owner</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Transactions</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Deployed</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
              </Table.Row>
            </Table.Header>

            <Table.Body>
              {contracts.map((contract) => (
                <Table.Row key={contract.address}>
                  <Table.Cell>
                    <Text weight="bold">{contract.name}</Text>
                  </Table.Cell>
                  <Table.Cell>
                    <Text className="font-mono" color="gray">
                      {contract.address}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>{getTypeBadge(contract.type)}</Table.Cell>
                  <Table.Cell>{getStatusBadge(contract.status)}</Table.Cell>
                  <Table.Cell>
                    <Text className="font-mono" color="gray">
                      {contract.owner}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>{contract.transactions}</Table.Cell>
                  <Table.Cell>{formatTimeAgo(contract.deployedAt)}</Table.Cell>
                  <Table.Cell>
                    <Flex gap="2">
                      <Button 
                        size="1" 
                        variant="soft" 
                        color="blue"
                        onClick={() => handleDetailsClick(contract.address)}
                      >
                        Details
                      </Button>
                      <Button 
                        size="1" 
                        variant="soft" 
                        color="green"
                        onClick={() => router.push(`/contracts/${contract.address}/interact`)}
                      >
                        Interact
                      </Button>
                    </Flex>
                  </Table.Cell>
                </Table.Row>
              ))}

              {contracts.length === 0 && (
                <Table.Row>
                  <Table.Cell colSpan={8}>
                    <Text align="center" color="gray" className="py-4">
                      No contracts found matching your filters
                    </Text>
                  </Table.Cell>
                </Table.Row>
              )}
            </Table.Body>
          </Table.Root>
        )}
      </Flex>
    </div>
  );
} 