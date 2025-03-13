'use client';

import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Grid, Badge } from '@radix-ui/themes';
import { useParams } from 'next/navigation';
import Link from 'next/link';

interface TransactionDetail {
  hash: string;
  blockNumber: number;
  from: string;
  to: string;
  value: string;
  type: 'transfer' | 'contract' | 'validator';
  timestamp: number;
  gasPrice: string;
  gasLimit: string;
  gasUsed: string;
  nonce: number;
  status: 'success' | 'failed' | 'pending';
  data?: string;
}

// Mock data - will be replaced with real API call
const mockTransactionDetail: TransactionDetail = {
  hash: '0xabcd...1234',
  blockNumber: 1234567,
  from: '0x1111...2222',
  to: '0x3333...4444',
  value: '0.5 ETH',
  type: 'transfer',
  timestamp: Date.now() - 30000,
  gasPrice: '20 Gwei',
  gasLimit: '21,000',
  gasUsed: '21,000',
  nonce: 42,
  status: 'success',
  data: '0x',
};

export default function TransactionDetailPage() {
  const params = useParams();
  const txHash = params.hash;

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    return `${minutes} minutes ago`;
  };

  const getStatusBadge = (status: TransactionDetail['status']) => {
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
            <h1 className="text-3xl font-bold mb-4">Transaction Details</h1>
            <Text className="text-gray-400">
              Transaction Hash: {txHash}
            </Text>
          </div>

          <Card>
            <Grid columns="2" gap="4" className="p-4">
              <Flex direction="column" gap="2">
                <Text weight="bold">Transaction Hash</Text>
                <Text className="font-mono" color="gray">{mockTransactionDetail.hash}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Status</Text>
                {getStatusBadge(mockTransactionDetail.status)}
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Block</Text>
                <Link 
                  href={`/explorer/block/${mockTransactionDetail.blockNumber}`}
                  className="text-blue-500 hover:text-blue-400"
                >
                  #{mockTransactionDetail.blockNumber}
                </Link>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Timestamp</Text>
                <Text>{formatTimeAgo(mockTransactionDetail.timestamp)}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">From</Text>
                <Text className="font-mono" color="gray">{mockTransactionDetail.from}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">To</Text>
                <Text className="font-mono" color="gray">{mockTransactionDetail.to}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Value</Text>
                <Text>{mockTransactionDetail.value}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Transaction Type</Text>
                <Text className="capitalize">{mockTransactionDetail.type}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Gas Price</Text>
                <Text>{mockTransactionDetail.gasPrice}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Gas Limit</Text>
                <Text>{mockTransactionDetail.gasLimit}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Gas Used</Text>
                <Text>{mockTransactionDetail.gasUsed}</Text>
              </Flex>

              <Flex direction="column" gap="2">
                <Text weight="bold">Nonce</Text>
                <Text>{mockTransactionDetail.nonce}</Text>
              </Flex>
            </Grid>
          </Card>

          {mockTransactionDetail.data && (
            <Card>
              <Text size="3" weight="bold" className="p-4">Input Data</Text>
              <div className="p-4 bg-gray-900 rounded-lg font-mono">
                <Text color="gray">{mockTransactionDetail.data}</Text>
              </div>
            </Card>
          )}
        </Flex>
      </div>
    </main>
  );
} 