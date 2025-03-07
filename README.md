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

### Geliştirme Aşamasındaki Özellikler
- 🔄 Validator sıralama sistemi (Round-Robin)
- 🔄 HTTP API desteği
- 🔄 P2P ağ desteği
- 🔄 Akıllı kontrat desteği
- 🔄 Validator oylama sistemi
- 🔄 Web arayüzü

## Kurulum

### Gereksinimler
- Go 1.24 veya üzeri

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
go run .
```

## Proje Yapısı

- `authority.go`: Validator yapısı ve imzalama işlemleri
- `block.go`: Blok yapısı ve ilgili metodlar
- `blockchain.go`: Blockchain yapısı ve temel işlemler
- `main.go`: Örnek kullanım ve test kodu

## Nasıl Çalışır?

1. **Validator Sistemi**
   - Her validator için ECDSA public/private anahtar çifti oluşturulur
   - Validatorlar blokları kendi private key'leri ile imzalar
   - İmzalar diğer validatorlar tarafından doğrulanır

2. **Blok Yapısı**
   - Timestamp
   - İşlem verisi
   - Önceki blok hash'i
   - Mevcut blok hash'i
   - Validator imzası
   - Validator adresi

3. **Konsensüs Mekanizması**
   - Proof of Authority kullanılır
   - Sadece yetkili validatorlar blok oluşturabilir
   - Her blok, oluşturan validator tarafından imzalanır
   - Blok zinciri sürekli doğrulanır

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