'use client';

import { Card, Text, Tabs, Box, Button, Flex } from '@radix-ui/themes';
import { CopyIcon, CheckIcon } from '@radix-ui/react-icons';
import { useState } from 'react';
import type { BlockchainContract } from '@/types/blockchain';
import Link from 'next/link';

interface ContractCodeProps {
  contract: BlockchainContract;
}

export default function ContractCode({ contract }: ContractCodeProps) {
  const [copied, setCopied] = useState<'source' | 'bytecode' | null>(null);

  const handleCopy = async (type: 'source' | 'bytecode') => {
    const text = type === 'source' ? contract.sourceCode : contract.bytecode;
    if (text) {
      await navigator.clipboard.writeText(text);
      setCopied(type);
      setTimeout(() => setCopied(null), 2000);
    }
  };

  if (!contract.sourceCode && !contract.bytecode) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[200px] gap-4">
        <Text color="gray">Contract code is not available</Text>
        <Link href={`/contracts/${contract.address}/verify`}>
          <Button variant="soft" color="blue">
            Verify Contract
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div>
      <Text size="3" weight="bold" mb="4">
        Contract Code
      </Text>

      <Tabs.Root defaultValue={contract.sourceCode ? 'source' : 'bytecode'}>
        <Tabs.List>
          <Tabs.Trigger 
            value="source" 
            disabled={!contract.sourceCode}
          >
            Source Code
          </Tabs.Trigger>
          <Tabs.Trigger 
            value="bytecode"
            disabled={!contract.bytecode}
          >
            Bytecode
          </Tabs.Trigger>
        </Tabs.List>

        <Box mt="4">
          <Tabs.Content value="source">
            {contract.sourceCode ? (
              <Card>
                <Flex justify="between" align="center" mb="2">
                  <Text size="2" color="gray">Solidity Source Code</Text>
                  <Button 
                    variant="soft" 
                    color={copied === 'source' ? 'green' : 'gray'}
                    onClick={() => handleCopy('source')}
                  >
                    {copied === 'source' ? (
                      <CheckIcon />
                    ) : (
                      <CopyIcon />
                    )}
                    {copied === 'source' ? 'Copied!' : 'Copy Code'}
                  </Button>
                </Flex>
                <Box 
                  className="bg-gray-900 p-4 rounded-lg overflow-x-auto"
                  style={{ maxHeight: '500px' }}
                >
                  <pre className="font-mono text-sm whitespace-pre-wrap">
                    {contract.sourceCode}
                  </pre>
                </Box>
              </Card>
            ) : (
              <div className="flex flex-col items-center justify-center min-h-[200px] gap-4">
                <Text color="gray">Source code is not verified</Text>
                <Link href={`/contracts/${contract.address}/verify`}>
                  <Button variant="soft" color="blue">
                    Verify Contract
                  </Button>
                </Link>
              </div>
            )}
          </Tabs.Content>

          <Tabs.Content value="bytecode">
            {contract.bytecode ? (
              <Card>
                <Flex justify="between" align="center" mb="2">
                  <Text size="2" color="gray">Contract Bytecode</Text>
                  <Button 
                    variant="soft" 
                    color={copied === 'bytecode' ? 'green' : 'gray'}
                    onClick={() => handleCopy('bytecode')}
                  >
                    {copied === 'bytecode' ? (
                      <CheckIcon />
                    ) : (
                      <CopyIcon />
                    )}
                    {copied === 'bytecode' ? 'Copied!' : 'Copy Bytecode'}
                  </Button>
                </Flex>
                <Box 
                  className="bg-gray-900 p-4 rounded-lg overflow-x-auto"
                  style={{ maxHeight: '200px' }}
                >
                  <Text className="font-mono text-sm break-all">
                    {contract.bytecode}
                  </Text>
                </Box>
              </Card>
            ) : (
              <div className="flex justify-center items-center min-h-[200px]">
                <Text color="gray">Bytecode is not available</Text>
              </div>
            )}
          </Tabs.Content>
        </Box>
      </Tabs.Root>
    </div>
  );
} 