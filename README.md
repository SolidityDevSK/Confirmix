# Confirmix Blockchain Explorer

A modern blockchain explorer and smart contract management platform for the Confirmix PoA (Proof of Authority) blockchain.

## Project Structure

```
blockchain/
├── cmd/
│   └── node/               # Blockchain node entry point
├── internal/
│   └── validator/          # Validator implementation
├── pkg/
│   ├── api/               # REST API implementation
│   ├── blockchain/        # Core blockchain implementation
│   ├── consensus/         # PoA consensus implementation
│   └── testing/          # Testing utilities
└── web/                   # Frontend application
    ├── app/              # Next.js pages and routing
    ├── components/       # React components
    ├── contexts/         # React contexts
    ├── lib/              # Utility functions
    ├── services/         # API service layer
    └── types/            # TypeScript type definitions
```

## Working Components

### Backend (Go)

1. **Blockchain Core**
   - Block creation and validation
   - Transaction processing
   - Smart contract deployment and execution
   - PoA consensus mechanism

2. **Validator System**
   - Validator management
   - Block signing and verification
   - Authority rotation

3. **API Server**
   - RESTful endpoints for blockchain interaction
   - WebSocket support for real-time updates
   - Contract deployment and interaction

### Frontend (Next.js)

1. **Overview Dashboard**
   - Blockchain statistics
   - Recent blocks and transactions
   - Network health indicators

2. **Explorer**
   - Block explorer
   - Transaction viewer
   - Address lookup

3. **Smart Contracts**
   - Contract deployment interface
   - Contract verification
   - Contract interaction UI
   - Event monitoring

4. **Validator Management**
   - Validator list and status
   - Performance metrics
   - Authority management

## Setup and Running

1. Start the blockchain node:
```bash
go run cmd/node/main.go
```

2. Install web dependencies:
```bash
cd web
npm install
```

3. Start the web application:
```bash
npm run dev
```

## Environment Configuration

### Backend
Default configuration:
- API Port: 8080
- WebSocket Port: 8081

### Frontend
Configure in `.env.local`:
```
NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://localhost:8080
NEXT_PUBLIC_BLOCKCHAIN_WS_URL=ws://localhost:8081
```

## Future Development Plans

### Short-term

1. **Smart Contract Improvements**
   - [ ] Contract template library
   - [ ] Batch deployment support
   - [ ] Contract upgradeability

2. **Explorer Enhancements**
   - [ ] Advanced transaction filtering
   - [ ] Token tracking
   - [ ] Gas analytics

3. **Validator Features**
   - [ ] Validator performance metrics
   - [ ] Automated health checks
   - [ ] Slashing conditions

### Long-term

1. **Scalability**
   - [ ] Sharding support
   - [ ] Layer 2 solutions
   - [ ] State pruning

2. **Security**
   - [ ] Audit logging
   - [ ] Permission management
   - [ ] Multi-signature support

3. **Integration**
   - [ ] Cross-chain bridges
   - [ ] Oracle integration
   - [ ] API SDK development

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 