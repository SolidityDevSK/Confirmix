'use client';

import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Grid, Tabs, Box, Button, Badge } from '@radix-ui/themes';
import Link from 'next/link';
import { useState, useEffect } from 'react';
import ContractInteraction from '@/components/contracts/ContractInteraction';
import ContractEvents from '@/components/contracts/ContractEvents';
import ContractCode from '@/components/contracts/ContractCode';
import ContractTransactions from '@/components/contracts/ContractTransactions';
import ContractAnalytics from '@/components/contracts/ContractAnalytics';
import { blockchainService } from '@/services/blockchain';
import { websocketService } from '@/lib/websocket';
import type { BlockchainContract, ContractEvent } from '@/types/blockchain';
import TransferOwnershipDialog from '@/components/contracts/TransferOwnershipDialog';

interface ContractActivity {
  last24hTransactions: number;
  last24hVolume: string;
}

interface ContractDetailPageProps {
  params: {
    address: string;
  };
}

export default function ContractDetailPage({ params }: ContractDetailPageProps) {
  const [contract, setContract] = useState<BlockchainContract | null>(null);
  const [activity, setActivity] = useState<ContractActivity>({
    last24hTransactions: 0,
    last24hVolume: '0',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [transferDialogOpen, setTransferDialogOpen] = useState(false);

  async function fetchContract() {
    try {
      setLoading(true);
      setError(null);
      const data = await blockchainService.getContract(params.address);
      setContract(data);
    } catch (err) {
      console.error('Error fetching contract:', err);
      setError('Failed to load contract details');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchContract();
  }, [params.address]);

  useEffect(() => {
    if (contract) {
      fetchActivity();
      
      // Subscribe to real-time events for activity updates
      const handleNewEvent = (event: ContractEvent) => {
        setActivity(prev => ({
          last24hTransactions: prev.last24hTransactions + 1,
          last24hVolume: updateVolume(prev.last24hVolume, event),
        }));
      };

      websocketService.subscribeToContractEvents(
        contract.address,
        contract.abi,
        handleNewEvent
      );

      return () => {
        websocketService.unsubscribeFromContractEvents(
          contract.address,
          handleNewEvent
        );
      };
    }
  }, [contract]);

  const fetchActivity = async () => {
    try {
      // Get last 24h transactions
      const oneDayAgo = Math.floor(Date.now() / 1000) - 24 * 60 * 60;
      const events = await blockchainService.getContractEvents(
        params.address,
        oneDayAgo
      );

      let volume = '0';
      events.forEach(event => {
        volume = updateVolume(volume, event);
      });

      setActivity({
        last24hTransactions: events.length,
        last24hVolume: volume,
      });
    } catch (err) {
      console.error('Error fetching activity:', err);
    }
  };

  const updateVolume = (currentVolume: string, event: ContractEvent): string => {
    // Sum up transfer amounts for common token events
    if (event.name === 'Transfer' && event.args.value) {
      const currentValue = BigInt(currentVolume);
      const newValue = BigInt(event.args.value);
      return (currentValue + newValue).toString();
    }
    return currentVolume;
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

  if (loading) {
    return (
      <main>
        <Navigation />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex justify-center items-center min-h-[400px]">
            <Text>Loading contract details...</Text>
          </div>
        </div>
      </main>
    );
  }

  if (error || !contract) {
    return (
      <main>
        <Navigation />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex justify-center items-center min-h-[400px]">
            <Text color="red">{error || 'Contract not found'}</Text>
          </div>
        </div>
      </main>
    );
  }

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div>
            <Link 
              href="/contracts" 
              className="text-blue-500 hover:text-blue-400 mb-4 inline-block"
            >
              ← Back to Contracts
            </Link>
            <div className="flex items-center justify-between mb-4">
              <div>
                <h1 className="text-3xl font-bold">{contract.name}</h1>
                <Text className="text-gray-400 font-mono">
                  {params.address}
                </Text>
              </div>
              <div className="flex items-center gap-4">
                {contract.sourceCode ? (
                  <Badge color="green" size="2">Verified</Badge>
                ) : (
                  <Link href={`/contracts/${params.address}/verify`}>
                    <Button color="blue" variant="soft">
                      Verify Contract
                    </Button>
                  </Link>
                )}
                <Button 
                  color="red" 
                  variant="soft"
                  onClick={() => setTransferDialogOpen(true)}
                >
                  Transfer Ownership
                </Button>
              </div>
            </div>
          </div>

          <Grid columns="2" gap="4">
            <Card>
              <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
                Contract Overview
              </Text>
              <Box p="4">
                <Grid columns="2" gap="4">
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Type</Text>
                    <Text className="capitalize">{contract.type}</Text>
                  </Flex>
                  
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Status</Text>
                    <Text>{contract.status}</Text>
                  </Flex>

                  <Flex direction="column" gap="2">
                    <Text weight="bold">Owner</Text>
                    <Text className="font-mono" color="gray">{contract.owner}</Text>
                  </Flex>

                  <Flex direction="column" gap="2">
                    <Text weight="bold">Balance</Text>
                    <Text>{contract.balance}</Text>
                  </Flex>

                  <Flex direction="column" gap="2">
                    <Text weight="bold">Transactions</Text>
                    <Text>{contract.transactions}</Text>
                  </Flex>

                  <Flex direction="column" gap="2">
                    <Text weight="bold">Deployed</Text>
                    <Text>{formatTimeAgo(contract.deployedAt)}</Text>
                  </Flex>
                </Grid>
              </Box>
            </Card>

            <Card>
              <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
                Contract Activity
              </Text>
              <Box p="4">
                <Grid columns="2" gap="4">
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Last 24h Transactions</Text>
                    <Text size="6">{activity.last24hTransactions}</Text>
                  </Flex>
                  
                  <Flex direction="column" gap="2">
                    <Text weight="bold">Last 24h Volume</Text>
                    <Text size="6">{activity.last24hVolume} CFX</Text>
                  </Flex>
                </Grid>
              </Box>
            </Card>
          </Grid>

          <Card>
            <Tabs.Root defaultValue="interact">
              <Tabs.List>
                <Tabs.Trigger value="interact">Interact</Tabs.Trigger>
                <Tabs.Trigger value="transactions">Transactions</Tabs.Trigger>
                <Tabs.Trigger value="events">Events</Tabs.Trigger>
                <Tabs.Trigger value="analytics">Analytics</Tabs.Trigger>
                <Tabs.Trigger value="code">Code</Tabs.Trigger>
              </Tabs.List>

              <Box p="4">
                <Tabs.Content value="interact">
                  <ContractInteraction contract={contract} />
                </Tabs.Content>

                <Tabs.Content value="transactions">
                  <ContractTransactions contract={contract} />
                </Tabs.Content>

                <Tabs.Content value="events">
                  <ContractEvents contract={contract} />
                </Tabs.Content>

                <Tabs.Content value="analytics">
                  <ContractAnalytics contract={contract} />
                </Tabs.Content>

                <Tabs.Content value="code">
                  <ContractCode contract={contract} />
                </Tabs.Content>
              </Box>
            </Tabs.Root>
          </Card>
        </Flex>
      </div>
      <TransferOwnershipDialog
        contract={contract}
        open={transferDialogOpen}
        onOpenChange={setTransferDialogOpen}
        onSuccess={() => {
          // Refresh contract details after transfer
          fetchContract();
        }}
      />
    </main>
  );
} 