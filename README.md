# Confirmix

Confirmix, Go programlama dili ile geliştirilmiş, Proof of Authority (PoA) konsensüs mekanizmasını kullanan bir blockchain implementasyonudur.

## Özellikler

### Mevcut Özellikler
- ✅ Proof of Authority (PoA) konsensüs mekanizması
- ✅ Çoklu validator desteği
- ✅ ECDSA tabanlı dijital imza sistemi
- ✅ Blok zinciri doğrulama
- ✅ Genesis blok oluşturma
- ✅ Blok imzalama ve doğrulama
- ✅ Round-Robin validator sıralama
- ✅ HTTP API desteği
- ✅ P2P ağ desteği
- ✅ Akıllı kontrat desteği

### Geliştirme Aşamasındaki Özellikler
- 🔄 Validator oylama sistemi
- 🔄 Web arayüzü

## Kurulum

### Gereksinimler
- Go 1.24 veya üzeri
- Gin web framework
- libp2p
- go-ethereum

### Kurulum Adımları
1. Repoyu klonlayın:
```bash
git clone https://github.com/SolidityDevSK/confirmix.git
```

2. Proje dizinine gidin:
```bash
cd confirmix
```

3. Bağımlılıkları yükleyin:
```bash
go mod download
```

4. Projeyi çalıştırın:
```bash
# İlk node'u başlat
go run cmd/confirmix/main.go -api-port 8080 -p2p-port 9000

# İkinci node'u başlat ve ilk node'a bağlan
go run cmd/confirmix/main.go -api-port 8081 -p2p-port 9001 -bootstrap /ip4/127.0.0.1/tcp/9000/p2p/FIRST_NODE_ID
```

## HTTP API

API dokümantasyonu için [API README](pkg/api/README.md) dosyasına bakın.

### Örnek API Kullanımı

1. Blockchain bilgisini al:
```bash
curl http://localhost:8080/info
```

2. Yeni bir işlem gönder:
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "data": "Alice'den Bob'a 50 coin transfer",
    "validator": "[VALIDATOR_ADDRESS]"
  }'
```

## P2P Ağ

### Node Başlatma
```bash
# Bootstrap node
go run cmd/confirmix/main.go -p2p-port 9000

# Diğer node'lar
go run cmd/confirmix/main.go -p2p-port 9001 -bootstrap BOOTSTRAP_NODE_ADDR
```

### Özellikler
- Otomatik peer keşfi
- Blockchain senkronizasyonu
- Blok ve validator duyuruları
- Güvenli P2P iletişim

## Akıllı Kontratlar

### Özellikler
- EVM (Ethereum Virtual Machine) uyumlu
- Solidity kontratlarını destekler
- Kontrat yönetimi (deploy, execute, disable/enable)
- Kontrat sahipliği ve yetkilendirme
- Gas limiti ve kontrol

### Örnek Kontrat Deploy
```bash
curl -X POST http://localhost:8080/contracts \
  -H "Content-Type: application/json" \
  -d '{
    "code": "608060405234801561001057600080fd5b50...",
    "owner": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
    "name": "MyToken",
    "version": "1.0.0"
  }'
```

### Örnek Kontrat Çağrısı
```bash
curl -X POST http://localhost:8080/contracts/0x1234.../execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": "a9059cbb000000000000000000000000..."
  }'
```

## Proje Yapısı

```
confirmix/
├── cmd/                    # Uygulamanın giriş noktaları
│   └── confirmix/         # Ana uygulama
├── pkg/                    # Dışa açık paketler
│   ├── api/               # HTTP API implementasyonu
│   ├── blockchain/        # Blockchain çekirdek yapısı
│   ├── consensus/         # Konsensüs mekanizmaları
│   ├── contracts/        # Akıllı kontrat sistemi
│   ├── network/          # P2P ağ implementasyonu
│   └── utils/            # Yardımcı fonksiyonlar
├── internal/              # Sadece içeride kullanılan paketler
│   └── validator/        # Validator işlemleri
├── docs/                  # Dokümantasyon
└── tests/                 # Test dosyaları
```

## Nasıl Çalışır?

1. **Validator Sistemi**
   - Her validator için ECDSA public/private anahtar çifti oluşturulur
   - Validatorlar blokları kendi private key'leri ile imzalar
   - İmzalar diğer validatorlar tarafından doğrulanır

2. **Round-Robin Konsensüs**
   - Validatorlar sırayla blok oluşturur
   - Her blok arasında minimum süre (5 saniye) beklenir
   - Sadece sırası gelen validator blok oluşturabilir

3. **P2P Ağ**
   - libp2p tabanlı P2P iletişim
   - Otomatik peer keşfi ve bağlantı
   - Blockchain senkronizasyonu
   - Blok ve validator duyuruları

4. **HTTP API**
   - RESTful API ile blockchain yönetimi
   - Blok ve validator işlemleri
   - İşlem gönderme ve sorgulama

4. **Akıllı Kontratlar**
   - EVM tabanlı akıllı kontrat çalıştırma ortamı
   - Kontrat deploy ve yönetimi
   - Güvenli kontrat yürütme
   - Kontrat sahipliği kontrolü

## Katkıda Bulunma

1. Bu repoyu fork edin
2. Feature branch'i oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add some amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için `LICENSE` dosyasına bakın.

## İletişim

GitHub: [SolidityDevSK](https://github.com/SolidityDevSK) 