import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, Table, ScrollArea } from '@radix-ui/themes';
import BlockList from '@/components/explorer/BlockList';
import TransactionList from '@/components/explorer/TransactionList';

export default function ExplorerPage() {
  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div>
            <h1 className="text-3xl font-bold mb-4">Blockchain Explorer</h1>
            <Text className="text-gray-400">
              Explore blocks, transactions, and smart contract interactions
            </Text>
          </div>

          <Flex gap="4">
            {/* Left column - Block List */}
            <Card className="flex-1">
              <Text size="3" weight="bold" mb="4">Latest Blocks</Text>
              <BlockList />
            </Card>

            {/* Right column - Transaction List */}
            <Card className="flex-1">
              <Text size="3" weight="bold" mb="4">Recent Transactions</Text>
              <TransactionList />
            </Card>
          </Flex>
        </Flex>
      </div>
    </main>
  );
} 