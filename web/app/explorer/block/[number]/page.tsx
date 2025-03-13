'use client';

import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Grid, Table } from '@radix-ui/themes';
import { useParams } from 'next/navigation';
import Link from 'next/link';

interface BlockDetail {
  number: number;
  hash: string;
  timestamp: number;
  transactions: number;
  validator: string;
  parentHash: string;
  size: string;
  gasUsed: string;
  gasLimit: string;
}

// Mock data - will be replaced with real API call
const mockBlockDetail: BlockDetail = {
  number: 1234567,
  hash: '0x1234...5678',
  timestamp: Date.now() - 60000,
  transactions: 5,
  validator: '0xabc...def',
  parentHash: '0x9876...5432',
  size: '12.5 KB',
  gasUsed: '21,000',
  gasLimit: '30,000'
};

export default function BlockDetailPage() {
  const params = useParams();
  const blockNumber = params.number;

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    return `${minutes} minutes ago`;
  };

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div>
            <Link 
              href="/explorer" 
              className="text-blue-500 hover:text-blue-400 mb-4 inline-block"
            >
              ← Back to Explorer
            </Link>
            <h1 className="text-3xl font-bold mb-4">Block #{blockNumber}</h1>
            <Text className="text-gray-400">
              Block Details and Transactions
            </Text>
          </div>

          <Card>
            <Grid columns="2" gap="4" className="p-4">
              <Flex direction="column" gap="2">
                <Text weight="bold">Block Number</Text>
                <Text>{mockBlockDetail.number}</Text>
              </Flex>
              
              <Flex direction="column" gap="2">
                <Text weight="bold">Timestamp</Text>
                <Text>{formatTimeAgo(mockBlockDetail.timestamp)}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Transactions</Text>
                <Text>{mockBlockDetail.transactions}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Validator</Text>
                <Text color="gray">{mockBlockDetail.validator}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Hash</Text>
                <Text className="font-mono" color="gray">{mockBlockDetail.hash}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Parent Hash</Text>
                <Text className="font-mono" color="gray">{mockBlockDetail.parentHash}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Size</Text>
                <Text>{mockBlockDetail.size}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Gas Used</Text>
                <Text>{mockBlockDetail.gasUsed}</Text>
              </Flex>
            </Grid>
          </Card>

          <Card>
            <Text size="3" weight="bold" className="p-4">Block Transactions</Text>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeaderCell>Hash</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>From</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>To</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>Value</Table.ColumnHeaderCell>
                </Table.Row>
              </Table.Header>

              <Table.Body>
                <Table.Row className="cursor-pointer hover:bg-gray-800">
                  <Table.Cell>
                    <Link href="/explorer/tx/0xabcd" className="text-blue-500 hover:text-blue-400">
                      0xabcd...1234
                    </Link>
                  </Table.Cell>
                  <Table.Cell>
                    <Text color="gray">0x1111...2222</Text>
                  </Table.Cell>
                  <Table.Cell>
                    <Text color="gray">0x3333...4444</Text>
                  </Table.Cell>
                  <Table.Cell>0.5 ETH</Table.Cell>
                </Table.Row>
              </Table.Body>
            </Table.Root>
          </Card>
        </Flex>
      </div>
    </main>
  );
} 