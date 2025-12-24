# 🔗 DropFiles API (Сервис обмена файлами на Go)

Простой, но полнофункциональный сервис по типу **Dropmefiles.com**, разработанный на **Go** с использованием фреймворка **Gin** и базы данных **PostgreSQL**.

## ✨ Особенности

* **Загрузка без регистрации** — получаешь ссылку на скачивание мгновенно
* **Случайная генерация UUID** — уникальные коды для файлов
* **TTL (Time-to-live)** — автоудаление файлов (макс 24ч)
* **Swagger UI** — полная документация API
* **Архитектура by feature** — впервые пробовал такой подход!
* **Автомиграции** — БД обновляется при каждом запуске (для простоты разработки)
* **File storage** — файлы сохраняются локально в `UPLOAD_DIR`

## 🛠️ Технологии

* **Язык:** Go 1.25+
* **Фреймворк:** Gin + Swagger
* **База данных:** PostgreSQL (pgx/v5)
* **Конфигурация:** `.env` + godotenv
* **Docker:** Полная поддержка compose

-----

## 🚀 Начало работы

### 1. Требования

* Docker + Docker Compose
* Go (только для генерации Swagger)

### 2. Конфигурация

1. Создай файл `.env`, описание ниже

### 3. Запуск

1. **Установи зависимости:**
```bash
go mod tidy
```

2. **Сгенерируй Swagger docs:**
```bash
swag init
```

3. **Запусти Docker:**
```bash
docker-compose up --build
```

**Готово!** Открой: `http://localhost:8080/swagger/index.html`

## 📋 Конфигурация (.env)

```env
# Сервер
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# База данных (PostgreSQL)
DB_HOST=db          # Docker / localhost (нативно)
DB_PORT=5432
DB_USER=dropfiles
DB_PASSWORD=dropfiles
DB_NAME=dropfiles
DB_SSL_MODE=disable

# Файлы
MAX_FILE_SIZE_MB=10
MAX_TTL_HOURS=24
DEFAULT_TTL_HOURS=24
MIN_FILENAME_LEN=1
MAX_FILENAME_LEN=255

# Хранилище
UPLOAD_DIR=/app/uploads
```

-----

## 📋 API Эндпоинты

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/ping` | Проверка работоспособности |
| `POST` | `/api/v1/upload` | Загрузка файла |
| `GET` | `/api/v1/f/{id}` | Скачать файл |

### **Загрузка**
```
POST /api/v1/upload
multipart/form-data:
- file (обязательно)
- ttl? (часы, макс 24)
```
**Ответ:** `{"url": "/f/abc123"}`

### **Скачивание**
```
GET /api/v1/f/{id}
```
**Ответ:** Файл (`application/octet-stream`)

-----

### **Улучшение: Удаление файлов из storage**
В `file/service.go:67` можно добавить физическое удаление файла:
```go
// При истечении TTL или ошибке
os.Remove(filepath.Join(h.config.UploadDir, filePath))
```

-----

## 💡 Будущие улучшения

**Текущая реализация** простая и надёжная, но на будущее:

1. **Redis кэш** — для быстрого поиска файлов по ID
2. **S3/MinIO** — облачное хранилище вместо локальных файлов
3. **Rate limiting** — защита от DDoS
4. **Файрвол по MIME** — запрет исполняемых файлов

-----

## 🚨 Дисклеймер: начинающий разработчик

**Postman** — работает идеально на любых файлах (маленьких/больших).

**Swagger UI** — маленькие файлы ок, но большие файлы (>MAX_FILE_SIZE_MB) дают:

```text
Failed to fetch.
Possible Reasons:
CORS
Network Failure
URL scheme must be "http" or "https" for CORS request.
```

>Я пытался решить, не получилось, каюсь
