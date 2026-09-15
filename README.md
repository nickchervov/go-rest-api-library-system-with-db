# Library REST API

Учебный REST API для библиотечной системы на Go: авторы, книги, читатели и выдача книг.

Проект написан, чтобы разобраться, как устроено backend-приложение целиком — от HTTP-слоя
до SQL и миграций, — а не на уровне отдельных упражнений. Дизайн слоёв и обработка ошибок
продумывались осознанно, поэтому код разделён по ответственности и не смешивает SQL с бизнес-логикой.

> Проект учебный: я начинающий Go-разработчик. Что уже сделано, а что осознанно отложено —
> честно описано в разделах «Что реализовано» и «Что пока не сделано».

## Стек

| Инструмент | Зачем выбран |
|---|---|
| `net/http` + `chi/v5` | стандартная библиотека + лёгкий роутер с поддержкой URL-параметров, без «магии» фреймворка |
| `modernc.org/sqlite` | чистый Go-драйвер SQLite — работает без CGO, сборка не требует компилятора C |
| `golang-migrate/migrate/v4` | версионирование схемы БД: миграции применяются автоматически при старте |
| `go-playground/validator/v10` | декларативная валидация входных данных через теги структур |
| `air` | live-reload при разработке |

## Архитектура

Слои общаются строго по направлению вниз, каждый знает только о своём уровне:

```
cmd/server/main.go   — точка входа: сборка зависимостей, запуск сервера, graceful shutdown
        │
internal/connectors  — HTTP-слой: роутинг и хендлеры, парсинг запроса, коды ответов
        │
internal/service     — бизнес-логика, валидация, зависимости от БД через интерфейс Storage
        │
internal/repository  — SQL-слой: запросы к SQLite, миграции, транзакции
        │
      SQLite
```

`internal/domain` — общие модели и доменные ошибки, `internal/dto` — структуры входных/выходных
данных, `pkg/` — переиспользуемые хелперы (`render`, `validator`).

```
cmd/server/main.go              запуск, DI, graceful shutdown
internal/
  connectors/
    routes.go                   регистрация маршрутов, middleware
    handlers.go                 хендлеры
  service/service.go            бизнес-логика + интерфейс Storage
  repository/store.go           SQL-запросы к SQLite
  domain/
    models.go                   Author, Book, Reader, Borrowing
    errors.go                   типизированные ошибки с HTTP-кодом
  dto/                          Input/Output-структуры
pkg/
  render/render.go              JSON-ответ
  validator/validator.go        кастомные правила валидации
  migrations/                   SQL-миграции
```

### Как я работаю с ошибками

Доменные ошибки вынесены в отдельный тип, который сам знает свой HTTP-код:

```go
type LibraryError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrBookAlreadyOccupied = &LibraryError{Code: 400, Message: "book is already occupied"}
```

Хендлер проверяет их через `errors.As` и отдаёт клиенту корректный статус, а всё остальное
считает внутренней ошибкой:

```go
if err := h.svc.CreateBook(input); err != nil {
	var targetErr *domain.LibraryError
	if errors.As(err, &targetErr) {
		render.JSONResponse(w, targetErr.Code, targetErr)
		return
	}
	log.Println(err)
	render.JSONResponse(w, http.StatusInternalServerError, map[string]string{"message": "create book internal error"})
	return
}
```

Ошибки оборачиваются через `%w`, чтобы сохранить цепочку для `errors.Is`/`errors.As`.

## Запуск

Требуется Go 1.25+.

```bash
git clone <repo-url>
cd go-rest-api-library-system-with-db

# обычный запуск
go run ./cmd/server

# или с live-reload
air
```

Сервер поднимется на `http://localhost:8080`.

При первом старте рядом с проектом создастся файл `library.db`, а миграции из `pkg/migrations`
применятся автоматически. **Важно:** запускать нужно из корня проекта — путь к миграциям
относительный (`file://pkg/migrations`).

## API

### Авторы

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/authors` | список авторов |
| `GET` | `/api/authors/{id}` | автор по id |
| `POST` | `/api/authors` | создать автора |
| `PUT` | `/api/authors/{id}` | обновить автора |
| `DELETE` | `/api/authors/{id}` | удалить автора |

### Книги

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/books` | список книг с вложенным автором |
| `GET` | `/api/books/{id}` | книга по id |
| `POST` | `/api/books` | создать книгу |
| `PUT` | `/api/books/{id}` | обновить книгу |
| `DELETE` | `/api/books/{id}` | удалить книгу |

### Читатели

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/readers` | список читателей |
| `POST` | `/api/readers` | создать читателя |
| `DELETE` | `/api/readers/{id}` | удалить читателя |

### Выдача книг

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/borrowings/active` | все активные выдачи |
| `POST` | `/api/borrowings` | выдать книгу читателю |
| `PUT` | `/api/borrowings/{id}/return` | вернуть книгу |

### Примеры

```bash
# создать автора
curl -X POST localhost:8080/api/authors \
  -H 'Content-Type: application/json' \
  -d '{"name":"George Orwell","country":"United Kingdom"}'
# => 201 {"id":1}

# создать книгу
curl -X POST localhost:8080/api/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"1984","isbn":"978-0451524935","year":1949,"authors_id":1}'
# => 201 {"id":1}

# создать читателя
curl -X POST localhost:8080/api/readers \
  -H 'Content-Type: application/json' \
  -d '{"name":"John Doe","email":"john@example.com","phone":"+1 555 0100"}'
# => 201 {"id":1}

# выдать книгу (дата выдачи подставляется сервером)
curl -X POST localhost:8080/api/borrowings \
  -H 'Content-Type: application/json' \
  -d '{"book_id":1,"readers_id":1}'
# => 201 {"id":1}

# вернуть книгу
curl -X PUT localhost:8080/api/borrowings/1/return \
  -H 'Content-Type: application/json' \
  -d '{"return_date":"2026-01-20"}'
# => 204

# получить список книг — автор вложен в ответ
curl localhost:8080/api/books
# => 200 {"books":[{"id":1,"title":"1984","isbn":"978-0451524935","year":1949,
#                  "authors":{"id":1,"name":"George Orwell","country":"United Kingdom"}}]}
```

## Что реализовано

**Слоистая архитектура с инверсией зависимости.** `service` не знает о SQLite — он работает
с интерфейсом `Storage`. Значит, репозиторий можно заменить или подменить в тестах, не трогая
бизнес-логику.

**Миграции.** Схема БД описана в `pkg/migrations` и версионируется. Миграции применяются
автоматически при старте приложения — ручной `CREATE TABLE` не нужен.

**Транзакции.** Выдача книги — критичная операция: нужно проверить, что книга и читатель
существуют и что книга ещё не выдана, а затем вставить запись. Это делается внутри транзакции
с `defer tx.Rollback()`, чтобы данные не разъехались:

```go
tx, err := s.db.Begin()
if err != nil {
	return 0, fmt.Errorf("begin tx issues: %w", err)
}
defer tx.Rollback()
// ... проверка занятости книги, INSERT в borrowing ...
if err := tx.Commit(); err != nil {
	return 0, fmt.Errorf("commit issues: %w", err)
}
```

**Бизнес-правила целостности.** Нельзя удалить автора, у которого есть книги; книгу или
читателя с историей выдач; нельзя выдать уже выданную книгу. Проверки живут в SQL-слое
через `EXISTS`-подзапросы.

**Валидация входных данных.** Помимо стандартных правил валидатора (`required`, `email`,
`datetime=2006-01-02`) написаны два своих: `alphaspace` (только буквы и пробелы — для имён
и стран) и `numericspecs` (цифры и пунктуация — для ISBN), через `RegisterValidation`
и регулярные выражения с Unicode-классами `\p{L}`.

**Graceful shutdown.** Сервер слушает сигналы `SIGINT`/`SIGTERM` через `signal.NotifyContext`,
после чего даёт активным запросам 5 секунд на завершение и только потом закрывает соединение с БД:

```go
ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
defer cancel()
// ...
ctxShutdown, stop := context.WithTimeout(context.Background(), 5*time.Second)
defer stop()
if err := server.Shutdown(ctxShutdown); err != nil { ... }
```

**Обработка ошибок без паник.** Каждая ошибка возвращается наверх, оборачивается через `%w`
и логируется в хендлере. Паник в коде нет.

## Что пока не сделано

Перечисляю осознанно — это не забытые места, а следующий фронт работы:

- **Тестов нет.** Это главный пробел проекта: бизнес-логику и репозиторий нужно покрыть тестами
  (`httptest` + `sqlmock` или отдельный файл БД).
- **Нет аутентификации и авторизации** — API полностью открыт.
- **Нет пагинации и фильтрации** в списках: `GET /api/books` отдаёт всё сразу.
- **Логирование** через стандартный `log`, без структурных логов и request-id.
- **Нет Docker и CI** — приложение запускается только локально.
- **Конфигурация зашита в код** (порт `:8080`, путь к БД `library.db`) вместо переменных окружения.
- **Путь к миграциям относительный** — свободно решается через `embed.FS`.

## Планы

Ближайшие шаги в порядке приоритета:

1. Покрыть `service` и `repository` тестами — без них проект сложно развивать.
2. Вынести конфигурацию в переменные окружения.
3. Добавить пагинацию в списковые эндпоинты.
4. Встроить миграции в бинарник через `embed`.
5. Структурное логирование (`log/slog`) и middleware с request-id.
