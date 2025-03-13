'use client';

import { useState } from 'react';
import { Dialog, Button, Flex, Text } from '@radix-ui/themes';
import { blockchainService } from '@/services/blockchain';
import type { BlockchainContract } from '@/types/blockchain';

interface TransferOwnershipDialogProps {
  contract: BlockchainContract;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export default function TransferOwnershipDialog({ 
  contract, 
  open, 
  onOpenChange,
  onSuccess 
}: TransferOwnershipDialogProps) {
  const [newOwner, setNewOwner] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);

      // Call the contract's transferOwnership method
      await blockchainService.callContractMethod(
        contract.address,
        'transferOwnership',
        [newOwner]
      );

      onOpenChange(false);
      onSuccess?.();
    } catch (err) {
      console.error('Error transferring ownership:', err);
      setError(err instanceof Error ? err.message : 'Failed to transfer ownership');
    } finally {
      setLoading(false);
    }
  };

  const validateAddress = (address: string) => {
    // Basic Ethereum address validation
    return /^0x[a-fA-F0-9]{40}$/.test(address);
  };

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content>
        <Dialog.Title>Transfer Contract Ownership</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Transfer ownership of this contract to another address. This action cannot be undone.
        </Dialog.Description>

        {error && (
          <Text color="red" size="2" mb="4">
            {error}
          </Text>
        )}

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Current Owner
              </Text>
              <Text className="font-mono text-gray-400">
                {contract.owner}
              </Text>
            </div>

            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                New Owner Address
              </Text>
              <input
                type="text"
                className="w-full px-3 py-2 bg-gray-900 border border-gray-800 rounded-md text-white font-mono"
                placeholder="0x..."
                value={newOwner}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewOwner(e.target.value)}
                disabled={loading}
              />
              {newOwner && !validateAddress(newOwner) && (
                <Text color="red" size="1" mt="1">
                  Please enter a valid Ethereum address
                </Text>
              )}
            </div>

            <Flex gap="3" mt="4" justify="end">
              <Dialog.Close>
                <Button variant="soft" color="gray" disabled={loading}>
                  Cancel
                </Button>
              </Dialog.Close>
              <Button 
                type="submit" 
                disabled={loading || !validateAddress(newOwner)}
                color="red"
              >
                {loading ? 'Transferring...' : 'Transfer Ownership'}
              </Button>
            </Flex>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  );
} 