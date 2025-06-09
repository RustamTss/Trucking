# Business Schedule - Система управления траковым бизнесом

Полнофункциональная система для управления траковым бизнесом с возможностью отслеживания компаний, транспорта, кредитов и финансовых показателей.

## 🚀 Быстрый старт

### Требования
- **Docker** и **Docker Compose**
- **Go 1.21+** (для локальной разработки)
- **Node.js 18+** (для локальной разработки)

### 🐳 Запуск через Docker (Рекомендуется)

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/RustamTss/Trucking.git
   cd Trucking
   ```

2. **Создайте .env файл в корне проекта**
   ```bash
   cp .env.example .env
   ```
   
   Или создайте `.env` файл с содержимым:
   ```env
   MONGODB_URI=mongodb://mongodb:27017/business_schedule
   JWT_SECRET=your_very_secret_jwt_key_here_make_it_long_and_secure
   PORT=8080
   REACT_APP_API_URL=http://localhost:8080
   ```

3. **Запустите все сервисы**
   ```bash
   docker-compose up -d
   ```

4. **Проверьте статус контейнеров**
   ```bash
   docker-compose ps
   ```

Приложение будет доступно на:
- **Фронтенд**: http://localhost:3000
- **Бэкенд API**: http://localhost:8080
- **MongoDB**: localhost:27017

### 💻 Локальная разработка

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/RustamTss/Trucking.git
   cd Trucking
   ```

2. **Запустите MongoDB**
   ```bash
   docker-compose up -d mongodb
   ```

3. **Запустите бэкенд**
   ```bash
   cd backend
   go mod download
   go run main.go
   ```
   Бэкенд будет доступен на `http://localhost:8080`

4. **Запустите фронтенд**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   Фронтенд будет доступен на `http://localhost:5174`

### Первый вход

**Готовый пользователь:**
- Email: `admin@test.com`
- Пароль: `123456`

Или зарегистрируйте нового пользователя через интерфейс.

## 📋 Функциональность

### MVP возможности
- ✅ **Аутентификация** - регистрация и вход пользователей
- ✅ **Управление компаниями** - создание и управление несколькими компаниями
- ✅ **Управление транспортом** - трак и трейлер менеджмент
- ✅ **Кредитное управление** - займы с автоматическими расчетами
- ✅ **Платежная система** - отслеживание платежей и остатков
- ✅ **Финансовые отчеты**:
  - График долгов по компаниям
  - Амортизационный график
  - График амортизации активов

### Дополнительные возможности
- ✅ **Dashboard статистика** - общие показатели бизнеса
- ✅ **Multi-tenant архитектура** - изоляция данных пользователей
- ✅ **Автоматические финансовые расчеты**
- ✅ **Responsive UI** - адаптивный дизайн
- ✅ **TypeScript поддержка** - типизированный фронтенд
- ✅ **CRUD пользователей** - управление профилем и аккаунтом

## 🛠 Технологический стек

### Бэкенд
- **Go 1.21** с Fiber framework
- **MongoDB** для хранения данных
- **JWT** аутентификация
- **bcrypt** для хеширования паролей
- **Clean Architecture** структура проекта

### Фронтенд  
- **React 19** с TypeScript
- **Vite** для сборки
- **Tailwind CSS** для стилизации
- **Axios** для API запросов
- **React Router** для навигации
- **Recharts** для графиков

## 📡 API Endpoints

### Аутентификация
- `POST /api/auth/register` - Регистрация пользователя
- `POST /api/auth/login` - Вход в систему

### Компании
- `GET /api/companies` - Список компаний
- `POST /api/companies` - Создание компании
- `PUT /api/companies/:id` - Обновление компании
- `DELETE /api/companies/:id` - Удаление компании

### Транспорт
- `GET /api/vehicles` - Список транспорта
- `POST /api/vehicles` - Добавление транспорта
- `PUT /api/vehicles/:id` - Обновление транспорта
- `DELETE /api/vehicles/:id` - Удаление транспорта

### Кредиты
- `GET /api/loans` - Список кредитов
- `POST /api/loans` - Создание кредита
- `PUT /api/loans/:id` - Обновление кредита
- `DELETE /api/loans/:id` - Удаление кредита

### Платежи
- `GET /api/payments` - Все платежи пользователя
- `GET /api/payments/loan/:loanId` - Платежи по кредиту
- `POST /api/payments` - Создание платежа

### Финансовые отчеты
- `GET /api/schedules/debt` - График долгов
- `GET /api/schedules/amortization` - Амортизационный график
- `GET /api/schedules/depreciation` - График амортизации активов

### Статистика
- `GET /api/stats/dashboard` - Статистика дашборда

### Пользователи
- `GET /api/users` - Список пользователей (текущий пользователь)
- `GET /api/users/profile` - Профиль текущего пользователя
- `GET /api/users/:id` - Получить пользователя по ID
- `PUT /api/users/:id` - Обновить данные пользователя
- `DELETE /api/users/:id` - Удалить пользователя (и все связанные данные)

## 🔒 Безопасность

- JWT токены для аутентификации
- bcrypt хеширование паролей  
- CORS защита
- Изоляция данных пользователей
- Валидация всех входящих данных

## 📁 Структура проекта

```
Business Schedule/
├── backend/                 # Go бэкенд
│   ├── config/             # Конфигурация
│   ├── database/           # Подключение к БД
│   ├── handlers/           # HTTP обработчики
│   ├── middleware/         # Middleware (Auth, CORS)
│   ├── models/            # Модели данных
│   ├── routes/            # Маршруты API
│   ├── utils/             # Утилиты (JWT, Hash, Math)
│   └── main.go            # Точка входа
├── frontend/               # React фронтенд
│   ├── src/
│   │   ├── components/    # React компоненты
│   │   ├── context/       # React Context
│   │   ├── services/      # API сервисы
│   │   └── types/         # TypeScript типы
│   └── package.json
├── docker-compose.yml      # MongoDB в Docker
└── README.md
```

## 📚 Документация

- [Подробная установка](./SETUP.md) - детальные инструкции для всех ОС
- [Примеры API](./API_EXAMPLES.md) - полные примеры использования API

## 🎯 Roadmap

- [ ] Интеграция с внешними банковскими API
- [ ] Система уведомлений
- [ ] Экспорт отчетов в PDF/Excel
- [ ] Мобильное приложение
- [ ] Продвинутая аналитика и прогнозирование

## 🤝 Вклад в проект

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект распространяется под MIT лицензией. См. файл `LICENSE` для подробностей.

---

**Business Schedule** - профессиональная система управления траковым бизнесом 🚛💼 