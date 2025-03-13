import { NextResponse } from 'next/server';
import * as solc from 'solc';

export async function POST(request: Request) {
  try {
    const { sourceCode } = await request.json();

    if (!sourceCode) {
      return NextResponse.json(
        { error: 'Source code is required' },
        { status: 400 }
      );
    }

    // Prepare input for solc
    const input = {
      language: 'Solidity',
      sources: {
        'contract.sol': {
          content: sourceCode,
        },
      },
      settings: {
        outputSelection: {
          '*': {
            '*': ['*'],
          },
        },
      },
    };

    // Compile
    const output = JSON.parse(solc.compile(JSON.stringify(input)));

    // Check for errors
    if (output.errors) {
      const errors = output.errors.filter((e: any) => e.severity === 'error');
      if (errors.length > 0) {
        return NextResponse.json(
          { error: 'Compilation failed', details: errors[0].message },
          { status: 400 }
        );
      }
    }

    // Get the contract
    const contractFile = Object.keys(output.contracts['contract.sol'])[0];
    const contract = output.contracts['contract.sol'][contractFile];

    return NextResponse.json({
      abi: contract.abi,
      bytecode: '0x' + contract.evm.bytecode.object,
    });
  } catch (error) {
    console.error('API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
} 