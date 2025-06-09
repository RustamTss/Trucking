# Инструкции по установке и запуску Business Schedule

## Предварительные требования

1. **Go 1.21+** - для backend
2. **Node.js 18+** - для frontend  
3. **MongoDB 6.0+** - база данных

## Установка MongoDB

### macOS (с Homebrew):
```bash
brew tap mongodb/brew
brew install mongodb-community@6.0
brew services start mongodb/brew/mongodb-community
```

### Ubuntu/Debian:
```bash
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
```

### Windows:
Скачайте и установите MongoDB Community Server с официального сайта: https://www.mongodb.com/try/download/community

## Запуск приложения

### 1. Клонирование и настройка

```bash
# Переходим в директорию проекта
cd "Business Schedule"

# Устанавливаем зависимости backend
cd backend
go mod tidy

# Устанавливаем зависимости frontend
cd ../frontend
npm install
```

### 2. Настройка переменных окружения (опционально)

Создайте файл `backend/.env` для кастомных настроек:

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
DB_NAME=business_schedule
```

### 3. Запуск backend сервера

```bash
cd backend
go run main.go
```

Backend будет доступен по адресу: http://localhost:8080

### 4. Запуск frontend (в новом терминале)

```bash
cd frontend
npm run dev
```

Frontend будет доступен по адресу: http://localhost:5173

## Первый запуск

1. Откройте браузер и перейдите на http://localhost:5173
2. Вы увидите страницу входа
3. Поскольку это первый запуск, нужно зарегистрировать пользователя через API

### Регистрация пользователя через API:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "name": "Администратор"
  }'
```

После регистрации вы можете войти в систему используя эти учетные данные.

## Структура API

### Аутентификация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход

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
- `GET /api/payments/loan/:loanId` - Платежи по кредиту
- `POST /api/payments` - Запись платежа

### Финансовые отчеты
- `GET /api/schedules/debt` - График задолженности
- `GET /api/schedules/amortization` - График амортизации
- `GET /api/schedules/depreciation` - График амортизации активов

## Функциональность

### Реализованные возможности:
- ✅ Аутентификация пользователей
- ✅ Управление компаниями
- ✅ Управление транспортом (траки и трейлеры)
- ✅ Управление кредитами с автоматическим расчетом платежей
- ✅ Запись платежей с обновлением остатков
- ✅ Финансовые дашборды:
  - Business Debt Schedule
  - Amortization Schedule  
  - Depreciation Schedule
- ✅ Адаптивный интерфейс
- ✅ Переключение между компаниями

### Планируемые улучшения:
- 📋 Формы для добавления/редактирования данных
- 📊 Графики и диаграммы
- 📄 Экспорт в Excel/PDF
- 📱 Мобильная версия
- 🔔 Уведомления о платежах
- 🚛 Интеграция с GPS трекерами

## Устранение неполадок

### Backend не запускается:
1. Проверьте что MongoDB запущен: `mongosh` или `mongo`
2. Проверьте что порт 8080 свободен: `lsof -i :8080`
3. Проверьте логи в терминале

### Frontend не запускается:
1. Проверьте что Node.js установлен: `node --version`
2. Переустановите зависимости: `rm -rf node_modules && npm install`
3. Проверьте что порт 5173 свободен

### Ошибки подключения к API:
1. Убедитесь что backend запущен на порту 8080
2. Проверьте CORS настройки в `backend/main.go`
3. Проверьте URL API в `frontend/src/services/api.ts`

## Разработка

### Добавление новых функций:
1. Backend: добавьте модели в `models/`, handlers в `handlers/`, маршруты в `routes/`
2. Frontend: добавьте типы в `types/`, API методы в `services/`, компоненты в `components/`

### Тестирование API:
Используйте Postman или curl для тестирования API endpoints. Не забудьте добавить JWT токен в заголовок Authorization. 