'use client';

import { Card, Text, Grid, Flex, Box, Button } from '@radix-ui/themes';
import { useState, useEffect } from 'react';
import { blockchainService } from '@/services/blockchain';
import type { BlockchainContract } from '@/types/blockchain';
import { ReloadIcon } from '@radix-ui/react-icons';

interface ContractAnalytics {
  dailyTransactions: {
    date: string;
    count: number;
    volume: string;
  }[];
  topInteractions: {
    address: string;
    transactions: number;
    volume: string;
  }[];
  methodStats: {
    name: string;
    calls: number;
    gasUsed: string;
  }[];
  gasStats: {
    totalGasUsed: string;
    averageGasUsed: string;
    maxGasUsed: string;
    minGasUsed: string;
  };
}

interface ContractAnalyticsProps {
  contract: BlockchainContract;
}

export default function ContractAnalytics({ contract }: ContractAnalyticsProps) {
  const [analytics, setAnalytics] = useState<ContractAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<'24h' | '7d' | '30d'>('7d');

  useEffect(() => {
    fetchAnalytics();
  }, [contract.address, timeRange]);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await blockchainService.getContractAnalytics(
        contract.address,
        timeRange
      );
      setAnalytics(data);
    } catch (err) {
      console.error('Error fetching analytics:', err);
      setError('Failed to fetch contract analytics');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Text>Loading analytics...</Text>
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

  if (!analytics) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Text color="gray">No analytics data available</Text>
      </div>
    );
  }

  return (
    <div>
      <Flex justify="between" align="center" mb="6">
        <Text size="3" weight="bold">
          Contract Analytics
        </Text>
        <Flex gap="4" align="center">
          <Flex gap="2">
            <Button 
              variant="soft" 
              color={timeRange === '24h' ? 'blue' : 'gray'}
              onClick={() => setTimeRange('24h')}
            >
              24h
            </Button>
            <Button 
              variant="soft" 
              color={timeRange === '7d' ? 'blue' : 'gray'}
              onClick={() => setTimeRange('7d')}
            >
              7d
            </Button>
            <Button 
              variant="soft" 
              color={timeRange === '30d' ? 'blue' : 'gray'}
              onClick={() => setTimeRange('30d')}
            >
              30d
            </Button>
          </Flex>
          <Button
            variant="soft"
            onClick={fetchAnalytics}
            disabled={loading}
          >
            <ReloadIcon className="mr-1" />
            Refresh
          </Button>
        </Flex>
      </Flex>

      <Grid columns="2" gap="4">
        <Card>
          <Text size="2" weight="bold" className="p-4 border-b border-gray-800">
            Gas Usage Statistics
          </Text>
          <Box p="4">
            <Grid columns="2" gap="4">
              <Flex direction="column" gap="1">
                <Text size="1" color="gray">Total Gas Used</Text>
                <Text size="4" weight="bold">{analytics.gasStats.totalGasUsed}</Text>
              </Flex>
              <Flex direction="column" gap="1">
                <Text size="1" color="gray">Average Gas per Tx</Text>
                <Text size="4" weight="bold">{analytics.gasStats.averageGasUsed}</Text>
              </Flex>
              <Flex direction="column" gap="1">
                <Text size="1" color="gray">Max Gas Used</Text>
                <Text size="4" weight="bold">{analytics.gasStats.maxGasUsed}</Text>
              </Flex>
              <Flex direction="column" gap="1">
                <Text size="1" color="gray">Min Gas Used</Text>
                <Text size="4" weight="bold">{analytics.gasStats.minGasUsed}</Text>
              </Flex>
            </Grid>
          </Box>
        </Card>

        <Card>
          <Text size="2" weight="bold" className="p-4 border-b border-gray-800">
            Method Usage Statistics
          </Text>
          <Box p="4">
            <div className="space-y-4">
              {analytics.methodStats.map((method) => (
                <div key={method.name}>
                  <Flex justify="between" mb="1">
                    <Text weight="bold">{method.name}</Text>
                    <Text>{method.calls} calls</Text>
                  </Flex>
                  <Flex justify="between">
                    <Text size="1" color="gray">Gas Used</Text>
                    <Text size="1">{method.gasUsed}</Text>
                  </Flex>
                </div>
              ))}
            </div>
          </Box>
        </Card>

        <Card>
          <Text size="2" weight="bold" className="p-4 border-b border-gray-800">
            Daily Transaction Volume
          </Text>
          <Box p="4">
            <div className="space-y-2">
              {analytics.dailyTransactions.map((day) => (
                <div key={day.date}>
                  <Flex justify="between" mb="1">
                    <Text>{day.date}</Text>
                    <Text>{day.count} transactions</Text>
                  </Flex>
                  <Flex justify="between">
                    <Text size="1" color="gray">Volume</Text>
                    <Text size="1">{day.volume} CFX</Text>
                  </Flex>
                </div>
              ))}
            </div>
          </Box>
        </Card>

        <Card>
          <Text size="2" weight="bold" className="p-4 border-b border-gray-800">
            Top Interactions
          </Text>
          <Box p="4">
            <div className="space-y-4">
              {analytics.topInteractions.map((interaction) => (
                <div key={interaction.address}>
                  <Flex justify="between" mb="1">
                    <Text className="font-mono" color="gray">
                      {interaction.address}
                    </Text>
                    <Text>{interaction.transactions} txs</Text>
                  </Flex>
                  <Flex justify="between">
                    <Text size="1" color="gray">Volume</Text>
                    <Text size="1">{interaction.volume} CFX</Text>
                  </Flex>
                </div>
              ))}
            </div>
          </Box>
        </Card>
      </Grid>
    </div>
  );
} 