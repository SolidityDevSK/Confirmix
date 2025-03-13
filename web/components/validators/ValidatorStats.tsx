'use client';

import { Card, Grid, Text, Flex } from '@radix-ui/themes';

interface NetworkStats {
  totalValidators: number;
  activeValidators: number;
  averageBlockTime: string;
  networkUptime: string;
  lastConsensus: string;
}

const mockStats: NetworkStats = {
  totalValidators: 5,
  activeValidators: 4,
  averageBlockTime: '5 seconds',
  networkUptime: '99.99%',
  lastConsensus: '2 minutes ago',
};

export default function ValidatorStats() {
  return (
    <Grid columns="5" gap="4">
      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Total Validators</Text>
          <Text size="6" weight="bold">{mockStats.totalValidators}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Active Validators</Text>
          <Text size="6" weight="bold" color="green">
            {mockStats.activeValidators}
          </Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Average Block Time</Text>
          <Text size="6" weight="bold">{mockStats.averageBlockTime}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Network Uptime</Text>
          <Text size="6" weight="bold">{mockStats.networkUptime}</Text>
        </Flex>
      </Card>

      <Card>
        <Flex direction="column" gap="1" className="p-4">
          <Text size="2" color="gray">Last Consensus</Text>
          <Text size="6" weight="bold">{mockStats.lastConsensus}</Text>
        </Flex>
      </Card>
    </Grid>
  );
} 