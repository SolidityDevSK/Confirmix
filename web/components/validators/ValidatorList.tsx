'use client';

import { Table, Text, Badge, Button, Flex } from '@radix-ui/themes';
import { useState } from 'react';
import { useRouter } from 'next/navigation';

interface Validator {
  address: string;
  name: string;
  status: 'active' | 'inactive' | 'pending';
  blocksProduced: number;
  lastActive: number;
  uptime: string;
}

const mockValidators: Validator[] = [
  {
    address: '0x1234...5678',
    name: 'Validator 1',
    status: 'active',
    blocksProduced: 1234,
    lastActive: Date.now() - 30000,
    uptime: '99.98%',
  },
  {
    address: '0x8765...4321',
    name: 'Validator 2',
    status: 'active',
    blocksProduced: 1156,
    lastActive: Date.now() - 60000,
    uptime: '99.95%',
  },
  {
    address: '0xabcd...efgh',
    name: 'Validator 3',
    status: 'inactive',
    blocksProduced: 985,
    lastActive: Date.now() - 3600000,
    uptime: '98.50%',
  },
];

export default function ValidatorList() {
  const [validators] = useState<Validator[]>(mockValidators);
  const router = useRouter();

  const formatTimeAgo = (timestamp: number) => {
    const seconds = Math.floor((Date.now() - timestamp) / 1000);
    if (seconds < 60) return `${seconds} seconds ago`;
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes} minutes ago`;
    const hours = Math.floor(minutes / 60);
    return `${hours} hours ago`;
  };

  const getStatusBadge = (status: Validator['status']) => {
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

  const handleDetailsClick = (address: string) => {
    router.push(`/validators/${address}`);
  };

  const handleStatusChange = (validator: Validator) => {
    // Here we'll add the API call to change validator status
    console.log('Changing status for validator:', validator.address);
  };

  return (
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeaderCell>Name</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Address</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Status</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Blocks Produced</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Last Active</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Uptime</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Actions</Table.ColumnHeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {validators.map((validator) => (
          <Table.Row key={validator.address}>
            <Table.Cell>
              <Text weight="bold">{validator.name}</Text>
            </Table.Cell>
            <Table.Cell>
              <Text className="font-mono" color="gray">
                {validator.address}
              </Text>
            </Table.Cell>
            <Table.Cell>{getStatusBadge(validator.status)}</Table.Cell>
            <Table.Cell>{validator.blocksProduced}</Table.Cell>
            <Table.Cell>{formatTimeAgo(validator.lastActive)}</Table.Cell>
            <Table.Cell>
              <Text color={parseFloat(validator.uptime) > 99 ? 'green' : 'yellow'}>
                {validator.uptime}
              </Text>
            </Table.Cell>
            <Table.Cell>
              <Flex gap="2">
                <Button 
                  size="1" 
                  variant="soft" 
                  color={validator.status === 'active' ? 'red' : 'green'}
                  onClick={() => handleStatusChange(validator)}
                >
                  {validator.status === 'active' ? 'Deactivate' : 'Activate'}
                </Button>
                <Button 
                  size="1" 
                  variant="soft" 
                  color="blue"
                  onClick={() => handleDetailsClick(validator.address)}
                >
                  Details
                </Button>
              </Flex>
            </Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
} 