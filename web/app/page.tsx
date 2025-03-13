'use client';

import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Grid } from '@radix-ui/themes';
import { useBlockchain } from '@/contexts/BlockchainContext';

export default function Home() {
  const { info, loading, error } = useBlockchain();
  
  console.log('Overview page rendering with info:', info);

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-3xl font-bold mb-8">Blockchain Overview</h1>
        
        {error ? (
          <Card>
            <Text color="red" size="2">
              {error}
            </Text>
          </Card>
        ) : loading ? (
          <Card>
            <Text>Loading blockchain information...</Text>
          </Card>
        ) : (
          <Grid columns="4" gap="4">
            <Card>
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">Latest Block</Text>
                <Text size="6" weight="bold">#{info?.current_block}</Text>
              </Flex>
            </Card>
            
            <Card>
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">Active Validators</Text>
                <Text size="6" weight="bold" color="green">
                  {info?.active_validators}/{info?.validator_count}
                </Text>
              </Flex>
            </Card>
            
            <Card>
              <Flex direction="column" gap="2">
                <Text size="2" weight="bold">Pending Transactions</Text>
                <Text size="6" weight="bold">{info?.pending_transactions}</Text>
              </Flex>
            </Card>
          </Grid>
        )}

        <div className="mt-8">
          <Card className="w-full">
            <Text size="2" weight="bold" className="mb-4">Recent Transactions</Text>
            <div className="text-center text-gray-500 py-8">
              Loading transactions...
            </div>
          </Card>
        </div>
      </div>
    </main>
  );
}
