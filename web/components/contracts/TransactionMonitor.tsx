import { Card, Text, Flex, Badge } from '@radix-ui/themes';
import { useEffect, useState } from 'react';
import { blockchainService } from '@/services/blockchain';

interface TransactionMonitorProps {
  hash: string;
  onComplete?: (success: boolean) => void;
}

interface TransactionStatus {
  confirmations: number;
  status: 'pending' | 'success' | 'failed';
  blockNumber?: number;
  gasUsed?: string;
}

export default function TransactionMonitor({ hash, onComplete }: TransactionMonitorProps) {
  const [status, setStatus] = useState<TransactionStatus>({
    confirmations: 0,
    status: 'pending'
  });
  const [checkCount, setCheckCount] = useState(0);

  useEffect(() => {
    let interval: NodeJS.Timeout;
    
    const checkStatus = async () => {
      try {
        setCheckCount(prev => prev + 1);
        const txStatus = await blockchainService.getTransactionStatus(hash);
        
        if (txStatus) {
          setStatus({
            confirmations: txStatus.confirmations,
            status: txStatus.status as 'success' | 'failed',
            blockNumber: txStatus.blockNumber,
            gasUsed: txStatus.gasUsed
          });

          if (txStatus.confirmations >= 1) {
            clearInterval(interval);
            onComplete?.(txStatus.status === 'success');
          }
        }
      } catch (error) {
        console.error('Error checking transaction status:', error);
      }
    };

    // Check immediately and then every 2 seconds
    checkStatus();
    interval = setInterval(checkStatus, 2000);

    return () => clearInterval(interval);
  }, [hash, onComplete]);

  const getStatusMessage = () => {
    if (status.status === 'pending') {
      return checkCount > 5 
        ? 'Still waiting for confirmation... This might take a few more moments.'
        : 'Waiting for blockchain confirmation...';
    }
    return status.status === 'success' 
      ? 'Transaction confirmed successfully!'
      : 'Transaction failed. Please check the details below.';
  };

  const getStatusBadge = () => {
    const colors = {
      pending: 'yellow',
      success: 'green',
      failed: 'red'
    };

    return (
      <Badge color={colors[status.status] as any} size="1">
        {status.status === 'pending' ? 'Pending...' : status.status.charAt(0).toUpperCase() + status.status.slice(1)}
      </Badge>
    );
  };

  return (
    <Card>
      <Flex direction="column" gap="2" p="4">
        <Flex justify="between" align="center">
          <Text weight="bold">Transaction Status</Text>
          {getStatusBadge()}
        </Flex>

        <Text size="2" color={status.status === 'failed' ? 'red' : undefined}>
          {getStatusMessage()}
        </Text>

        <Text className="font-mono" size="1" color="gray">
          Transaction Hash: {hash}
        </Text>

        {status.blockNumber && (
          <Text size="1">
            Block Number: #{status.blockNumber}
          </Text>
        )}

        {status.gasUsed && (
          <Text size="1">
            Gas Used: {status.gasUsed}
          </Text>
        )}

        <Text size="1">
          Confirmations: {status.confirmations}
          {status.status === 'pending' && status.confirmations === 0 && ' (Waiting for first confirmation...)'}
        </Text>
      </Flex>
    </Card>
  );
} 