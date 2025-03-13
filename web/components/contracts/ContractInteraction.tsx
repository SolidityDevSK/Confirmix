'use client';

import { useState, useMemo } from 'react';
import { Card, Text, Button, Select, Flex, Badge } from '@radix-ui/themes';
import type { BlockchainContract, ContractABI } from '@/types/blockchain';
import { blockchainService } from '@/services/blockchain';
import { ethers } from 'ethers';

interface ContractInteractionProps {
  contract: BlockchainContract;
}

interface FunctionInput {
  name: string;
  type: string;
  value: string;
  error?: string;
}

interface ValidationRule {
  type: string;
  validate: (value: string) => boolean;
  format: (value: string) => string;
  error: string;
}

export default function ContractInteraction({ contract }: ContractInteractionProps) {
  const [selectedFunction, setSelectedFunction] = useState('');
  const [inputs, setInputs] = useState<FunctionInput[]>([]);
  const [result, setResult] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const validationRules: Record<string, ValidationRule> = {
    address: {
      type: 'address',
      validate: (value) => /^0x[a-fA-F0-9]{40}$/.test(value),
      format: (value) => value.toLowerCase(),
      error: 'Invalid Ethereum address format',
    },
    uint256: {
      type: 'uint256',
      validate: (value) => {
        try {
          // Check if the value is a valid number and non-negative
          return /^\d+$/.test(value) && value !== '';
        } catch {
          return false;
        }
      },
      format: (value) => value,
      error: 'Invalid uint256 value (must be a non-negative number)',
    },
    bool: {
      type: 'bool',
      validate: (value) => value === 'true' || value === 'false',
      format: (value) => value.toLowerCase(),
      error: 'Value must be true or false',
    },
    string: {
      type: 'string',
      validate: () => true,
      format: (value) => value,
      error: 'Invalid string value',
    },
    bytes32: {
      type: 'bytes32',
      validate: (value) => /^0x[a-fA-F0-9]{64}$/.test(value),
      format: (value) => value.toLowerCase(),
      error: 'Invalid bytes32 format',
    },
  };

  const handleFunctionChange = (functionName: string) => {
    setSelectedFunction(functionName);
    const func = contract.abi.find(f => f.name === functionName);
    if (func) {
      setInputs(func.inputs.map((input) => ({
        name: input.name,
        type: input.type,
        value: '',
      })));
    }
    setResult(null);
    setError(null);
  };

  const validateInput = (input: FunctionInput): FunctionInput => {
    const baseType = input.type.replace(/\[\]$/, ''); // Remove array suffix if present
    const isArray = input.type.endsWith('[]');
    
    if (isArray) {
      try {
        const values = JSON.parse(input.value);
        if (!Array.isArray(values)) {
          return { ...input, error: 'Invalid array format' };
        }
        
        const rule = validationRules[baseType];
        if (!rule) {
          return { ...input, error: `Unsupported type: ${input.type}` };
        }

        const invalidValue = values.find(v => !rule.validate(v));
        if (invalidValue) {
          return { ...input, error: rule.error };
        }

        return { ...input, error: undefined };
      } catch {
        return { ...input, error: 'Invalid array format' };
      }
    } else {
      const rule = validationRules[baseType];
      if (!rule) {
        return { ...input, error: `Unsupported type: ${input.type}` };
      }

      return {
        ...input,
        error: rule.validate(input.value) ? undefined : rule.error,
      };
    }
  };

  const handleInputChange = (index: number, value: string) => {
    const newInputs = [...inputs];
    newInputs[index] = validateInput({ ...newInputs[index], value });
    setInputs(newInputs);
  };

  const formatInputValue = (input: FunctionInput): any => {
    const baseType = input.type.replace(/\[\]$/, '');
    const isArray = input.type.endsWith('[]');
    const rule = validationRules[baseType];

    if (!rule) {
      throw new Error(`Unsupported type: ${input.type}`);
    }

    if (isArray) {
      const values = JSON.parse(input.value);
      return values.map(rule.format);
    }

    return rule.format(input.value);
  };

  const getFunctionStateMutability = (functionName: string): string | undefined => {
    const func = contract.abi.find(f => f.name === functionName);
    return func?.stateMutability;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setLoading(true);
      setError(null);
      setResult(null);

      // Validate all inputs
      const validatedInputs = inputs.map(validateInput);
      const hasErrors = validatedInputs.some(input => input.error);
      if (hasErrors) {
        setInputs(validatedInputs);
        throw new Error('Please fix input errors before submitting');
      }

      // Format input values
      const formattedValues = inputs.map(formatInputValue);

      const result = await blockchainService.callContractMethod(
        contract.address,
        selectedFunction,
        formattedValues
      );

      setResult(JSON.stringify(result, null, 2));
    } catch (err) {
      console.error('Error calling contract method:', err);
      setError(err instanceof Error ? err.message : 'Failed to call contract method');
    } finally {
      setLoading(false);
    }
  };

  const getMethodBadge = (functionName: string) => {
    const stateMutability = getFunctionStateMutability(functionName);
    const colors: Record<string, any> = {
      view: 'blue',
      pure: 'blue',
      nonpayable: 'orange',
      payable: 'green',
    };
    return stateMutability ? (
      <Badge color={colors[stateMutability]} size="1">
        {stateMutability}
      </Badge>
    ) : null;
  };

  const getInputPlaceholder = (type: string): string => {
    const baseType = type.replace(/\[\]$/, '');
    const isArray = type.endsWith('[]');
    
    const placeholders: Record<string, string> = {
      address: '0x...',
      uint256: '1000000',
      bool: 'true or false',
      string: 'text...',
      bytes32: '0x...',
    };

    const basePlaceholder = placeholders[baseType] || type;
    return isArray ? `[${basePlaceholder}, ...]` : basePlaceholder;
  };

  // Get contract methods from ABI
  const methods = useMemo(() => {
    if (!contract?.abi) return [];
    
    return contract.abi
      .filter(item => item?.type === 'function')
      .map(method => ({
        name: method.name,
        inputs: method.inputs || [],
        outputs: method.outputs || [],
        stateMutability: method.stateMutability
      }));
  }, [contract?.abi]);

  // Group methods by type
  const groupedMethods = useMemo(() => {
    if (!methods.length) return {
      view: [],
      nonpayable: [],
      payable: []
    };

    return methods.reduce((acc, method) => {
      const type = method.stateMutability || 'nonpayable';
      acc[type] = [...(acc[type] || []), method];
      return acc;
    }, {} as Record<string, typeof methods>);
  }, [methods]);

  return (
    <div>
      <Text size="3" weight="bold" mb="4">
        Contract Functions
      </Text>

      <Card>
        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <div>
              <Text as="div" size="2" mb="1" weight="bold">
                Select Function
              </Text>
              <Select.Root
                value={selectedFunction}
                onValueChange={handleFunctionChange}
                disabled={loading}
              >
                <Select.Trigger placeholder="Choose a function" />
                <Select.Content>
                  {contract.abi
                    .filter((item) => item.type === 'function')
                    .map((func) => (
                      <Select.Item key={func.name} value={func.name}>
                        <Flex gap="2" align="center">
                          {func.name}
                          {getMethodBadge(func.name)}
                        </Flex>
                      </Select.Item>
                    ))}
                </Select.Content>
              </Select.Root>
            </div>

            {inputs.map((input, index) => (
              <div key={input.name}>
                <Text as="div" size="2" mb="1" weight="bold">
                  {input.name} ({input.type})
                </Text>
                <input
                  type="text"
                  className={`w-full px-3 py-2 bg-gray-900 border rounded-md text-white font-mono ${
                    input.error ? 'border-red-500' : 'border-gray-800'
                  }`}
                  placeholder={getInputPlaceholder(input.type)}
                  value={input.value}
                  onChange={(e) => handleInputChange(index, e.target.value)}
                  disabled={loading}
                />
                {input.error && (
                  <Text color="red" size="1" mt="1">
                    {input.error}
                  </Text>
                )}
              </div>
            ))}

            {error && (
              <Text color="red" size="2">
                {error}
              </Text>
            )}

            {selectedFunction && (
              <Button 
                type="submit" 
                mt="4" 
                disabled={loading || inputs.some(input => input.error)}
              >
                {loading ? 'Executing...' : 'Execute Function'}
              </Button>
            )}
          </Flex>
        </form>
      </Card>

      {result && (
        <Card mt="4">
          <Text size="2" weight="bold" mb="2">
            Result
          </Text>
          <pre className="font-mono text-sm whitespace-pre-wrap bg-gray-900 p-4 rounded-lg overflow-x-auto">
            {result}
          </pre>
        </Card>
      )}
    </div>
  );
} 