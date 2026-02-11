# Traktors Backend API

REST API для управления тракторами с MongoDB и загрузкой медиа.

## Возможности

- CRUD операции для тракторов
- Загрузка изображений (JPEG, PNG, GIF, WebP, HEIC)
- MongoDB для хранения данных
- Автоматический деплой через GitHub Actions

## API Endpoints

- `GET /tractors` - получить все тракторы
- `GET /tractors/{id}` - получить трактор по ID
- `POST /tractors` - создать трактор
- `PUT /tractors/{id}` - обновить трактор
- `DELETE /tractors/{id}` - удалить трактор
- `POST /media` - загрузить изображение
- `GET /media/{filename}` - получить изображение

## Локальная разработка

```bash
# Запустить MongoDB
# Убедитесь что MongoDB работает на localhost:27017

# Запустить сервер
go run .
```

Сервер будет доступен на `http://localhost:8080`

## Деплой на VM

### Настройка GitHub Secrets

Перед первым деплоем настройте секреты в GitHub:

1. Откройте https://github.com/zeppelin-tt/traktors_be/settings/secrets/actions
2. Добавьте следующие secrets:

- **SSH_PRIVATE_KEY** - ваш приватный SSH ключ для подключения к VM
- **VM_HOST** - `213.171.30.253`
- **VM_USER** - `user1`

### Автоматический деплой

При каждом push в `main` GitHub Actions автоматически:
1. Подключается к VM через SSH
2. Устанавливает Go и MongoDB (если не установлены)
3. Копирует код на сервер
4. Собирает приложение
5. Перезапускает сервис

### Ручной деплой

```bash
./deploy.sh 213.171.30.253 user1 ~/.ssh/id_rsa
```

## Технологии

- Go 1.22+
- MongoDB 7.0
- GitHub Actions для CI/CD
