'use client';

import { Table, Text, Badge, Button, Flex, Select, Card, Grid, Box, Dialog } from '@radix-ui/themes';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ReloadIcon, MagnifyingGlassIcon, ArrowUpIcon, ArrowDownIcon } from '@radix-ui/react-icons';
import { blockchainService } from '@/services/blockchain';
import type { BlockchainContract, ContractTransaction } from '@/types/blockchain';
import { ethers } from 'ethers';

interface ContractTransactionsProps {
  contract: BlockchainContract;
}

interface TransactionFilters {
  search?: string;
  method?: string;
  status?: 'all' | 'success' | 'failed';
  timeRange: '1h' | '24h' | '7d' | '30d' | 'all';
  sortBy: 'timestamp' | 'value' | 'gasUsed';
  sortDirection: 'asc' | 'desc';
}

interface TransactionStats {
  totalValue: string;
  averageGasUsed: string;
  methodFrequency: Record<string, number>;
  successRate: number;
}

interface TransactionReceipt {
  blockHash: string;
  blockNumber: number;
  contractAddress: string | null;
  cumulativeGasUsed: string;
  effectiveGasPrice: string;
  gasUsed: string;
  logs: Array<{
    address: string;
    topics: string[];
    data: string;
    blockNumber: number;
    transactionHash: string;
    logIndex: number;
  }>;
  status: boolean;
  transactionHash: string;
  transactionIndex: number;
}

export default function ContractTransactions({ contract }: ContractTransactionsProps) {
  const router = useRouter();
  const [transactions, setTransactions] = useState<ContractTransaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState<TransactionFilters>({
    status: 'all',
    timeRange: '24h',
    sortBy: 'timestamp',
    sortDirection: 'desc',
  });
  const [stats, setStats] = useState<TransactionStats>({
    totalValue: '0',
    averageGasUsed: '0',
    methodFrequency: {},
    successRate: 0,
  });
  const [selectedTx, setSelectedTx] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<TransactionReceipt | null>(null);

  useEffect(() => {
    fetchTransactions();
  }, [contract.address, page, pageSize, filters]);

  useEffect(() => {
    if (transactions.length > 0) {
      calculateStats();
    }
  }, [transactions]);

  const fetchTransactions = async () => {
    try {
      setLoading(true);
      setError(null);

      // Calculate fromTimestamp based on timeRange
      let fromTimestamp: number | undefined;
      const now = Math.floor(Date.now() / 1000);
      switch (filters.timeRange) {
        case '1h':
          fromTimestamp = now - 3600;
          break;
        case '24h':
          fromTimestamp = now - 86400;
          break;
        case '7d':
          fromTimestamp = now - 604800;
          break;
        case '30d':
          fromTimestamp = now - 2592000;
          break;
        default:
          fromTimestamp = undefined;
      }

      const data = await blockchainService.getContractTransactions(
        contract.address,
        page,
        pageSize,
        {
          fromTimestamp,
          search: filters.search,
          method: filters.method,
          status: filters.status === 'all' ? undefined : filters.status,
          sortBy: filters.sortBy,
          sortDirection: filters.sortDirection,
        }
      );

      setTransactions(data.transactions);
      setTotalPages(Math.ceil(data.total / pageSize));
    } catch (err) {
      console.error('Error fetching transactions:', err);
      setError('Failed to fetch contract transactions');
    } finally {
      setLoading(false);
    }
  };

  const calculateStats = () => {
    const totalValue = transactions.reduce(
      (sum, tx) => sum + BigInt(tx.value),
      BigInt(0)
    ).toString();

    const totalGasUsed = transactions.reduce(
      (sum, tx) => sum + BigInt(tx.gasUsed),
      BigInt(0)
    );
    const averageGasUsed = (totalGasUsed / BigInt(transactions.length)).toString();

    const methodFrequency: Record<string, number> = {};
    transactions.forEach(tx => {
      if (tx.method) {
        methodFrequency[tx.method] = (methodFrequency[tx.method] || 0) + 1;
      }
    });

    const successCount = transactions.filter(tx => tx.status === 'success').length;
    const successRate = (successCount / transactions.length) * 100;

    setStats({
      totalValue,
      averageGasUsed,
      methodFrequency,
      successRate,
    });
  };

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes} minutes ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours} hours ago`;
    const days = Math.floor(hours / 24);
    return `${days} days ago`;
  };

  const formatValue = (value: string): string => {
    // Convert Wei to ETH for better readability
    try {
      const eth = parseFloat(value) / 1e18;
      return `${eth.toFixed(6)} ETH`;
    } catch {
      return value;
    }
  };

  const getStatusBadge = (status: ContractTransaction['status']) => {
    const colors: Record<ContractTransaction['status'], 'green' | 'red'> = {
      success: 'green',
      failed: 'red',
    };
    return (
      <Badge color={colors[status]} size="1">
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  const handleTransactionClick = (hash: string) => {
    router.push(`/explorer/tx/${hash}`);
  };

  const handleSortChange = (field: TransactionFilters['sortBy']) => {
    setFilters(prev => ({
      ...prev,
      sortBy: field,
      sortDirection: prev.sortBy === field && prev.sortDirection === 'desc' ? 'asc' : 'desc',
    }));
  };

  const getSortIcon = (field: TransactionFilters['sortBy']) => {
    if (filters.sortBy !== field) return null;
    return filters.sortDirection === 'desc' ? <ArrowDownIcon /> : <ArrowUpIcon />;
  };

  const getMethodOptions = () => {
    const methods = new Set<string>();
    contract.abi
      .filter(item => item.type === 'function')
      .forEach(func => methods.add(func.name));
    return Array.from(methods);
  };

  const fetchTransactionReceipt = async (hash: string) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/transactions/${hash}/receipt`);
      if (!response.ok) throw new Error('Failed to fetch transaction receipt');
      const data = await response.json();
      setReceipt(data);
    } catch (err) {
      console.error('Error fetching transaction receipt:', err);
      setError('Failed to fetch transaction receipt');
    }
  };

  const handleViewReceipt = async (hash: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedTx(hash);
    await fetchTransactionReceipt(hash);
  };

  const formatHex = (value: string) => {
    return value.startsWith('0x') ? value : `0x${value}`;
  };

  const decodeLogs = (logs: TransactionReceipt['logs']) => {
    return logs.map(log => {
      const eventAbi = contract.abi.find(item => {
        if (item.type !== 'event') return false;
        // Calculate event signature hash
        const signature = `${item.name}(${item.inputs.map(input => input.type).join(',')})`;
        const hash = ethers.id(signature);
        return hash === log.topics[0];
      });

      if (!eventAbi) return log;

      try {
        const decodedData = ethers.AbiCoder.defaultAbiCoder().decode(
          eventAbi.inputs.map(input => input.type),
          log.data
        );

        return {
          ...log,
          decodedName: eventAbi.name,
          decodedData: decodedData.map((value, index) => ({
            name: eventAbi.inputs[index].name,
            type: eventAbi.inputs[index].type,
            value: value.toString(),
          })),
        };
      } catch {
        return log;
      }
    });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Text>Loading transactions...</Text>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Text color="red">{error}</Text>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <Text size="3" weight="bold">
          Contract Transactions
        </Text>
        <Button
          variant="soft"
          onClick={fetchTransactions}
          disabled={loading}
        >
          <ReloadIcon className="mr-1" />
          Refresh
        </Button>
      </div>

      <Card className="mb-6">
        <Flex direction="column" gap="4" p="4">
          <Flex gap="4" align="end">
            <div className="flex-1">
              <Text as="div" size="2" mb="1" weight="bold">
                Search Transactions
              </Text>
              <div className="relative">
                <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400">
                  <MagnifyingGlassIcon height={16} width={16} />
                </span>
                <input
                  className="w-full pl-9 pr-3 py-2 bg-gray-900 border border-gray-800 rounded-md text-white"
                  placeholder="Search by hash, from, to..."
                  value={filters.search || ''}
                  onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
                />
              </div>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Method
              </Text>
              <Select.Root 
                value={filters.method || 'all'}
                onValueChange={(value) => setFilters(prev => ({ 
                  ...prev, 
                  method: value === 'all' ? undefined : value 
                }))}
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="all">All Methods</Select.Item>
                  {getMethodOptions().map(method => (
                    <Select.Item key={method} value={method}>{method}</Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Status
              </Text>
              <Select.Root 
                value={filters.status}
                onValueChange={(value: string) => 
                  setFilters(prev => ({ 
                    ...prev, 
                    status: value as TransactionFilters['status']
                  }))
                }
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="all">All Status</Select.Item>
                  <Select.Item value="success">Success</Select.Item>
                  <Select.Item value="failed">Failed</Select.Item>
                </Select.Content>
              </Select.Root>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Time Range
              </Text>
              <Select.Root 
                value={filters.timeRange}
                onValueChange={(value: TransactionFilters['timeRange']) => 
                  setFilters(prev => ({ ...prev, timeRange: value }))
                }
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="1h">Last Hour</Select.Item>
                  <Select.Item value="24h">Last 24 Hours</Select.Item>
                  <Select.Item value="7d">Last 7 Days</Select.Item>
                  <Select.Item value="30d">Last 30 Days</Select.Item>
                  <Select.Item value="all">All Time</Select.Item>
                </Select.Content>
              </Select.Root>
            </div>
          </Flex>
        </Flex>
      </Card>

      <Grid columns="4" gap="4" mb="6">
        <Card>
          <Box p="4">
            <Text size="1" color="gray" mb="1">Total Value</Text>
            <Text size="4" weight="bold">{formatValue(stats.totalValue)}</Text>
          </Box>
        </Card>
        <Card>
          <Box p="4">
            <Text size="1" color="gray" mb="1">Average Gas Used</Text>
            <Text size="4" weight="bold">{stats.averageGasUsed}</Text>
          </Box>
        </Card>
        <Card>
          <Box p="4">
            <Text size="1" color="gray" mb="1">Success Rate</Text>
            <Text size="4" weight="bold">{stats.successRate.toFixed(1)}%</Text>
          </Box>
        </Card>
        <Card>
          <Box p="4">
            <Text size="1" color="gray" mb="1">Most Used Method</Text>
            <Text size="4" weight="bold">
              {Object.entries(stats.methodFrequency)
                .sort(([, a], [, b]) => b - a)
                .map(([method]) => method)[0] || 'N/A'}
            </Text>
          </Box>
        </Card>
      </Grid>

      {transactions.length === 0 ? (
        <div className="flex justify-center items-center min-h-[200px]">
          <Text color="gray">No transactions found matching your filters</Text>
        </div>
      ) : (
        <>
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.ColumnHeaderCell>Hash</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Method</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>From</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>To</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell 
                  className="cursor-pointer"
                  onClick={() => handleSortChange('value')}
                >
                  <Flex align="center" gap="1">
                    Value
                    {getSortIcon('value')}
                  </Flex>
                </Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell 
                  className="cursor-pointer"
                  onClick={() => handleSortChange('timestamp')}
                >
                  <Flex align="center" gap="1">
                    Age
                    {getSortIcon('timestamp')}
                  </Flex>
                </Table.ColumnHeaderCell>
                <Table.ColumnHeaderCell 
                  className="cursor-pointer"
                  onClick={() => handleSortChange('gasUsed')}
                >
                  <Flex align="center" gap="1">
                    Gas Used
                    {getSortIcon('gasUsed')}
                  </Flex>
                </Table.ColumnHeaderCell>
              </Table.Row>
            </Table.Header>

            <Table.Body>
              {transactions.map((tx) => (
                <Table.Row 
                  key={tx.hash}
                  className="cursor-pointer hover:bg-gray-800"
                  onClick={() => handleTransactionClick(tx.hash)}
                >
                  <Table.Cell>
                    <Text className="font-mono" color="blue">
                      {tx.hash}
                      {tx.status === 'success' && (
                        <Button
                          size="1"
                          variant="ghost"
                          ml="2"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewReceipt(tx.hash, e);
                          }}
                        >
                          View Receipt
                        </Button>
                      )}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>
                    {tx.method ? (
                      <Badge color="purple" size="1">
                        {tx.method}
                      </Badge>
                    ) : (
                      <Text color="gray">Transfer</Text>
                    )}
                  </Table.Cell>
                  <Table.Cell>
                    <Text className="font-mono" color="gray">
                      {tx.from}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>
                    <Text className="font-mono" color="gray">
                      {tx.to}
                    </Text>
                  </Table.Cell>
                  <Table.Cell>{formatValue(tx.value)}</Table.Cell>
                  <Table.Cell>{getStatusBadge(tx.status)}</Table.Cell>
                  <Table.Cell>{formatTimeAgo(tx.timestamp)}</Table.Cell>
                  <Table.Cell>{tx.gasUsed}</Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Root>

          <Flex justify="between" align="center" mt="4">
            <Flex align="center" gap="4">
              <Text size="2" color="gray">
                Showing {(page - 1) * pageSize + 1} to {Math.min(page * pageSize, contract.transactions)} of {contract.transactions} transactions
              </Text>
              <Select.Root 
                value={pageSize.toString()} 
                onValueChange={(value) => {
                  setPageSize(parseInt(value));
                  setPage(1);
                }}
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="10">10 per page</Select.Item>
                  <Select.Item value="25">25 per page</Select.Item>
                  <Select.Item value="50">50 per page</Select.Item>
                  <Select.Item value="100">100 per page</Select.Item>
                </Select.Content>
              </Select.Root>
            </Flex>

            <Flex gap="2">
              <Button
                variant="soft"
                disabled={page === 1}
                onClick={() => setPage(p => Math.max(1, p - 1))}
              >
                Previous
              </Button>
              <Button
                variant="soft"
                disabled={page === totalPages}
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              >
                Next
              </Button>
            </Flex>
          </Flex>
        </>
      )}

      <Dialog.Root open={selectedTx !== null} onOpenChange={() => setSelectedTx(null)}>
        <Dialog.Content style={{ maxWidth: 600 }}>
          <Dialog.Title>Transaction Receipt</Dialog.Title>
          
          {receipt ? (
            <div className="space-y-4">
              <div>
                <Text weight="bold">Transaction Hash</Text>
                <Text className="font-mono">{receipt.transactionHash}</Text>
              </div>

              <div>
                <Text weight="bold">Status</Text>
                <Badge color={receipt.status ? 'green' : 'red'}>
                  {receipt.status ? 'Success' : 'Failed'}
                </Badge>
              </div>

              <div>
                <Text weight="bold">Block</Text>
                <Flex gap="2">
                  <Text>#{receipt.blockNumber}</Text>
                  <Text className="font-mono" color="gray">{receipt.blockHash}</Text>
                </Flex>
              </div>

              {receipt.contractAddress && (
                <div>
                  <Text weight="bold">Created Contract</Text>
                  <Text className="font-mono">{receipt.contractAddress}</Text>
                </div>
              )}

              <div>
                <Text weight="bold">Gas Used</Text>
                <Text>{receipt.gasUsed} ({receipt.cumulativeGasUsed} cumulative)</Text>
              </div>

              <div>
                <Text weight="bold">Effective Gas Price</Text>
                <Text>{formatValue(receipt.effectiveGasPrice)}</Text>
              </div>

              {receipt.logs.length > 0 && (
                <div>
                  <Text weight="bold" mb="2">Event Logs</Text>
                  <div className="space-y-2">
                    {decodeLogs(receipt.logs).map((log, index) => (
                      <Card key={index}>
                        <Box p="3">
                          <Flex justify="between" mb="2">
                            <Text weight="bold">
                              {'decodedName' in log ? log.decodedName : 'Unknown Event'}
                            </Text>
                            <Text size="1" color="gray">Log #{log.logIndex}</Text>
                          </Flex>
                          
                          {'decodedData' in log ? (
                            <div className="space-y-1">
                              {log.decodedData.map((param, i) => (
                                <Flex key={i} justify="between">
                                  <Text size="2">{param.name} ({param.type})</Text>
                                  <Text size="2" className="font-mono">{param.value}</Text>
                                </Flex>
                              ))}
                            </div>
                          ) : (
                            <div className="space-y-1">
                              <Text size="2">Topics:</Text>
                              {log.topics.map((topic, i) => (
                                <Text key={i} size="2" className="font-mono">{topic}</Text>
                              ))}
                              <Text size="2">Data:</Text>
                              <Text size="2" className="font-mono">{log.data}</Text>
                            </div>
                          )}
                        </Box>
                      </Card>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <Flex justify="center" py="4">
              <Text>Loading receipt...</Text>
            </Flex>
          )}

          <Flex gap="3" mt="4" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Close
              </Button>
            </Dialog.Close>
          </Flex>
        </Dialog.Content>
      </Dialog.Root>
    </div>
  );
} 