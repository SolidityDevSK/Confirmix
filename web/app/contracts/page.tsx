'use client';

import { useState } from 'react';
import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Button, Grid } from '@radix-ui/themes';
import ContractList from '@/components/contracts/ContractList';
import ContractStats from '@/components/contracts/ContractStats';
import DeployContractDialog from '@/components/contracts/DeployContractDialog';

export default function ContractsPage() {
  const [deployDialogOpen, setDeployDialogOpen] = useState(false);

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold mb-4">Smart Contracts</h1>
              <Text className="text-gray-400">
                Deploy and manage smart contracts on the blockchain
              </Text>
            </div>
            <Button 
              size="3" 
              color="blue"
              onClick={() => setDeployDialogOpen(true)}
            >
              Deploy New Contract
            </Button>
          </div>

          <ContractStats />

          <Grid columns="1" gap="6">
            <Card>
              <Text size="3" weight="bold" className="p-4 border-b border-gray-800">
                Deployed Contracts
              </Text>
              <ContractList />
            </Card>
          </Grid>
        </Flex>
      </div>

      <DeployContractDialog 
        open={deployDialogOpen}
        onOpenChange={setDeployDialogOpen}
      />
    </main>
  );
} 