# Exoplanets_backend

Учебный REST API на Go (gin + GORM) для сервиса каталога экзопланет. Данные хранятся в PostgreSQL, изображения и видео — в MinIO.

## Стек

| Компонент | Назначение |
|---|---|
| Go, gin-gonic/gin | HTTP-сервер и маршрутизация |
| gorm.io/gorm + driver/postgres | ORM и работа с PostgreSQL |
| minio-go/v7 | Загрузка файлов в MinIO |
| viper, godotenv | Чтение `config/config.toml` и `.env` |
| logrus | Логирование |
| Docker Compose | Запуск MinIO |

## Структура проекта

```
cmd/exoplanets/main.go     точка входа сервера
cmd/migrate/main.go        миграция схемы БД
config/config.toml         порт и адрес сервера
.env                       переменные окружения (БД, MinIO)
docker-compose.yml         сервис MinIO
internal/app/ds/           структуры данных (таблицы)
internal/app/dsn/          формирование строки подключения к БД
internal/app/config/       загрузка конфигурации
internal/app/repository/   работа с БД и MinIO
internal/app/handler/      HTTP-обработчики
internal/pkg/app.go        запуск приложения
```

## Запуск

1. Запустить MinIO:
   ```bash
   docker compose up -d minio
   ```
   Консоль MinIO: http://localhost:9001. Создать бакет `logo`.

2. Заполнить `.env`:
   ```
   DB_HOST=localhost
   DB_PORT=5433
   DB_NAME=mydb
   DB_USER=myuser
   DB_PASS=mypassword

   MINIO_ENDPOINT=localhost:9000
   MINIO_ACCESS_KEY=root
   MINIO_SECRET_KEY=rootpassword
   MINIO_BUCKET_NAME=logo
   ```
   Логин и пароль должны совпадать с `MINIO_ROOT_USER` и `MINIO_ROOT_PASSWORD` из `docker-compose.yml`.

3. Создать таблицы:
   ```bash
   go run ./cmd/migrate
   ```

4. Запустить сервер (по умолчанию http://localhost:8080):
   ```bash
   go run ./cmd/exoplanets
   ```

## Таблицы БД

Таблицы создаёт GORM (`AutoMigrate`) по структурам из `internal/app/ds`.

### Таблица `exoplanets`

| Поле | Тип в БД | Тип в Go | JSON | Описание |
|---|---|---|---|---|
| `id` | serial, PK | `int` | `id` | Идентификатор |
| `title` | varchar(30) | `string` | `title` | Название |
| `status` | varchar(30) | `string` | `status` | `draft`, `published` или `deleted` |
| `distance` | integer | `int` | `distance` | Расстояние |
| `weight` | numeric(6,2) | `float64` | `weight` | Масса |
| `description` | varchar(300) | `string` | `description` | Описание |
| `image` | varchar(100) | `string` | `image` | URL изображения в MinIO |
| `video` | varchar(100) | `string` | `video` | URL видео в MinIO |
| `created_at` | timestamp | `time.Time` | скрыто | Заполняется автоматически при создании |
| `creator_id` | integer, индекс | `uint` | скрыто | ID создателя (внешний ключ на `users.id`) |
| `formed_at` | timestamp | `time.Time` | скрыто | Обновляется автоматически при изменении |

### Таблица `users`

| Поле | Тип в БД | Тип в Go | JSON | Описание |
|---|---|---|---|---|
| `id` | serial, PK | `uint` | `id` | Идентификатор |
| `login` | varchar(40) | `string` | `login` | Логин |
| `password` | varchar(40) | `string` | скрыто | Пароль (хранится открытым текстом, см. замечания) |

### Таблица `likes`

| Поле | Тип в БД | Тип в Go | Описание |
|---|---|---|---|
| `id` | serial, PK | `uint` | Идентификатор |
| `user_id` | integer, NOT NULL, уникальный индекс `idx_render_server_unit_user` | `uint` | Кто поставил лайк |
| `exoplanet_id` | integer, NOT NULL, тот же уникальный индекс | `uint` | Какой экзопланете |

Уникальный индекс по паре `(user_id, exoplanet_id)` не даёт поставить два лайка одним пользователем на одну планету.

### Связи

- `exoplanets.creator_id` → `users.id` (многие к одному)
- `likes.user_id` → `users.id`
- `likes.exoplanet_id` → `exoplanets.id`

## Методы API

Базовый адрес: `http://localhost:8080`

Все пути написаны точно так, как в коде. Обратите внимание: у маршрута списка экзопланет в пути два слеша подряд (`/api//Exoplanets`).

| Метод | Путь | Тело запроса | Успех | Описание |
|---|---|---|---|---|
| GET | `/api//Exoplanets` | — | 200 | Список опубликованных экзопланет. Фильтр: `?query=<вес>` |
| GET | `/api/Exoplanets/:id` | — | 200 | Одна опубликованная экзопланета. Параметр `?next=true` делает редирект на следующую |
| GET | `/api/DraftExoplanets` | — | 200 | Черновик текущего пользователя (или `null`) |
| POST | `/api/AddExoplanets` | `multipart/form-data` | 201 | Создать черновик. Поля: `title`, `image`, `video` |
| PUT | `/api/PublishExoplanets` | `multipart/form-data` | 200 | Опубликовать черновик. Поля: `distance`, `weight`, `description` |
| DELETE | `/api/DeleteExoplanets/:id` | — | 200 | Мягкое удаление (статус `deleted`) |
| POST | `/api/LikeExoplanets/:id/like` | JSON `{"like": true}` | 201 / 200 | Поставить или снять лайк |
| POST | `/api/SignUp` | `form-data` / `x-www-form-urlencoded`: `login`, `password` | 201 | Регистрация |
| POST | `/api/SignIn` | — | 200 | Заглушка, возвращает сообщение «авторизация» |
| POST | `/api/SignOut` | — | 200 | Заглушка, возвращает сообщение «деавторизация» |

### Подробно по методам

#### GET `/api//Exoplanets`

- Без параметров возвращает все записи со статусом `published`.
- `?query=7.5` возвращает записи с весом `weight = 7.5`. Запятая тоже допускается (`7,5`).
- Для каждой записи добавляются поля `likes_number` (число лайков) и `is_mine` (создал ли её текущий пользователь).

Ответ:
```json
{
  "exoplanets": [
    {
      "id": 1,
      "title": "Суперземля",
      "status": "published",
      "distance": 120,
      "weight": 4.5,
      "description": "…",
      "image": "http://127.0.0.1:9000/logo/1_image.jpg",
      "video": "",
      "likes_number": 3,
      "is_mine": true
    }
  ],
  "query": ""
}
```

Ошибки: 400, если `query` не число.

#### GET `/api/Exoplanets/:id`

- Возвращает одну опубликованную запись с `likes_number` и `is_liked`.
- `?next=true` перенаправляет (302) на запись с большим `id`, у которой статус `published`.
- Значение `next` кроме `true` и пустого — ошибка 400.

#### GET `/api/DraftExoplanets`

Возвращает черновик текущего пользователя (статус `draft`). Если черновика нет, `exoplanet` равен `null`.

#### POST `/api/AddExoplanets`

Создаёт черновик со статусом `draft` и `creator_id` текущего пользователя.

| Поле | Тип | Обязательно | Допустимые значения |
|---|---|---|---|
| `title` | text | нет | строка до 30 символов |
| `image` | file | нет | `image/jpeg`, `image/png`, `image/gif`, `image/webp` |
| `video` | file | нет | `video/mp4`, `video/webm`, `video/quicktime`, `video/x-msvideo`, `video/x-matroska` |

Тип файла определяется по первым 512 байтам содержимого, а не по расширению. Файл сохраняется в MinIO как `<id>_image.<ext>` или `<id>_video.<ext>`, а в поле `image` или `video` записывается URL вида `http://127.0.0.1:9000/logo/<имя_файла>`.

Ошибки:
- 400 `некорректный файл`: содержимое не совпадает с допустимым типом.
- 400 `request Content-Type isn't multipart/form-data`: запрос отправлен не как form-data.
- 500: ошибка БД или MinIO (например, нет бакета `logo`).

#### PUT `/api/PublishExoplanets`

Публикует черновик текущего пользователя.

| Поле | Тип | Описание |
|---|---|---|
| `distance` | text (целое) | Записывается в `distance`, если передано |
| `weight` | text (число) | Записывается в `weight`, запятая заменяется на точку |
| `description` | text | Записывается в `description` |

После успешного запроса статус становится `published`.

Ошибки: 404, если черновика нет; 400, если `distance` или `weight` не разбираются как числа.

#### DELETE `/api/DeleteExoplanets/:id`

Выполняет мягкое удаление: статус меняется на `deleted`, запись остаётся в таблице. Удалить может только создатель (иначе 403). Удаляются только опубликованные записи.

#### POST `/api/LikeExoplanets/:id/like`

Тело: `{"like": true}` — поставить лайк (201). `{"like": false}` — снять лайк (200). Повторный лайк не создаёт дубликат.

#### POST `/api/SignUp`

Создаёт пользователя. Возвращает `{"user": {"id": ..., "login": ...}}` со статусом 201.

#### POST `/api/SignIn`, `/api/SignOut`

Заглушки. Они только возвращают сообщение и не проверяют пароль и не выдают токен или cookie.

## Внутренние методы

### Repository (`internal/app/repository`)

| Метод | Файл | Что делает |
|---|---|---|
| `New(settings)` | repository.go | Подключается к PostgreSQL через GORM и создаёт клиент MinIO |
| `GetPublishedExoplanet(id)` | exoplanet.go | Одна запись со статусом `published` |
| `GetPublishedExoplanets()` | exoplanet.go | Все записи со статусом `published` |
| `GetExoplanetsByWeight(weight)` | exoplanet.go | Опубликованные записи с заданным весом |
| `GetDraftExoplanet(creatorID)` | exoplanet.go | Черновик пользователя или `nil` |
| `GetNextPublishedExoplanetTo(id)` | exoplanet.go | Следующая запись с большим `id` (без проверки статуса) |
| `AddDraftExoplanet(exoplanet)` | exoplanet.go | Создаёт запись |
| `UpdateExoplanet(exoplanet)` | exoplanet.go | Сохраняет запись целиком (`Save`) |
| `DeleteExoplanet(id)` | exoplanet.go | Ставит статус `deleted` |
| `AddOrReplaceExoplanetImage(id, header, ctx)` | exoplanet.go | Загружает картинку в MinIO и пишет URL в `image` |
| `AddOrReplaceExoplanetVideo(id, header, ctx)` | exoplanet.go | Загружает видео в MinIO и пишет URL в `video` |
| `GetLikesNumber(id)` | Likes.go | Число лайков |
| `IsLikedByUser(exoplanetID, userID)` | Likes.go | Поставил ли пользователь лайк |
| `AddLike(like)` | Likes.go | Создаёт лайк |
| `DeleteLike(userID, exoplanetID)` | Likes.go | Удаляет лайк |
| `AddUser(user)` | User.go | Создаёт пользователя |

### Handler (`internal/app/handler`)

| Функция | Файл | Что делает |
|---|---|---|
| `GetUserID()` | exoplanet.go | Возвращает ID текущего пользователя (сейчас всегда `1`) |
| `isImage(contentType)` | exoplanet.go | Проверяет, что тип — картинка |
| `isVideo(contentType)` | exoplanet.go | Проверяет, что тип — видео |
| `validateFileUpload(header, check)` | exoplanet.go | Читает первые 512 байт и проверяет тип |
| `errorHandler(ctx, code, err)` | handler.go | Возвращает JSON `{"status": "error", "description": ...}` |
| `RegisterHandler(router)` | handler.go | Регистрирует маршруты API |
| `RegisterStatic(router)` | handler.go | Подключает шаблоны и статику |

## Замечания и известные ограничения

- **Авторизация не реализована.** `GetUserID()` всегда возвращает `1`. `SignIn` и `SignOut` ничего не проверяют.
- **Пароли хранятся открытым текстом** в поле `users.password`.
- **Каждый POST `/api/AddExoplanets` создаёт новый черновик.** Проверки на единственность черновика нет.
- **В `GetExoplanet` ошибка загрузки записи игнорируется.** Для несуществующего `id` ответ будет 200 с пустой записью, а не 404.
- **`DeleteExoplanet`** при несуществующем или неопубликованном `id` возвращает 500 (`record not found`).
- **Ссылки на файлы жёстко заданы** как `http://127.0.0.1:9000/...`, поэтому при другом адресе MinIO их нужно менять в коде.
- **Структура `PublishExoplanetRequest`** в handler объявлена, но не используется: данные читаются через `FormValue`.
- **Поле `title`** ограничено 30 символами на уровне БД.
