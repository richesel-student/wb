# 🚀 Запуск

```bash
docker-compose up --build
```
# 🧪 Проверка
```bash
curl http://localhost:8080/api/health
```

# ➕ Создать событие
```bash
ID=$(curl -s -X POST http://localhost:8080/api/events \
  -d "name=test&date=2026-01-01&capacity=2&ttl=5" | jq -r .id)

echo $ID
```
# 📄 Получить событие
```bash
curl http://localhost:8080/api/events/$ID
```

# 🎟 Забронировать

```bash 
curl -X POST http://localhost:8080/api/events/$ID/book \
  -H "X-User-ID: user1"
  ```

# 🔁 Проверить места
```bash 
curl http://localhost:8080/api/events/$ID
```

# 💳 Подтвердить

```bash
curl -X POST http://localhost:8080/api/events/$ID/confirm \
  -H "X-User-ID: user1"
```