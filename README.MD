# 🌐 SocialNet — Fullstack Social Network & Messenger

**SocialNet** — это амбициозный проект современной социальной сети с интегрированным real-time мессенджером. Приложение объединяет в себе лучшие черты популярных платформ: алгоритмическую ленту, систему подписок и мгновенный обмен сообщениями.

---

## 🚀 Технологический стек

### Backend (The Engine)
- **Язык:** Go (Golang) — высокая производительность и эффективная многопоточность
- **Фреймворк:** Gin Gonic — быстрый HTTP-роутинг
- **База данных:** PostgreSQL с партиционированием (Range по времени для постов, Hash по ID чата для сообщений)
- **Кэширование и сессии:** Redis — статусы «онлайн» и кеширование сессий
- **Real-time:** WebSockets — мгновенная доставка сообщений
- **Инфраструктура:** Docker Compose, VPS, GitHub Actions (CI/CD)

### Frontend (The Interface)
| Платформа | Стек | State Management |
|-----------|------|-----------------|
| **Web** | React + Vite → TypeScript → Tailwind CSS v4 | Zustand |
| **Mobile** | Flutter (Dart) → BLoC Pattern → Dio | flutter_bloc |

---

## 🛠 Основной функционал

### 👥 Социальная сеть
- **Авторизация:** Регистрация и вход через JWT (JSON Web Tokens)
- **Лента постов:**
    - Создание постов с поддержкой медиа
    - Глобальная лента и персональная лента (только подписки)
    - Бесконечный скролл (Infinite Scroll)
- **Взаимодействия:** Лайки и многоуровневые комментарии
- **Подписки:** Модель «Подписчик → Подписка» для персонализации контента

### 💬 Мессенджер
- **Приватные чаты:** Создание и поиск комнат между пользователями
- **Real-time доставка:** Мгновенный обмен через WebSocket-хаб
- **История переписки:** Загрузка из PostgreSQL при открытии чата
- **Персистентность:** Все сообщения сохраняются в базе данных

---

## 📦 Установка и запуск

### 1. Backend

```bash
# Клонируйте репозиторий
git clone https://github.com/yourname/social-network.git
cd social-network/backend
```

#### Создайте файл .env в корне бэкенда:
```
PORT=8080
DB_URL=postgres://user:pass@localhost:5432/social_db
REDIS_URL=localhost:6379
JWT_SECRET=your_secret_key
```

```bash
# Запустите базу данных (PostgreSQL и Redis), затем сервер
go run cmd/api/main.go
```

### 2. Web Frontend

```bash
cd frontend-web
npm install
npm run dev
```

### 3. Mobile App

```bash
cd frontend_mobile
flutter pub get
flutter run
```

## 📐 Архитектура базы данных
Проект использует продвинутые техники оптимизации БД:
- **Партиционирование постов:** Таблица posts разделена по времени (RANGE) для молниеносной обработки архивов
- **Хеширование сообщений:** Таблица messages разделена по chat_id (HASH) для распределения нагрузки между физическими сегментами диска
- **Индексация:** Оптимизированные индексы на username, email и author_id для минимизации времени отклика

## 📈 Roadmap
- **Статусы пользователей («Печатает...», «В сети»)**
- **Групповые чаты и администрирование групп**
- **Push-уведомления (Firebase Cloud Messaging)**
- **Загрузка медиа через S3-совместимое хранилище**
- **Сквозное шифрование (E2EE) для секретных чатов**

## 📄 Лицензия
Distributed under the MIT License. See LICENSE for more information.

<div align="center">
Developed with ❤️ by <a href="https://github.com/1Yurko1">Yurko</a>
</div>
