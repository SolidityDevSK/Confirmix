# Confirmix HTTP API

## Endpoints

### Blockchain Bilgisi

#### GET /info
Blockchain hakkında genel bilgi alır.

```bash
curl http://localhost:8080/info
```

#### GET /blocks
Tüm blokları listeler.

```bash
curl http://localhost:8080/blocks
```

#### GET /blocks/:hash
Hash değeri ile belirli bir bloğu getirir.

```bash
curl http://localhost:8080/blocks/[BLOCK_HASH]
```

### Validator İşlemleri

#### GET /validators
Tüm validatorları listeler.

```bash
curl http://localhost:8080/validators
```

#### GET /validators/current
Şu anki aktif validator'ı gösterir.

```bash
curl http://localhost:8080/validators/current
```

#### POST /validators
Yeni bir validator ekler.

```bash
curl -X POST http://localhost:8080/validators
```

#### DELETE /validators/:address
Belirtilen validator'ı siler.

```bash
curl -X DELETE http://localhost:8080/validators/[VALIDATOR_ADDRESS]
```

### İşlem Gönderme

#### POST /transactions
Yeni bir işlem gönderir.

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "data": "Alice'den Bob'a 50 coin transfer",
    "validator": "[VALIDATOR_ADDRESS]"
  }'
```

## Smart Contract Endpoints

### Deploy Contract
```bash
POST /contracts
```

Deploy a new smart contract.

**Request Body:**
```json
{
    "code": "608060405234801561001057600080fd5b50...",  // Contract bytecode (hex)
    "owner": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
    "name": "MyToken",
    "version": "1.0.0"
}
```

**Response:**
```json
{
    "address": "0x1234...",
    "owner": "0x742d...",
    "name": "MyToken",
    "version": "1.0.0",
    "timestamp": 1647123456,
    "is_enabled": true
}
```

### List Contracts
```bash
GET /contracts
```

Returns a list of all deployed contracts.

**Response:**
```json
[
    {
        "address": "0x1234...",
        "owner": "0x742d...",
        "name": "MyToken",
        "version": "1.0.0",
        "timestamp": 1647123456,
        "is_enabled": true
    }
]
```

### Get Contract
```bash
GET /contracts/:address
```

Returns details of a specific contract.

**Response:**
```json
{
    "address": "0x1234...",
    "owner": "0x742d...",
    "name": "MyToken",
    "version": "1.0.0",
    "timestamp": 1647123456,
    "is_enabled": true
}
```

### Execute Contract
```bash
POST /contracts/:address/execute
```

Execute a smart contract method.

**Request Body:**
```json
{
    "input": "a9059cbb000000000000000000000000..."  // Method call data (hex)
}
```

**Response:**
```json
{
    "result": "0000000000000000000000000000000000000000000000000000000000000001"
}
```

### Disable Contract
```bash
POST /contracts/:address/disable
```

Disable a smart contract (only contract owner).

**Request Body:**
```json
{
    "owner": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
}
```

### Enable Contract
```bash
POST /contracts/:address/enable
```

Enable a disabled smart contract (only contract owner).

**Request Body:**
```json
{
    "owner": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
}
```

## Örnek Kullanım

1. Blockchain bilgisini al:
```bash
curl http://localhost:8080/info
```

2. Mevcut validator'ı kontrol et:
```bash
curl http://localhost:8080/validators/current
```

3. Yeni bir işlem gönder:
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "data": "Alice'den Bob'a 50 coin transfer",
    "validator": "[VALIDATOR_ADDRESS]"
  }'
```

4. Tüm blokları listele:
```bash
curl http://localhost:8080/blocks
```

## Notlar

- Tüm POST istekleri için `Content-Type: application/json` header'ı gereklidir
- Validator adresleri, validator oluşturulduğunda console'da görüntülenir
- Round-Robin konsensüs nedeniyle, işlemler sadece sırası gelen validator tarafından eklenebilir
- Bloklar arası minimum 5 saniyelik bekleme süresi vardır
- All POST requests must include the `Content-Type: application/json` header
- Contract code and input data must be hex-encoded
- Contract addresses are automatically generated based on the code, owner, and timestamp
- Only the contract owner can disable or enable a contract
- Contract execution follows the EVM specification 