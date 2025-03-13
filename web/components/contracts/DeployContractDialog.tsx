'use client';

import { useState, useEffect } from 'react';
import { Dialog, Button, Flex, Text, TextField, TextArea, Select, Tabs, Grid, Card } from '@radix-ui/themes';
import { CodeIcon, UploadIcon, ReloadIcon, CheckCircledIcon, CrossCircledIcon } from '@radix-ui/react-icons';
import { blockchainService } from '@/services/blockchain';
import { websocketService } from '@/services/websocket';
import type { ContractTemplate, ContractType, BlockchainContract, ContractEvent } from '@/types/blockchain';
import ContractTemplates from './ContractTemplates';
import TransactionMonitor from './TransactionMonitor';

declare global {
  interface Window {
    ethereum?: {
      selectedAddress: string;
      request: (args: { method: string; params?: any[] }) => Promise<any>;
    };
  }
}

interface DeployContractDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

interface FormData {
  name: string;
  type: ContractType;
  code: string;
  constructorArgs?: Record<string, string>;
}

interface GasEstimate {
  gasLimit: string;
  gasPrice: string;
  totalCost: string;
}

interface DeploymentStatus {
  step: 'compiling' | 'estimating' | 'deploying' | 'verifying';
  message: string;
  txHash?: string;
}

function getConstructorInputs(sourceCode: string): Array<{ name: string; type: string; id: string }> {
  // Find the constructor in the source code using a simpler regex
  const constructorMatch = sourceCode.match(/constructor\s*\(([\s\S]*?)\)/);
  if (!constructorMatch) return [];

  const params = constructorMatch[1].split(',').filter(Boolean);
  return params.map((param, index) => {
    const [type, name] = param.trim().split(' ');
    return {
      type: type.replace('memory', '').trim(),
      name: name.replace('_', '').trim(),
      id: `${index}-${name.replace('_', '').trim()}-${type.replace('memory', '').trim()}`
    };
  });
}

export default function DeployContractDialog({ open, onOpenChange }: DeployContractDialogProps) {
  const [selectedTemplate, setSelectedTemplate] = useState<ContractTemplate | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [constructorArgs, setConstructorArgs] = useState<Record<string, string>>({});
  const [isDeploying, setIsDeploying] = useState(false);
  const [deploymentSuccess, setDeploymentSuccess] = useState(false);

  useEffect(() => {
    if (isDeploying) {
      const unsubscribe = websocketService.subscribeToContractStatus('global', (status) => {
        if (status.type === 'CONTRACT_DEPLOY_SUCCESS' && 
            status.data.name === selectedTemplate?.name) {
          setDeploymentSuccess(true);
          setTimeout(() => {
            onOpenChange(false);
          }, 2000);
        }
      });

      return () => unsubscribe();
    }
  }, [isDeploying, selectedTemplate?.name]);

  const handleDeploy = async () => {
    if (!selectedTemplate) return;
    
    try {
      setError(null);
      setIsDeploying(true);
      setDeploymentSuccess(false);

      // Validate contract code
      if (!selectedTemplate.code || selectedTemplate.code.trim() === '') {
        throw new Error('Contract code is empty');
      }

      // Convert contract code to hex
      const encoder = new TextEncoder();
      const codeBytes = encoder.encode(selectedTemplate.code.trim());
      const codeHex = Array.from(codeBytes)
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');

      const deploymentData = {
        name: selectedTemplate.name,
        type: selectedTemplate.type,
        code: codeHex,
        constructorArgs: Object.values(constructorArgs).map(arg => arg.trim()).filter(Boolean),
        version: '1.0.0',
        owner: window.ethereum?.selectedAddress || '0x0000000000000000000000000000000000000000',
        abi: []
      };

      console.log('Selected template:', selectedTemplate);
      console.log('Deploying contract with data:', deploymentData);
      const response = await blockchainService.deployContract(deploymentData);
      console.log('Deploy response:', response);
      
    } catch (err: any) {
      console.error('Deploy error details:', err);
      const errorMessage = err?.response?.data?.error || err?.message || 'Failed to deploy contract. Please try again.';
      setError(errorMessage);
      setIsDeploying(false);
      setDeploymentSuccess(false);
    }
  };

  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Content size="4">
        <Dialog.Title>
          {deploymentSuccess ? 'Deployment Successful' : (isDeploying ? 'Deploying Contract' : 'Deploy New Contract')}
        </Dialog.Title>
        
        {isDeploying ? (
          <div className="py-8 animate-fadeIn">
            <Flex direction="column" gap="4" align="center">
              <div className="rounded-full bg-blue-500/10 p-4">
                {deploymentSuccess ? (
                  <CheckCircledIcon className="w-12 h-12 text-green-500 animate-scaleIn" />
                ) : (
                  <ReloadIcon className="w-12 h-12 text-blue-500 animate-spin" />
                )}
              </div>
              
              <div className="text-center space-y-2">
                <Text size="4" weight="bold" className={deploymentSuccess ? "text-green-500" : "text-blue-500"}>
                  {deploymentSuccess ? 'Contract Deployed Successfully!' : 'Deploying Your Contract'}
                </Text>
                <Text size="2" color="gray">
                  {deploymentSuccess 
                    ? 'Your contract has been deployed. The dialog will close automatically.'
                    : 'Please wait while your contract is being deployed. This may take a few minutes.'}
                </Text>
              </div>

              <Button 
                size="3"
                className="mt-4 animate-slideUp"
                onClick={() => onOpenChange(false)}
              >
                Close
              </Button>
            </Flex>
          </div>
        ) : (
          <>
            {!selectedTemplate ? (
              <>
                <Text size="2" mb="4">
                  Select a contract template to deploy
                </Text>
                <ContractTemplates onSelect={setSelectedTemplate} />
              </>
            ) : (
              <div className="space-y-4">
                <Flex justify="between" align="center">
                  <div>
                    <Text size="3" weight="bold">{selectedTemplate.name}</Text>
                    <Text size="2" color="gray">{selectedTemplate.description}</Text>
                  </div>
                  <Button variant="soft" onClick={() => setSelectedTemplate(null)}>
                    Change Template
                  </Button>
                </Flex>

                {/* Constructor Arguments Form */}
                <div className="space-y-3">
                  <Text size="2" weight="bold">Constructor Arguments</Text>
                  {getConstructorInputs(selectedTemplate.code).map((input) => (
                    <div key={input.id}>
                      <Text as="div" size="2" mb="1">
                        {input.name} ({input.type})
                      </Text>
                      <input
                        type="text"
                        className="w-full px-3 py-2 bg-gray-900 border border-gray-800 rounded-md"
                        placeholder={`Enter ${input.type}`}
                        value={constructorArgs[input.name] || ''}
                        onChange={(e) => setConstructorArgs(prev => ({
                          ...prev,
                          [input.name]: e.target.value
                        }))}
                      />
                    </div>
                  ))}
                </div>

                {error && (
                  <Text color="red" size="2">
                    {error}
                  </Text>
                )}

                <Flex gap="3" mt="4">
                  <Button onClick={handleDeploy}>
                    Deploy Contract
                  </Button>
                  <Button 
                    variant="soft" 
                    color="gray" 
                    onClick={() => onOpenChange(false)}
                  >
                    Cancel
                  </Button>
                </Flex>
              </div>
            )}
          </>
        )}
      </Dialog.Content>
    </Dialog.Root>
  );
} 