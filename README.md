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

### Geliştirme Aşamasındaki Özellikler
- 🔄 P2P ağ desteği
- 🔄 Akıllı kontrat desteği
- 🔄 Validator oylama sistemi
- 🔄 Web arayüzü

## Kurulum

### Gereksinimler
- Go 1.24 veya üzeri
- Gin web framework

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
go run cmd/confirmix/main.go
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

## Proje Yapısı

```
confirmix/
├── cmd/                    # Uygulamanın giriş noktaları
│   └── confirmix/         # Ana uygulama
├── pkg/                    # Dışa açık paketler
│   ├── api/               # HTTP API implementasyonu
│   ├── blockchain/        # Blockchain çekirdek yapısı
│   ├── consensus/         # Konsensüs mekanizmaları
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

3. **HTTP API**
   - RESTful API ile blockchain yönetimi
   - Blok ve validator işlemleri
   - İşlem gönderme ve sorgulama

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