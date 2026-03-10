# LearnFlow

**LearnFlow** — это платформа (маркетплейс) для автоматизации образовательного процесса. Система позволяет ученикам находить репетиторов по предметам, просматривать свободные слоты, бронировать занятия и общаться во встроенном чате в реальном времени.

Проект разрабатывается в рамках изучения дисциплин «Web-программирование» и «Конструирование программного обеспечения».

## 🛠 Технологический стек

- **Язык:** Go (Golang) 1.22+
- **Архитектура:** Clean Architecture / Layered Architecture
- **Роутинг:** `net/http` (Standard Library)
- **Конфигурация:** `viper` (чтение из `.env`)
- **Логирование:** `log/slog` (Standard Library)
- **Шаблонизатор:** `html/template` (SSR)
- **База данных:** PostgreSQL + `pgx` (в планах)

## 📁 Структура проекта (Standard Go Project Layout)

```text
learningflow/
├── cmd/
│   └── learnflow/
│       └── main.go         # Точка входа в приложение
├── internal/
│   ├── app/                # Инициализация приложения и DI (Dependency Injection)
│   ├── models/             # Доменные сущности (User, Tutor, Subject)
│   ├── repository/         # Слой работы с данными (БД)
│   ├── service/            # Бизнес-логика (Use cases)
│   └── transport/
│       ├── rest/           # Обработчики REST API (JSON)
│       └── ssr/            # Обработчики Server-Side Rendering (HTML)
├── pkg/
│   ├── config/             # Модуль загрузки конфигурации
│   └── logger/             # Обёртка над slog для работы с контекстом
├── web/
│   ├── static/             # Статические файлы (CSS, JS, изображения)
│   └── templates/          # HTML-шаблоны
├── .env                    # Переменные окружения
└── go.mod                  # Зависимости
```
🚀 Запуск проекта локально
Склонируйте репозиторий:
```Bash
git clone https://github.com/yourusername/learningflow.git
cd learningflow
```
Установите зависимости:
```Bash
go mod tidy
```
Настройте переменные окружения. Создайте файл .env в корне проекта:
```Env
APP_ENV=local
HTTP_PORT=8080
```
Запустите сервер:
```Bash
go run cmd/learnflow/main.go
```
Откройте в браузере:
http://localhost:8080