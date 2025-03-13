'use client';

import { Card, Grid, Text, Flex } from '@radix-ui/themes';

interface ContractStats {
  totalContracts: number;
  activeContracts: number;
  totalTransactions: number;
  totalValue: string;
  averageGasUsed: string;
}

const mockStats: ContractStats = {
  totalContracts: 23,
  activeContracts: 20,
  totalTransactions: 1234,
  totalValue: '50,000 CFX',
  averageGasUsed: '150,000',
};

export default function ContractStats() {
  return (
    <Grid columns="5" gap="4">
      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Total Contracts</Text>
          <Text size="6" weight="bold">{mockStats.totalContracts}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Active Contracts</Text>
          <Text size="6" weight="bold" color="green">
            {mockStats.activeContracts}
          </Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Total Transactions</Text>
          <Text size="6" weight="bold">{mockStats.totalTransactions}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Total Value Locked</Text>
          <Text size="6" weight="bold">{mockStats.totalValue}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Average Gas Used</Text>
          <Text size="6" weight="bold">{mockStats.averageGasUsed}</Text>
        </Flex>
      </Card>
    </Grid>
  );
} 