'use client';

import { useState } from 'react';
import { Dialog, Button, Flex, Text, TextField } from '@radix-ui/themes';

interface AddValidatorDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function AddValidatorDialog({ open, onOpenChange }: AddValidatorDialogProps) {
  const [formData, setFormData] = useState({
    name: '',
    address: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Here we'll add the API call to add a new validator
    console.log('Adding validator:', formData);
    onOpenChange(false);
  };

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content>
        <Dialog.Title>Add New Validator</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Add a new validator to the network. Make sure you have the correct address.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Validator Name
              </Text>
              <TextField.Root>
                <TextField.Input
                  placeholder="Enter validator name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
              </TextField.Root>
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Validator Address
              </Text>
              <TextField.Root>
                <TextField.Input
                  placeholder="0x..."
                  value={formData.address}
                  onChange={(e) => setFormData({ ...formData, address: e.target.value })}
                />
              </TextField.Root>
            </label>

            <Flex gap="3" mt="4" justify="end">
              <Dialog.Close>
                <Button variant="soft" color="gray">
                  Cancel
                </Button>
              </Dialog.Close>
              <Button type="submit">Add Validator</Button>
            </Flex>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  );
} 