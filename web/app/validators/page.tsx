'use client';

import { useState } from 'react';
import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Button } from '@radix-ui/themes';
import ValidatorList from '@/components/validators/ValidatorList';
import ValidatorStats from '@/components/validators/ValidatorStats';
import AddValidatorDialog from '@/components/validators/AddValidatorDialog';

export default function ValidatorsPage() {
  const [addDialogOpen, setAddDialogOpen] = useState(false);

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold mb-4">Validators</h1>
              <Text className="text-gray-400">
                Manage network validators and view their performance
              </Text>
            </div>
            <Button 
              size="3" 
              color="blue"
              onClick={() => setAddDialogOpen(true)}
            >
              Add Validator
            </Button>
          </div>

          <ValidatorStats />

          <Card>
            <Text size="3" weight="bold" className="p-4">Active Validators</Text>
            <ValidatorList />
          </Card>
        </Flex>
      </div>

      <AddValidatorDialog 
        open={addDialogOpen}
        onOpenChange={setAddDialogOpen}
      />
    </main>
  );
} 