# Примеры API запросов для Business Schedule

## Аутентификация

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "name": "Администратор"
  }'
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

**Ответ содержит JWT токен, который нужно использовать в заголовке Authorization для всех последующих запросов:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "...",
    "email": "admin@example.com",
    "name": "Администратор"
  }
}
```

## Управление компаниями

### Создание компании
```bash
curl -X POST http://localhost:8080/api/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "ABC Trucking LLC",
    "ein": "12-3456789",
    "address": "123 Main St, City, State 12345"
  }'
```

### Получение списка компаний
```bash
curl -X GET http://localhost:8080/api/companies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Обновление компании
```bash
curl -X PUT http://localhost:8080/api/companies/COMPANY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "ABC Trucking LLC (Updated)",
    "ein": "12-3456789",
    "address": "456 New St, City, State 12345"
  }'
```

## Управление транспортом

### Добавление трака
```bash
curl -X POST http://localhost:8080/api/vehicles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "company_id": "COMPANY_ID",
    "type": "truck",
    "vin": "1HGBH41JXMN109186",
    "make": "Freightliner",
    "model": "Cascadia",
    "year": 2020,
    "purchase_price": 120000,
    "purchase_date": "2020-01-15T00:00:00Z",
    "status": "active"
  }'
```

### Добавление трейлера
```bash
curl -X POST http://localhost:8080/api/vehicles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "company_id": "COMPANY_ID",
    "type": "trailer",
    "vin": "1M2AA18C1WM123456",
    "make": "Great Dane",
    "model": "Everest",
    "year": 2019,
    "purchase_price": 45000,
    "purchase_date": "2019-06-10T00:00:00Z",
    "status": "active"
  }'
```

### Получение списка транспорта
```bash
# Весь транспорт
curl -X GET http://localhost:8080/api/vehicles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Транспорт конкретной компании
curl -X GET "http://localhost:8080/api/vehicles?company_id=COMPANY_ID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Управление кредитами

### Создание кредита
```bash
curl -X POST http://localhost:8080/api/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "vehicle_id": "VEHICLE_ID",
    "company_id": "COMPANY_ID",
    "lender": "Bank of America",
    "principal_amount": 100000,
    "interest_rate": 5.5,
    "term_months": 60,
    "start_date": "2023-01-01T00:00:00Z",
    "status": "active"
  }'
```

### Получение списка кредитов
```bash
# Все кредиты
curl -X GET http://localhost:8080/api/loans \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Кредиты конкретной компании
curl -X GET "http://localhost:8080/api/loans?company_id=COMPANY_ID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Управление платежами

### Запись платежа
```bash
curl -X POST http://localhost:8080/api/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "loan_id": "LOAN_ID",
    "payment_date": "2023-02-01T00:00:00Z",
    "total_paid": 1887.12
  }'
```

### Получение платежей по кредиту
```bash
curl -X GET http://localhost:8080/api/payments/loan/LOAN_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Финансовые отчеты

### График задолженности (Debt Schedule)
```bash
curl -X GET http://localhost:8080/api/schedules/debt \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Пример ответа:**
```json
[
  {
    "company_name": "ABC Trucking LLC",
    "total_debt": 95000.50,
    "monthly_payment": 1887.12,
    "vehicles_count": 3
  }
]
```

### График амортизации (Amortization Schedule)
```bash
# Для всех кредитов
curl -X GET http://localhost:8080/api/schedules/amortization \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Для конкретного кредита
curl -X GET "http://localhost:8080/api/schedules/amortization?loan_id=LOAN_ID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### График амортизации активов (Depreciation Schedule)
```bash
# Для всех компаний
curl -X GET http://localhost:8080/api/schedules/depreciation \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Для конкретной компании
curl -X GET "http://localhost:8080/api/schedules/depreciation?company_id=COMPANY_ID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Полный пример workflow

### 1. Регистрация и вход
```bash
# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123", "name": "Test User"}'

# Вход (сохраните токен из ответа)
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}' | jq -r '.token')
```

### 2. Создание компании
```bash
COMPANY_ID=$(curl -s -X POST http://localhost:8080/api/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Test Trucking", "ein": "12-3456789", "address": "123 Test St"}' | jq -r '.id')
```

### 3. Добавление транспорта
```bash
VEHICLE_ID=$(curl -s -X POST http://localhost:8080/api/vehicles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"company_id\": \"$COMPANY_ID\", \"type\": \"truck\", \"vin\": \"TEST123456789\", \"make\": \"Freightliner\", \"model\": \"Cascadia\", \"year\": 2020, \"purchase_price\": 120000, \"purchase_date\": \"2020-01-15T00:00:00Z\", \"status\": \"active\"}" | jq -r '.id')
```

### 4. Создание кредита
```bash
LOAN_ID=$(curl -s -X POST http://localhost:8080/api/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"vehicle_id\": \"$VEHICLE_ID\", \"company_id\": \"$COMPANY_ID\", \"lender\": \"Test Bank\", \"principal_amount\": 100000, \"interest_rate\": 5.5, \"term_months\": 60, \"start_date\": \"2023-01-01T00:00:00Z\", \"status\": \"active\"}" | jq -r '.id')
```

### 5. Запись платежа
```bash
curl -X POST http://localhost:8080/api/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"loan_id\": \"$LOAN_ID\", \"payment_date\": \"2023-02-01T00:00:00Z\", \"total_paid\": 1887.12}"
```

### 6. Получение отчетов
```bash
# График задолженности
curl -X GET http://localhost:8080/api/schedules/debt \
  -H "Authorization: Bearer $TOKEN"

# График амортизации
curl -X GET http://localhost:8080/api/schedules/amortization \
  -H "Authorization: Bearer $TOKEN"

# График амортизации активов
curl -X GET http://localhost:8080/api/schedules/depreciation \
  -H "Authorization: Bearer $TOKEN"
```

## Коды ошибок

- `400` - Неверный формат данных
- `401` - Не авторизован (неверный или отсутствующий токен)
- `403` - Нет доступа к ресурсу
- `404` - Ресурс не найден
- `500` - Внутренняя ошибка сервера

## Примечания

1. Замените `YOUR_JWT_TOKEN` на реальный токен, полученный при входе
2. Замените `COMPANY_ID`, `VEHICLE_ID`, `LOAN_ID` на реальные ID из ответов API
3. Все даты должны быть в формате ISO 8601 (RFC 3339)
4. Суммы указываются в долларах США
5. Процентные ставки указываются в процентах (например, 5.5 для 5.5%) 