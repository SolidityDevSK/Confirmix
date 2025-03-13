'use client';

import { useParams } from 'next/navigation';
import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Grid, Table, Badge } from '@radix-ui/themes';
import Link from 'next/link';
import { useState } from 'react';

interface ValidatorDetail {
  address: string;
  name: string;
  status: 'active' | 'inactive' | 'pending';
  blocksProduced: number;
  lastActive: number;
  uptime: string;
  stake: string;
  rewards: string;
  delegators: number;
  performance: {
    last24h: string;
    last7d: string;
    last30d: string;
  };
  recentBlocks: {
    number: number;
    timestamp: number;
    transactions: number;
    reward: string;
  }[];
}

// Mock data - will be replaced with real API call
const mockValidatorDetail: ValidatorDetail = {
  address: '0x1234...5678',
  name: 'Validator 1',
  status: 'active',
  blocksProduced: 1234,
  lastActive: Date.now() - 30000,
  uptime: '99.98%',
  stake: '100,000 CFX',
  rewards: '1,234 CFX',
  delegators: 5,
  performance: {
    last24h: '100%',
    last7d: '99.95%',
    last30d: '99.98%',
  },
  recentBlocks: [
    {
      number: 1234567,
      timestamp: Date.now() - 60000,
      transactions: 5,
      reward: '0.5 CFX',
    },
    {
      number: 1234566,
      timestamp: Date.now() - 120000,
      transactions: 3,
      reward: '0.5 CFX',
    },
  ],
};

export default function ValidatorDetailPage() {
  const params = useParams();
  const [validator] = useState<ValidatorDetail>(mockValidatorDetail);

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes} minutes ago`;
    const hours = Math.floor(minutes / 60);
    return `${hours} hours ago`;
  };

  const getStatusBadge = (status: ValidatorDetail['status']) => {
    const colors = {
      active: 'green',
      inactive: 'red',
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
              href="/validators" 
              className="text-blue-500 hover:text-blue-400 mb-4 inline-block"
            >
              ← Back to Validators
            </Link>
            <div className="flex items-center gap-4 mb-4">
              <h1 className="text-3xl font-bold">{validator.name}</h1>
              {getStatusBadge(validator.status)}
            </div>
            <Text className="text-gray-400 font-mono">
              {params.address}
            </Text>
          </div>

          <Grid columns="2" gap="4">
            <Card>
              <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
                Overview
              </Text>
              <div className="p-4">
                <Grid columns="2" gap="4">
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Total Stake</Text>
                    <Text size="5">{validator.stake}</Text>
                  </Flex>
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Total Rewards</Text>
                    <Text size="5">{validator.rewards}</Text>
                  </Flex>
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Blocks Produced</Text>
                    <Text size="5">{validator.blocksProduced}</Text>
                  </Flex>
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Delegators</Text>
                    <Text size="5">{validator.delegators}</Text>
                  </Flex>
                </Grid>
              </div>
            </Card>

            <Card>
              <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
                Performance
              </Text>
              <div className="p-4">
                <Grid columns="3" gap="4">
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Last 24h</Text>
                    <Text size="5" color="green">{validator.performance.last24h}</Text>
                  </Flex>
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Last 7d</Text>
                    <Text size="5" color="green">{validator.performance.last7d}</Text>
                  </Flex>
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Last 30d</Text>
                    <Text size="5" color="green">{validator.performance.last30d}</Text>
                  </Flex>
                </Grid>
              </div>
            </Card>
          </Grid>

          <Card>
            <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
              Recent Blocks
            </Text>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeaderCell>Block</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>Age</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>Transactions</Table.ColumnHeaderCell>
                  <Table.ColumnHeaderCell>Reward</Table.ColumnHeaderCell>
                </Table.Row>
              </Table.Header>

              <Table.Body>
                {validator.recentBlocks.map((block) => (
                  <Table.Row key={block.number}>
                    <Table.Cell>
                      <Link 
                        href={`/explorer/block/${block.number}`}
                        className="text-blue-500 hover:text-blue-400"
                      >
                        #{block.number}
                      </Link>
                    </Table.Cell>
                    <Table.Cell>{formatTimeAgo(block.timestamp)}</Table.Cell>
                    <Table.Cell>{block.transactions}</Table.Cell>
                    <Table.Cell>{block.reward}</Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table.Root>
          </Card>
        </Flex>
      </div>
    </main>
  );
} 