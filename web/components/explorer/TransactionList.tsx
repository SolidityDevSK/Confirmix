'use client';

import { Table, Text, Badge } from '@radix-ui/themes';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useBlockchain } from '@/contexts/BlockchainContext';
import { blockchainService } from '@/services/blockchain';

interface Transaction {
  hash: string;
  from: string;
  to: string;
  value: string;
  type: 'transfer' | 'contract' | 'validator';
  timestamp: number;
  status: 'success' | 'failed' | 'pending';
}

export default function TransactionList() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const { info } = useBlockchain();

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        setLoading(true);
        const data = await blockchainService.getTransactions();
        console.log('Fetched transactions:', data);
        setTransactions(data);
      } catch (error) {
        console.error('Error fetching transactions:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchTransactions();
  }, [info?.current_block]); // Refresh when new blocks are added

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    return `${minutes} minutes ago`;
  };

  const getTransactionBadge = (type: Transaction['type']) => {
    const colors = {
      transfer: 'blue',
      contract: 'purple',
      validator: 'green',
    };
    return (
      <Badge color={colors[type] as any} size="1">
        {type.charAt(0).toUpperCase() + type.slice(1)}
      </Badge>
    );
  };

  const getStatusBadge = (status: Transaction['status']) => {
    const colors = {
      success: 'green',
      failed: 'red',
      pending: 'yellow',
    };
    return (
      <Badge color={colors[status] as any} size="1">
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    );
  };

  const handleTransactionClick = (hash: string) => {
    router.push(`/explorer/tx/${hash}`);
  };

  if (loading) {
    return (
      <div className="text-center py-8">
        <Text>Loading transactions...</Text>
      </div>
    );
  }

  return (
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeaderCell>Hash</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Type</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>From</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>To</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Value</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Age</Table.ColumnHeaderCell>
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
              <Text color="blue">{tx.hash}</Text>
            </Table.Cell>
            <Table.Cell>{getTransactionBadge(tx.type)}</Table.Cell>
            <Table.Cell>{getStatusBadge(tx.status)}</Table.Cell>
            <Table.Cell>
              <Text color="gray">{tx.from}</Text>
            </Table.Cell>
            <Table.Cell>
              <Text color="gray">{tx.to}</Text>
            </Table.Cell>
            <Table.Cell>{tx.value}</Table.Cell>
            <Table.Cell>{formatTimeAgo(tx.timestamp)}</Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
} 