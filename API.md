# Traktors Backend API

Base URL: `http://213.171.30.253:8080`

All responses use `Content-Type: application/json`. CORS разрешён для всех origins.

---

## Модель трактора

```json
{
  "id": "string",
  "name": "string",
  "images": ["string"],
  "brand": "string",
  "color": "string",
  "engineType": "diesel" | "gasoline" | "gas",
  "horsepower": 120,
  "year": 2020,
  "mileage": 1500,
  "vin": "string",
  "pts": "string",
  "ptsOwners": 1,
  "location": "string",
  "phone": "string",
  "ownerName": "string",
  "description": "string",
  "price": 2300000
}
```

Все поля кроме `id`, `name` и `images` — опциональные.

---

## Тракторы (MongoDB)

### GET /tractors

Возвращает все тракторы из базы данных.

**Response 200**
```json
[{ ...tractor }, ...]
```

---

### GET /tractors/{id}

**Response 200**
```json
{ ...tractor }
```

**Response 404**
```json
{ "error": "tractor not found" }
```

---

### POST /tractors

Создаёт трактор. Поле `id` опционально — если не передать, генерируется автоматически.

**Request**
```json
{
  "name": "John Deere 6120",
  "brand": "John Deere",
  "horsepower": 120,
  "price": 2300000
}
```

**Response 201**
```json
{ ...tractor }
```

**Response 400** — если `name` не передан или тело невалидно.

**Response 409** — если `id` уже существует.

---

### PUT /tractors/{id}

Полная замена документа. Поле `name` обязательно.

**Request**
```json
{ ...tractor }
```

**Response 200**
```json
{ ...tractor }
```

**Response 404** — если трактор не найден.

---

### DELETE /tractors/{id}

**Response 204** — успешно удалён.

**Response 404** — если трактор не найден.

---

## Статичные тракторы (feature flag)

### GET /get_tractors

Возвращает статичный список тракторов. Набор данных зависит от feature flag:
- flag = `false` (по умолчанию) → базовый датасет
- flag = `true` → альтернативный датасет

**Response 200**
```json
[{ ...tractor }, ...]
```

---

## Feature Flag

Feature flag хранится в MongoDB (коллекция `features`). Управляет поведением `GET /get_tractors`.

### GET /feature_check

Возвращает текущее значение флага.

**Response 200**
```json
true
```

---

### POST /feature_set

Устанавливает значение флага.

**Request**
```json
{ "value": true }
```

**Response 200**
```json
true
```

---

### POST /set_feature

Алиас для `POST /feature_set`.

**Request**
```json
{ "value": false }
```

**Response 200**
```json
false
```

---

## Медиа

### POST /media

Загружает изображение. Принимает `multipart/form-data` с полем `image`.

Поддерживаемые форматы: `jpeg`, `png`, `gif`, `webp`, `heic`.
Максимальный размер: **10 MB**.

**Response 201**
```json
{ "url": "http://213.171.30.253:8080/media/64f1a2b3c4d5e6f7a8b9c0d1.jpg" }
```

**Response 400** — неподдерживаемый формат или ошибка чтения файла.

---

### GET /media/{filename}

Отдаёт загруженный файл напрямую.

**Response 200** — файл в бинарном виде.
