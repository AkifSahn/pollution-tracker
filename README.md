# Gerçek Zamanlı Hava Kirliliği İzleme Platformu

Bu proje, dünya genelinde hava kirliliği verilerini toplayan, analiz eden ve görselleştiren gerçek zamanlı bir web platformudur

---

## Projenin Amacı ve Kapsamı

Bu proje, dünya genelindeki hava kirliliği seviyelerinin anlık olarak izlenebilmesini amaçlar. Çeşitli bölgelerden gelen kirlilik verileri API aracılığıyla sisteme iletir. Sistem belirli kriterlere göre bu verileri işler, anomalileri otomatik olarak tespit eder ve gerçek zamanlı olarak kullanıcıları bilgilendirir.

Projenin temel hedefleri:

- Farklı konumlardan gelen hava kirliliği verilerinin toplanması.
- İstatistiksel yöntemlerle anomali tespiti yapılması.
- Anormal durumlarda kullanıcıların `WebSocket` üzerinden anlık olarak bilgilendirilmesi.
- Harita ve grafikler üzerinden bölgesel analizlerin yapılabilmesi.

---

## Teknoloji

### Go (Fiber)
- Kolay kullanımı ve yüksek performansı ile gerçek zamanlı hava kirliliği izleme sistemi için ideal bir dildir.

### TimescaleDB (PostgreSQL)
- Hava kirliliği verileri zaman serisi verileri olduğundan `TimescaleDB`, bu tür veriler için `PostgreSQL` üzerinde geliştirilmiş güçlü bir veritabanıdır.
- Zaman bazlı sorgular, histogramlar ve bucket işlemleri gibi analizler kolay bir şekilde yapılabilir.

### RabbitMQ
- Verilerin anlık olarak kuyruğa alınıp farklı bileşenler tarafından asenkron şekilde işlenmesini sağlar.
  
### Swagger (OpenAPI)
- API’lerin kullanıcı ve geliştiriciler tarafından anlaşılabilir olması için otomatik dokümantasyon sağlar.
- `Swagger` kullanıcı arayüzünden endpointler test edilebilir. 

---

## Sistem Mimarisi ve Komponentleri

Sistem monolitik bir mimari ile geliştirilmiştir. Bileşenler arasında veri akışı `RabbitMQ` aracılığıyla sağlanır. Zaman serisi verileri TimescaleDB üzerinde saklanırken, anomaliler anlık olarak WebSocket ile yayınlanır.

### Temel Komponentler

#### 1. API Sunucusu (Go + Fiber)
- API aracılığıyla veri alma, gönderme işlemleri yapılır.
- `Swagger` üzerinden dökümantasyon sunar.

#### 2. Veri İşleme Servisi (ingest)
- RabbitMQ `ingest_queue` üzerinden gelen ham verileri dinler.
- Gelen verileri işler ve belli kriterlere(Z-score, yüzde artış) göre anomali olup olmadığını belirler.
- Sonuçları `TimescaleDB`’ye kaydeder.
- Anomali varsa `notification_queue` kuyruğuna bir bildirim gönderir.

#### 3. WebSocket Bildirim Servisi(Notification)
- Merkezi bir `Hub`, RabbitMQ `notification_queue` kuyruğunu dinler.
- Her kullanıcı için backendde bir `Client` yapısı oluşturulur.
- Her `Client` için bir Websocket bağlantısı oluşur.
- `Hub` kendisine register olan tüm `Client`lara oluşan anomali bildirimleri iletir.
- Bildirimleri alan `Client`lar Websocket bağlantısı ile frontend tarafına iletilir  

#### 4. Veritabanı
- `TimeScaleDB` veritabanı ile olan bağlantıyı sağlar.
  
#### 5. Kuyruklama Sistemi (RabbitMQ)
- `RabbitMQ` bağlantısını oluşturur.
- `RabbitMQ` kuyruklarını -hali hazırda yoksa- oluşturur


### Veri Akışı

1. API aracılığıyla bir ölçüm verisi sisteme gönderilir.
2. Veri `RabbitMQ` kuyruğuna alınır (`ingest_queue`).
3. Veri işlenir, anomali tespiti yapılır ve veritabanına kaydedilir.
4. Eğer anomali varsa, sistem `RabbitMQ` bildirim kuyruğuna(`notification_queue`) bir mesaj gönderir.
5. Bildirim **WebSocket Bildirim Servisi** ile kullanıcılara anlık olarak iletilir.

---

## Kurulum

### 1. GitHub Deposunu Klonlayın

```bash
git clone https://github.com/AkifSahn/pollution-tracker.git
cd pollution-tracker
```

### 2. Ortam Değişkenlerini Tanımlayın

`pollution-tracker/backend` dizininde `.env` adında bir dosya oluşturun ve aşağıdaki örneğe uygun şekilde doldurun:

```env
DB_USER=timescale
DB_PASSWORD=root1234
DB_NAME=pollution
DB_HOST=127.0.0.1
DB_PORT=5432

SERVER_HOST=127.0.0.1
SERVER_PORT=3000

AMQP_USER=guest
AMQP_PASSWORD=guest
AMQP_HOST=127.0.0.1
AMQP_PORT=5672
```

> Not: Docker Compose içerisindeki servisler, `DB_HOST` ve `AMQP_HOST` değerlerini `db` ve `rabbitmq` olarak otomatik değiştirecektir.

### 3. Docker Compose ile Uygulamayı Başlatın

Aşağıdaki komutu kök dizinde çalıştırarak sistemdeki tüm bileşenleri ayağa kaldırabilirsiniz:

```
docker-compose --env-file backend/.env up --build
```
> Not: Docker konteynerlarını arkaplanda çalıştırmak için `-d` flagı eklenebilir

### 4. Servislerin Durumunu Kontrol Edin

| Servis               | Adres                                           |
|----------------------|-------------------------------------------------|
| Frontend UI          | http://localhost:5173/                          |
| API Sunucusu         | http://localhost:3000                           |
| Swagger Dokümantasyonu | http://localhost:3000/swagger/index.html      |
| RabbitMQ UI          | http://localhost:15672 (guest / guest)          |

> `RabbitMQ` arayüzüne girerek kuyruğa gelen mesajları gözlemleyebilir ve sistemin düzgün çalışıp çalışmadığını kontrol edebilirsiniz.

### 5. Sistemi Durdurmak

Tüm ilgili docker servislerini durdurmak ve konteylerlarla ilişkili `volume`ları silmek için:

```
docker-compose --env-file backend/.env down -v
```

---

## Kullanım Rehberi

- Kullanıcı arayüzü aşağıdaki gibi gözükür, bu arayüzde görünen ana komponentler şunlardır:
  * Isı Haritası
  * Yoğunluk Grafiği
  * Navbar ve bildirim paneli
![Kullanıcı Arayüzü](screenshots/1.png?raw=true "Optional Title")

- Haritada 1 saat içinde oluşan anomali işaretçileri aşağıdaki gibi gözükür ve işaretçiye tıklayarak anomali hakkında detaylı bilgi görülebilir
![Marker](screenshots/2.png?raw=true "Optional Title")

- Harita komponentinin üstünde bulunan parametre filtreleri ile ısı haritası üzerinde gösterilen değerler filtrelenebilir.
- Seçilen filtreye ikinci kez tıklayarak filtre kaldırılabilir.
- Hiçbir filtre seçili olmadığı zaman bütün parametre verileri haritada gözükür.
- Haritanın aşağısında bulunan **Zaman Aralığı** girdileri ayarlanarak harita üzerindeki veriler belirli zaman aralığından seçilebillir.
- Varsayılan olarak 23 saat içerisindeki veriler gösterilir.
![Filtre](screenshots/3.png?raw=true "Optional Title")

- Sağ üst köşede bulunan **Bildirimler** butonu ile bildirim paneli açılabilir.
- Kullanıcıya iletilen bildirimler anlık olarak burada gözükür ve ısı haritasına işaretçi olarak da eklenir.
![Bildirim](screenshots/4.png?raw=true "Optional Title")

- Isı haritası üzerinde görünen dikdörtgen alanın boyu ve konumu iki köşesinde bulunan imleçler değiştirilebilir.
- Bu seçilen alanın altında kalan bölge için **Yoğunluk Grafiği** gösterilir.
- Alan değiştirildiğinde grafik de otomatik olarak güncellenir
![Alan Seçme](screenshots/5.png?raw=true "Optional Title")

- Grafik seçilen bölgedeki kirlilik yoğunluğuklarını gösterir.
- Grafik üzerindeki yoğunluk verileri ilk veriden başlanarak 5dk değerlerin ortalaması alınarak backend tarafından hesaplanır.
- Harita üzerindeki parametre filtreleri aynı zamanda grafiğe de uygulanır.
![Grafik](screenshots/6.png?raw=true "Optional Title")



---

## API Dökümantasyonu

- [POST `/api/pollutions`](#post-apipollutions)
- [GET `/api/pollution/density/rect`](#get-apipollutionsdensityrect)
- [GET `/api/pollutions/{latitude}/{longitude}`](#get-apipollutionslatitudelongitude)
- [GET `/api/anomalies`](#get-apianomalies)
- [GET `/api/pollutions`](#get-apipollutions)
- [GET `/api/pollutants`](#get-apipollutants)
- [GET `/ws`](#get-ws)

* ### Swagger Arayüzü

- Adres: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)

Swagger üzerinde aşağıdaki endpointleri görebilir ve deneyebilirsiniz:


* ### POST `/api/pollutions`

Yeni bir kirlilik verisi gönderir.

**Body (JSON):**

```json
{
  "latitude": 41.0,
  "longitude": 29.0,
  "pollutant": "PM10",
  "value": 85.2
}
```

* ### GET `/api/pollutions/density/rect`

Belirtilen dikdörtgen alanda belirli zaman aralığında ortalama kirlilik yoğunluklarını verir.

**Query Parametreleri:**
- `latFrom`, `latTo`
- `longFrom`, `longTo`
- `from`, `to` (ISO 8601)
- `pollutant`


* ### GET `/api/pollutions/{latitude}/{longitude}`

Verilen konum ve zaman aralığına göre tüm kirlilik ölçümlerini getirir.

**Path Parametreleri:**
- `latitude`
- `longitude`

**Query Parametreleri:**
- `from`
- `to`


* ### GET `/api/anomalies`

Belirtilen zaman aralığındaki tespit edilen anomalileri getirir.

**Query Parametreleri:**
- `from`
- `to`


* ### GET `/api/pollutions`

Zaman aralığına göre tüm ölçüm verilerini getirir (isteğe bağlı olarak kirleten parametresi filtresi uygulanabilir).
`pollutant` parametresi sağlanmadığı durumda, bütün parametreler için verileri getirir.

**Query Parametreleri:**
- `from`
- `to`
- `pollutant` (opsiyonel)


* ### GET `/api/pollutants`

Veritabanındaki tüm farklı kirleten parametreleri (PM2.5, NO2, SO2, vb.) listeler.


* ### GET `/ws`

WebSocket bağlantı noktasıdır. Anomali tespit edildikçe bağlı istemcilere anlık mesaj gönderilir.

---

## Scriptler

Sistem iki script ile test edilebilir:
  * `manual-input.sh`
  * `auto-test.sh`

### 1. manual-input.sh  

Bu script kullanıcının manuel olarak sisteme kirlilik değeri ekleyebilmesini sağlar

```
Usage: ./manual-input.sh <latitude> <longitude> <parameter> <value>
```

### 2. auto-test.sh

Bu script ile sisteme otomatik bir şekilde kirlilik değerleri eklenebilir
```
Usage: auto-test.sh [--anomaly-chance=<0-100>] [--duration=<seconds>] [--rate=<number>] [--lat-min=] [--lat-max=] [--long-min=] [--long-max=]
```
**Parametreler**:
  * `--anomaly-chance=<0-100>`: 0-100 arası verilen bir değer ile üretilen kirlilik değerinin anomali olma olasılığını belirtir. Varsayılan: 20  
  * `--duration=<seconds>`: Scriptin kaç saniye boyunca çalışacağını belirtir. Varsayılan: 30
  * `--rate=<number>`: Sisteme saniyede kaç tane kirlilik verisi gönderileceğini belirtir. Bu kirlilik verileri rastgele bir şekilde üretilir. Varsayılan 5
  * `--lat-min=`: Oluşturulan kirlilik verisinin alabileceği en küçük **enlem** değerini belirtir. Varsayılan -90
  * `--lat-max=`: Oluşturulan kirlilik verisinin alabileceği en büyük **enlem** değerini belirtir. Varsayılan 90
  * `--long-min=`: Oluşturulan kirlilik verisinin alabileceği en küçük **boylam** değerini belirtir. Varsayılan -180
  * `--long-max=`: Oluşturulan kirlilik verisinin alabileceği en büyük **boylam** değerini belirtir. Varsayılan 180

* Örnek kullanım:
  - Türkiye sınırlarında %20 şans ile anomali üretecek ve 2 dakika boyunca saniyede 2 kirlilik verisi üretecek kullanım örneği:
    ```
     ./auto-test.sh --anomaly-chance=20 --duration=120 --rate=2 \
                    --lat-min=36.0 --lat-max=42.1 \
                     --long-min=26.0 --long-max=45.0
    ```
    > Not: Bu örnek 2 dakika boyunca toplam 240 kirlilik verisi üretir. Ve kabaca 48 tane anomali değeri oluşturur.


## Sorun Giderme

TODO

