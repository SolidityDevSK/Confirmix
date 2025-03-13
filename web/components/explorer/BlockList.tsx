'use client';

import { Table, Text } from '@radix-ui/themes';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useBlockchain } from '@/contexts/BlockchainContext';
import { blockchainService } from '@/services/blockchain';

interface Block {
  height: number;
  hash: string;
  timestamp: number;
  transactions: number;
  validator: string;
}

export default function BlockList() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const { info } = useBlockchain();

  useEffect(() => {
    const fetchBlocks = async () => {
      try {
        setLoading(true);
        const data = await blockchainService.getBlocks();
        console.log('Fetched blocks:', data);
        setBlocks(data);
      } catch (error) {
        console.error('Error fetching blocks:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchBlocks();
  }, [info?.current_block]); // Refresh when new blocks are added

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    return `${minutes} minutes ago`;
  };

  const handleBlockClick = (blockNumber: number) => {
    router.push(`/explorer/block/${blockNumber}`);
  };

  if (loading) {
    return (
      <div className="text-center py-8">
        <Text>Loading blocks...</Text>
      </div>
    );
  }

  return (
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeaderCell>Block</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Hash</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Age</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Txs</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Validator</Table.ColumnHeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {blocks.map((block) => (
          <Table.Row 
            key={block.height} 
            className="cursor-pointer hover:bg-gray-800"
            onClick={() => handleBlockClick(block.height)}
          >
            <Table.Cell>
              <Text weight="bold" color="blue">
                #{block.height}
              </Text>
            </Table.Cell>
            <Table.Cell>
              <Text color="gray">{block.hash}</Text>
            </Table.Cell>
            <Table.Cell>{formatTimeAgo(block.timestamp)}</Table.Cell>
            <Table.Cell>{block.transactions}</Table.Cell>
            <Table.Cell>
              <Text color="gray">{block.validator}</Text>
            </Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
} 