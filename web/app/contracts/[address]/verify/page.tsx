'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Navigation from '@/components/shared/Navigation';
import { Card, Flex, Text, TextField, Button } from '@radix-ui/themes';
import Link from 'next/link';
import { blockchainService } from '@/lib/blockchain';

export default function VerifyContractPage() {
  const params = useParams();
  const router = useRouter();
  const [sourceCode, setSourceCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);

      // First compile the contract to verify it compiles correctly
      const compileResponse = await fetch('/api/compile', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ sourceCode }),
      });

      if (!compileResponse.ok) {
        const error = await compileResponse.json();
        throw new Error(error.details || 'Failed to compile contract');
      }

      // If compilation succeeds, verify the contract
      await blockchainService.verifyContract(params.address as string, sourceCode);

      // Redirect back to contract details page
      router.push(`/contracts/${params.address}`);
    } catch (err) {
      console.error('Error verifying contract:', err);
      setError(err instanceof Error ? err.message : 'Failed to verify contract');
    } finally {
      setLoading(false);
    }
  };

  return (
    <main>
      <Navigation />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Flex direction="column" gap="6">
          <div>
            <Link 
              href={`/contracts/${params.address}`}
              className="text-blue-500 hover:text-blue-400 mb-4 inline-block"
            >
              ← Back to Contract
            </Link>
            <h1 className="text-3xl font-bold mb-4">Verify Contract</h1>
            <Text className="text-gray-400">
              Verify and publish your contract&apos;s source code
            </Text>
          </div>

          <Card>
            <form onSubmit={handleSubmit}>
              <Flex direction="column" gap="4" p="4">
                <div>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Contract Address
                  </Text>
                  <Text className="font-mono text-gray-400">
                    {params.address}
                  </Text>
                </div>

                <div>
                  <Text as="div" size="2" mb="1" weight="bold">
                    Source Code
                  </Text>
                  <TextField.Root>
                    <TextField.Input
                      placeholder="Paste your Solidity source code here"
                      value={sourceCode}
                      onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => 
                        setSourceCode(e.target.value)
                      }
                      multiline
                      style={{ height: 400, fontFamily: 'monospace' }}
                      disabled={loading}
                    />
                  </TextField.Root>
                  <Text size="1" color="gray" mt="1">
                    Make sure the source code matches exactly with the deployed contract
                  </Text>
                </div>

                {error && (
                  <Text color="red" size="2">
                    {error}
                  </Text>
                )}

                <Flex gap="3" mt="4" justify="end">
                  <Link href={`/contracts/${params.address}`}>
                    <Button variant="soft" color="gray" disabled={loading}>
                      Cancel
                    </Button>
                  </Link>
                  <Button type="submit" disabled={loading}>
                    {loading ? 'Verifying...' : 'Verify Contract'}
                  </Button>
                </Flex>
              </Flex>
            </form>
          </Card>
        </Flex>
      </div>
    </main>
  );
} 