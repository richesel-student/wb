# 🚀 Crypto Analytics

Сервис для учёта крипто-сделок с аналитикой и графиком.

---

## ⚙️ Запуск

```bash
docker-compose up --build
```

**Открыть:**
[http://localhost:8080/web/](http://localhost:8080/web/)

---

## 📌 Возможности

* ➕ Добавление сделок (buy / sell)
* 📋 Просмотр списка
* 📊 Аналитика (sum, avg, median, p90)
* 📈 График
* 📅 Фильтр по дате
* 📥 Экспорт CSV

---

## 🧪 API

```bash
# создать
curl -X POST http://localhost:8080/items \
-H "Content-Type: application/json" \
-d '{"asset_name":"BTC","type":"buy","amount_usd":1000}'

# получить список
curl http://localhost:8080/items

# аналитика
curl http://localhost:8080/analytics

# график
curl http://localhost:8080/chart
```

---

## ⚠️ Примечание

Bybit API может возвращать `403`, поэтому используется fallback цена:

```go
price = 70000
```

---

## 🗄 БД

PostgreSQL (Docker) + dump.sql

---

## 👨‍💻 Автор

Rinat