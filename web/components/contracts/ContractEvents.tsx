'use client';

import { Table, Text, Card, Badge, Button, Flex, Select } from '@radix-ui/themes';
import type { BlockchainContract, ContractEvent } from '@/types/blockchain';
import { blockchainService } from '@/services/blockchain';
import { websocketService } from '@/lib/websocket';
import { useState, useEffect } from 'react';
import { ReloadIcon, MagnifyingGlassIcon } from '@radix-ui/react-icons';

interface ContractEventsProps {
  contract: BlockchainContract;
}

interface EventFilters {
  eventName: string;
  fromBlock?: number;
  toBlock?: number;
  search?: string;
  timeRange: '1h' | '24h' | '7d' | '30d' | 'all';
}

export default function ContractEvents({ contract }: ContractEventsProps) {
  const [events, setEvents] = useState<ContractEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [liveUpdates, setLiveUpdates] = useState(true);
  const [filters, setFilters] = useState<EventFilters>({
    eventName: 'all',
    timeRange: '24h',
  });

  useEffect(() => {
    fetchEvents();

    // Subscribe to real-time events
    if (liveUpdates) {
      const handleNewEvent = (event: ContractEvent) => {
        if (shouldIncludeEvent(event)) {
          setEvents(prev => [event, ...prev]);
        }
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
  }, [contract.address, contract.abi, liveUpdates, filters]);

  const fetchEvents = async () => {
    try {
      setLoading(true);
      setError(null);

      // Calculate fromBlock based on timeRange
      let fromTimestamp: number | undefined;
      const now = Math.floor(Date.now() / 1000);
      switch (filters.timeRange) {
        case '1h':
          fromTimestamp = now - 3600;
          break;
        case '24h':
          fromTimestamp = now - 86400;
          break;
        case '7d':
          fromTimestamp = now - 604800;
          break;
        case '30d':
          fromTimestamp = now - 2592000;
          break;
        default:
          fromTimestamp = undefined;
      }

      const data = await blockchainService.getContractEvents(
        contract.address,
        fromTimestamp || 0
      );

      // Apply filters
      const filteredEvents = data.filter(shouldIncludeEvent);
      setEvents(filteredEvents);
    } catch (err) {
      console.error('Error fetching events:', err);
      setError('Failed to fetch contract events');
    } finally {
      setLoading(false);
    }
  };

  const shouldIncludeEvent = (event: ContractEvent): boolean => {
    // Filter by event name
    if (filters.eventName !== 'all' && event.name !== filters.eventName) {
      return false;
    }

    // Filter by search term
    if (filters.search) {
      const searchTerm = filters.search.toLowerCase();
      const eventString = JSON.stringify(event).toLowerCase();
      if (!eventString.includes(searchTerm)) {
        return false;
      }
    }

    return true;
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

  const formatValue = (value: string) => {
    // Convert Wei to ETH for better readability
    if (value.length > 18) {
      const eth = parseInt(value) / 1e18;
      return `${eth} ETH`;
    }
    return value;
  };

  const getEventNames = () => {
    const names = new Set<string>();
    contract.abi
      .filter(item => item.type === 'event')
      .forEach(event => names.add(event.name));
    return Array.from(names);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Text>Loading events...</Text>
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

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <Text size="3" weight="bold">
          Contract Events
        </Text>
        <div className="flex items-center gap-4">
          <Button
            variant="soft"
            color={liveUpdates ? 'green' : 'gray'}
            onClick={() => setLiveUpdates(!liveUpdates)}
          >
            {liveUpdates ? 'Live Updates: On' : 'Live Updates: Off'}
          </Button>
          <Button
            variant="soft"
            onClick={fetchEvents}
            disabled={loading}
          >
            <ReloadIcon className="mr-1" />
            Refresh
          </Button>
        </div>
      </div>

      <Card className="mb-6">
        <Flex direction="column" gap="4" p="4">
          <Flex gap="4" align="end">
            <div className="flex-1">
              <Text as="div" size="2" mb="1" weight="bold">
                Search Events
              </Text>
              <div className="relative">
                <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400">
                  <MagnifyingGlassIcon height={16} width={16} />
                </span>
                <input
                  className="w-full pl-9 pr-3 py-2 bg-gray-900 border border-gray-800 rounded-md text-white"
                  placeholder="Search by transaction hash, arguments..."
                  value={filters.search || ''}
                  onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
                />
              </div>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Event Type
              </Text>
              <Select.Root 
                value={filters.eventName} 
                onValueChange={(value) => setFilters(prev => ({ ...prev, eventName: value }))}
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="all">All Events</Select.Item>
                  {getEventNames().map(name => (
                    <Select.Item key={name} value={name}>{name}</Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Time Range
              </Text>
              <Select.Root 
                value={filters.timeRange}
                onValueChange={(value: EventFilters['timeRange']) => 
                  setFilters(prev => ({ ...prev, timeRange: value }))
                }
              >
                <Select.Trigger />
                <Select.Content>
                  <Select.Item value="1h">Last Hour</Select.Item>
                  <Select.Item value="24h">Last 24 Hours</Select.Item>
                  <Select.Item value="7d">Last 7 Days</Select.Item>
                  <Select.Item value="30d">Last 30 Days</Select.Item>
                  <Select.Item value="all">All Time</Select.Item>
                </Select.Content>
              </Select.Root>
            </div>
          </Flex>
        </Flex>
      </Card>

      {events.length === 0 ? (
        <div className="flex justify-center items-center min-h-[200px]">
          <Text color="gray">No events found matching your filters</Text>
        </div>
      ) : (
        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeaderCell>Event</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Block</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Age</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Transaction</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Arguments</Table.ColumnHeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {events.map((event) => (
              <Table.Row key={`${event.transactionHash}-${event.name}`}>
                <Table.Cell>
                  <Badge color="blue" size="1">
                    {event.name}
                  </Badge>
                </Table.Cell>
                <Table.Cell>
                  <Text className="font-mono">
                    #{event.blockNumber}
                  </Text>
                </Table.Cell>
                <Table.Cell>
                  {formatTimeAgo(event.timestamp)}
                </Table.Cell>
                <Table.Cell>
                  <Text className="font-mono" color="gray">
                    {event.transactionHash}
                  </Text>
                </Table.Cell>
                <Table.Cell>
                  <Card variant="classic" size="1">
                    {Object.entries(event.args).map(([key, value]) => (
                      <div key={key} className="flex gap-2">
                        <Text size="1" color="gray">{key}:</Text>
                        <Text size="1" className="font-mono">
                          {typeof value === 'string' ? formatValue(value) : String(value)}
                        </Text>
                      </div>
                    ))}
                  </Card>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      )}
    </div>
  );
} 