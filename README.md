# Сервис сокращения ссылок

## Запуск сервера сервиса

```bash
$> ./cmd/shortener/shortener
```

### Флаги для запуска

```
Usage of cmd/shortener/shortener/main:
  -a, --address string             address of shortener service server (default: localhost:8080)
  -b, --basepath string            address of short link basepath (default: http://localhost:8080)
  -c, --config string              path to config file in json format
  -d, --database_dsn string        db connection string
  -s, --enable_https               use HTTPS connection
  -f, --file_storage_path string   path to file storage of URLs
  -e, --server_cert_path string    path to server certificate file
  -k, --server_key_path string     path to server key file
```

### Переменные окружения (повторяют ф-нал флагов)

```
SERVER_ADDRESS      // address of shortener service server 
BASE_URL            // address of short link basepath
CONFIG              // path to config file in json format
DATABASE_DSN        // db connection string
ENABLE_HTTPS        // use HTTPS connection
FILE_STORAGE_PATH   // default "/tmp/short-url-db.json"
SERVER_CERT_PATH    // path to server certificate file
SERVER_KEY_PATH     // path to server key file
```

### Конфиг из файла

Пример файла кофигурации в файле `config.sample.json`

```
{
    "server_address": "localhost:8080",
    "base_url": "http://localhost:8080",
    "file_storage_path": "",
    "database_dsn": "host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable",
    "enable_https": false,
    "server_key_path": "",
    "server_cert_path": ""
}
```

## Запуск базы сервиса в контейнере

```bash
$> docker compose up
```

База будет запущена вот тут

```
host=0.0.0.0
port=5433
user=postgres
password=postgres
dbname=short
sslmode=disable
```

## HTTPS

Для запуска в режиме `HTTPS` необходимо получить сертификат и ключ, либо сгенерировать самоподписанные:

### Generate private key (.key)

```bash
# Key considerations for algorithm "RSA" ≥ 2048-bit
$> openssl genrsa -out server.key 2048

# OR

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
$> openssl ecparam -genkey -name secp384r1 -out server.key
```

### Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

```bash
$> openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```