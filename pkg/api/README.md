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