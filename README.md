# Marketplace API

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Docker](https://img.shields.io/badge/Docker-20.10+-2496ED?style=for-the-badge&logo=docker)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)
![Kafka](https://img.shields.io/badge/Apache%20Kafka-3.2+-231F20?style=for-the-badge&logo=apachekafka)
![Fiber](https://img.shields.io/badge/Fiber-v2-000000?style=for-the-badge&logo=go)

Бэкенд-сервис для простого маркетплейса, построенный на Go с использованием фреймворка Fiber. Проект включает в себя аутентификацию пользователей по JWT, ролевую модель (покупатель/продавец), полный CRUD для товаров и изображений, а также интеграцию с Kafka для асинхронной обработки событий.

## ✨ Возможности

- **Аутентификация**: Регистрация и вход пользователей с использованием JWT.
- **Ролевая модель**: Разделение пользователей на `покупателей` (buyer) и `продавцов` (seller).
- **Управление товарами**: Полный CRUD (Create, Read, Update, Delete) для товаров. Только продавцы могут управлять своими товарами.
- **Управление изображениями**: Загрузка, скачивание, удаление и привязка изображений к товарам.
- **Публичный API**: Возможность просматривать товары и их изображения без аутентификации.
- **Конфигурация**: Гибкая настройка через YAML-файл и переменные окружения.
- **Асинхронность**: Интеграция с кластером Apache Kafka.
- **Контейнеризация**: Полная настройка для запуска всего стека через Docker Compose.

## 🛠️ Технологический стек

- **Язык**: Go
- **Веб-фреймворк**: [Fiber v2](https://gofiber.io/)
- **База данных**: PostgreSQL 16
- **Драйвер БД**: [pgx/v5](https://github.com/jackc/pgx)
- **Конфигурация**: [Viper](https://github.com/spf13/viper)
- **Логирование**: [Zap](https://github.com/uber-go/zap)
- **Брокер сообщений**: Apache Kafka
- **Контейнеризация**: Docker, Docker Compose

## 🚀 Начало работы

### Предварительные требования

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- `git`

### Установка и запуск

1.  **Клонируйте репозиторий:**
    ```sh
    git clone https://github.com/BOBAvov/Marketplace.git
    cd Marketplace
    ```

2.  **Настройте окружение:**
    Создайте файл `.env` в корне проекта для переопределения секретных значений.
    ```env
    # .env
    AUTH_JWTSECRET="your_super_secret_key_for_development"
    ```

3.  **(Опционально) Раскомментируйте сервисы `db` и `api` в `docker-compose.yml`**, если они закомментированы.

4.  **Запустите все сервисы с помощью Docker Compose:**
    Эта команда соберёт образ API, поднимет контейнеры с PostgreSQL и кластером Kafka. Миграции базы данных применяются автоматически.
    ```sh
    docker-compose up -d --build
    ```

5.  **Сервис готов к работе!**
    - API доступен по адресу: `http://localhost:8080`
    - Kafka UI (для просмотра топиков) доступен по адресу: `http://localhost:9020`

---
## 📖 API Reference & cURL Примеры

Ниже — подборка cURL-примеров для каждого эндпоинта с типовыми телами запросов и примерами ответов. Подставьте свой токен и ID из ваших данных.

#### Подготовка переменных:
```bash
BASE=http://localhost:8080/api/v1
# Пример: получите токен через /auth/login и подставьте сюда
TOKEN="eyJhbGciOi..."
```

### Auth

#### 1) `POST /auth/register`
- **Описание**: регистрация пользователя как `buyer` или `seller`.
- **Запрос**:
```bash
curl -X POST "$BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seller@example.com",
    "password": "123456",
    "role": "seller"
  }'
```
- **Успешный ответ `201`**:
```json
{
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
  "role": "seller"
}
```

#### 2) `POST /auth/login`
- **Описание**: вход и получение JWT.
- **Запрос**:
```bash
curl -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "seller@example.com",
    "password": "123456"
  }'
```
- **Успешный ответ `200`**:
```json
{
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
  "role": "seller"
}
```

### Products (публичные эндпоинты)

#### 3) `GET /products?q=&limit=&offset=`
- **Описание**: список товаров с поиском по имени (`LIKE`) и пагинацией.
- **Запрос**:
```bash
curl "$BASE/products?q=phone&limit=2&offset=0"
```
- **Успешный ответ `200`**:
```json
[
  {
    "id": 2, "seller_id": 1, "name": "Phone X", "description": "Nice",
    "price_cents": 99900, "stock": 5, "cover_picture_id": 10,
    "created_at": "2025-01-01T12:00:00Z", "updated_at": "2025-01-01T12:00:00Z"
  },
  {
    "id": 1, "seller_id": 1, "name": "Phone Mini", "description": "",
    "price_cents": 59900, "stock": 10, "cover_picture_id": null,
    "created_at": "2025-01-01T11:00:00Z", "updated_at": "2025-01-01T11:00:00Z"
  }
]
```

#### 4) `GET /products/:id`
- **Описание**: получить один товар по ID.
- **Запрос**: `curl "$BASE/products/2"`
- **Успешный ответ `200`**:
```json
{
  "id": 2, "seller_id": 1, "name": "Phone X", "description": "Nice",
  "price_cents": 99900, "stock": 5, "cover_picture_id": 10,
  "created_at": "2025-01-01T12:00:00Z", "updated_at": "2025-01-01T12:00:00Z"
}
```

#### 5) `GET /products/:id/pictures`
- **Описание**: список картинок товара (с позициями).
- **Запрос**: `curl "$BASE/products/2/pictures"`
- **Успешный ответ `200`**:
```json
[
  {
    "id": 10, "mime_type": "image/jpeg", "size_bytes": 34567,
    "created_at": "2025-01-01T12:00:05Z", "position": 1
  },
  {
    "id": 11, "mime_type": "image/png", "size_bytes": 28765,
    "created_at": "2025-01-01T12:01:10Z", "position": 2
  }
]
```

#### 6) `GET /pictures/:id`
- **Описание**: скачать бинарный файл картинки.
- **Запрос для сохранения в файл**: `curl -L "$BASE/pictures/10" -o out.jpg`
- **Успешный ответ `200`**: Тело ответа — бинарные данные с заголовком `Content-Type: image/jpeg`.

### Products (только для `seller` с Bearer JWT)

#### 7) `POST /products`
- **Описание**: создать товар.
- **Запрос**:
```bash
curl -X POST "$BASE/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Phone X",
    "description": "Nice",
    "price_cents": 99900,
    "stock": 5
  }'
```
- **Успешный ответ `201`**:
```json
{
  "id": 2, "seller_id": 1, "name": "Phone X", "description": "Nice",
  "price_cents": 99900, "stock": 5, "cover_picture_id": null,
  "created_at": "2025-01-01T12:00:00Z", "updated_at": "2025-01-01T12:00:00Z"
}
```

#### 8) `PUT /products/:id`
- **Описание**: обновить товар.
- **Запрос**:
```bash
curl -X PUT "$BASE/products/2" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Phone X Pro", "description": "Better",
    "price_cents": 109900, "stock": 3, "cover_picture_id": 10
  }'
```
- **Успешный ответ `200`**:
```json
{
  "id": 2, "seller_id": 1, "name": "Phone X Pro", "description": "Better",
  "price_cents": 109900, "stock": 3, "cover_picture_id": 10,
  "created_at": "2025-01-01T12:00:00Z", "updated_at": "2025-01-01T12:05:00Z"
}
```

#### 9) `DELETE /products/:id`
- **Описание**: удалить товар.
- **Запрос**: `curl -X DELETE "$BASE/products/2" -H "Authorization: Bearer $TOKEN"`
- **Успешный ответ**: `204 No Content`

### Pictures (только для `seller`)

#### 10) `POST /products/:id/pictures` (multipart)
- **Описание**: загрузить картинку и привязать к товару.
- **Запрос**: `curl -X POST "$BASE/products/2/pictures" -H "Authorization: Bearer $TOKEN" -F "file=@./photo.jpg"`
- **Успешный ответ `201`**:
```json
{
  "id": 10, "mime_type": "image/jpeg", "size_bytes": 34567,
  "position": 1, "created_at": "2025-01-01T12:00:05Z"
}
```

#### 11) `DELETE /products/:id/pictures/:pid?hard=1`
- **Описание**: отвязать картинку от товара. При `?hard=1` картинка удаляется из хранилища.
- **Запрос (отвязать и удалить)**: `curl -X DELETE "$BASE/products/2/pictures/10?hard=1" -H "Authorization: Bearer $TOKEN"`
- **Успешный ответ**: `204 No Content`

#### 12) `PUT /products/:id/cover/:pid`
- **Описание**: назначить картинку обложкой товара.
- **Запрос**: `curl -X PUT "$BASE/products/2/cover/10" -H "Authorization: Bearer $TOKEN"`
- **Успешный ответ**: `204 No Content`

### Шаблоны ошибок

Сервис возвращает ошибки в формате JSON `{"error":"<сообщение>"}`.

| Код | Описание                | Примеры сообщений                                                                                                       |
|:----|:------------------------|:------------------------------------------------------------------------------------------------------------------------|
| 400 | **Bad Request**         | `invalid json`, `email and password required`, `invalid role`, `invalid id`, `product not found`, `missing file`      |
| 401 | **Unauthorized**        | `missing bearer token`, `invalid token`, `invalid credentials`                                                          |
| 403 | **Forbidden**           | `seller role required`, `forbidden: not owner`                                                                          |
| 404 | **Not Found**           | `product not found`, `<текст ошибки БД>`                                                                                |
| 500 | **Internal Server Error** | `internal server error`                                                                                                 |
